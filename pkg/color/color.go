package color

import "fmt"

// ANSI escape codes for text colors
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

// Colorize formats the input value with the specified color.
func Colorize[T any](color string, value T) string {
	return color + fmt.Sprint(value) + Reset
}
