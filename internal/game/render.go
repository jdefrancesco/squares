package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// whitePixel supplies a solid white texture for vertex-colored circle and ring
// meshes. It is initialized once and shared by the drawing helpers.
var whitePixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()

// drawSquareAA draws a filled, axis-aligned square using a temporary image
// positioned around (x, y). Despite its name, it does not enable antialiasing.
func drawSquareAA(dst *ebiten.Image, x, y, size float64, c color.RGBA) {
	img := ebiten.NewImage(int(size), int(size))
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-size/2, y-size/2)
	dst.DrawImage(img, op)
}

// drawRotatedSquare draws a filled square around (x, y), rotated by angle radians.
func drawRotatedSquare(dst *ebiten.Image, x, y, size, angle float64, c color.RGBA) {
	img := ebiten.NewImage(int(size), int(size))
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-size/2, -size/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x, y)
	dst.DrawImage(img, op)
}

// drawFilledCircle fills a circle centered at (x, y) with radius r, using a
// triangulated vector path tinted with c.
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
	dst.DrawTriangles(vs, is, whitePixel, nil)
}

// drawRing strokes a circle centered at (x, y) with path radius r and the supplied
// thickness, using vertex colors to tint the stroke with c.
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
	dst.DrawTriangles(vs, is, whitePixel, nil)
}
