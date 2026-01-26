package bomber

import (
	"fmt"
	"log"
)

const (
	logPrefix = "[BOMBER] "
	// ANSI color codes
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBlue  = "\033[34m"
	colorCyan  = "\033[36m"
)

func green(format string, v ...any) string {
	message := fmt.Sprintf(format, v...)
	return fmt.Sprintf("%s%s%s", colorGreen, message, colorReset)
}

func blue(format string, v ...any) string {
	message := fmt.Sprintf(format, v...)
	return fmt.Sprintf("%s%s%s", colorBlue, message, colorReset)
}

func red(format string, v ...any) string {
	message := fmt.Sprintf(format, v...)
	return fmt.Sprintf("%s%s%s", colorRed, message, colorReset)
}

// success logs a formatted success message in green
func success(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Printf("%s%s%s%s\n", logPrefix, colorGreen, message, colorReset)
}

// error logs a formatted error message in red
func error(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Printf("%s%s%s%s\n", logPrefix, colorRed, message, colorReset)
}

// info logs a formatted info message in blue
func info(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Printf("%s%s%s%s\n", logPrefix, colorBlue, message, colorReset)
}

// debug logs a formatted debug message in cyan
func debug(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Printf("%s%s%s%s\n", logPrefix, colorCyan, message, colorReset)
}
