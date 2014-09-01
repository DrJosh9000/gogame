package game

import (
	"time"
	
	"sdl"
)

const (
	keyGuideFrameWidth, keyGuideFrameHeight = 64, 64
	keyGuideUpdateInterval = 500 * time.Millisecond
	keyGuideTextureFile = "assets/keys.png"
	
)

type KeyGuide int
const (
	KeyGuideA KeyGuide = iota
	KeyGuideD
	KeyGuideE
)

type Guide struct {
	tex *sdl.Texture
	x, y, frame int
	KeyGuide
	
	Updater *time.Ticker
}

func NewGuide(ctx *sdl.Context, k KeyGuide) (*Guide, error) {
	tex, err := ctx.GetTexture(keyGuideTextureFile)
	if err != nil {
		return nil, err
	}
	
	g := &Guide{
		tex: tex,
		KeyGuide: k,
		Updater: time.NewTicker(keyGuideUpdateInterval),
	}
	go g.life()
	return g, nil
}

func (g *Guide) Draw(r *sdl.Renderer) error {
	return r.Copy(g.tex, sdl.Rect(g.frame*keyGuideFrameWidth, int(g.KeyGuide)*keyGuideFrameHeight, keyGuideFrameWidth, keyGuideFrameHeight), sdl.Rect(g.x, g.y, keyGuideFrameWidth, keyGuideFrameHeight))
}

func (g *Guide) Destroy() {
	g.Updater.Stop()
}

func (g *Guide) life() {
	for _ = range g.Updater.C {
		g.frame = (g.frame+1)%2
	}
}
