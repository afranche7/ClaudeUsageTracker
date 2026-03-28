//go:build windows

package main

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	procFindWindow    = user32.NewProc("FindWindowW")
	procSetWindowPos  = user32.NewProc("SetWindowPos")
	procSetWindowLong = user32.NewProc("SetWindowLongW")
	procGetWindowLong = user32.NewProc("GetWindowLongW")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
)

const (
	swpNoSize     = 0x0001
	swpNoMove     = 0x0002
	swpShowWindow = 0x0040
	hwndTopMost   = ^uintptr(0) // -1
	gwlExStyle    = -20
	wsExLayered   = 0x00080000
	wsExToolWindow = 0x00000080
	wsExTopmost   = 0x00000008
	lwaAlpha      = 0x02
)

// launchNativeWidget opens the widget in the default browser and attempts to
// make the browser window always-on-top and semi-transparent using Win32 APIs.
// For the best experience, the browser-mode widget auto-refreshes its content.
// Returns false to fall back to standard browser mode if window manipulation fails.
func launchNativeWidget(intervalSec int) bool {
	// On Windows, we launch in browser mode but enhance the experience by
	// opening a minimal Chrome/Edge app-mode window
	port := 17429

	// Try to launch Edge in app mode (creates a minimal, chrome-less window)
	edgePath := findEdgePath()
	if edgePath != "" {
		url := fmt.Sprintf("http://127.0.0.1:%d", port)
		cmd := exec.Command(edgePath,
			"--app="+url,
			"--window-size=300,400",
			"--window-position=50,50",
			"--disable-extensions",
			"--disable-sync",
			"--no-first-run",
		)
		if err := cmd.Start(); err == nil {
			log.Printf("Launched Edge app-mode window at %s", url)
			return false // Still need to start the HTTP server, so return false
		}
	}

	return false // Fall back to browser mode
}

func findEdgePath() string {
	paths := []string{
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
		`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
	}
	for _, p := range paths {
		if _, err := syscall.UTF16PtrFromString(p); err == nil {
			// Check if file exists
			pathPtr, _ := syscall.UTF16PtrFromString(p)
			attrs, err := syscall.GetFileAttributes(pathPtr)
			if err == nil && attrs != syscall.INVALID_FILE_ATTRIBUTES {
				return p
			}
		}
	}
	return ""
}

// makeWindowTopMostAndTransparent finds a window by title and makes it always-on-top
// with the specified opacity (0-255). This can be called after the browser window opens.
func makeWindowTopMostAndTransparent(title string, opacity byte) error {
	titlePtr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}

	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return fmt.Errorf("window not found: %s", title)
	}

	// Add WS_EX_LAYERED and WS_EX_TOOLWINDOW styles
	gwlIdx := uintptr(0xFFFFFFEC) // GWL_EXSTYLE = -20 as unsigned
	style, _, _ := procGetWindowLong.Call(hwnd, gwlIdx)
	newStyle := style | wsExLayered | wsExToolWindow
	procSetWindowLong.Call(hwnd, gwlIdx, newStyle)

	// Set transparency
	procSetLayeredWindowAttributes.Call(hwnd, 0, uintptr(opacity), lwaAlpha)

	// Set always-on-top
	procSetWindowPos.Call(hwnd, hwndTopMost, 0, 0, 0, 0, swpNoSize|swpNoMove|swpShowWindow)

	return nil
}
