package scenes

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/tducasse/ebiten-camera"
	collisions "github.com/tducasse/ebiten-collisions"
	manager "github.com/tducasse/ebiten-manager"
)

type ContextType struct {
	Camera         *camera.Camera
	World          *ebiten.Image
	CollisionWorld *collisions.World
	AssetsFolder   *embed.FS
	Manager        *manager.Manager
}

var Context *ContextType
