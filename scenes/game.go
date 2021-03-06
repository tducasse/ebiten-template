package scenes

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/tducasse/ebiten-input"
	ldtk "github.com/tducasse/ebiten-ldtk"
	manager "github.com/tducasse/ebiten-manager"
	signals "github.com/tducasse/ebiten-signals"
	"github.com/tducasse/ebiten-template/entities"
)

var levels *ldtk.Ldtk

var player *entities.Player

var message string

var GameScene *manager.Scene = &manager.Scene{

	// Init is run on switch to scene
	Init: func(setReady func()) {
		var err error
		levels, err = ldtk.Load(
			"sample.ldtk",
			&ldtk.Options{
				Aseprite:    true,
				EmbedFolder: Context.AssetsFolder,
				Root:        "assets/maps",
				CollidesWith: map[int]bool{
					0: true,
				},
				OnCollisionAdd: Context.CollisionWorld.AddNewBox,
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		entityOptions := entities.EntityOptions{
			EmbedFolder: Context.AssetsFolder,
			Root:        "assets/images",
			World:       Context.CollisionWorld,
		}

		keys := map[string][]ebiten.Key{
			"up":    {ebiten.KeyArrowUp, ebiten.KeyW},
			"down":  {ebiten.KeyArrowDown, ebiten.KeyS},
			"right": {ebiten.KeyArrowRight, ebiten.KeyD},
			"left":  {ebiten.KeyArrowLeft, ebiten.KeyA},
		}
		input.Init(keys)

		player = new(entities.Player)
		player.Init(levels, &entityOptions)

		setReady()
	},

	// Update runs once per tick
	Update: func(setReady func()) error {
		player.Update()
		Context.Camera.Follow.W, Context.Camera.Follow.H = player.Sprite.CurrentFrame.Image.Size()
		Context.Camera.X, Context.Camera.Y = player.X, player.Y

		signals.Connect("collided", func(i []interface{}) {
			message = ""
			for _, part := range i {
				message += part.(string) + " "
			}
		})

		setReady()

		return nil
	},

	// Draw runs once per frame
	Draw: func(screen *ebiten.Image) {
		Context.World.Fill(color.RGBA{R: 10, G: 10, B: 30, A: 255})
		levels.Draw(Context.World)
		player.Draw(Context.World)
		if message != "" {
			ebitenutil.DebugPrint(Context.World, message)
		}
		Context.Camera.Draw(Context.World, screen)
	},
}
