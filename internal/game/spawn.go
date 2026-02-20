package game

import (
	"image/color"
	"math"
	"math/rand"
)

func (g *Game) spawnEntityWithDifficulty(d float64) {
	edge := rand.Intn(4)
	margin := 50.0

	var x, y float64
	switch edge {
	case 0:
		x = rand.Float64() * ScreenWidth
		y = -margin
	case 1:
		x = rand.Float64() * ScreenWidth
		y = ScreenHeight + margin
	case 2:
		x = -margin
		y = rand.Float64() * ScreenHeight
	case 3:
		x = ScreenWidth + margin
		y = rand.Float64() * ScreenHeight
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
		col = color.RGBA{15, 15, 15, 255}
	case KindCircleBoost:
		col = color.RGBA{40, 150, 165, 255}
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
