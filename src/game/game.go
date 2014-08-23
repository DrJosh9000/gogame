// Package game contains the game logic and stuff.
package game

import (
	"fmt"
	"time"

	"sdl"
)

const (
	gameTickerDuration = 10 * time.Millisecond
)

type Game struct {
	Base
	ctx *sdl.Context
	t0     time.Time
	ticker *time.Ticker

	player *Player
}

func NewGame(ctx *sdl.Context) (*Game, error) {
	p, err  := NewPlayer(ctx)
	if err != nil {
		return nil, err
	}
	t, err := NewTerrain(ctx)
	if err != nil {
		return nil, err
	}
	
	g := &Game{
		ctx: ctx,
		t0:     time.Now(),
		ticker: time.NewTicker(gameTickerDuration),
		player: p,
	}
	g.AddChild(t)
	g.AddChild(p)
	go g.tickLoop()
	return g, nil
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
		fmt.Println("jump")
	case 'w':
		fmt.Println("up")
	case 'a':
		fmt.Println("left")
	case 's':
		fmt.Println("down")
	case 'd':
		fmt.Println("right")
	default:
		fmt.Println("other")
	}
	return nil
}

