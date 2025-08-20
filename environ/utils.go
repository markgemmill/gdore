package environ

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenFile opens a file with the default
// application based on the operating system
func OpenFile(filepath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filepath)
	case "darwin": // macOS
		cmd = exec.Command("open", filepath)
	case "linux":
		cmd = exec.Command("xdg-open", filepath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
