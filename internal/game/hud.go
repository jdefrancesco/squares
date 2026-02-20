package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

var (
	hudFace   = basicfont.Face7x13
	hudColor  = color.RGBA{20, 20, 20, 255}
	hudPanel  = color.RGBA{255, 255, 255, 220}
	hudBorder = color.RGBA{0, 0, 0, 70}
)

type hudLine struct {
	label string
	value string
}

func drawHUD(screen *ebiten.Image, x, y int, score int, invLeft float64) {
	lines := []hudLine{
		{label: "Squares eaten:", value: fmt.Sprintf("%d", score)},
		{label: "Invincible:", value: fmt.Sprintf("%.1fs", math.Max(0, invLeft))},
		{label: "Pause:", value: "p"},
		{label: "Quit:", value: "q"},
		{},
		{label: "Green circle:", value: "invincibility"},
		{label: "Black circle:", value: "loss"},
	}

	lineHeight := hudFace.Metrics().Height.Ceil()
	ascent := hudFace.Metrics().Ascent.Ceil()

	maxW := 0
	for _, line := range lines {
		if line.label == "" && line.value == "" {
			continue
		}
		full := line.label
		if line.value != "" {
			full = line.label + " " + line.value
		}
		b := text.BoundString(hudFace, full)
		if b.Dx() > maxW {
			maxW = b.Dx()
		}
	}

	pad := 10
	w := pad*2 + maxW
	h := pad*2 + len(lines)*lineHeight

	panel := ebiten.NewImage(w, h)
	panel.Fill(hudPanel)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(panel, op)

	drawRect(screen, x, y, w, 1, hudBorder)
	drawRect(screen, x, y+h-1, w, 1, hudBorder)
	drawRect(screen, x, y, 1, h, hudBorder)
	drawRect(screen, x+w-1, y, 1, h, hudBorder)

	labelCol := color.RGBA{5, 5, 5, 255}
	valueCol := color.RGBA{35, 35, 35, 255}
	shadowCol := color.RGBA{0, 0, 0, 70}

	baseY := y + pad + ascent
	for i, line := range lines {
		if line.label == "" && line.value == "" {
			continue
		}

		ly := baseY + i*lineHeight
		x0 := x + pad

		// Label: simulate bold by drawing twice with a 1px offset.
		text.Draw(screen, line.label, hudFace, x0+1, ly+1, shadowCol)
		text.Draw(screen, line.label, hudFace, x0, ly, labelCol)
		text.Draw(screen, line.label, hudFace, x0+1, ly, labelCol)

		if line.value != "" {
			labelW := text.BoundString(hudFace, line.label+" ").Dx()
			text.Draw(screen, line.value, hudFace, x0+labelW, ly, valueCol)
		}
	}
}

func drawPauseOverlay(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 70})
	screen.DrawImage(overlay, nil)

	text.Draw(screen, "PAUSED", hudFace, ScreenWidth/2-21, ScreenHeight/2-8, color.RGBA{255, 255, 255, 230})
	text.Draw(screen, "Press P or Esc to resume", hudFace, ScreenWidth/2-77, ScreenHeight/2+12, color.RGBA{255, 255, 255, 220})
}

func drawRect(dst *ebiten.Image, x, y, w, h int, c color.RGBA) {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(img, op)
}
