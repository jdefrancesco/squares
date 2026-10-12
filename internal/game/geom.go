package game

import "math"

// Number permits integer and floating-point types, including named types with
// those underlying types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// clamp returns value limited to the inclusive range [low, high].
// The caller must ensure low <= high.
func clamp[T Number](value, low, high T) T {
	return max(low, min(value, high))
}

// squareIntersectsSquare reports whether two axis-aligned squares overlap or
// touch. Entity.size is the side length; visual rotation is ignored.
func squareIntersectsSquare(a, b Entity) bool {
	ah := a.size / 2
	bh := b.size / 2
	return math.Abs(a.x-b.x) <= (ah+bh) && math.Abs(a.y-b.y) <= (ah+bh)
}

// circleIntersectsSquare reports whether a circle overlaps or touches an
// axis-aligned square by testing the nearest point on the square to the circle's
// center. The circle's size is its diameter; the square's size is its side length.
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
