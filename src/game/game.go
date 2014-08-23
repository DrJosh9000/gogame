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
	ctx    *sdl.Context
	t0     time.Time
	ticker *time.Ticker

	player *Player
}

func NewGame(ctx *sdl.Context) (*Game, error) {
	p, err := NewPlayer(ctx)
	if err != nil {
		return nil, err
	}
	t, err := NewTerrain(ctx)
	if err != nil {
		return nil, err
	}

	g := &Game{
		ctx:    ctx,
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
	g.player.Controller <- Quit
	g.ticker.Stop()
	g.Base.Destroy()
}

func (g *Game) HandleEvent(ev interface{}) error {
	switch v := ev.(type) {
	case sdl.KeyEvent:
		switch v.KeyCode {
			case ' ':
				if v.Type == sdl.KeyDown {
					g.player.Controller <- StartJump
				} else {
					g.player.Controller <- StopJump
				}
			case 'w':
				fmt.Println("up")
			case 'a':
				if v.Type == sdl.KeyDown {
					g.player.Controller <- StartWalkLeft
				} else {
					g.player.Controller <- StopWalkLeft
				}
			case 's':
				fmt.Println("down")
			case 'd':
				if v.Type == sdl.KeyDown {
					g.player.Controller <- StartWalkRight
				} else {
					g.player.Controller <- StopWalkRight
				}
			case 'e':
				fmt.Println("use")
			default:
				fmt.Println("other")
			}
	}
	return nil
}
