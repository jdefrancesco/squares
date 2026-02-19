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

func drawHUD(screen *ebiten.Image, x, y int, score int, invLeft float64) {
	textStr := fmt.Sprintf(
		"Squares eaten: %d\nInvincible: %.1fs\n\nGreen circle: invincibility\nBlue circle: instant death",
		score,
		math.Max(0, invLeft),
	)

	lines := 5
	maxLen := 0
	cur := 0
	for _, r := range textStr {
		if r == '\n' {
			if cur > maxLen {
				maxLen = cur
			}
			cur = 0
		} else {
			cur++
		}
	}
	if cur > maxLen {
		maxLen = cur
	}

	pad := 10
	w := pad*2 + maxLen*7
	h := pad*2 + lines*13

	panel := ebiten.NewImage(w, h)
	panel.Fill(hudPanel)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(panel, op)

	drawRect(screen, x, y, w, 1, hudBorder)
	drawRect(screen, x, y+h-1, w, 1, hudBorder)
	drawRect(screen, x, y, 1, h, hudBorder)
	drawRect(screen, x+w-1, y, 1, h, hudBorder)

	text.Draw(screen, textStr, hudFace, x+pad+1, y+pad+13+1, color.RGBA{0, 0, 0, 90})
	text.Draw(screen, textStr, hudFace, x+pad, y+pad+13, hudColor)
}

func pauseButtonRect() (x, y, w, h int) {
	pad := 12
	w = 96
	h = 26
	x = ScreenWidth - pad - w
	y = pad
	return
}

func drawPauseButton(screen *ebiten.Image, paused bool) {
	x, y, w, h := pauseButtonRect()

	mx, my := ebiten.CursorPosition()
	mx = clampInt(mx, 0, ScreenWidth-1)
	my = clampInt(my, 0, ScreenHeight-1)
	hover := mx >= x && mx < x+w && my >= y && my < y+h

	bg := color.RGBA{255, 255, 255, 235}
	if hover {
		bg = color.RGBA{250, 250, 250, 245}
	}

	drawRect(screen, x, y, w, h, bg)
	drawRect(screen, x, y, w, 1, hudBorder)
	drawRect(screen, x, y+h-1, w, 1, hudBorder)
	drawRect(screen, x, y, 1, h, hudBorder)
	drawRect(screen, x+w-1, y, 1, h, hudBorder)

	label := "Pause"
	if paused {
		label = "Resume"
	}
	text.Draw(screen, label, hudFace, x+12, y+18, hudColor)
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
