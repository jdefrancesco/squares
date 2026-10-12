package game

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// renderCheck runs the real renderer through active, pause, and retry views.
type renderCheck struct {
	t        *testing.T
	g        *Game
	view     int
	prepared int
	dir      string
}

// testLoop runs tests within Ebiten's graphics loop while RunGame remains on the
// OS main thread, as required by macOS window creation.
type testLoop struct {
	m    *testing.M
	code int
}

// Update runs the test suite with graphics available, then ends the loop.
func (g *testLoop) Update() error {
	g.code = g.m.Run()
	return ebiten.Termination
}

// Draw is empty because individual rendering tests draw to off-screen images.
func (g *testLoop) Draw(screen *ebiten.Image) {}

// Layout gives the test window the same logical dimensions as the game.
func (g *testLoop) Layout(width, height int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// TestMain runs the graphics loop on the main OS thread and reports suite failures.
func TestMain(m *testing.M) {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Squares - UI check")
	g := &testLoop{m: m, code: 1}
	if err := ebiten.RunGame(g); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	os.Exit(g.code)
}

// Update prepares each view once and exits after all screenshots are captured.
func (r *renderCheck) Update() error {
	if r.view >= 4 {
		return ebiten.Termination
	}
	if r.prepared == r.view {
		return nil
	}
	r.prepared = r.view
	if r.view == 0 {
		// Use a stable scene so screenshots do not depend on the real cursor.
		r.g.step(0.01, inputState{x: 320, y: 240})
		r.g.player = Entity{x: 320, y: 280, size: 32}
		r.g.angle, r.g.score, r.g.best = 0.3, 7, 12
		r.g.eatPulse, r.g.scorePopupLeft = 0.09, 0.4
		r.g.scorePopupX, r.g.scorePopupY = 320, 240
		r.g.ents = []Entity{
			{x: 150, y: 180, size: 12},
			{x: 440, y: 180, size: 48},
			{x: 320, y: -120, size: 35, vy: 100},
		}
		r.g.playEatSound()
	}
	r.g.paused = r.view == 1
	r.g.gameOver = r.view >= 2
	r.g.newBest = r.view == 3
	if r.view == 3 {
		r.g.best = r.g.score
		r.g.playEatSound() // Reuse the same sound player on successive eats.
	}
	return nil
}

// Draw captures the rendered pixels and checks visible food, danger, and prompts.
func (r *renderCheck) Draw(screen *ebiten.Image) {
	r.g.Draw(screen)
	img := image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))
	screen.ReadPixels(img.Pix)
	if r.view == 0 {
		if got := img.RGBAAt(150, 180); got != foodColor {
			r.t.Errorf("food pixel=%v; want %v", got, foodColor)
		}
		if got := img.RGBAAt(440, 180); got != dangerColor {
			r.t.Errorf("danger pixel=%v; want %v", got, dangerColor)
		}
		if img.RGBAAt(320, 5) == (color.RGBA{245, 246, 248, 255}) {
			r.t.Error("incoming danger marker should be visible")
		}
		if countDarkPixels(img, image.Rect(20, 13, 100, 48)) < 30 {
			r.t.Error("Score and Best should be visible")
		}
		if countDarkPixels(img, image.Rect(20, 448, 140, 470)) < 30 {
			r.t.Error("EAT and AVOID legend should be visible")
		}
	} else if r.view >= 2 {
		if countLightPixels(img, image.Rect(240, 273, 400, 288)) < 20 {
			r.t.Error("click-to-retry prompt should be visible")
		}
	}
	names := []string{"active", "paused", "game-over", "new-best"}
	f, err := os.Create(filepath.Join(r.dir, names[r.view]+".png"))
	if err != nil {
		r.t.Error(err)
	} else {
		if err := png.Encode(f, img); err != nil {
			r.t.Error(err)
		}
		if err := f.Close(); err != nil {
			r.t.Error(err)
		}
	}
	r.view++
}

// Layout delegates to the game's fixed-size layout.
func (r *renderCheck) Layout(width, height int) (int, int) {
	return r.g.Layout(width, height)
}

// countDarkPixels measures visible text and markers against the light playfield.
func countDarkPixels(img *image.RGBA, bounds image.Rectangle) int {
	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.R < 150 && c.G < 150 && c.B < 150 {
				count++
			}
		}
	}
	return count
}

// countLightPixels measures prompt text over the dark game-over overlay.
func countLightPixels(img *image.RGBA, bounds image.Rectangle) int {
	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.R > 220 && c.G > 220 && c.B > 220 {
				count++
			}
		}
	}
	return count
}

// TestArcadeRendering validates the actual graphics and audio path and captures
// the visible score, legend, warning, pause, and retry surfaces for inspection.
func TestArcadeRendering(t *testing.T) {
	dir := os.Getenv("SQUARES_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	g := New()
	g.best = 0
	g.scorePath = filepath.Join(t.TempDir(), "best.json")
	r := &renderCheck{t: t, g: g, prepared: -1, dir: dir}
	width, height := r.Layout(800, 600)
	if width != ScreenWidth || height != ScreenHeight {
		t.Fatal("the playfield should keep its logical dimensions")
	}
	screen := ebiten.NewImage(width, height)
	for r.view < 4 {
		if err := r.Update(); err != nil {
			t.Fatal(err)
		}
		r.Draw(screen)
	}
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
	if r.view != 4 {
		t.Fatal("not all game views were rendered")
	}
}
