package game

import (
	"fmt"
	"math"
	"math/rand/v2"
	"path/filepath"
	"testing"
)

// newTestGame creates a deterministic round without touching the user's record.
func newTestGame(t *testing.T) *Game {
	t.Helper()
	g := &Game{rng: rand.New(rand.NewPCG(1, 2)), scorePath: filepath.Join(t.TempDir(), "best.json")}
	g.reset()
	return g
}

// mouseAt supplies an input snapshot with the cursor at the player's position.
func mouseAt(g *Game) inputState {
	return inputState{x: g.player.x, y: g.player.y}
}

// TestEatSquareScoresAndGrows checks the reward and immediate visual feedback.
func TestEatSquareScoresAndGrows(t *testing.T) {
	g := newTestGame(t)
	g.ents = []Entity{{x: g.player.x, y: g.player.y, size: 12}}
	g.step(0.01, mouseAt(g))
	if g.gameOver || g.score != 1 || g.best != 1 || len(g.ents) != 0 {
		t.Fatalf("eating failed: over=%v score=%d best=%d entities=%d", g.gameOver, g.score, g.best, len(g.ents))
	}
	if g.player.size <= initialPlayerSize || g.eatPulse <= 0 || !g.pendingEatSound {
		t.Fatal("eating should grow the player and trigger feedback")
	}
}

// TestEqualOrLargerSquareEndsRound checks fatal collisions and saved records.
func TestEqualOrLargerSquareEndsRound(t *testing.T) {
	for _, size := range []float64{initialPlayerSize, initialPlayerSize + 20} {
		t.Run(fmt.Sprintf("size_%g", size), func(t *testing.T) {
			g := newTestGame(t)
			g.score, g.best = 7, 7
			g.ents = []Entity{{x: g.player.x, y: g.player.y, size: size}}
			in := mouseAt(g)
			in.mouseDown = true // Clicking during a round never provides protection.
			g.step(0.01, in)
			if !g.gameOver || g.score != 7 {
				t.Fatalf("danger collision: over=%v score=%d", g.gameOver, g.score)
			}
			best, err := loadBestScore(g.scorePath)
			if err != nil || best != 7 {
				t.Fatalf("saved record = %d, %v; want 7", best, err)
			}
		})
	}
}

// TestGrowthTapersAtLimit checks continued scoring without unbounded size.
func TestGrowthTapersAtLimit(t *testing.T) {
	g := newTestGame(t)
	firstGrowth := grownSize(g.player.size, 12) - g.player.size
	lateGrowth := grownSize(maxPlayerSize-1, 12) - (maxPlayerSize - 1)
	if lateGrowth <= 0 || lateGrowth >= firstGrowth {
		t.Fatalf("growth should taper: first=%v late=%v", firstGrowth, lateGrowth)
	}
	for i := 0; i < 1000; i++ {
		g.ents = []Entity{{x: g.player.x, y: g.player.y, size: 12}}
		g.step(0.001, mouseAt(g))
	}
	if g.score != 1000 || g.player.size > maxPlayerSize || g.player.size < maxPlayerSize-0.1 {
		t.Fatalf("long run: score=%d size=%v", g.score, g.player.size)
	}
}

// TestClickRetryRequiresFreshPress checks one-click retry and held-button safety.
func TestClickRetryRequiresFreshPress(t *testing.T) {
	g := newTestGame(t)
	g.score, g.best, g.gameOver = 3, 9, true
	g.prevMouseBtn = true
	in := mouseAt(g)
	in.mouseDown = true
	g.step(0.01, in)
	if !g.gameOver {
		t.Fatal("a held click should not restart")
	}
	in.mouseDown = false
	g.step(0.01, in)
	in.mouseDown = true
	g.step(0.01, in)
	if g.gameOver || g.score != 0 || g.best != 9 || len(g.ents) != 0 {
		t.Fatal("fresh click should start a clean round and retain the best score")
	}
	if !g.prevMouseBtn || g.elapsed != 0 {
		t.Fatal("retry should preserve the held input and start at time zero")
	}
}

// TestPauseFreezesRound checks movement, difficulty, feedback, and key edges.
func TestPauseFreezesRound(t *testing.T) {
	g := newTestGame(t)
	g.eatPulse = 0.1
	g.ents = []Entity{{x: 100, y: 100, size: 12, vx: 50}}
	in := mouseAt(g)
	in.pauseDown = true
	g.step(0.01, in)
	in.x = 10
	g.step(1, in)
	if !g.paused || g.elapsed != 0 || g.player.x != ScreenWidth/2 || g.ents[0].x != 100 || g.eatPulse != 0.1 {
		t.Fatal("holding pause should leave the complete round frozen")
	}
	in.pauseDown = false
	g.step(0.01, in)
	in.escapeDown = true
	g.step(0.01, in)
	if g.paused {
		t.Fatal("a fresh Escape press should resume")
	}
}

// TestSquareColorsFollowCurrentSize checks that newly edible threats change color.
func TestSquareColorsFollowCurrentSize(t *testing.T) {
	if squareColor(30, 22) != dangerColor || squareColor(22, 22) != dangerColor {
		t.Fatal("equal and larger squares should be marked as dangerous")
	}
	if squareColor(30, 31) != foodColor || foodColor == dangerColor {
		t.Fatal("edible squares should have a distinct food color")
	}
}

// TestDifficultyRampsAndCaps checks the early game and bounded late pressure.
func TestDifficultyRampsAndCaps(t *testing.T) {
	if difficultyLevel(0, 0) != 0 || difficultyLevel(60, 10) != 1 || difficultyLevel(600, 1000) != 1 {
		t.Fatal("difficulty should start at zero and cap around a minute of play")
	}
	if difficultyLevel(10, 0) >= difficultyLevel(10, 10) {
		t.Fatal("score should also increase difficulty")
	}
}

// TestSweptCollision detects crossings even when neither endpoint overlaps.
func TestSweptCollision(t *testing.T) {
	a := Entity{x: 0, y: 0, size: 10, vx: 100}
	b := Entity{x: 50, y: 0, size: 10}
	if !sweptSquareCollision(a, b, 1) || sweptSquareCollision(a, b, 0.1) {
		t.Fatal("crossing collision should depend on the sweep duration")
	}
	b.y = 11
	if sweptSquareCollision(a, b, 1) {
		t.Fatal("separated parallel paths should stay safe")
	}
	b.y = 10
	if !sweptSquareCollision(a, b, 1) {
		t.Fatal("tangent paths should count as collisions")
	}
}

// TestMouseSweepCannotTunnelThroughDanger checks fast cursor movement.
func TestMouseSweepCannotTunnelThroughDanger(t *testing.T) {
	g := newTestGame(t)
	g.player.x, g.player.y = 100, 240
	g.ents = []Entity{{x: 200, y: 240, size: 40}}
	g.step(0.01, inputState{x: 300, y: 240})
	if !g.gameOver {
		t.Fatal("moving across a danger square should end the round")
	}
}

// TestEarlierFoodMakesLaterSquareEdible checks chronological growth in a sweep.
func TestEarlierFoodMakesLaterSquareEdible(t *testing.T) {
	g := newTestGame(t)
	g.player.x, g.player.y, g.player.size = 100, 240, 31
	// Deliberately reverse slice order so collision chronology decides the result.
	g.ents = []Entity{{x: 180, y: 240, size: 31.5}, {x: 150, y: 240, size: 23}}
	g.step(0.01, inputState{x: 300, y: 240})
	if g.gameOver || g.score != 2 || len(g.ents) != 0 {
		t.Fatalf("earlier food should make the later square edible: over=%v score=%d", g.gameOver, g.score)
	}
}

// TestSimultaneousDangerWins checks that food cannot cancel same-time death.
func TestSimultaneousDangerWins(t *testing.T) {
	g := newTestGame(t)
	g.ents = []Entity{{x: g.player.x, y: g.player.y, size: 12}, {x: g.player.x, y: g.player.y, size: 23}}
	g.step(0.01, mouseAt(g))
	if !g.gameOver || g.score != 0 {
		t.Fatal("danger at the same contact time should be fatal")
	}
}

// TestSpawnLeavesReactionTime checks incoming threats at corners and the center.
func TestSpawnLeavesReactionTime(t *testing.T) {
	for _, point := range [][2]float64{{32, 32}, {608, 448}, {320, 240}} {
		g := newTestGame(t)
		g.player.x, g.player.y, g.player.size, g.elapsed = point[0], point[1], maxPlayerSize, 60
		for i := 0; i < 200; i++ {
			g.ents = nil
			g.spawnsSinceEdible = 0
			g.spawnEntityWithDifficulty(1)
			e := g.ents[0]
			if e.size >= g.player.size {
				playerHit := g.player
				playerHit.size *= playerHitScale
				if sweptSquareCollision(e, playerHit, spawnReactionTime-0.001) {
					t.Fatalf("threat spawned without reaction time: %+v", e)
				}
			}
			if speed := math.Hypot(e.vx, e.vy); speed <= 0 || speed > maxEntitySpeed {
				t.Fatalf("spawn speed=%v", speed)
			}
		}
	}
}

// TestBlockedEscapeRejectsThreat checks the spawn safety gate and food exclusion.
func TestBlockedEscapeRejectsThreat(t *testing.T) {
	g := newTestGame(t)
	candidate := Entity{x: 320, y: -100, size: 50, vy: 200}
	if !g.hasEscapeRoute(candidate) {
		t.Fatal("an empty playfield should offer an escape route")
	}
	g.ents = []Entity{{x: g.player.x, y: g.player.y, size: 200}}
	if g.hasEscapeRoute(candidate) {
		t.Fatal("a blocked playfield should reject another threat")
	}
	g.ents[0].size = 10
	if !g.hasEscapeRoute(candidate) {
		t.Fatal("food should not block an escape route")
	}
}

// TestQuitAndKeyboardRetry checks the retained convenience controls.
func TestQuitAndKeyboardRetry(t *testing.T) {
	g := newTestGame(t)
	g.gameOver = true
	in := mouseAt(g)
	in.restartDown = true
	g.step(0.01, in)
	if g.gameOver {
		t.Fatal("R should restart after death")
	}
	in.quitDown = true
	if !g.step(0.01, in) {
		t.Fatal("Q should request termination")
	}
}
