package game

import (
	"sdl"
)

var tileTemplate = &spriteTemplate{
	sheetFile:   "assets/tiles.png",
	frameWidth:  32,
	frameHeight: 32,
	framesX:     8,
	framesY:     8,
}

type tile struct {
	*sprite
}

func (t tile) destroy() {}

type layer struct {
	complexBase
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
				l.addChild(tile{sprite: s})
			}
		}
	}
	return l, nil
}

type terrain struct {
	complexBase
	exits []*exit
}

func newTerrain(ctx *sdl.Context, lev *level) (*terrain, error) {
	t := &terrain{}
	for _, m := range lev.levelMap {
		l, err := newLayer(ctx, m)
		if err != nil {
			return nil, err
		}
		t.addChild(l)
	}

	if lev.hasExit {
		e, err := newExit(ctx)
		if err != nil {
			return nil, err
		}
		e.x, e.y = tileTemplate.frameWidth*lev.exitX, tileTemplate.frameHeight*lev.exitY-16 // hax
		t.exits = append(t.exits, e)
		t.addChild(e)
	}
	return t, nil
}
