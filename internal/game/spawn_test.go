package game

import (
	"math"
	"testing"
)

// TestSpawnForceEdibleAlwaysSquare checks the guaranteed food spawn and counter.
func TestSpawnForceEdibleAlwaysSquare(t *testing.T) {
	g := newTestGame(t)
	g.player.size, g.elapsed = 30, 60
	g.spawnsSinceEdible = maxSpawnsWithoutEdible
	g.spawnEntityWithDifficulty(1)
	if len(g.ents) != 1 || g.spawnsSinceEdible != 0 {
		t.Fatal("guaranteed food should append one square and reset the counter")
	}
	e := g.ents[0]
	if e.size < 10 || e.size >= g.player.size || squareColor(e.size, g.player.size) != foodColor {
		t.Fatalf("expected smaller amber square, got %+v", e)
	}
}

// TestSpawnFoodAndDangerAppear checks both square roles at late difficulty.
func TestSpawnFoodAndDangerAppear(t *testing.T) {
	g := newTestGame(t)
	g.elapsed = 60
	seenFood, seenDanger := false, false
	for i := 0; i < 200; i++ {
		g.ents, g.spawnsSinceEdible = nil, 0
		g.spawnEntityWithDifficulty(1)
		e := g.ents[0]
		seenFood = seenFood || e.size < g.player.size
		seenDanger = seenDanger || e.size >= g.player.size
	}
	if !seenFood || !seenDanger {
		t.Fatal("both food and danger should appear")
	}
}

// TestSpawnVelocityHasReasonableMagnitude checks bounded positive speed.
func TestSpawnVelocityHasReasonableMagnitude(t *testing.T) {
	g := newTestGame(t)
	g.spawnEntityWithDifficulty(100)
	e := g.ents[0]
	if speed := math.Hypot(e.vx, e.vy); speed <= 0 || speed > maxEntitySpeed {
		t.Fatalf("spawn speed should be positive and bounded: %v", speed)
	}
}

// TestSpawnLimitsCrowdingAndStartsGently checks food-only openings and threat cap.
func TestSpawnLimitsCrowdingAndStartsGently(t *testing.T) {
	g := newTestGame(t)
	for i := 0; i < 20; i++ {
		g.spawnEntityWithDifficulty(1)
		if g.ents[i].size >= g.player.size {
			t.Fatal("the opening seconds should spawn only food")
		}
	}
	g.elapsed = 60
	g.ents = make([]Entity, maxActiveDangers)
	for i := range g.ents {
		g.ents[i] = Entity{x: 30, y: 30, size: 40}
	}
	g.spawnEntityWithDifficulty(1)
	if g.ents[len(g.ents)-1].size >= g.player.size {
		t.Fatal("a crowded board should spawn food instead of another danger")
	}
}

// TestPlayfieldEntry checks edge warnings, on-screen squares, and departing paths.
func TestPlayfieldEntry(t *testing.T) {
	cases := []struct {
		name     string
		e        Entity
		x, y     float64
		incoming bool
	}{
		{"top", Entity{x: 100, y: -100, vy: 50}, 100, 0, true},
		{"left", Entity{x: -50, y: 100, vx: 50}, 0, 100, true},
		{"right", Entity{x: 740, y: 100, vx: -50}, 640, 100, true},
		{"bottom", Entity{x: 100, y: 580, vy: -50}, 100, 480, true},
		{"inside", Entity{x: 100, y: 100}, 0, 0, false},
		{"leaving", Entity{x: -10, y: 100, vx: -50}, 0, 0, false},
		{"parallel", Entity{x: -10, y: 100, vy: 50}, 0, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			x, y, incoming := playfieldEntry(tc.e)
			if incoming != tc.incoming || (incoming && (x != tc.x || y != tc.y)) {
				t.Fatalf("entry=(%v,%v,%v); want (%v,%v,%v)", x, y, incoming, tc.x, tc.y, tc.incoming)
			}
		})
	}
}
