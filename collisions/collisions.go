package collisions

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Box struct {
	Data interface{}
	W    float64
	H    float64
	X    float64
	Y    float64
}

type World struct {
	Items []*Box
}

type Collision struct {
	Other      *Box
	CollidingX bool
	CollidingY bool
}

func MakeBox(x float64, y float64, w float64, h float64) *Box {
	return &Box{
		W: w,
		H: h,
		X: x,
		Y: y,
	}
}

func MakeWorld() *World {
	return &World{Items: make([]*Box, 0)}
}

func (w *World) Add(box *Box) {
	w.Items = append(w.Items, box)
}

func (b *Box) AddData(data interface{}) {
	b.Data = data
}

func (world *World) AddNewBox(x float64, y float64, w float64, h float64) {
	box := MakeBox(x, y, w, h)
	world.Add(box)
}

func (w *World) Remove(box *Box) {
	for i, item := range w.Items {
		if item == box {
			w.Items[i] = w.Items[len(w.Items)-1]
			w.Items = w.Items[:len(w.Items)-1]
		}
	}
}

func (b *Box) Collides(o *Box, filter func(self *Box, other *Box) bool) bool {
	if b == o {
		return false
	}
	isColliding := b.X < o.X+o.W && b.X+b.W > o.X && b.Y < o.Y+o.H && b.Y+b.H > o.Y
	if isColliding && (filter != nil) {
		isColliding = filter(b, o)
	}
	return isColliding
}

func (w *World) Move(box *Box, dx float64, dy float64, filterFunc func(self *Box, other *Box) bool) (float64, float64, []*Collision) {
	collisions := make([]*Collision, 0)
	var collidingX, collidingY bool
	for _, other := range w.Items {
		box.Y += dy
		onY := box.Collides(other, filterFunc)
		box.Y -= dy
		box.X += dx
		onX := box.Collides(other, filterFunc)
		box.X -= dx
		if onX || onY {
			collisions = append(
				collisions,
				&Collision{
					CollidingX: onX,
					CollidingY: onY,
					Other:      other,
				},
			)
			collidingX = collidingX || onX
			collidingY = collidingY || onY
		}
	}
	if !collidingX {
		box.X = math.Round(box.X + dx)
	}
	if !collidingY {
		box.Y = math.Round(box.Y + dy)
	}
	if (math.Abs(dx) > 1 || math.Abs(dy) > 1) && (collidingX || collidingY) {
		return w.Move(box, dx/2, dy/2, filterFunc)
	}
	return box.X, box.Y, collisions
}

func (w *World) Debug(screen *ebiten.Image) {
	for _, item := range w.Items {
		ebitenutil.DrawRect(screen, item.X, item.Y, item.W, item.H, color.White)
	}
}
