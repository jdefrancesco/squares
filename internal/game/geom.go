package game

import "math"

// Number is a constraint that permits any numeric type
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64
}

// Clamp returns value clamped to the range [low, high]
func clamp[T Number](value, low, high T) T {
    return max(low, min(value, high))
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
