package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tducasse/ebiten-template/manager"
)

var MenuScene *manager.Scene = &manager.Scene{

	// Init is run on switch to scene
	Init: func(setReady func()) {
		setReady()
	},

	// Update runs once per tick
	Update: func(setReady func()) error {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			Context.Manager.SwitchTo("game")
		}

		setReady()
		return nil
	},

	// Draw runs once per frame
	Draw: func(screen *ebiten.Image) {
		ebitenutil.DebugPrint(screen, "Press ENTER to move to the game")
	},
}
