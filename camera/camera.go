package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	Width  int
	Height int
	X      float64
	Y      float64
	Follow *FollowData
}

type FollowData struct {
	W int
	H int
}

func (c *Camera) GetTransform() ebiten.GeoM {
	geom := ebiten.GeoM{}
	geom.Translate(
		float64(c.Width)/2-float64(c.X)-float64(c.Follow.W)/2,
		float64(c.Height)/2-float64(c.Y)-float64(c.Follow.H)/2,
	)
	return geom
}

func (c *Camera) Draw(world *ebiten.Image, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.GetTransform(),
	})
}

func Init(width int, height int) *Camera {
	cam := &Camera{
		Width:  width,
		Height: height,
		Follow: &FollowData{},
	}
	return cam
}
