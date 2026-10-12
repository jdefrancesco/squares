package game

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// bestScorePath resolves the per-user record location without creating files.
func bestScorePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find user config directory: %w", err)
	}
	return filepath.Join(dir, "Squares", "best-score.json"), nil
}

// loadBestScore reads a bounded JSON record. A missing file means a new player;
// malformed or negative records return an error rather than becoming a score.
func loadBestScore(path string) (int, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("open personal best: %w", err)
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, 4097))
	if err != nil {
		return 0, fmt.Errorf("read personal best: %w", err)
	}
	var record struct {
		Best int `json:"best"`
	}
	if len(data) > 4096 {
		return 0, fmt.Errorf("personal best file is too large")
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return 0, fmt.Errorf("decode personal best: %w", err)
	}
	if record.Best < 0 {
		return 0, fmt.Errorf("personal best must be nonnegative")
	}
	return record.Best, nil
}

// saveBestScore writes through a temporary file and atomically replaces the
// record, preserving a higher existing score and repairing damaged records.
func saveBestScore(path string, best int) error {
	if best < 0 {
		return fmt.Errorf("personal best must be nonnegative")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create personal best directory: %w", err)
	}
	// Keep the lock file's inode stable: removing it could let another process
	// lock a different file while a waiting writer still holds the original.
	lock, err := openScoreLock(path + ".lock")
	if err != nil {
		return err
	}
	defer lock.Close() // Closing the descriptor releases the OS lock after crashes too.
	if previous, err := loadBestScore(path); err == nil {
		best = max(best, previous)
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".best-*.json")
	if err != nil {
		return fmt.Errorf("create personal best file: %w", err)
	}
	defer os.Remove(f.Name()) // Remove any leftover temporary file after a failure.
	if _, err := fmt.Fprintf(f, "{\"best\":%d}\n", best); err != nil {
		f.Close()
		return fmt.Errorf("write personal best: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close personal best file: %w", err)
	}
	if err := os.Rename(f.Name(), path); err != nil {
		return fmt.Errorf("replace personal best: %w", err)
	}
	return nil
}

// openScoreLock opens a persistent advisory lock file and serializes record
// updates across game windows. The caller releases the lock by closing the file.
func openScoreLock(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open personal best lock: %w", err)
	}
	if err := lockScoreFile(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("lock personal best: %w", err)
	}
	return f, nil
}

// SaveBestScore persists the personal best. Without a storage path it retains
// the in-memory record, allowing play when configuration storage is unavailable.
func (g *Game) SaveBestScore() error {
	if g.scorePath == "" || g.best == 0 {
		return nil
	}
	return saveBestScore(g.scorePath, g.best)
}
