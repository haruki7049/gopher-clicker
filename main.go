package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/haruki7049/gopher-clicker/assets"
)

func gopherColor() color.RGBA {
	return color.RGBA{R: 0x00, G: 0xad, B: 0xd8, A: 0xff}
}

type Gopher struct {
	image  *ebiten.Image
	x      int
	y      int
	scaleX float64
	scaleY float64
}

type Game struct {
	isTitle bool
	ticks   int
	gopher  Gopher
}

func newGame() (*Game, error) {
	g := &Game{}

	// Load gopher image
	gopher_img, _, err := ebitenutil.NewImageFromFileSystem(assets.Assets, "images/gopher.png")
	if err != nil {
		return nil, err
	}
	g.gopher.image = gopher_img

	// Set initial gopher scale
	g.gopher.scaleX = 0.5
	g.gopher.scaleY = 0.5

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
	// Fill the screen with Cyan Blue (Gopher's color!!)
	screen.Fill(gopherColor())

	// Draw Gopher image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.gopher.scaleX, g.gopher.scaleY)
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
