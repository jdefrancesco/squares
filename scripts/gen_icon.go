package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

func main() {
	outPath := flag.String("o", "assets/icon.png", "output PNG path (should be 1024x1024)")
	size := flag.Int("s", 1024, "icon size in pixels")
	flag.Parse()

	if *size < 256 {
		fmt.Fprintln(os.Stderr, "size too small; use >= 256")
		os.Exit(2)
	}

	img := image.NewNRGBA(image.Rect(0, 0, *size, *size))

	bg := color.NRGBA{R: 245, G: 246, B: 248, A: 255}
	ink := color.NRGBA{R: 15, G: 18, B: 22, A: 255}
	accent := color.NRGBA{R: 42, G: 199, B: 105, A: 255}
	shadow := color.NRGBA{R: 0, G: 0, B: 0, A: 22}
	soft := color.NRGBA{R: 0, G: 0, B: 0, A: 10}

	fill(img, bg)

	// Subtle vignette / depth.
	addRadialShade(img, soft)

	// Central rotated square (the player).
	cx, cy := float64(*size)/2, float64(*size)/2
	playerSize := float64(*size) * 0.44
	drawRotSquare(img, cx, cy, playerSize*1.03, math.Pi/6, shadow) // soft shadow
	drawRotSquare(img, cx, cy, playerSize, math.Pi/6, ink)

	// A couple of smaller “food” squares.
	drawRotSquare(img, float64(*size)*0.72, float64(*size)*0.32, float64(*size)*0.16, -math.Pi/12, color.NRGBA{R: 90, G: 125, B: 255, A: 255})
	drawRotSquare(img, float64(*size)*0.28, float64(*size)*0.72, float64(*size)*0.12, math.Pi/10, color.NRGBA{R: 255, G: 152, B: 85, A: 255})

	// Green power-up ring.
	drawRing(img, float64(*size)*0.76, float64(*size)*0.72, float64(*size)*0.11, float64(*size)*0.02, accent)

	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	f, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(f, img); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("wrote", *outPath)
}

func fill(img *image.NRGBA, c color.NRGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		o := (y - b.Min.Y) * img.Stride
		for x := 0; x < b.Dx(); x++ {
			i := o + x*4
			img.Pix[i+0] = c.R
			img.Pix[i+1] = c.G
			img.Pix[i+2] = c.B
			img.Pix[i+3] = c.A
		}
	}
}

func addRadialShade(img *image.NRGBA, shade color.NRGBA) {
	b := img.Bounds()
	cx := float64(b.Dx()) / 2
	cy := float64(b.Dy()) / 2
	maxR := math.Hypot(cx, cy)

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			d := math.Hypot(float64(x)-cx, float64(y)-cy) / maxR
			// Stronger toward edges.
			a := uint8(float64(shade.A) * clamp01(math.Pow(d, 1.9)))
			blendOver(img, x, y, color.NRGBA{R: shade.R, G: shade.G, B: shade.B, A: a})
		}
	}
}

func drawRotSquare(img *image.NRGBA, cx, cy, size, angle float64, c color.NRGBA) {
	h := size / 2
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)

	// Conservative bounding box.
	r := h * math.Sqrt2
	minX := int(math.Floor(cx - r))
	maxX := int(math.Ceil(cx + r))
	minY := int(math.Floor(cy - r))
	maxY := int(math.Ceil(cy + r))

	b := img.Bounds()
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX > b.Dx()-1 {
		maxX = b.Dx() - 1
	}
	if maxY > b.Dy()-1 {
		maxY = b.Dy() - 1
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			// Rotate point into square's local space.
			lx := dx*cosA + dy*sinA
			ly := -dx*sinA + dy*cosA
			if math.Abs(lx) <= h && math.Abs(ly) <= h {
				blendOver(img, x, y, c)
			}
		}
	}
}

func drawRing(img *image.NRGBA, cx, cy, radius, thickness float64, c color.NRGBA) {
	rOuter := radius
	rInner := math.Max(0, radius-thickness)

	minX := int(math.Floor(cx - rOuter))
	maxX := int(math.Ceil(cx + rOuter))
	minY := int(math.Floor(cy - rOuter))
	maxY := int(math.Ceil(cy + rOuter))

	b := img.Bounds()
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX > b.Dx()-1 {
		maxX = b.Dx() - 1
	}
	if maxY > b.Dy()-1 {
		maxY = b.Dy() - 1
	}

	rOuter2 := rOuter * rOuter
	rInner2 := rInner * rInner
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			d2 := dx*dx + dy*dy
			if d2 <= rOuter2 && d2 >= rInner2 {
				blendOver(img, x, y, c)
			}
		}
	}
}

func blendOver(img *image.NRGBA, x, y int, src color.NRGBA) {
	i := y*img.Stride + x*4
	dr := img.Pix[i+0]
	dg := img.Pix[i+1]
	db := img.Pix[i+2]
	da := img.Pix[i+3]

	sa := float64(src.A) / 255.0
	daF := float64(da) / 255.0
	outA := sa + daF*(1-sa)
	if outA <= 0 {
		img.Pix[i+0] = 0
		img.Pix[i+1] = 0
		img.Pix[i+2] = 0
		img.Pix[i+3] = 0
		return
	}

	toByte := func(v float64) uint8 {
		v = math.Round(v)
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}

	// Non-premultiplied alpha blend.
	outR := (float64(src.R)*sa + float64(dr)*daF*(1-sa)) / outA
	outG := (float64(src.G)*sa + float64(dg)*daF*(1-sa)) / outA
	outB := (float64(src.B)*sa + float64(db)*daF*(1-sa)) / outA

	img.Pix[i+0] = toByte(outR)
	img.Pix[i+1] = toByte(outG)
	img.Pix[i+2] = toByte(outB)
	img.Pix[i+3] = toByte(outA * 255.0)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
