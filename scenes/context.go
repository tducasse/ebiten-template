package scenes

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tducasse/ebiten-template/camera"
	"github.com/tducasse/ebiten-template/collision"
	"github.com/tducasse/ebiten-template/manager"
)

type ContextType struct {
	Camera         *camera.Camera
	World          *ebiten.Image
	CollisionWorld *collision.World
	AssetsFolder   *embed.FS
	Manager        *manager.Manager
}

var Context *ContextType
