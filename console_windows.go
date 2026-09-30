//go:build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	enableExtendedFlags = 0x0080
	enableQuickEditMode = 0x0040
)

// disableConsoleQuickEdit turns off the Windows console QuickEdit mode. With it
// on, clicking/selecting text in the console window puts it into selection mode
// and any process writing to stdout blocks until the selection is cleared
// (e.g. by a keypress), which makes a server look like it has hung.
func disableConsoleQuickEdit() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	handle := uintptr(os.Stdin.Fd())

	var mode uint32
	if r, _, _ := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode))); r == 0 {
		return // no console attached (redirected output or running as a service)
	}

	// QuickEdit can only be changed while ENABLE_EXTENDED_FLAGS is set.
	mode |= enableExtendedFlags
	mode &^= enableQuickEditMode
	setConsoleMode.Call(handle, uintptr(mode))
}
