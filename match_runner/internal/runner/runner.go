package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/N3moAhead/bomberman/match_runner/internal/history"
	"github.com/N3moAhead/bomberman/match_runner/internal/match"
	"github.com/N3moAhead/bomberman/match_runner/pkg/logger"
	"github.com/google/uuid"
)

var log = logger.New("[Runner]")

// Runner handles the execution of a single match.
type Runner struct{}

// New creates a new Runner.
func New() *Runner {
	return &Runner{}
}

// RunMatch executes a full match lifecycle: creates a pod,
// runs containers, waits for completion, and cleans up
func (r *Runner) RunMatch(ctx context.Context, details *match.Details) (*match.Result, error) {
	podName := fmt.Sprintf("bomberman-match-%s", details.MatchID)
	log.Info("Starting match %s in pod %s", details.MatchID, podName)

	client1AuthToken := uuid.NewString()
	client2AuthToken := uuid.NewString()

	// Cleanup is deferred to ensure it runs even if errors occur
	defer r.cleanupPod(podName)

	if err := r.createPod(ctx, podName); err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}

	serverContainerName := "server"
	if err := r.runServer(ctx, podName, serverContainerName, details.ServerImage); err != nil {
		return nil, fmt.Errorf("failed to run server: %w", err)
	}

	// Run clients concurrently
	clientErrCh := make(chan error, 2)
	go func() {
		clientErrCh <- r.runClient(ctx, podName, "client1", details.Client1Image, client1AuthToken)
	}()
	go func() {
		clientErrCh <- r.runClient(ctx, podName, "client2", details.Client2Image, client2AuthToken)
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

	serverLogs, err := r.getContainerLogs(context.Background(), serverContainerName)
	if err != nil {
		log.Warn("Could not get server logs after match completion: %v", err)
	} else {
		// Try to parse the game history from the logs
		gameHistory, err := parseGameHistory(serverLogs)
		if err != nil {
			log.Warn("Failed to parse game history from server logs: %v", err)
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
			log.Success("Successfully parsed game history with %d ticks.", len(gameHistory.Ticks))
		}
	}

	go r.removeImage(context.Background(), details.Client1Image)
	go r.removeImage(context.Background(), details.Client2Image)

	return result, nil
}

func (r *Runner) createPod(ctx context.Context, podName string) error {
	log.Debug("Creating pod: %s", podName)
	// TODO Rething if --network ="host" is the correct decision here...
	cmd := exec.CommandContext(ctx, "podman", "pod", "create", "--name", podName, "--network=host")
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
		return fmt.Errorf("podman pull of '%s' failed: %s: %w", image, string(output), err)
	}
	return nil
}

func (r *Runner) runServer(ctx context.Context, podName, containerName, image string) error {
	log.Info("Starting server container '%s' with image '%s'", containerName, image)
	if err := r.pullImage(ctx, image); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "podman", "run", "--pod", podName, "--name", containerName, "--detach", image)
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
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("podman wait for '%s' failed: %w", containerName, err)
	}
	return nil
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

func (r *Runner) cleanupPod(podName string) {
	// Use a background context with a timeout for cleanup,
	// as the original match context might have been cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

func parseGameHistory(logs string) (*history.GameHistory, error) {
	const prefix = "GameHistory:"
	scanner := bufio.NewScanner(strings.NewReader(logs))
	for scanner.Scan() {
		line := scanner.Text()
		if after, ok := strings.CutPrefix(line, prefix); ok {
			jsonBody := after
			var gameHistory history.GameHistory
			if err := json.Unmarshal([]byte(jsonBody), &gameHistory); err != nil {
				return nil, fmt.Errorf("failed to unmarshal game history JSON: %w", err)
			}
			return &gameHistory, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading server logs: %w", err)
	}

	return nil, fmt.Errorf("game history prefix not found in server logs")
}
