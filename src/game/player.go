package game

import (
	"math/rand"

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
	return &Player{
		x: int(rand.Int31() % (1024-playerWidth)),
		y: int(rand.Int31() % (768-playerHeight)),
		tex: textureSpaceperson,
	}
}

func (p *Player) Draw(r *sdl.Renderer) {
	if p.tex == nil {
		panic("player texture not initialised")
	}
	r.Copy(p.tex, sdl.Rect(0, 0, playerWidth, playerHeight), sdl.Rect(p.x, p.y, playerWidth, playerHeight))
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
