package entities

import (
	"bytes"
	"embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	collisions "github.com/tducasse/ebiten-collisions"
)

type EntityOptions struct {
	EmbedFolder *embed.FS
	Root        string
	World       *collisions.World
}

func NewImageFromEmbed(path string, folder *embed.FS) *ebiten.Image {
	imgData, err := folder.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	imgFile, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(imgFile)
}
