//go:build !windows

package main

// initSystemTray is a no-op on non-Windows platforms
func initSystemTray() {}
