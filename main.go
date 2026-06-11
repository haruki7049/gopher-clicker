package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	ticks int
}

func (g *Game) Update() error {
	// Increment ticks and reset at 120 to prevent overflow
	g.ticks++
	if g.ticks >= 120 {
		g.ticks = 0
	}

	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with the calculated color
	screen.Fill(color.RGBA{R: 0x00, G: 0xad, B: 0xd8, A: 0xff})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
