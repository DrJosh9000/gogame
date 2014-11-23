package game

import (
	"sdl"
)

var tileTemplate = &spriteTemplate{
	name:        "tile",
	sheetFile:   "assets/tiles.png",
	frameWidth:  32,
	frameHeight: 32,
	framesX:     8,
	framesY:     8,
}

type tile struct {
	*sprite
}

type layer struct {
	ComplexBase
}

func newLayer(ctx *sdl.Context, m levelLayer) (*layer, error) {
	l := &layer{}
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if m[i][j].index != 0 {
				s, err := tileTemplate.new(ctx)
				if err != nil {
					return nil, err
				}
				s.x = j * tileTemplate.frameWidth
				s.y = i * tileTemplate.frameHeight
				s.frame = m[i][j].index
				l.AddChild(tile{sprite: s})
			}
		}
	}
	return l, nil
}

type terrain struct {
	ComplexBase
	exit *exit
}

func newTerrain(ctx *sdl.Context, lev *level) (*terrain, error) {
	t := &terrain{}
	for _, m := range lev.levelMap {
		l, err := newLayer(ctx, m)
		if err != nil {
			return nil, err
		}
		t.AddChild(l)
	}

	if lev.hasExit {
		e, err := newExit(ctx)
		if err != nil {
			return nil, err
		}
		e.x, e.y = tileTemplate.frameWidth*lev.exitX, tileTemplate.frameHeight*lev.exitY-16 // hax
		t.exit = e
		t.AddChild(e)
	}
	return t, nil
}
