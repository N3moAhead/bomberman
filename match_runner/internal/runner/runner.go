package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/N3moAhead/bombahead/match_runner/internal/history"
	"github.com/N3moAhead/bombahead/match_runner/internal/match"
	"github.com/N3moAhead/bombahead/match_runner/pkg/logger"
	"github.com/google/uuid"
)

var log = logger.New("[Runner]")

const (
	errCodeImagePullRateLimit = "ERR_IMAGE_PULL_RATE_LIMIT"
	errCodeImageNotPullable   = "ERR_IMAGE_NOT_PULLABLE"
)

// Runner handles the execution of a single match.
type Runner struct{}

// New creates a new Runner.
func New() *Runner {
	return &Runner{}
}

// RunMatch executes a full match lifecycle: creates a pod,
// runs containers, waits for completion, and cleans up
func (r *Runner) RunMatch(ctx context.Context, details *match.Details) (*match.Result, error) {
	runID := uuid.NewString()[:8]
	podName := fmt.Sprintf("bomberman-match-%s-%s", details.MatchID, runID)
	serverContainerName := fmt.Sprintf("%s-server", podName)
	client1ContainerName := fmt.Sprintf("%s-client1", podName)
	client2ContainerName := fmt.Sprintf("%s-client2", podName)

	log.Info("Starting match %s in pod %s", details.MatchID, podName)

	client1AuthToken := uuid.NewString()
	client2AuthToken := uuid.NewString()

	historyDir := os.Getenv("MATCH_HISTORY_DIR")
	if historyDir != "" {
		if err := os.MkdirAll(historyDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create history directory '%s': %w", historyDir, err)
		}
	}
	historyFile, err := os.CreateTemp(historyDir, "bombahead-match-history-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary history file: %w", err)
	}
	historyFilePath := historyFile.Name()
	if closeErr := historyFile.Close(); closeErr != nil {
		return nil, fmt.Errorf("failed to close temporary history file: %w", closeErr)
	}
	// Allow the container process to write to the bind-mounted file.
	if chmodErr := os.Chmod(historyFilePath, 0666); chmodErr != nil {
		return nil, fmt.Errorf("failed to set permissions on temporary history file: %w", chmodErr)
	}
	defer func() {
		if removeErr := os.Remove(historyFilePath); removeErr != nil && !os.IsNotExist(removeErr) {
			log.Warn("Failed to remove temporary history file '%s': %v", historyFilePath, removeErr)
		}
	}()

	// Ensure no stale resources from previous runs can interfere with this match.
	r.cleanupResources(context.Background(), podName, serverContainerName, client1ContainerName, client2ContainerName)

	// Cleanup is deferred to ensure it runs even if errors occur
	defer r.cleanupResources(context.Background(), podName, serverContainerName, client1ContainerName, client2ContainerName)

	if err := r.createPod(ctx, podName); err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}

	if err := r.runServer(ctx, podName, serverContainerName, details.ServerImage, historyFilePath); err != nil {
		return nil, fmt.Errorf("failed to run server: %w", err)
	}

	// Run clients concurrently
	clientErrCh := make(chan error, 2)
	go func() {
		clientErrCh <- r.runClient(ctx, podName, client1ContainerName, details.Client1Image, client1AuthToken)
	}()
	go func() {
		clientErrCh <- r.runClient(ctx, podName, client2ContainerName, details.Client2Image, client2AuthToken)
	}()

	for range 2 {
		if err := <-clientErrCh; err != nil {
			return nil, fmt.Errorf("failed to run a client: %w", err)
		}
	}

	log.Info("All containers started for match %s. Waiting for server to complete...", details.MatchID)

	if err := r.waitForContainer(ctx, serverContainerName); err != nil {
		// Attempt to get logs even if wait fails, as they might contain error info
		serverLogs, _ := r.getContainerLogs(context.Background(), serverContainerName)
		log.Error("Server logs on wait error: %s", serverLogs)
		return nil, fmt.Errorf("error waiting for server container: %w", err)
	}

	log.Info("Server container exited. Match %s finished.", details.MatchID)

	result := &match.Result{
		MatchID:       details.MatchID,
		Winner:        "",
		Client1GameID: client1AuthToken,
		Client2GameID: client2AuthToken,
	}

	gameHistory, err := r.readGameHistoryFromFile(historyFilePath)
	if err != nil {
		log.Warn("Failed to read game history from file '%s': %v", historyFilePath, err)
	} else {
		if gameHistory.WinnerAuthToken != "" {
			switch gameHistory.WinnerAuthToken {
			case client1AuthToken:
				result.Winner = details.Client1Image
			case client2AuthToken:
				result.Winner = details.Client2Image
			}
		}
		result.Log = gameHistory
		log.Success("Successfully read game history with %d ticks from file.", len(gameHistory.Ticks))
	}

	go r.removeImage(context.Background(), details.Client1Image)
	go r.removeImage(context.Background(), details.Client2Image)

	return result, nil
}

func (r *Runner) createPod(ctx context.Context, podName string) error {
	log.Debug("Creating pod: %s", podName)
	cmd := exec.CommandContext(ctx, "podman", "pod", "create", "--name", podName, "--network=none")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("podman pod create failed: %s: %w", string(output), err)
	}
	log.Success("Pod '%s' created successfully.", podName)
	return nil
}

func (r *Runner) pullImage(ctx context.Context, image string) error {
	log.Info("Pulling image: %s", image)
	cmd := exec.CommandContext(ctx, "podman", "pull", image)
	if output, err := cmd.CombinedOutput(); err != nil {
		rawOutput := strings.TrimSpace(string(output))
		if code := classifyImagePullError(rawOutput); code != "" {
			return fmt.Errorf("%s: podman pull of '%s' failed: %s: %w", code, image, rawOutput, err)
		}
		return fmt.Errorf("podman pull of '%s' failed: %s: %w", image, rawOutput, err)
	}
	return nil
}

func (r *Runner) runServer(ctx context.Context, podName, containerName, image, hostHistoryFilePath string) error {
	log.Info("Starting server container '%s' with image '%s'", containerName, image)
	if err := r.pullImage(ctx, image); err != nil {
		return err
	}
	cmd := exec.CommandContext(
		ctx,
		"podman",
		"run",
		"--pod", podName,
		"--name", containerName,
		"--detach",
		"--mount", fmt.Sprintf("type=bind,src=%s,dst=/tmp/match-history.json,relabel=shared", hostHistoryFilePath),
		"--env", "BOMBERMAN_MATCH_HISTORY_PATH=/tmp/match-history.json",
		image,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("podman run (server) failed: %s: %w", string(output), err)
	}
	log.Success("Server container '%s' started.", containerName)
	return nil
}

func (r *Runner) runClient(ctx context.Context, podName, containerName, image, clientAuthToken string) error {
	log.Info("Starting client container '%s' with image '%s'", containerName, image)
	if err := r.pullImage(ctx, image); err != nil {
		return err
	}
	// Secure the client containers
	cmd := exec.CommandContext(ctx, "podman", "run", "--pod", podName, "--name", containerName,
		"--detach",
		"--cap-drop=all",
		"--security-opt=no-new-privileges",
		"--memory=500m", // Yeah well so that shall be my security haha let's see how far this will bring me XD
		"--cpus=0.5",
		"--env", "BOMBERMAN_CLIENT_AUTH_TOKEN="+clientAuthToken,
		image)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("podman run (client: %s) failed: %s: %w", containerName, string(output), err)
	}
	log.Success("Client container '%s' started.", containerName)
	return nil
}

func (r *Runner) waitForContainer(ctx context.Context, containerName string) error {
	log.Debug("Waiting for container '%s' to stop...", containerName)
	cmd := exec.CommandContext(ctx, "podman", "wait", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("podman wait for '%s' failed: %w", containerName, err)
	}

	// podman wait may include extra whitespace/lines; parse the first token safely.
	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		exitCode, inspectErr := r.inspectContainerExitCode(ctx, containerName)
		if inspectErr != nil {
			return fmt.Errorf("podman wait for '%s' returned empty output and inspect failed: %w", containerName, inspectErr)
		}
		if exitCode != 0 {
			return fmt.Errorf("container '%s' exited with code %d", containerName, exitCode)
		}
		return nil
	}

	exitCode, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("podman wait for '%s' returned non-integer exit code token %q (raw: %q): %w", containerName, fields[0], strings.TrimSpace(string(output)), err)
	}
	if exitCode != 0 {
		return fmt.Errorf("container '%s' exited with code %d", containerName, exitCode)
	}
	return nil
}

func (r *Runner) inspectContainerExitCode(ctx context.Context, containerName string) (int, error) {
	cmd := exec.CommandContext(ctx, "podman", "inspect", "--format", "{{.State.ExitCode}}", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("podman inspect for '%s' failed: %s: %w", containerName, strings.TrimSpace(string(output)), err)
	}

	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		return 0, fmt.Errorf("podman inspect for '%s' returned empty output", containerName)
	}

	exitCode, parseErr := strconv.Atoi(fields[0])
	if parseErr != nil {
		return 0, fmt.Errorf("podman inspect for '%s' returned non-integer exit code token %q", containerName, fields[0])
	}
	return exitCode, nil
}

func (r *Runner) getContainerLogs(ctx context.Context, containerName string) (string, error) {
	log.Debug("Getting logs for container '%s'", containerName)
	cmd := exec.CommandContext(ctx, "podman", "logs", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("podman logs for '%s' failed: %w", containerName, err)
	}
	return string(output), nil
}

func (r *Runner) removeImage(ctx context.Context, image string) {
	log.Info("Attempting to remove image: %s", image)
	cmd := exec.CommandContext(ctx, "podman", "rmi", "--force", image)
	if err := cmd.Run(); err != nil {
		log.Warn("Failed to remove image '%s' (this may not be an error): %v", image, err)
	} else {
		log.Success("Successfully removed image '%s'", image)
	}
}

func (r *Runner) cleanupResources(ctx context.Context, podName string, containerNames ...string) {
	// Use a timeout for cleanup operations so we don't block forever on a bad engine state.
	cleanupCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for _, containerName := range containerNames {
		r.forceRemoveContainer(cleanupCtx, containerName)
	}
	r.cleanupPod(cleanupCtx, podName)
}

func (r *Runner) forceRemoveContainer(ctx context.Context, containerName string) {
	log.Info("Ensuring container '%s' is removed...", containerName)
	cmd := exec.CommandContext(ctx, "podman", "rm", "-f", containerName)
	if output, err := cmd.CombinedOutput(); err != nil {
		combined := strings.TrimSpace(string(output))
		if strings.Contains(combined, "no container with name") || strings.Contains(combined, "no such container") {
			log.Debug("Container '%s' does not exist, nothing to remove.", containerName)
			return
		}

		// Ignore cancellations/timeouts from parent context during best-effort cleanup.
		if errorsIsContextDone(err) {
			log.Warn("Container cleanup for '%s' stopped by context: %v", containerName, err)
			return
		}

		log.Warn("Failed to force remove container '%s': %s: %v", containerName, string(output), err)
		return
	}
	log.Success("Container '%s' removed.", containerName)
}

func (r *Runner) cleanupPod(ctx context.Context, podName string) {
	// Use a background context with a timeout for cleanup,
	// as the original match context might have been cancelled
	log.Info("Cleaning up resources for pod '%s'", podName)
	cmd := exec.CommandContext(ctx, "podman", "pod", "exists", podName)
	if err := cmd.Run(); err != nil {
		// Pod does not exist, nothing to clean up.
		log.Info("Pod '%s' does not exist, no cleanup needed.", podName)
		return
	}

	log.Info("Stopping and removing pod '%s'...", podName)
	rmCmd := exec.CommandContext(ctx, "podman", "pod", "rm", "-f", podName)
	if output, err := rmCmd.CombinedOutput(); err != nil {
		log.Warn("Failed to remove pod '%s': %s: %v", podName, string(output), err)
	} else {
		log.Success("Successfully removed pod '%s'", podName)
	}
}

func errorsIsContextDone(err error) bool {
	return err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded))
}

func (r *Runner) readGameHistoryFromFile(filePath string) (*history.GameHistory, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return nil, fmt.Errorf("history file is empty")
	}

	var gameHistory history.GameHistory
	if err := json.Unmarshal([]byte(trimmed), &gameHistory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal history JSON: %w", err)
	}

	return &gameHistory, nil
}

func classifyImagePullError(output string) string {
	out := strings.ToLower(output)

	// Docker Hub / OCI registry rate limit signatures.
	if strings.Contains(out, "toomanyrequests") ||
		strings.Contains(out, "pull rate limit") ||
		strings.Contains(out, "you have reached your unauthenticated pull rate limit") ||
		strings.Contains(out, "too many requests") {
		return errCodeImagePullRateLimit
	}

	// Common signatures for non-pullable images (missing/private/invalid reference).
	if strings.Contains(out, "manifest unknown") ||
		strings.Contains(out, "not found") ||
		strings.Contains(out, "name unknown") ||
		strings.Contains(out, "pull access denied") ||
		strings.Contains(out, "requested access to the resource is denied") ||
		strings.Contains(out, "repository does not exist") ||
		strings.Contains(out, "insufficient_scope") {
		return errCodeImageNotPullable
	}

	return ""
}
