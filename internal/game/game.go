package game

import (
	// Standard libraries
	"bytes"
	"image/color"
	"math/rand/v2"
	"strconv"

	// Externals
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"

	// Internals
	"github.com/haruki7049/gopher-clicker/assets"
)

const GAME_TITLE = "Gopher Clicker"
const GAME_HEIGHT = 480
const GAME_WIDTH = 640

// CLICK_TIMEOUT_SECONDS is the time limit (in seconds) the player can go
// without clicking the gopher before the score is reset.
const CLICK_TIMEOUT_SECONDS = 1.5

const gaugeWidth float32 = 200
const gaugeHeight float32 = 12
const gaugeY float32 = 16

type gopher struct {
	image  *ebiten.Image
	x      float64
	y      float64
	scaleX float64
	scaleY float64
}

type inputPosition struct {
	x int
	y int
}

type successSe struct {
	player *audio.Player
}

func (se *successSe) Play() {
	se.player.Rewind()
	se.player.Play()
}

type missSe struct {
	player *audio.Player
}

func (se *missSe) Play() {
	se.player.Rewind()
	se.player.Play()
}

type Se interface {
	Play()
}

func newSuccessSe(g *Game) (*successSe, error) {
	var successSe successSe

	// Open the WAV file from embedded assets
	file, err := assets.Assets.Open("se/success.wav")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the WAV file
	decoded, err := wav.DecodeWithSampleRate(44100, file)
	if err != nil {
		return nil, err
	}

	// Create a new audio player
	player, err := g.audioCtx.NewPlayer(decoded)
	if err != nil {
		return nil, err
	}

	successSe.player = player
	return &successSe, nil
}

func newMissSe(g *Game) (*missSe, error) {
	var missSe missSe

	// Open the WAV file from embedded assets
	file, err := assets.Assets.Open("se/miss.wav")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the WAV file
	decoded, err := wav.DecodeWithSampleRate(44100, file)
	if err != nil {
		return nil, err
	}

	// Create a new audio player
	player, err := g.audioCtx.NewPlayer(decoded)
	if err != nil {
		return nil, err
	}

	missSe.player = player
	return &missSe, nil
}

type Game struct {
	states    states
	gopher    gopher
	fontFace  *text.GoTextFaceSource
	audioCtx  *audio.Context
	successSe successSe
	missSe    missSe
}

type states struct {
	inTitle         bool
	ticks           int
	score           int
	ticksSinceClick int
}

func NewGame() (*Game, error) {
	g := &Game{}

	g.newAudioCtx()

	if err := g.newGameGopher(); err != nil {
		return nil, err
	}

	if err := g.newGameFont(); err != nil {
		return nil, err
	}

	// Success SE
	successSe, err := newSuccessSe(g)
	if err != nil {
		return nil, err
	}
	g.successSe = *successSe

	// Miss SE
	missSe, err := newMissSe(g)
	if err != nil {
		return nil, err
	}
	g.missSe = *missSe

	g.newStates()

	return g, nil
}

func (g *Game) newAudioCtx() {
	// Initialize audio context with 44100 sample rate
	g.audioCtx = audio.NewContext(44100)

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
		g.states.ticksSinceClick = 0

		// Success SE play
		g.successSe.Play()
	}

	g.updateClickTimeout()

	return nil
}

func (g *Game) updateTicks() {
	// Increment ticks and reset at 120 to prevent overflow
	g.states.ticks++
	if g.states.ticks >= 120 {
		g.states.ticks = 0
	}
}

func (g *Game) updateClickTimeout() {
	if g.states.inTitle {
		return
	}

	g.states.ticksSinceClick += 1

	timeoutTicks := g.clickTimeoutTicks()
	if g.states.ticksSinceClick >= timeoutTicks {
		if g.states.score != 0 {
			// Miss SE play
			g.missSe.Play()
		}

		g.states.score = 0
	}
}

func (g *Game) clickTimeoutTicks() int {
	return int(CLICK_TIMEOUT_SECONDS * float64(ebiten.TPS()))
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

func (g *Game) justPressedPositions() []inputPosition {
	var positions []inputPosition

	// Mouse button
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		positions = append(positions, inputPosition{x: cx, y: cy})
	}

	// Touchpad (for mobile devices)
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		tx, ty := ebiten.TouchPosition(id)
		positions = append(positions, inputPosition{x: tx, y: ty})
	}

	return positions
}

func (g *Game) isPositionOnGopher(cx, cy int) bool {
	// Adjust position to relative coordinates inside the image
	relX := float64(cx) - g.gopher.x
	relY := float64(cy) - g.gopher.y

	bounds := g.gopher.image.Bounds()
	w := float64(bounds.Dx()) * g.gopher.scaleX
	h := float64(bounds.Dy()) * g.gopher.scaleY

	// Check if the position is within the bounding box
	if relX < 0 || relY < 0 || relX >= w || relY >= h {
		return false
	}

	// Convert relative coordinates to local image coordinates
	localX := int(relX / g.gopher.scaleX)
	localY := int(relY / g.gopher.scaleY)

	// Get RGBA value of the pixel
	_, _, _, a := g.gopher.image.At(localX, localY).RGBA()

	// Check if the alpha value is not zero (not transparent)
	return a > 0
}

func (g *Game) isGopherClicked() bool {
	for _, pos := range g.justPressedPositions() {
		if g.isPositionOnGopher(pos.x, pos.y) {
			return true
		}
	}

	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with Cyan Blue (Gopher's color!!)
	screen.Fill(gopherColor())

	g.drawGopher(screen)
	g.drawClickGauge(screen)

	if g.states.inTitle {
		g.drawTitle(screen)
	} else {
		g.drawScore(screen)
	}
}

func (g *Game) drawClickGauge(screen *ebiten.Image) {
	if g.states.inTitle {
		return
	}

	timeoutTicks := g.clickTimeoutTicks()
	remainingTicks := timeoutTicks - g.states.ticksSinceClick
	if remainingTicks <= 0 {
		return
	}

	ratio := float32(remainingTicks) / float32(timeoutTicks)
	gaugeX := (float32(GAME_WIDTH) - gaugeWidth) / 2

	vector.FillRect(screen, gaugeX, gaugeY, gaugeWidth, gaugeHeight, color.RGBA{A: 0x80}, false)

	vector.FillRect(screen, gaugeX, gaugeY, gaugeWidth*ratio, gaugeHeight, color.White, false)
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
	op.ColorScale.ScaleWithColor(color.Black)
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
	op.ColorScale.ScaleWithColor(color.Black)
	op.LayoutOptions = text.LayoutOptions{LineSpacing: h, PrimaryAlign: text.AlignCenter, SecondaryAlign: text.AlignCenter}
	op.GeoM.Translate(GAME_WIDTH/2, GAME_HEIGHT/3*2)

	// Draw
	text.Draw(screen, GAME_TITLE+"\nClick and Start!!", face, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GAME_WIDTH, GAME_HEIGHT
}
