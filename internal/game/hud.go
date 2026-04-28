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
	hudFace = basicfont.Face7x13
)

type hudLine struct {
	label string
	value string
}

func drawHUD(screen *ebiten.Image, x, y int, score int, invLeft float64) {
	lines := []hudLine{
		{label: "Score:", value: fmt.Sprintf("%d", score)},
		{label: "Invincible:", value: fmt.Sprintf("%.1fs", math.Max(0, invLeft))},
		{},
		{label: "Green:", value: "invincibility"},
		{label: "Black:", value: "loss"},
	}

	lineHeight := hudFace.Metrics().Height.Ceil()
	ascent := hudFace.Metrics().Ascent.Ceil()

	pad := 10

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
