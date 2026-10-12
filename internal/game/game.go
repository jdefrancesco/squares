package game

import (
	"image/color"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Game holds a mouse-controlled survival round and its persistent personal best.
// Update, Draw, and Layout implement ebiten.Game.
type Game struct {
	player              Entity
	angle               float64
	score, best         int
	gameOver, newBest   bool
	ents                []Entity
	spawnTimer, elapsed float64
	spawnsSinceEdible   int
	rng                 *rand.Rand

	eatPulse, scorePopupLeft float64
	scorePopupX, scorePopupY float64
	pendingEatSound          bool
	eatSound                 *audio.Player
	scorePath                string
	last                     time.Time

	paused                                                 bool
	prevMouseBtn, prevPKey, prevEscKey, prevRKey, prevQKey bool
}

// New returns a fresh round, loading the personal best when storage is available.
// Missing or damaged storage never prevents playing.
func New() *Game {
	g := &Game{}
	path, err := bestScorePath()
	if err != nil {
		log.Printf("Squares: personal best unavailable: %v", err)
	} else {
		g.scorePath = path
		if g.best, err = loadBestScore(path); err != nil {
			log.Printf("Squares: could not load personal best: %v", err)
		}
	}
	g.reset()
	return g
}

// reset clears round state while retaining the record, random source, and held
// controls. Preserving input edges prevents a held click from repeatedly retrying.
func (g *Game) reset() {
	g.player = Entity{x: ScreenWidth / 2, y: ScreenHeight / 2, size: initialPlayerSize}
	g.angle, g.score, g.elapsed, g.spawnTimer = 0, 0, 0, 0
	g.gameOver, g.paused, g.newBest = false, false, false
	g.ents = nil
	g.spawnsSinceEdible = 0
	g.eatPulse, g.scorePopupLeft = 0, 0
	g.pendingEatSound = false
	g.last = time.Now()
}

// Update reads the mouse and convenience keys, advances one bounded simulation
// step, and plays eating feedback. Q requests normal Ebiten termination.
func (g *Game) Update() error {
	now := time.Now()
	dt := clamp(now.Sub(g.last).Seconds(), 0, 0.05)
	g.last = now
	mx, my := ebiten.CursorPosition()
	quit := g.step(dt, inputState{
		x: float64(mx), y: float64(my),
		mouseDown:   ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		pauseDown:   ebiten.IsKeyPressed(ebiten.KeyP),
		escapeDown:  ebiten.IsKeyPressed(ebiten.KeyEscape),
		restartDown: ebiten.IsKeyPressed(ebiten.KeyR),
		quitDown:    ebiten.IsKeyPressed(ebiten.KeyQ),
	})
	if quit {
		return ebiten.Termination
	}
	if g.pendingEatSound {
		g.playEatSound()
	}
	return nil
}

// Draw renders color-coded squares, incoming warnings, eating feedback, and the
// score and record, followed by pause or retry prompts when needed.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{245, 246, 248, 255})
	for _, e := range g.ents {
		drawSquareAA(screen, e.x, e.y, e.size, squareColor(e.size, g.player.size))
		drawSpawnWarning(screen, e, g.player.size)
	}
	pulse := clamp(g.eatPulse/eatPulseDuration, 0, 1)
	playerCol := playerColor
	playerCol.R += uint8(100 * pulse)
	playerCol.G += uint8(65 * pulse)
	playerCol.B += uint8(100 * pulse)
	drawRotatedSquare(screen, g.player.x, g.player.y, g.player.size, g.angle, playerCol)
	if pulse > 0 {
		drawRing(screen, float32(g.player.x), float32(g.player.y), float32(g.player.size*0.8+12*(1-pulse)), 2, color.RGBA{235, 166, 48, uint8(180 * pulse)})
	}
	if g.scorePopupLeft > 0 {
		alpha := g.scorePopupLeft / scorePopupDuration
		y := int(g.scorePopupY - 12*(1-alpha))
		text.Draw(screen, "+1", hudFace, int(g.scorePopupX)-7, y, color.RGBA{155, 99, 13, uint8(255 * alpha)})
	}
	drawHUD(screen, 12, 12, g.score, g.best, pulse)
	drawLegend(screen)
	if g.gameOver {
		drawGameOverOverlay(screen, g.score, g.best, g.newBest)
	} else if g.paused {
		drawPauseOverlay(screen)
	}
}

// Layout returns the fixed logical playfield size regardless of window size.
func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Close saves the personal best and releases sound resources when the window
// closes or the player quits. It is safe to call when sound was never initialized.
func (g *Game) Close() error {
	if g.eatSound != nil {
		if err := g.eatSound.Close(); err != nil {
			log.Printf("Squares: could not close eating sound: %v", err)
		}
		g.eatSound = nil
	}
	return g.SaveBestScore()
}
