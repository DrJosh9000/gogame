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

func newLayer(ctx *sdl.Context, m LevelLayer) (*layer, error) {
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

type Terrain struct {
	ComplexBase
	exit *Exit
}

func NewTerrain(ctx *sdl.Context, lev *Level) (*Terrain, error) {
	t := &Terrain{}
	for _, m := range lev.Map {
		l, err := newLayer(ctx, m)
		if err != nil {
			return nil, err
		}
		t.AddChild(l)
	}

	if lev.HasExit {
		e, err := NewExit(ctx)
		if err != nil {
			return nil, err
		}
		e.x, e.y = tileTemplate.frameWidth*lev.ExitX, tileTemplate.frameHeight*lev.ExitY-16 // hax
		t.exit = e
		t.AddChild(e)
	}
	return t, nil
}
