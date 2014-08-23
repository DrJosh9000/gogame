// Package game contains the game logic and stuff.
package game

import (
	"time"
	
	"sdl"
)

const (
	gameTickerDuration = 10 * time.Millisecond
)
type Game struct {
	Base
	t0 time.Time
	ticker *time.Ticker
		
	player *Player
}

func NewGame(r *sdl.Renderer) *Game {
	// Load some initial stuff
	g := &Game {
		t0: time.Now(),
		ticker: time.NewTicker(gameTickerDuration),
		player: NewPlayer(),
	}
	g.AddChild(g.player)
	g.AddChild(NewTerrain())
	go g.tickLoop()
	return g
}

func (g *Game) tickLoop() {
	for t := range g.ticker.C {
		g.Update(t.Sub(g.t0))
	}
}

func (g *Game) Destroy() {
	g.ticker.Stop()
	g.Base.Destroy()
}

func (g *Game) HandleKey(k uint32) error {
	switch k {
	case ' ':
		return nil
	case 'w':
		return nil
	case 'a':
		return nil
	case 's':
		return nil
	case 'd':
		return nil
	default:
		return nil
	}
	return nil
}

func InitAll(ctx *sdl.Context) error {
	if err := InitPlayerTexture(ctx.Renderer); err != nil {
		return err
	}
	if err := InitTerrainTextures(ctx.Renderer); err != nil {
		return err
	}
	return nil
}