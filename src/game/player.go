package game

import (
	"time"

	"sdl"
)

const (
	textureSpacepersonFile    = "assets/spacepsn.png"
	playerWidth, playerHeight = 32, 32
)

var (
	textureSpaceperson *sdl.Texture
)

type Player struct {
	Base
	x, y int
	tex  *sdl.Texture
}

func NewPlayer() *Player {
	return &Player{tex: textureSpaceperson}
}

func (p *Player) Draw(r *sdl.Renderer) {
	if p.tex == nil {
		panic("player texture not initialised")
	}
	r.Copy(p.tex, nil, sdl.Rect(p.x, p.y, playerWidth, playerHeight))
}

func (p *Player) Update(t time.Duration) {
	// TODO: implement

	p.Base.Update(t)
}

func InitPlayerTexture(r *sdl.Renderer) error {
	if textureSpaceperson == nil {
		t, err := r.LoadImage(textureSpacepersonFile)
		if err != nil {
			return err
		}
		textureSpaceperson = t
	}
	return nil
}
