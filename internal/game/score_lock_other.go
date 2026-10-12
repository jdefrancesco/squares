//go:build !(darwin || linux || freebsd || netbsd || openbsd || dragonfly || windows)

package game

import (
	"fmt"
	"os"
)

// lockScoreFile reports unsupported persistent storage while allowing the game
// to retain its in-memory best and continue on other platforms.
func lockScoreFile(f *os.File) error {
	return fmt.Errorf("personal best file locking is unavailable on this platform")
}
