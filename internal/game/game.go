package game

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Game holds the player, entities, score, input state, and simulation timers.
// Its Update, Draw, and Layout methods implement ebiten.Game.
type Game struct {
	player Entity
	angle  float64

	score    int
	gameOver bool

	ents       []Entity
	spawnTimer float64

	elapsed           float64
	spawnsSinceEdible int

	dashCDLeft   float64
	dashInvLeft  float64
	prevMouseBtn bool

	invincibleLeft float64
	popupText      string
	popupLeft      float64

	paused     bool
	prevPKey   bool
	prevEscKey bool
	prevQKey   bool

	last time.Time
}

// New returns a game initialized for a new round with the player at the center.
func New() *Game {
	g := &Game{}
	g.reset()
	return g
}

// reset starts a new round, clearing entities, score, power-ups, and input state.
func (g *Game) reset() {
	g.player = Entity{
		kind: KindSquare,
		x:    ScreenWidth / 2,
		y:    ScreenHeight / 2,
		size: 22,
		col:  color.RGBA{35, 145, 85, 255},
	}
	g.angle = 0
	g.score = 0
	g.gameOver = false

	g.ents = nil
	g.spawnTimer = 0

	g.elapsed = 0
	g.spawnsSinceEdible = 0

	g.dashCDLeft = 0
	g.dashInvLeft = 0
	g.prevMouseBtn = false

	g.invincibleLeft = 0
	g.popupText = ""
	g.popupLeft = 0

	g.paused = false
	g.prevPKey = false
	g.prevEscKey = false
	g.prevQKey = false

	g.last = time.Now()
}

// Update handles input and advances timers, spawning, movement, and collisions.
// Pausing freezes the simulation; after game over, only restart and quit are
// handled. Pressing Q returns ebiten.Termination to end the game loop.
func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.last).Seconds()
	// Limit each simulation step so a delayed frame cannot cause a large jump.
	if dt > 0.05 {
		dt = 0.05
	}
	g.last = now

	qDown := ebiten.IsKeyPressed(ebiten.KeyQ)
	pressedQuitKey := qDown && !g.prevQKey
	g.prevQKey = qDown
	if pressedQuitKey {
		return ebiten.Termination
	}

	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.reset()
		}
		return nil
	}

	mx, my := ebiten.CursorPosition()
	mx = clamp(mx, 0, ScreenWidth-1)
	my = clamp(my, 0, ScreenHeight-1)

	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	// Trigger dash and pause on press transitions rather than while held.
	clicked := mouseDown && !g.prevMouseBtn
	g.prevMouseBtn = mouseDown

	pDown := ebiten.IsKeyPressed(ebiten.KeyP)
	escDown := ebiten.IsKeyPressed(ebiten.KeyEscape)
	pressedPauseKey := (pDown && !g.prevPKey) || (escDown && !g.prevEscKey)
	g.prevPKey = pDown
	g.prevEscKey = escDown

	if pressedPauseKey {
		g.paused = !g.paused
		g.last = now
		return nil
	}

	if g.paused {
		g.last = now
		return nil
	}

	g.elapsed += dt

	if g.dashCDLeft > 0 {
		g.dashCDLeft = math.Max(0, g.dashCDLeft-dt)
	}
	if g.dashInvLeft > 0 {
		g.dashInvLeft = math.Max(0, g.dashInvLeft-dt)
	}
	if g.invincibleLeft > 0 {
		g.invincibleLeft = math.Max(0, g.invincibleLeft-dt)
	}
	if g.popupLeft > 0 {
		g.popupLeft = math.Max(0, g.popupLeft-dt)
	}

	g.player.x = float64(mx)
	g.player.y = float64(my)

	if clicked && g.dashCDLeft <= 0 {
		// Dash briefly protects against squares; movement still follows the mouse.
		g.dashCDLeft = dashCooldown
		g.dashInvLeft = dashInvDuration
	}

	g.angle += playerRotationRate * dt

	difficulty := 0.12*g.elapsed + 0.8*float64(g.score)
	spawnInterval := clamp(0.85-0.0035*difficulty, 0.25, 0.85)

	g.spawnTimer += dt
	for g.spawnTimer >= spawnInterval {
		g.spawnTimer -= spawnInterval
		g.spawnEntityWithDifficulty(difficulty)
	}

	playerHit := g.player
	// Collisions use a slightly smaller, axis-aligned square than the rotating art.
	playerHit.size *= 0.90

	// Capture protection for this collision pass; collected boosts affect later ticks.
	inv := (g.invincibleLeft > 0) || (g.dashInvLeft > 0)

	// Reuse the entity slice, removing consumed entities and those beyond the
	// off-screen margin.
	alive := g.ents[:0]
	for _, e := range g.ents {
		e.x += e.vx * dt
		e.y += e.vy * dt

		if e.x < -160 || e.x > ScreenWidth+160 || e.y < -160 || e.y > ScreenHeight+160 {
			continue
		}

		switch e.kind {
		case KindSquare:
			if squareIntersectsSquare(playerHit, e) {
				// Equal-sized squares are fatal too unless protection is active.
				if inv || g.player.size > e.size {
					g.score++
					g.player.size += growthScale*e.size + growthFlat
					g.spawnsSinceEdible = 0
					continue
				}
				g.gameOver = true
			}

		case KindCircleHazard:
			// Black circles ignore both dash protection and boost invincibility.
			if circleIntersectsSquare(e, playerHit) {
				g.gameOver = true
			}

		case KindCircleBoost:
			if circleIntersectsSquare(e, playerHit) {
				g.invincibleLeft = invincibleDuration
				g.dashInvLeft = 0
				g.popupText = "INVINCIBLE"
				g.popupLeft = invinciblePopupDur
				continue
			}
		}

		alive = append(alive, e)
	}
	g.ents = alive

	return nil
}

// Draw renders the entities, rotating player, protection ring, and HUD, followed
// by the pause or game-over message when appropriate.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{245, 245, 245, 255})

	for _, e := range g.ents {
		if e.kind == KindSquare {
			drawSquareAA(screen, e.x, e.y, e.size, e.col)
		} else {
			drawFilledCircle(screen, float32(e.x), float32(e.y), float32(e.size/2), e.col)
		}
	}

	if g.invincibleLeft > 0 {
		drawRing(screen, float32(g.player.x), float32(g.player.y), float32(g.player.size*0.80), 4, color.RGBA{40, 150, 165, 220})
	} else if g.dashInvLeft > 0 {
		drawRing(screen, float32(g.player.x), float32(g.player.y), float32(g.player.size*0.75), 3, color.RGBA{120, 120, 120, 200})
	}

	drawRotatedSquare(screen, g.player.x, g.player.y, g.player.size, g.angle, g.player.col)

	invLeft := g.invincibleLeft
	if g.dashInvLeft > invLeft {
		invLeft = g.dashInvLeft
	}
	drawHUD(screen, 12, 12, g.score, invLeft)
	drawTopPopup(screen, g.popupText, g.popupLeft, invinciblePopupDur)
	if !g.gameOver {
		if g.paused {
			drawPauseOverlay(screen)
		}
	}

	if g.gameOver {
		text.Draw(screen, "GAME OVER\nPress R to restart", hudFace, ScreenWidth/2-90, ScreenHeight/2, color.RGBA{20, 20, 20, 255})
	}
}

// Layout returns the fixed logical playfield size regardless of window size.
func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return ScreenWidth, ScreenHeight
}
