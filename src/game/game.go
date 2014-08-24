// Package game contains the game logic and stuff.
package game

import (
	"fmt"
	"time"

	"sdl"
)

const (
	gameTickerDuration = 10 * time.Millisecond
	
	level0File = "assets/level0.txt"
)

var gameInstance *Game

type Game struct {
	Base
	ctx    *sdl.Context
	t0     time.Time
	ticker *time.Ticker

	player *Player
	level Level
}

func GetGame(ctx *sdl.Context) (*Game, error) {
	if gameInstance != nil {
		return gameInstance, nil
	}
	
	p, err := NewPlayer(ctx)
	if err != nil {
		return nil, err
	}
	
	m, err := LoadLevel(level0File)
	if err != nil {
		return nil, err
	}
	
	t, err := NewTerrain(ctx, m)
	if err != nil {
		return nil, err
	}

	gameInstance = &Game{
		ctx:    ctx,
		t0:     time.Now(),
		ticker: time.NewTicker(gameTickerDuration),
		player: p,
		level:  m,
	}
	gameInstance.AddChild(t)
	gameInstance.AddChild(p)
	go gameInstance.tickLoop()
	return gameInstance, nil
}

func (g *Game) tickLoop() {
	for t := range g.ticker.C {
		dt := t.Sub(g.t0)
		g.player.Updater <- dt
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
				if v.Type == sdl.KeyDown {
					g.player.Controller <- Teleport
				}
			}
	}
	return nil
}
