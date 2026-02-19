package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jdefrancesco/squares/internal/game"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Squares")

	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
