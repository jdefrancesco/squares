package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// whitePixel is a shared texture for rectangles and rotated player artwork.
var whitePixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()

// drawSquareRect draws a colored rectangle from the shared texture without
// allocating a new image for every entity and frame.
func drawSquareRect(dst *ebiten.Image, x, y, width, height float64, c color.Color) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(width, height)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	dst.DrawImage(whitePixel, op)
}

// drawSquareAA draws a filled axis-aligned square centered at (x, y). The shared
// texture preserves fractional sizes; its historical name does not imply AA.
func drawSquareAA(dst *ebiten.Image, x, y, size float64, c color.RGBA) {
	drawSquareRect(dst, x-size/2, y-size/2, size, size, c)
}

// drawRotatedSquare draws a filled square around (x, y), rotated by angle radians.
func drawRotatedSquare(dst *ebiten.Image, x, y, size, angle float64, c color.RGBA) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(size, size)
	op.GeoM.Translate(-size/2, -size/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	dst.DrawImage(whitePixel, op)
}

// drawRing strokes a circle centered at (x, y) for the short-lived eating pulse.
func drawRing(dst *ebiten.Image, x, y, r, thickness float32, c color.RGBA) {
	vector.StrokeCircle(dst, x, y, r, thickness, c, true)
}

// drawSpawnWarning marks the predicted entry point of an incoming square using
// its current food or danger color. Squares already inside need no marker.
func drawSpawnWarning(dst *ebiten.Image, e Entity, playerSize float64) {
	if x, y, incoming := playfieldEntry(e); incoming {
		c := squareColor(e.size, playerSize)
		c.A = 175
		drawSquareAA(dst, clamp(x, 5, ScreenWidth-5), clamp(y, 5, ScreenHeight-5), 7, c)
	}
}
