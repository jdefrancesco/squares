package game

import "image/color"

type Kind int

const (
	KindSquare Kind = iota
	KindCircleHazard
	KindCircleBoost
)

type Entity struct {
	kind Kind
	x    float64
	y    float64
	size float64
	vx   float64
	vy   float64
	col  color.RGBA
}
