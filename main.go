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
	// Determine the current step in the 120-tick cycle
	cycle := game.ticks
	var r uint8
	var g uint8
	var b uint8

	if cycle < 60 {
		// Fade from black to white (0 to 255)
		r = uint8((cycle * 255) / 60)
		g = uint8((cycle * 255) / 30)
		b = uint8((cycle * 255) / 20)
	} else {
		// Fade from white to black (255 to 0)
		r = uint8(((120 - cycle) * 255) / 60)
		g = uint8(((120 - cycle) * 255) / 30)
		b = uint8(((120 - cycle) * 255) / 20)
	}

	// Fill the screen with the calculated color
	screen.Fill(color.RGBA{R: r, G: g, B: b, A: 255})
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
