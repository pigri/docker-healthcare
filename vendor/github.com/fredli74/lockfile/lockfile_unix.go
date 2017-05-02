// Copyright 2015 Fredrik Lidström. All rights reserved.
// Use of this source code is governed by the standard MIT License (MIT)
// that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package lockfile

import (
	"os"
	"syscall"
)

// ProcessRunning is a cross-platform check to work on both Windows and Unix systems as the os.FindProcess() function works differently.
func ProcessRunning(pid int) bool {
	p, _ := os.FindProcess(pid)   // On unix the FindProcess never returns an error
	err := p.Signal(syscall.Signal(0)) // Returns error if process is not running
	return err == nil
}
