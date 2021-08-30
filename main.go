package main

import (
	"bytes"
	"embed"
	"errors"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tducasse/ebiten-template/ldtk"
)

var errQuit = errors.New("quit")

type Game struct {
}

//go:embed assets/*
var assetsFolder embed.FS

var img *ebiten.Image

var levels *ldtk.Ldtk

func init() {
	imgData, err := assetsFolder.ReadFile("assets/images/player.png")
	if err != nil {
		log.Fatal(err)
	}

	imgFile, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatal(err)
	}

	img = ebiten.NewImageFromImage(imgFile)

	mapData, err := assetsFolder.ReadFile("assets/maps/sample.ldtk")
	if err != nil {
		log.Fatal(err)
	}

	levels, err = ldtk.Load(
		mapData,
		&ldtk.Options{
			Aseprite:    true,
			EmbedFolder: &assetsFolder,
			FilePrefix:  "assets/maps",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (game *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errQuit
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	levels.Draw(screen)
}

func (game *Game) Layout(w, h int) (int, int) { return 256, 144 }

func main() {
	ebiten.SetWindowSize(1024, 576)
	if err := ebiten.RunGame(&Game{}); err != nil && err != errQuit {
		log.Fatal(err)
	}
}
