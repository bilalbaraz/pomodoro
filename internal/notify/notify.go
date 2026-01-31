package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Notify sends a macOS notification using osascript (macOS only).
func Notify(title, message string) {
	if runtime.GOOS != "darwin" {
		return
	}

	safeTitle := escapeAppleScript(title)
	safeMessage := escapeAppleScript(message)
	script := fmt.Sprintf(
		`display dialog "%s" with title "%s" buttons {"OK"} default button "OK"`,
		safeMessage,
		safeTitle,
	)

	cmd := exec.Command("osascript", "-e", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		if len(out) > 0 {
			fmt.Fprintf(os.Stderr, "notify error: %s\n", strings.TrimSpace(string(out)))
		}
	}
}

func escapeAppleScript(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
