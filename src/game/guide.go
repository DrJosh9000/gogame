package game

import (
	"time"
	
	"sdl"
)

type KeyGuide int
const (
	KeyGuideA KeyGuide = iota
	KeyGuideD
	KeyGuideE
)

type Guide struct {
	tex *sdl.Texture
	x, y int
	KeyGuide
	
	Updater *time.Ticker
}

func NewGuide(ctx *sdl.Context, k KeyGuide) (*Guide, error) {
	
	return nil, nil
}

