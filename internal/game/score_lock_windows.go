//go:build windows

package game

import (
	"os"

	"golang.org/x/sys/windows"
)

// lockScoreFile locks the first byte exclusively until file close or process
// exit. The lock file is separate from the atomically replaced JSON record.
func lockScoreFile(f *os.File) error {
	return windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &windows.Overlapped{})
}
