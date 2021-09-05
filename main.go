package main

import (
	"embed"
	"errors"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tducasse/ebiten-template/camera"
	"github.com/tducasse/ebiten-template/collision"
	"github.com/tducasse/ebiten-template/entities"
	"github.com/tducasse/ebiten-template/input"
	"github.com/tducasse/ebiten-template/ldtk"
)

const (
	screenWidth  = 256
	screenHeight = 144
	windowScale  = 4
)

const (
	worldWidth  = 512
	worldHeight = 512
)

var errQuit = errors.New("quit")

type Game struct {
	Camera         *camera.Camera
	World          *ebiten.Image
	CollisionWorld *collision.World
}

//go:embed assets/*
var assetsFolder embed.FS

var levels *ldtk.Ldtk

var player *entities.Player

func (g *Game) Init() {
	var err error
	levels, err = ldtk.Load(
		"sample.ldtk",
		&ldtk.Options{
			Aseprite:    true,
			EmbedFolder: &assetsFolder,
			Root:        "assets/maps",
			CollidesWith: map[int]bool{
				0: true,
			},
			OnCollisionAdd: g.CollisionWorld.AddNewBox,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	entityOptions := entities.EntityOptions{
		EmbedFolder: &assetsFolder,
		Root:        "assets/images",
		World:       g.CollisionWorld,
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
}

func (game *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errQuit
	}
	player.Update()
	game.Camera.Follow.W, game.Camera.Follow.H = player.Sprite.CurrentFrame.Image.Size()
	game.Camera.X, game.Camera.Y = player.X, player.Y

	if ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		toggleFullscreen()
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.World.Clear()
	levels.Draw(game.World)
	player.Draw(game.World)

	game.Camera.Draw(game.World, screen)
}

func (game *Game) Layout(w, h int) (int, int) { return screenWidth, screenHeight }

func toggleFullscreen() {
	ebiten.SetFullscreen(!ebiten.IsFullscreen())
}

func main() {
	g := &Game{
		Camera:         camera.Init(screenWidth, screenHeight),
		World:          ebiten.NewImage(worldWidth, worldHeight),
		CollisionWorld: collision.MakeWorld(),
	}
	g.Init()
	ebiten.SetWindowSize(screenWidth*windowScale, screenHeight*windowScale)
	if err := ebiten.RunGame(g); err != nil && err != errQuit {
		log.Fatal(err)
	}
}
