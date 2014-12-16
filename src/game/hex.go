package game

import (
	"sdl"
)

var hexTemplate = &spriteTemplate{
	sheetFile:   "assets/hex.png",
	framesX:     1,
	framesY:     1,
	frameWidth:  128,
	frameHeight: 128,
}

type hex struct {
	*sprite
}

func newHex(ctx *sdl.Context) (*hex, error) {
	s, err := hexTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	return &hex{sprite: s}, nil
}
