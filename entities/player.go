package entities

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tducasse/ebiten-template/aseprite"
	"github.com/tducasse/ebiten-template/ldtk"
)

type Player struct {
	X      float64
	Y      float64
	Sprite *aseprite.Animation
}

func (p *Player) Init(levels *ldtk.Ldtk, opt *EntityOptions) {
	playerEntity := levels.FindEntity("Player")
	if playerEntity == nil {
		log.Fatal("Could not find Player")
	}
	p.X = float64(playerEntity.Px[0])
	p.Y = float64(playerEntity.Px[1])
	p.Sprite = aseprite.Load(
		"player.json",
		&aseprite.Options{
			EmbedFolder: opt.EmbedFolder,
			Root:        opt.Root,
		},
		"walk",
		p,
	)
	p.Sprite.OnLoop(onAnimLoop)
}

func onAnimLoop(player interface{}, anim *aseprite.Animation) {
	p := player.(*Player)
	fmt.Println(anim.CurrentTag.Name)
	fmt.Println(p.X)
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Sprite.Draw(p.X, p.Y, screen)
}

func (p *Player) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		p.Sprite.SetTag("walk")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		p.Sprite.SetTag("red")
	}
	p.Sprite.Update()
}
