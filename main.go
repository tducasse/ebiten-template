package main

import (
	"embed"
	"errors"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tducasse/ebiten-template/camera"
	"github.com/tducasse/ebiten-template/collisions"
	"github.com/tducasse/ebiten-template/manager"
	"github.com/tducasse/ebiten-template/scenes"
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

type Game struct{}

var m *manager.Manager
var g *Game

//go:embed assets/*
var assetsFolder embed.FS

func (game *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errQuit
	}
	if ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		toggleFullscreen()
	}
	return m.Update()
}

func (game *Game) Draw(screen *ebiten.Image) {
	m.Draw(screen)
}

func (game *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func toggleFullscreen() {
	ebiten.SetFullscreen(!ebiten.IsFullscreen())
}

func main() {
	m = manager.MakeManager(map[string]*manager.Scene{
		"menu": scenes.MenuScene,
		"game": scenes.GameScene,
	}, "menu")

	g = &Game{}

	scenes.Context = &scenes.ContextType{
		Camera:         camera.Init(screenWidth, screenHeight),
		World:          ebiten.NewImage(worldWidth, worldHeight),
		CollisionWorld: collisions.MakeWorld(),
		AssetsFolder:   &assetsFolder,
		Manager:        m,
	}

	ebiten.SetWindowSize(screenWidth*windowScale, screenHeight*windowScale)
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(g); err != nil && err != errQuit {
		log.Fatal(err)
	}
}
