package game

import (
	"time"

	"sdl"
)

const (
	// TODO: replace with bigger sheet
	tileSheetFile                   = "assets/default.png"
	tileWidth, tileHeight           = 32, 32
	tileSheetWidth, tileSheetHeight = 1, 1
)

var (
	terrainTileSheet *sdl.Texture
)

type tile struct {
	// x and y position in pixels - woo arbitrary!
	// id index in tile sheet
	x, y, id int
}

// Tile needs basically no work.
func (t *tile) AddChild(Object)      {}
func (t *tile) Children() []Object   { return nil }
func (t *tile) Update(time.Duration) {}
func (t *tile) Destroy()             {}

func (t *tile) Draw(r *sdl.Renderer) {
	if terrainTileSheet == nil {
		panic("Tilesheet not initialized")
	}
	r.Copy(terrainTileSheet, sdl.Rect((t.id%tileSheetWidth)*32, (t.id/tileSheetWidth)*32, tileWidth, tileHeight),
		sdl.Rect(t.x, t.y, tileWidth, tileHeight))
}

type layer struct {
	Base
}

func newLayer() *layer {
	l := &layer{}
	for i := 0; i < 32; i++ {
		for j := 0; j < 24; j++ {
			l.AddChild(&tile{
				x:  i * tileWidth,
				y:  j * tileHeight,
				id: 0, // TODO: load from somewhere
			})
		}
	}
	return l
}

type Terrain struct {
	Base
}

func NewTerrain() *Terrain {
	t := &Terrain{}
	t.AddChild(newLayer())
	return t
}

func InitTerrainTextures(r *sdl.Renderer) error {
	if terrainTileSheet == nil {
		t, err := r.LoadImage("assets/default.png")
		if err != nil {
			return err
		}
		terrainTileSheet = t
	}
	return nil
}
