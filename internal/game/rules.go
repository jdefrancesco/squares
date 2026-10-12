package game

import (
	"image/color"
	"log"
	"math"
	"math/rand/v2"
)

var (
	playerColor = color.RGBA{35, 145, 85, 255}
	foodColor   = color.RGBA{235, 166, 48, 255}
	dangerColor = color.RGBA{40, 46, 56, 255}
)

// squareColor marks squares relative to the player's current side length, so a
// previously dangerous square changes to food as soon as the player outgrows it.
func squareColor(size, playerSize float64) color.RGBA {
	if size < playerSize {
		return foodColor
	}
	return dangerColor
}

// grownSize applies the eating reward, tapering it toward a 64-pixel side length.
// Score keeps increasing after growth becomes negligible.
func grownSize(size, foodSize float64) float64 {
	remaining := clamp((maxPlayerSize-size)/(maxPlayerSize-initialPlayerSize), 0, 1)
	return math.Min(maxPlayerSize, size+(growthScale*foodSize+growthFlat)*remaining)
}

// difficultyLevel combines active time and score into bounded pressure, reaching
// its maximum at 75 seconds of play minus 1.5 seconds for each point earned.
func difficultyLevel(elapsed float64, score int) float64 {
	return clamp((elapsed+1.5*float64(score))/75, 0, 1)
}

// random returns a value in [0, 1), using an injected source for repeatable tests.
func (g *Game) random() float64 {
	if g.rng != nil {
		return g.rng.Float64()
	}
	return rand.Float64()
}

// step advances a round from an input snapshot and returns whether Q was pressed.
// Clicks have no effect during play; after death, a fresh click or R retries.
func (g *Game) step(dt float64, in inputState) bool {
	clicked := in.mouseDown && !g.prevMouseBtn
	restart := in.restartDown && !g.prevRKey
	pause := (in.pauseDown && !g.prevPKey) || (in.escapeDown && !g.prevEscKey)
	quit := in.quitDown && !g.prevQKey
	g.prevMouseBtn, g.prevPKey, g.prevEscKey = in.mouseDown, in.pauseDown, in.escapeDown
	g.prevRKey, g.prevQKey = in.restartDown, in.quitDown
	g.pendingEatSound = false
	if quit {
		return true
	}
	if g.gameOver {
		if clicked || restart {
			g.reset()
			g.positionPlayer(in.x, in.y)
		}
		return false
	}
	if pause {
		g.paused = !g.paused
		return false
	}
	if g.paused {
		return false
	}
	dt = clamp(dt, 0, 0.05)
	g.elapsed += dt
	g.angle = math.Mod(g.angle+playerRotationRate*dt, 2*math.Pi)
	g.eatPulse = math.Max(0, g.eatPulse-dt)
	g.scorePopupLeft = math.Max(0, g.scorePopupLeft-dt)
	g.moveAndCollide(dt, in.x, in.y)
	if g.gameOver {
		return false
	}
	difficulty := difficultyLevel(g.elapsed, g.score)
	interval := 0.80 - 0.48*difficulty
	g.spawnTimer += dt
	for g.spawnTimer >= interval {
		g.spawnTimer -= interval
		g.spawnEntityWithDifficulty(difficulty)
	}
	return false
}

// moveAndCollide resolves contacts in chronological order along the mouse path.
// Growth can make a later square edible, so each eating event recomputes the
// remaining sweep. Larger squares win simultaneous contacts to avoid ambiguity.
func (g *Game) moveAndCollide(dt, mouseX, mouseY float64) {
	eaten := make([]bool, len(g.ents))
	elapsed := 0.0
	fatal := false
	for {
		from := g.player
		g.positionPlayer(mouseX, mouseY)
		target := g.player
		g.player = from
		remaining := dt - elapsed
		motion := from
		motion.size *= playerHitScale
		if remaining > 0 {
			motion.vx, motion.vy = (target.x-from.x)/remaining, (target.y-from.y)/remaining
		}
		first, hitIndex := math.Inf(1), -1
		for i, original := range g.ents {
			if eaten[i] {
				continue
			}
			e := original
			e.x, e.y = e.x+e.vx*elapsed, e.y+e.vy*elapsed
			if at, hit := squareCollisionTime(motion, e, remaining); hit {
				if at < first || (at == first && hitIndex >= 0 && e.size > g.ents[hitIndex].size) {
					first, hitIndex = at, i
				}
			}
		}
		if hitIndex < 0 {
			g.player, elapsed = target, dt
			break
		}
		elapsed += first
		g.player.x, g.player.y = from.x+motion.vx*first, from.y+motion.vy*first
		e := g.ents[hitIndex]
		if e.size >= g.player.size {
			fatal = true
			break
		}
		g.eatSquare(e)
		eaten[hitIndex] = true
	}

	// Reuse the entity slice, advancing survivors only as far as the round ran.
	alive := g.ents[:0]
	for i, e := range g.ents {
		if eaten[i] {
			continue
		}
		e.x, e.y = e.x+e.vx*elapsed, e.y+e.vy*elapsed
		margin := spawnMargin + e.size
		if e.x < -margin || e.x > ScreenWidth+margin || e.y < -margin || e.y > ScreenHeight+margin {
			continue
		}
		alive = append(alive, e)
	}
	g.ents = alive
	if fatal {
		g.finishRound()
	}
}

// positionPlayer follows the mouse while keeping the rotating square within the
// playfield, including at the maximum player size.
func (g *Game) positionPlayer(x, y float64) {
	margin := g.player.size / math.Sqrt2
	g.player.x = clamp(x, margin, ScreenWidth-margin)
	g.player.y = clamp(y, margin, ScreenHeight-margin)
}

// eatSquare awards a point, tapers growth, refreshes feedback, and records a new
// personal best in memory. Durable storage is updated at death or window close.
func (g *Game) eatSquare(e Entity) {
	g.score++
	if g.score > g.best {
		g.best, g.newBest = g.score, true
	}
	g.player.size = grownSize(g.player.size, e.size)
	g.positionPlayer(g.player.x, g.player.y)
	g.spawnsSinceEdible = 0
	g.eatPulse, g.scorePopupLeft = eatPulseDuration, scorePopupDuration
	g.scorePopupX, g.scorePopupY = g.player.x, math.Max(20, g.player.y-g.player.size)
	g.pendingEatSound = true
}

// finishRound freezes play and persists the current record, retaining the
// in-memory best if the filesystem is unavailable.
func (g *Game) finishRound() {
	g.gameOver = true
	if err := g.SaveBestScore(); err != nil {
		log.Printf("Squares: could not save personal best: %v", err)
	}
}
