package main

import (
	"embed"
	"errors"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tducasse/ebiten-template/entities"
	"github.com/tducasse/ebiten-template/input"
	"github.com/tducasse/ebiten-template/ldtk"
)

var errQuit = errors.New("quit")

type Game struct {
}

//go:embed assets/*
var assetsFolder embed.FS

var levels *ldtk.Ldtk

var player *entities.Player

func init() {
	var err error
	levels, err = ldtk.Load(
		"sample.ldtk",
		&ldtk.Options{
			Aseprite:    true,
			EmbedFolder: &assetsFolder,
			Root:        "assets/maps",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	entityOptions := entities.EntityOptions{
		EmbedFolder: &assetsFolder,
		Root:        "assets/images",
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
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	levels.Draw(screen)
	player.Draw(screen)
}

func (game *Game) Layout(w, h int) (int, int) { return 256, 144 }

func main() {
	ebiten.SetWindowSize(1024, 576)
	if err := ebiten.RunGame(&Game{}); err != nil && err != errQuit {
		log.Fatal(err)
	}
}
