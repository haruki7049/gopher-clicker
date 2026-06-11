package main

import (
	// Standard libraries
	"bytes"
	"image/color"
	"log"

	// Externals
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"

	// Internals
	"github.com/haruki7049/gopher-clicker/assets"
)

const GAME_TITLE = "Gopher Clicker"
const GAME_HEIGHT = 480
const GAME_WIDTH = 640

func gopherColor() color.RGBA {
	return color.RGBA{R: 0x00, G: 0xad, B: 0xd8, A: 0xff}
}

type gopher struct {
	image  *ebiten.Image
	x      float64
	y      float64
	scaleX float64
	scaleY float64
}

type game struct {
	isTitle  bool
	ticks    int
	gopher   gopher
	fontFace *text.GoTextFaceSource
}

func newGame() (*game, error) {
	g := &game{}

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

func (g *game) Update() error {
	// Increment ticks and reset at 120 to prevent overflow
	g.ticks++
	if g.ticks >= 120 {
		g.ticks = 0
	}

	if g.isGopherClicked() {
		g.gopher.x += 100.0
		g.gopher.y += 100.0
	}

	return nil
}

func (g *game) isGopherClicked() bool {
	// Check if the left mouse button is just pressed
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}

	cx, cy := ebiten.CursorPosition()

	// Adjust cursor position to relative coordinates inside the image
	relX := float64(cx) - g.gopher.x
	relY := float64(cy) - g.gopher.y

	bounds := g.gopher.image.Bounds()
	w := float64(bounds.Dx()) * g.gopher.scaleX
	h := float64(bounds.Dy()) * g.gopher.scaleY

	// Check if the cursor is within the bounding box
	if relX < 0 || relY < 0 || relX >= w || relY >= h {
		return false
	}

	// Convert relative coordinates to local image coordinates
	localX := int(relX / g.gopher.scaleX)
	localY := int(relY / g.gopher.scaleY)

	// Get the color of the pixel
	_, _, _, a := g.gopher.image.At(localX, localY).RGBA()

	// Check if the alpha value is not zero (not transparent)
	return a > 0
}
func (g *game) Draw(screen *ebiten.Image) {
	// Fill the screen with Cyan Blue (Gopher's color!!)
	screen.Fill(gopherColor())

	// Draw Gopher image
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.gopher.scaleX, g.gopher.scaleY)
		op.GeoM.Translate(g.gopher.x, g.gopher.y)
		screen.DrawImage(g.gopher.image, op)
	}

	// Draw title
	{
		face := &text.GoTextFace{
			Source: g.fontFace,
			Size:   24,
		}

		_, h := text.Measure(GAME_TITLE, face, face.Size)

		op := &text.DrawOptions{}
		op.LayoutOptions = text.LayoutOptions{LineSpacing: h, PrimaryAlign: text.AlignCenter, SecondaryAlign: text.AlignCenter}
		op.GeoM.Translate(GAME_WIDTH/2, GAME_HEIGHT/3*2)

		// Draw
		text.Draw(screen, GAME_TITLE+"\nClick and Start!!", face, op)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GAME_WIDTH, GAME_HEIGHT
}

func main() {
	ebiten.SetWindowSize(GAME_WIDTH, GAME_HEIGHT)
	ebiten.SetWindowTitle(GAME_TITLE)
	g, err := newGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
