package scenes

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tducasse/ebiten-template/camera"
	"github.com/tducasse/ebiten-template/collisions"
	"github.com/tducasse/ebiten-template/manager"
)

type ContextType struct {
	Camera         *camera.Camera
	World          *ebiten.Image
	CollisionWorld *collisions.World
	AssetsFolder   *embed.FS
	Manager        *manager.Manager
}

var Context *ContextType
