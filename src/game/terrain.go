package game

import (
	"sdl"
)

const (
	tileSheetFile                   = "assets/tiles.png"
	tileWidth, tileHeight           = 32, 32
	tileSheetWidth, tileSheetHeight = 8, 8
)

type tile struct {
	// x and y position in pixels - woo arbitrary!
	// id index in tile sheet
	x, y, id int
}

type layer struct {
	ComplexBase
	tiles []tile
	tex   *sdl.Texture
}

func newLayer(ctx *sdl.Context, m LevelLayer) (*layer, error) {
	tex, err := ctx.GetTexture(tileSheetFile)
	if err != nil {
		return nil, err
	}
	l := &layer{tex: tex}
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			if m[i][j].index != 0 {
				l.tiles = append(l.tiles, tile{
					x:  j * tileWidth,
					y:  i * tileHeight,
					id: m[i][j].index,
				})
			}
		}
	}
	return l, nil
}

func (l *layer) Draw(r *sdl.Renderer) error {
	for _, t := range l.tiles {
		if err := r.Copy(l.tex,
			sdl.Rect((t.id%tileSheetWidth)*32, (t.id/tileSheetWidth)*32, tileWidth, tileHeight),
			sdl.Rect(t.x, t.y, tileWidth, tileHeight)); err != nil {
			return err
		}
	}
	return nil
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
		e.x, e.y = tileWidth*lev.ExitX, tileHeight*lev.ExitY-16 // hax
		t.exit = e
		t.AddChild(e)
	}
	return t, nil
}
