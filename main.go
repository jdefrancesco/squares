package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenW = 1024
	screenH = 800

	// Click dash: brief invincibility burst (since player is cursor-locked)
	dashCooldown    = 0.90
	dashInvDuration = 0.25

	// Booster (green circle)
	invincibleDuration = 5.0

	// Fairness
	maxSpawnsWithoutEdible = 6

	// Visuals
	playerRotationRate = 6.0 // radians/sec constant

	// Growth tuning (slower)
	growthScale = 0.05
	growthFlat  = 0.8
)

var (
	hudFace   = basicfont.Face7x13
	hudColor  = color.RGBA{20, 20, 20, 255}
	hudPanel  = color.RGBA{255, 255, 255, 220}
	hudBorder = color.RGBA{0, 0, 0, 70}
)

type Kind int

const (
	KindSquare Kind = iota
	KindCircleHazard
	KindCircleBoost
)

type Entity struct {
	kind     Kind
	x, y     float64
	size     float64
	vx, vy   float64
	col      color.RGBA
}

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

	last time.Time
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func newGame() *Game {
	g := &Game{}
	g.reset()
	return g
}

func (g *Game) reset() {
	g.player = Entity{
		kind: KindSquare,
		x:    screenW / 2,
		y:    screenH / 2,
		size: 22,
		col:  color.RGBA{0, 0, 0, 255},
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

	g.last = time.Now()
}

// --- Collisions ---

func squareIntersectsSquare(a, b Entity) bool {
	ah := a.size / 2
	bh := b.size / 2
	return math.Abs(a.x-b.x) <= (ah+bh) && math.Abs(a.y-b.y) <= (ah+bh)
}

func circleIntersectsSquare(circle, square Entity) bool {
	r := circle.size / 2
	half := square.size / 2

	minX := square.x - half
	maxX := square.x + half
	minY := square.y - half
	maxY := square.y + half

	cx := clamp(circle.x, minX, maxX)
	cy := clamp(circle.y, minY, maxY)

	dx := circle.x - cx
	dy := circle.y - cy
	return (dx*dx + dy*dy) <= r*r
}

// --- Spawning ---

func (g *Game) spawnEntityWithDifficulty(d float64) {
	edge := rand.Intn(4)
	margin := 50.0

	var x, y float64
	switch edge {
	case 0:
		x = rand.Float64() * screenW
		y = -margin
	case 1:
		x = rand.Float64() * screenW
		y = screenH + margin
	case 2:
		x = -margin
		y = rand.Float64() * screenH
	case 3:
		x = screenW + margin
		y = rand.Float64() * screenH
	}

	hazardP := clamp(0.06+0.0009*d, 0.06, 0.22)
	boostP := clamp(0.06-0.00025*d, 0.02, 0.06)

	forceEdible := g.spawnsSinceEdible >= maxSpawnsWithoutEdible

	kind := KindSquare
	if !forceEdible {
		r := rand.Float64()
		switch {
		case r < hazardP:
			kind = KindCircleHazard
		case r < hazardP+boostP:
			kind = KindCircleBoost
		default:
			kind = KindSquare
		}
	}

	p := g.player.size
	var size float64

	if kind == KindSquare {
		edibleBias := clamp(0.78-0.0007*d, 0.45, 0.78)
		isEdible := forceEdible || (rand.Float64() < edibleBias)

		if isEdible {
			size = p * (0.45 + rand.Float64()*0.45)
			size = math.Max(10, size)
			g.spawnsSinceEdible = 0
		} else {
			threatMin := 1.02
			threatMax := clamp(1.35+0.0006*d, 1.35, 2.10)
			size = p*(threatMin+rand.Float64()*(threatMax-threatMin)) + 8
			g.spawnsSinceEdible++
		}
	} else {
		minS := math.Max(18, p*0.60)
		maxS := p*1.10 + 34
		size = minS + rand.Float64()*(maxS-minS)
		g.spawnsSinceEdible++
	}

	// Aim toward player
	dx := g.player.x - x
	dy := g.player.y - y
	dist := math.Hypot(dx, dy)
	if dist < 1 {
		dist = 1
	}
	dx /= dist
	dy /= dist

	baseK := 950.0 + 0.9*d
	speed := baseK/math.Sqrt(size) + (rand.Float64()*60 - 30)
	speed = clamp(speed, 75, 340)

	switch kind {
	case KindCircleHazard:
		speed *= 1.10
	case KindCircleBoost:
		speed *= 0.92
	}

	jitter := clamp(0.35-0.0002*d, 0.18, 0.35)
	dx += (rand.Float64()*2 - 1) * jitter
	dy += (rand.Float64()*2 - 1) * jitter
	nd := math.Hypot(dx, dy)
	if nd < 1 {
		nd = 1
	}
	dx /= nd
	dy /= nd

	var col color.RGBA
	switch kind {
	case KindSquare:
		c := uint8(60 + rand.Intn(150))
		col = color.RGBA{c, c, c, 255}
	case KindCircleHazard:
		col = color.RGBA{70, 110, 200, 255}
	case KindCircleBoost:
		col = color.RGBA{40, 180, 80, 255}
	}

	g.ents = append(g.ents, Entity{
		kind: kind,
		x:    x,
		y:    y,
		size: size,
		vx:   dx * speed,
		vy:   dy * speed,
		col:  col,
	})
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.last).Seconds()
	if dt > 0.05 {
		dt = 0.05
	}
	g.last = now

	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.reset()
		}
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

	// Cursor-locked player (clamped)
	mx, my := ebiten.CursorPosition()
	mx = clampInt(mx, 0, screenW-1)
	my = clampInt(my, 0, screenH-1)
	g.player.x = float64(mx)
	g.player.y = float64(my)

	// Click edge + dash inv burst
	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	clicked := mouseDown && !g.prevMouseBtn
	g.prevMouseBtn = mouseDown

	if clicked && g.dashCDLeft <= 0 {
		g.dashCDLeft = dashCooldown
		g.dashInvLeft = dashInvDuration
	}

	// Rotation
	g.angle += playerRotationRate * dt

	// Spawn
	difficulty := 0.12*g.elapsed + 0.8*float64(g.score)
	spawnInterval := clamp(0.85-0.0035*difficulty, 0.25, 0.85)

	g.spawnTimer += dt
	for g.spawnTimer >= spawnInterval {
		g.spawnTimer -= spawnInterval
		g.spawnEntityWithDifficulty(difficulty)
	}

	// Collision hitbox (slightly forgiving)
	playerHit := g.player
	playerHit.size *= 0.90

	inv := (g.invincibleLeft > 0) || (g.dashInvLeft > 0)

	// Move entities + collisions
	alive := g.ents[:0]
	for _, e := range g.ents {
		e.x += e.vx * dt
		e.y += e.vy * dt

		if e.x < -160 || e.x > screenW+160 || e.y < -160 || e.y > screenH+160 {
			continue
		}

		switch e.kind {
		case KindSquare:
			if squareIntersectsSquare(playerHit, e) {
				if inv || g.player.size > e.size {
					g.score++
					g.player.size += growthScale*e.size + growthFlat
					g.spawnsSinceEdible = 0
					continue
				}
				g.gameOver = true
			}

		case KindCircleHazard:
			if circleIntersectsSquare(e, playerHit) && !inv {
				g.gameOver = true
			}

		case KindCircleBoost:
			if circleIntersectsSquare(e, playerHit) {
				g.invincibleLeft = invincibleDuration
				g.dashInvLeft = 0
				continue
			}
		}

		alive = append(alive, e)
	}
	g.ents = alive

	return nil
}

// --- HUD (simplified) ---

func drawHUD(screen *ebiten.Image, x, y int, score int, invLeft float64) {
	textStr := fmt.Sprintf(
		"Squares eaten: %d\nInvincible: %.1fs\n\nGreen circle: invincibility\nBlue circle: instant death",
		score,
		math.Max(0, invLeft),
	)

	// simple sizing
	lines := 5
	maxLen := 0
	cur := 0
	for _, r := range textStr {
		if r == '\n' {
			if cur > maxLen {
				maxLen = cur
			}
			cur = 0
		} else {
			cur++
		}
	}
	if cur > maxLen {
		maxLen = cur
	}

	pad := 10
	w := pad*2 + maxLen*7
	h := pad*2 + lines*13

	panel := ebiten.NewImage(w, h)
	panel.Fill(hudPanel)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(panel, op)

	// border
	drawRect(screen, x, y, w, 1, hudBorder)
	drawRect(screen, x, y+h-1, w, 1, hudBorder)
	drawRect(screen, x, y, 1, h, hudBorder)
	drawRect(screen, x+w-1, y, 1, h, hudBorder)

	// text + shadow
	text.Draw(screen, textStr, hudFace, x+pad+1, y+pad+13+1, color.RGBA{0, 0, 0, 90})
	text.Draw(screen, textStr, hudFace, x+pad, y+pad+13, hudColor)
}

func drawRect(dst *ebiten.Image, x, y, w, h int, c color.RGBA) {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(img, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{245, 245, 245, 255})

	for _, e := range g.ents {
		if e.kind == KindSquare {
			drawSquareAA(screen, e.x, e.y, e.size, e.col)
		} else {
			drawFilledCircle(screen, float32(e.x), float32(e.y), float32(e.size/2), e.col)
		}
	}

	// rings for invincibility states
	if g.invincibleLeft > 0 {
		drawRing(screen, float32(g.player.x), float32(g.player.y), float32(g.player.size*0.80), 4, color.RGBA{40, 180, 80, 220})
	} else if g.dashInvLeft > 0 {
		drawRing(screen, float32(g.player.x), float32(g.player.y), float32(g.player.size*0.75), 3, color.RGBA{120, 120, 120, 200})
	}

	drawRotatedSquare(screen, g.player.x, g.player.y, g.player.size, g.angle, g.player.col)

	// show invincibility time remaining (include dash burst)
	invLeft := g.invincibleLeft
	if g.dashInvLeft > invLeft {
		invLeft = g.dashInvLeft
	}
	drawHUD(screen, 12, 12, g.score, invLeft)

	if g.gameOver {
		text.Draw(screen, "GAME OVER\nPress R to restart", hudFace, screenW/2-90, screenH/2, color.RGBA{20, 20, 20, 255})
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) { return screenW, screenH }

// --- Drawing helpers ---

func drawSquareAA(dst *ebiten.Image, x, y, size float64, c color.RGBA) {
	img := ebiten.NewImage(int(size), int(size))
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-size/2, y-size/2)
	dst.DrawImage(img, op)
}

func drawRotatedSquare(dst *ebiten.Image, x, y, size, angle float64, c color.RGBA) {
	img := ebiten.NewImage(int(size), int(size))
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-size/2, -size/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x, y)
	dst.DrawImage(img, op)
}

func drawFilledCircle(dst *ebiten.Image, x, y, r float32, c color.RGBA) {
	var p vector.Path
	p.Arc(x, y, r, 0, 2*math.Pi, vector.Clockwise)
	p.Close()

	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = float32(c.R) / 255
		vs[i].ColorG = float32(c.G) / 255
		vs[i].ColorB = float32(c.B) / 255
		vs[i].ColorA = float32(c.A) / 255
	}
	dst.DrawTriangles(vs, is, ebiten.NewImage(1, 1), nil)
}

func drawRing(dst *ebiten.Image, x, y, r float32, thickness float32, c color.RGBA) {
	var p vector.Path
	p.Arc(x, y, r, 0, 2*math.Pi, vector.Clockwise)

	vs, is := p.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: thickness,
	})
	for i := range vs {
		vs[i].ColorR = float32(c.R) / 255
		vs[i].ColorG = float32(c.G) / 255
		vs[i].ColorB = float32(c.B) / 255
		vs[i].ColorA = float32(c.A) / 255
	}
	dst.DrawTriangles(vs, is, ebiten.NewImage(1, 1), nil)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Squares")

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}

