package game

import (
	// Standard libraries
	"bytes"
	"image/color"
	"math/rand/v2"
	"strconv"

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

type Gopher struct {
	image  *ebiten.Image
	x      float64
	y      float64
	scaleX float64
	scaleY float64
}

type Game struct {
	states   states
	gopher   Gopher
	fontFace *text.GoTextFaceSource
}

type states struct {
	inTitle bool
	ticks   int
	score   int
}

func NewGame() (*Game, error) {
	g := &Game{}

	if err := g.newGameGopher(); err != nil {
		return nil, err
	}

	if err := g.newGameFont(); err != nil {
		return nil, err
	}

	g.newStates()

	return g, nil
}

func (g *Game) newStates() {
	g.states.inTitle = true
}

func (g *Game) newGameGopher() error {
	// Load gopher image
	gopher_img, _, err := ebitenutil.NewImageFromFileSystem(assets.Assets, "images/gopher.png")
	if err != nil {
		return err
	}
	g.gopher.image = gopher_img

	// Set initial gopher scale
	g.gopher.scaleX = 1.0
	g.gopher.scaleY = 1.0

	return nil
}

func (g *Game) newGameFont() error {
	// Set the standard Go font TTF data
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return err
	}

	// Set g.fontFace
	g.fontFace = s

	return nil
}

func gopherColor() color.RGBA {
	return color.RGBA{R: 0x00, G: 0xad, B: 0xd8, A: 0xff}
}

func (g *Game) Update() error {
	g.updateTicks()

	if g.isGopherClicked() {
		g.randomizeGopherPosition()
		g.states.inTitle = false
		g.states.score += 1
	}

	return nil
}

func (g *Game) updateTicks() {
	// Increment ticks and reset at 120 to prevent overflow
	g.states.ticks++
	if g.states.ticks >= 120 {
		g.states.ticks = 0
	}
}

func (g *Game) randomizeGopherPosition() {
	bounds := g.gopher.image.Bounds()
	w := float64(bounds.Dx()) * g.gopher.scaleX
	h := float64(bounds.Dy()) * g.gopher.scaleY

	maxX := float64(GAME_WIDTH) - w
	maxY := float64(GAME_HEIGHT) - h

	// Generate random coordinates within the screen bounds
	g.gopher.x = rand.Float64() * maxX
	g.gopher.y = rand.Float64() * maxY
}

func (g *Game) isGopherClicked() bool {
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

func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with Cyan Blue (Gopher's color!!)
	screen.Fill(gopherColor())

	g.drawGopher(screen)

	if g.states.inTitle {
		g.drawTitle(screen)
	} else {
		g.drawScore(screen)
	}
}

// Draw Gopher
func (g *Game) drawGopher(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.gopher.scaleX, g.gopher.scaleY)
	op.GeoM.Translate(g.gopher.x, g.gopher.y)
	screen.DrawImage(g.gopher.image, op)
}

func (g *Game) drawScore(screen *ebiten.Image) {
	face := &text.GoTextFace{
		Source: g.fontFace,
		Size:   24,
	}

	_, h := text.Measure(GAME_TITLE, face, face.Size)

	op := &text.DrawOptions{}
	op.LayoutOptions = text.LayoutOptions{LineSpacing: h, PrimaryAlign: text.AlignCenter, SecondaryAlign: text.AlignCenter}
	op.GeoM.Translate(GAME_WIDTH/2, GAME_HEIGHT/3*2)

	// Draw
	text.Draw(screen, strconv.Itoa(g.states.score), face, op)
}

// Draw title
func (g *Game) drawTitle(screen *ebiten.Image) {
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GAME_WIDTH, GAME_HEIGHT
}
