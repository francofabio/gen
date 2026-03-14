package clipboard

import (
	"bytes"
	"os/exec"
	"runtime"
)

// Write copies text to the system clipboard.
// Uses pbcopy (macOS), clip (Windows), or xclip/xsel/wl-copy (Linux).
// No special permissions are required; it uses the same mechanism as piping to pbcopy.
// Returns an error if the platform command is not available or fails.
func Write(text string) error {
	switch runtime.GOOS {
	case "darwin":
		return writeDarwin(text)
	case "windows":
		return writeWindows(text)
	default:
		return writeLinux(text)
	}
}

func writeDarwin(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}

func writeWindows(text string) error {
	// clip.exe reads from stdin
	cmd := exec.Command("cmd", "/c", "clip")
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}

func writeLinux(text string) error {
	// Try wl-copy (Wayland), then xclip, then xsel
	tries := []struct {
		name string
		args []string
	}{
		{"wl-copy", nil},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
	}
	var lastErr error
	for _, t := range tries {
		var cmd *exec.Cmd
		if len(t.args) == 0 {
			cmd = exec.Command(t.name)
		} else {
			cmd = exec.Command(t.name, t.args...)
		}
		cmd.Stdin = bytes.NewBufferString(text)
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return exec.ErrNotFound
}
