package game

import (
	"math"
	"math/rand"
	"testing"
)

func TestSpawnForceEdibleAlwaysSquare(t *testing.T) {
	rand.Seed(1)

	g := New()
	g.player.size = 30
	g.spawnsSinceEdible = maxSpawnsWithoutEdible

	before := len(g.ents)
	g.spawnEntityWithDifficulty(0)
	if len(g.ents) != before+1 {
		t.Fatalf("expected 1 entity appended, got %d", len(g.ents)-before)
	}

	e := g.ents[len(g.ents)-1]
	if e.kind != KindSquare {
		t.Fatalf("expected forced-edible spawn to be KindSquare, got %v", e.kind)
	}
	if g.spawnsSinceEdible != 0 {
		t.Fatalf("expected spawnsSinceEdible reset to 0, got %d", g.spawnsSinceEdible)
	}
	if e.col.A != 255 {
		t.Fatalf("expected alpha=255, got %d", e.col.A)
	}
	if e.col.R != e.col.G || e.col.G != e.col.B {
		t.Fatalf("expected square color to be grayscale, got %+v", e.col)
	}
	if e.col.R < 60 || e.col.R > 209 {
		t.Fatalf("expected grayscale value in [60,209], got %d", e.col.R)
	}

	// Forced-edible squares should be smaller than the player.
	if e.size < 10 {
		t.Fatalf("expected edible square size >= 10, got %v", e.size)
	}
	if e.size >= g.player.size {
		t.Fatalf("expected edible square smaller than player (player=%v, square=%v)", g.player.size, e.size)
	}
}

func TestSpawnCircleKindsAndColorsAppear(t *testing.T) {
	rand.Seed(2)

	g := New()
	g.player.size = 30

	seenBoost := false
	seenHazard := false

	for i := 0; i < 800; i++ {
		// Avoid being forced into edible-only squares.
		g.spawnsSinceEdible = 0
		g.spawnEntityWithDifficulty(0)
		e := g.ents[len(g.ents)-1]

		switch e.kind {
		case KindCircleBoost:
			seenBoost = true
			if e.col.R != 40 || e.col.G != 150 || e.col.B != 165 || e.col.A != 255 {
				t.Fatalf("unexpected boost color: %+v", e.col)
			}
		case KindCircleHazard:
			seenHazard = true
			if e.col.R != 15 || e.col.G != 15 || e.col.B != 15 || e.col.A != 255 {
				t.Fatalf("unexpected hazard color: %+v", e.col)
			}
		}

		if seenBoost && seenHazard {
			break
		}
	}

	if !seenBoost {
		t.Fatalf("expected to see at least one boost circle spawn")
	}
	if !seenHazard {
		t.Fatalf("expected to see at least one hazard circle spawn")
	}
}

func TestSpawnVelocityHasReasonableMagnitude(t *testing.T) {
	rand.Seed(3)

	g := New()
	g.player.size = 40
	g.spawnsSinceEdible = 0

	g.spawnEntityWithDifficulty(200)
	e := g.ents[len(g.ents)-1]

	speed := math.Hypot(e.vx, e.vy)
	if speed <= 0 {
		t.Fatalf("expected speed > 0")
	}
	// Internal clamp is [75,340] with small kind multipliers; allow a little headroom.
	if speed > 450 {
		t.Fatalf("expected speed to be bounded, got %v", speed)
	}
}
