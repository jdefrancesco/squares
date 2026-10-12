package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

var hudFace = basicfont.Face7x13

// drawHUD displays the current score and personal best, briefly tinting the
// score after eating to reinforce its increase.
func drawHUD(screen *ebiten.Image, x, y, score, best int, pulse float64) {
	scoreCol := color.RGBA{35, 35, 35, 255}
	if pulse > 0 {
		scoreCol = color.RGBA{155, 99, 13, 255}
	}
	text.Draw(screen, fmt.Sprintf("Score: %d", score), hudFace, x+10, y+17, scoreCol)
	text.Draw(screen, fmt.Sprintf("Best:  %d", best), hudFace, x+10, y+34, color.RGBA{95, 100, 110, 255})
}

// drawLegend explains the two square colors and the mouse and pause controls.
func drawLegend(screen *ebiten.Image) {
	y := ScreenHeight - 18
	drawSquareAA(screen, 26, float64(y-4), 8, foodColor)
	text.Draw(screen, "EAT", hudFace, 37, y, color.RGBA{110, 80, 25, 255})
	drawSquareAA(screen, 85, float64(y-4), 10, dangerColor)
	text.Draw(screen, "AVOID", hudFace, 98, y, dangerColor)
	text.Draw(screen, "Move mouse | P / Esc: pause", hudFace, ScreenWidth-205, y, color.RGBA{95, 100, 110, 255})
}

// drawCenteredText draws one line with a baseline at y, scaling around the
// logical center so prompts remain aligned as their text changes.
func drawCenteredText(screen *ebiten.Image, msg string, y int, scale float64, c color.Color) {
	bounds := text.BoundString(hudFace, msg)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((ScreenWidth-float64(bounds.Dx())*scale)/2, float64(y))
	op.ColorScale.ScaleWithColor(c)
	text.DrawWithOptions(screen, msg, hudFace, op)
}

// drawDimOverlay shades the playfield using the shared texture.
func drawDimOverlay(screen *ebiten.Image, alpha uint8) {
	drawSquareRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, alpha})
}

// drawPauseOverlay dims the playfield and displays the pause and resume prompt.
func drawPauseOverlay(screen *ebiten.Image) {
	drawDimOverlay(screen, 160)
	drawCenteredText(screen, "PAUSED", ScreenHeight/2-10, 2, color.White)
	drawCenteredText(screen, "Press P or Esc to resume", ScreenHeight/2+22, 1, color.White)
}

// drawGameOverOverlay shows the round result, highlights a new personal best,
// and offers a one-click retry without an additional menu.
func drawGameOverOverlay(screen *ebiten.Image, score, best int, newBest bool) {
	drawDimOverlay(screen, 185)
	drawCenteredText(screen, "GAME OVER", ScreenHeight/2-44, 2, color.White)
	drawCenteredText(screen, fmt.Sprintf("Score: %d    Best: %d", score, best), ScreenHeight/2-12, 1, color.White)
	if newBest {
		drawCenteredText(screen, "NEW BEST!", ScreenHeight/2+12, 1, foodColor)
	}
	drawCenteredText(screen, "Click to try again", ScreenHeight/2+45, 1, color.White)
	drawCenteredText(screen, "R: retry   Q: quit", ScreenHeight/2+67, 1, color.RGBA{205, 210, 220, 255})
}
