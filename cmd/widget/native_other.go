//go:build !windows

package main

// launchNativeWidget is a no-op on non-Windows platforms.
// Use browser mode instead: run with --browser flag or it will auto-fallback.
func launchNativeWidget(intervalSec int) bool {
	return false
}
