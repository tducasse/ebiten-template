package entities

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tducasse/ebiten-template/aseprite"
	"github.com/tducasse/ebiten-template/input"
	"github.com/tducasse/ebiten-template/ldtk"
)

type Player struct {
	X      float64
	Y      float64
	Sprite *aseprite.Animation
	Speed  float64
}

func (p *Player) Init(levels *ldtk.Ldtk, opt *EntityOptions) {
	playerEntity := levels.FindEntity("Player")
	if playerEntity == nil {
		log.Fatal("Could not find Player")
	}
	p.Speed = 100
	p.X = float64(playerEntity.Px[0])
	p.Y = float64(playerEntity.Px[1])
	p.Sprite = aseprite.Load(
		"player.json",
		&aseprite.Options{
			EmbedFolder: opt.EmbedFolder,
			Root:        opt.Root,
		},
		"idle",
		p,
	)
	p.Sprite.OnLoop(onAnimLoop)
}

func onAnimLoop(player interface{}, anim *aseprite.Animation) {
	// p := player.(*Player)
	// fmt.Println(anim.CurrentTag.Name)
	// fmt.Println(p.X)
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Sprite.Draw(p.X, p.Y, screen)
}

func (p *Player) Move() {
	var (
		dx float64
		dy float64
	)
	if input.IsPressed("up") {
		dy = -p.Speed
	} else if input.IsPressed("down") {
		dy = p.Speed
	}
	if input.IsPressed("left") {
		dx = -p.Speed
	} else if input.IsPressed("right") {
		dx = p.Speed
	}

	mag := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	if mag > 0 {
		dx = dx / mag
		dy = dy / mag
	}

	if dx != 0 || dy != 0 {
		p.Sprite.SetTag("walk")
	} else {
		p.Sprite.SetTag("idle")
	}
	p.X += dx
	p.Y += dy
}

func (p *Player) Update() {
	p.Move()
	p.Sprite.Update()
}
