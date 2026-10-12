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

// squareCollisionTime returns the first contact time of two moving axis-aligned
// squares during duration. Relative-motion slabs include edge contact and avoid
// tunneling when the mouse jumps across an entire square in one update.
func squareCollisionTime(a, b Entity, duration float64) (float64, bool) {
	enter, leave := 0.0, duration
	half := (a.size + b.size) / 2
	for _, axis := range [][2]float64{{a.x - b.x, a.vx - b.vx}, {a.y - b.y, a.vy - b.vy}} {
		pos, velocity := axis[0], axis[1]
		if velocity == 0 {
			if math.Abs(pos) > half {
				return 0, false
			}
			continue
		}
		first, last := (-half-pos)/velocity, (half-pos)/velocity
		if first > last {
			first, last = last, first
		}
		enter, leave = math.Max(enter, first), math.Min(leave, last)
		if enter > leave {
			return 0, false
		}
	}
	return enter, enter <= duration && leave >= 0
}

// sweptSquareCollision reports overlap anywhere along two squares' linear paths.
func sweptSquareCollision(a, b Entity, duration float64) bool {
	_, hit := squareCollisionTime(a, b, duration)
	return hit
}
