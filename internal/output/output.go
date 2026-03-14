package output

import (
	"fmt"
	"io"
	"os"
)

// PrintValue writes the value to w followed by a newline.
// Used for successful command output (stdout).
func PrintValue(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}

// Err writes the error message to w (typically stderr).
func Err(w io.Writer, msg string) {
	fmt.Fprintln(w, msg)
}

// Exit exits the process with the given code.
func Exit(code int) {
	os.Exit(code)
}
