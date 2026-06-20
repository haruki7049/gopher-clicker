package main

import (
	// Standard libraries
	"log"

	// Externals
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/haruki7049/gopher-clicker/internal/game"
)

func main() {
	ebiten.SetWindowSize(game.GAME_WIDTH, game.GAME_HEIGHT)
	ebiten.SetWindowTitle(game.GAME_TITLE)
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
