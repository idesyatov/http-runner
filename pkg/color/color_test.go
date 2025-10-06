package color

import (
	"testing"
)

// TestColorize checks that the Colorize function applies the ANSI escape codes correctly.
func TestColorize(t *testing.T) {
	tests := []struct {
		color  string
		input  string
		output string
	}{
		{Green, "Hello, World!", "\033[32mHello, World!\033[0m"},
		{Red, "Error!", "\033[31mError!\033[0m"},
		{Blue, "Info:", "\033[34mInfo:\033[0m"},
		{Yellow, "Warning!", "\033[33mWarning!\033[0m"},
		{Cyan, "Cyan text", "\033[36mCyan text\033[0m"},
		{Magenta, "Magenta text", "\033[35mMagenta text\033[0m"},
		{White, "White text", "\033[37mWhite text\033[0m"},
	}

	for _, tt := range tests {
		result := Colorize(tt.color, tt.input)
		if result != tt.output {
			t.Errorf("Colorize(%q, %q) = %q, want %q", tt.color, tt.input, result, tt.output)
		}
	}
}
