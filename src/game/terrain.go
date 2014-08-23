package game

import (
	"time"
	
	"sdl"
)

const (
	// TODO: replace with bigger sheet
	tileSheetFile = "assets/default.png"
	tileWidth, tileHeight = 32, 32
	tileSheetWidth, tileSheetHeight = 1, 1
)

var (
	terrainTileSheet *sdl.Texture
)

type tile struct {
	Base
	// x and y position in pixels - woo arbitrary!
	// id index in tile sheet
	x, y, id int
}

func (t *tile) Draw(r *sdl.Renderer) {
	if terrainTileSheet == nil {
		panic("Tilesheet not initialized")
	}
	r.Copy(terrainTileSheet, sdl.Rect((t.id % tileSheetWidth) * 32, (t.id / tileSheetWidth) * 32, tileWidth, tileHeight),
		sdl.Rect(t.x, t.y, tileWidth, tileHeight))
}

func (t *tile) Update(time.Duration) {
	// No updates here.
}

type layer struct {
	Base
}

func (l *layer) Update(time.Duration) {
	// No updates here.
}

type Terrain struct {
	Base
}

func NewTerrain() *Terrain {
	return &Terrain{}
}

func (t *Terrain) Update(time.Duration) {
	// No updates here.
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