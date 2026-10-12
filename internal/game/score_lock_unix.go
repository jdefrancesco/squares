//go:build darwin || linux || freebsd || netbsd || openbsd || dragonfly

package game

import (
	"os"
	"syscall"
)

// lockScoreFile acquires an exclusive lock released automatically on file close
// or process exit, including abnormal termination.
func lockScoreFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}
