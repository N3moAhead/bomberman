package logger

import (
	"fmt"
	"log"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorViolet = "\033[35m"
	colorCyan   = "\033[36m"
)

// Logger is a struct that can be configured with a prefix
type Logger struct {
	prefix string
}

// New creates a new Logger instance with a given prefix
// The prefix will be printed in violet with every log message from this instance
func New(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) print(level, color, format string, a ...any) {
	prefixStr := ""
	if l.prefix != "" {
		prefixStr = fmt.Sprintf("%s%s%s ", colorViolet, l.prefix, colorReset)
	}

	levelTag := fmt.Sprintf("%s[%-7s]%s", color, level, colorReset)
	message := fmt.Sprintf(format, a...)

	log.Printf("%s %s%s", levelTag, prefixStr, message)
}

// Success logs a message with a green 'SUCCESS' level
func (l *Logger) Success(format string, a ...any) {
	l.print("SUCCESS", colorGreen, format, a...)
}

// Error logs a message with a red 'ERROR' level
func (l *Logger) Error(format string, a ...any) {
	l.print("ERROR", colorRed, format, a...)
}

// Warn logs a message with a yellow 'WARN' level
func (l *Logger) Warn(format string, a ...any) {
	l.print("WARN", colorYellow, format, a...)
}

// Info logs a message with a blue 'INFO' level
func (l *Logger) Info(format string, a ...any) {
	l.print("INFO", colorBlue, format, a...)
}

// Debug logs a message with a cyan 'DEBUG' level
func (l *Logger) Debug(format string, a ...any) {
	l.print("DEBUG", colorCyan, format, a...)
}

// --- Global Logger ---

// std is the default logger instance, used by the global functions
var std = New("")

func init() {
	// Remove default timestamp and prefix from the standard logger,
	// as we are handling formatting ourselves
	log.SetFlags(0)
}

// Success logs a message using the default global logger
func Success(format string, a ...any) {
	std.Success(format, a...)
}

// Error logs a message using the default global logger
func Error(format string, a ...any) {
	std.Error(format, a...)
}

// Warn logs a message using the default global logger
func Warn(format string, a ...any) {
	std.Warn(format, a...)
}

// Info logs a message using the default global logger
func Info(format string, a ...any) {
	std.Info(format, a...)
}

// Debug logs a message using the default global logger
func Debug(format string, a ...any) {
	std.Debug(format, a...)
}
