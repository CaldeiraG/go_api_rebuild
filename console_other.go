//go:build !windows

package main

// disableConsoleQuickEdit is a no-op on non-Windows platforms.
func disableConsoleQuickEdit() {}
