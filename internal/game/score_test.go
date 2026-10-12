package game

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestBestScoreRoundTrip checks creation, loading, and preservation of records.
func TestBestScoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Squares", "best-score.json")
	if best, err := loadBestScore(path); err != nil || best != 0 {
		t.Fatalf("new player record=%d err=%v", best, err)
	}
	if err := saveBestScore(path, 27); err != nil {
		t.Fatal(err)
	}
	if err := saveBestScore(path, 12); err != nil {
		t.Fatal(err)
	}
	if best, err := loadBestScore(path); err != nil || best != 27 {
		t.Fatalf("stored best=%d err=%v; want 27", best, err)
	}
}

// TestBestScoreRejectsInvalidFiles checks malformed, negative, and oversized data.
func TestBestScoreRejectsInvalidFiles(t *testing.T) {
	for _, data := range []string{"not json", `{"best":-1}`, strings.Repeat("x", 4097)} {
		path := filepath.Join(t.TempDir(), "best.json")
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := loadBestScore(path); err == nil {
			t.Fatalf("invalid record accepted: %.30s", data)
		}
		if err := saveBestScore(path, 8); err != nil {
			t.Fatalf("could not recover damaged record: %v", err)
		}
		if best, err := loadBestScore(path); err != nil || best != 8 {
			t.Fatalf("recovered best=%d err=%v", best, err)
		}
	}
}

// TestBestScoreWriteError checks that storage failures are reported.
func TestBestScoreWriteError(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(parent, []byte("occupied"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := saveBestScore(filepath.Join(parent, "best.json"), 10); err == nil {
		t.Fatal("saving beneath a file should fail")
	}
	if err := saveBestScore(filepath.Join(t.TempDir(), "best.json"), -1); err == nil {
		t.Fatal("negative records should be rejected")
	}
}

// TestSaveBestScoreDoesNotLoseInMemoryRecord checks graceful storage failure.
func TestSaveBestScoreDoesNotLoseInMemoryRecord(t *testing.T) {
	g := newTestGame(t)
	g.best = 14
	g.scorePath = t.TempDir() // A directory cannot be replaced by the record file.
	if err := g.SaveBestScore(); err == nil || g.best != 14 {
		t.Fatal("storage failure should be reported without losing the record")
	}
}

// TestConcurrentBestScoresKeepMaximum checks serialized record updates.
func TestConcurrentBestScoresKeepMaximum(t *testing.T) {
	path := filepath.Join(t.TempDir(), "best.json")
	var group sync.WaitGroup
	start := make(chan struct{})
	for best := 1; best <= 24; best++ {
		group.Go(func() {
			<-start
			if err := saveBestScore(path, best); err != nil {
				t.Errorf("save %d: %v", best, err)
			}
		})
	}
	close(start)
	group.Wait()
	if best, err := loadBestScore(path); err != nil || best != 24 {
		t.Fatalf("concurrent record=%d err=%v; want 24", best, err)
	}
}
