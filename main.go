package main

import (
	// Standard libraries
	"bytes"
	"image/color"
	"log"

	// Externals
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"

	// Internals
	"github.com/haruki7049/gopher-clicker/assets"
)

const GAME_HEIGHT = 480
const GAME_WIDTH = 640

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
	isTitle  bool
	ticks    int
	gopher   Gopher
	fontFace *text.GoTextFaceSource
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
	g.gopher.scaleX = 1.0
	g.gopher.scaleY = 1.0

	// Set the standard Go font TTF data
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	g.fontFace = s

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
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.gopher.scaleX, g.gopher.scaleY)
		screen.DrawImage(g.gopher.image, op)
	}

	// Draw title
	{
		face := &text.GoTextFace{
			Source: g.fontFace,
			Size:   24,
		}
		op := &text.DrawOptions{}
		op.GeoM.Translate(GAME_WIDTH/2, GAME_HEIGHT/3*2)

		// Draw
		text.Draw(screen, "Gopher Clicker", face, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GAME_WIDTH, GAME_HEIGHT
}

func main() {
	ebiten.SetWindowSize(GAME_WIDTH, GAME_HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")
	g, err := newGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
