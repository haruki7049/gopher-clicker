package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/haruki7049/gopher-clicker/assets"
)

type Gopher struct {
	image *ebiten.Image
	x     int
	y     int
}

type Game struct {
	ticks  int
	gopher Gopher
}

func newGame() (*Game, error) {
	g := &Game{}

	// Load gopher image
	gopher_img, _, err := ebitenutil.NewImageFromFileSystem(assets.Assets, "images/gopher.png")
	if err != nil {
		return nil, err
	}
	g.gopher.image = gopher_img

	return g, nil
}

func (g *Game) Update() error {
	// Increment ticks and reset at 120 to prevent overflow
	g.ticks++
	if g.ticks >= 120 {
		g.ticks = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with the calculated color
	screen.Fill(color.RGBA{R: 0x00, G: 0xad, B: 0xd8, A: 0xff})

	// Draw Gopher image
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.gopher.image, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	g, err := newGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
