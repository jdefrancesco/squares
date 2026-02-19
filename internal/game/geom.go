package game

import "math"

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
