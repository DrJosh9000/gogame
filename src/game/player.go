package game

import (
	"math/rand"

	"sdl"
)

const (
	textureSpacepersonFile    = "assets/spacepsn.png"
	playerWidth, playerHeight = 32, 32
)

type Player struct {
	Base
	x, y int
	tex  *sdl.Texture
}

func NewPlayer(ctx *sdl.Context) (*Player, error) {
	tex, err := ctx.GetTexture(textureSpacepersonFile)
	if err != nil {
		return nil, err
	}
	return &Player{
		x: int(rand.Int31() % (1024-playerWidth)),
		y: int(rand.Int31() % (768-playerHeight)),
		tex: tex,
	}, nil
}

func (p *Player) Draw(r *sdl.Renderer) error {
	return r.Copy(p.tex, sdl.Rect(0, 0, playerWidth, playerHeight), sdl.Rect(p.x, p.y, playerWidth, playerHeight))
}

