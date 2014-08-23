package game

import (
	"sdl"
)

const (
	// TODO: replace with bigger sheet
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
	Base
	tiles []tile
	tex   *sdl.Texture
}

func newLayer(ctx *sdl.Context) (*layer, error) {
	tex, err := ctx.GetTexture(tileSheetFile)
	if err != nil {
		return nil, err
	}
	l := &layer{tex: tex}
	for i := 0; i < 32; i++ {
		for j := 0; j < 24; j++ {
			l.tiles = append(l.tiles, tile{
				x:  i * tileWidth,
				y:  j * tileHeight,
				id: 1, // TODO: load from somewhere
			})
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
	Base
}

func NewTerrain(ctx *sdl.Context) (*Terrain, error) {
	t := &Terrain{}
	l, err := newLayer(ctx)
	if err != nil {
		return nil, err
	}
	t.AddChild(l)
	return t, nil
}
