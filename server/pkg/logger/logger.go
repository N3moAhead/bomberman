package logger

import (
	"fmt"
	"log"
	"os"
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

// --- Printf Style ---

func (l *Logger) print(level, color, format string, a ...any) {
	prefixStr := ""
	if l.prefix != "" {
		prefixStr = fmt.Sprintf("%s%s%s ", colorViolet, l.prefix, colorReset)
	}

	levelTag := fmt.Sprintf("%s[%-7s]%s", color, level, colorReset) // Padded for alignment
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

// --- Println Style ---

func (l *Logger) println(level, color string, a ...any) {
	prefixStr := ""
	if l.prefix != "" {
		prefixStr = fmt.Sprintf("%s%s%s ", colorViolet, l.prefix, colorReset)
	}

	levelTag := fmt.Sprintf("%s[%-7s]%s", color, level, colorReset) // Padded for alignment
	message := fmt.Sprint(a...)

	log.Printf("%s %s%s", levelTag, prefixStr, message)
}

// Successln logs a message with a green 'SUCCESS' level
func (l *Logger) Successln(a ...any) {
	l.println("SUCCESS", colorGreen, a...)
}

// Errorln logs a message with a red 'ERROR' level
func (l *Logger) Errorln(a ...any) {
	l.println("ERROR", colorRed, a...)
}

// Warnln logs a message with a yellow 'WARN' level
func (l *Logger) Warnln(a ...any) {
	l.println("WARN", colorYellow, a...)
}

// Infoln logs a message with a blue 'INFO' level
func (l *Logger) Infoln(a ...any) {
	l.println("INFO", colorBlue, a...)
}

// Debugln logs a message with a cyan 'DEBUG' level
func (l *Logger) Debugln(a ...any) {
	l.println("DEBUG", colorCyan, a...)
}

func (l *Logger) Fatal(a ...any) {
	l.println("ERROR", colorRed, a...)
	os.Exit(1)
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

// Successln logs a message using the default global logger
func Successln(a ...any) {
	std.Successln(a...)
}

// Errorln logs a message using the default global logger
func Errorln(a ...any) {
	std.Errorln(a...)
}

// Warnln logs a message using the default global logger
func Warnln(a ...any) {
	std.Warnln(a...)
}

// Infoln logs a message using the default global logger
func Infoln(a ...any) {
	std.Infoln(a...)
}

// Debugln logs a message using the default global logger
func Debugln(a ...any) {
	std.Debugln(a...)
}

func Fatal(a ...any) {
	std.Errorln(a...)
}
