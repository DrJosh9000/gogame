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
	level1AFile = "assets/level1a.txt"
	level1BFile = "assets/level1b.txt"
)

var gameInstance *Game

type Game struct {
	Base
	ctx    *sdl.Context
	t0     time.Time
	ticker *time.Ticker

	player *Player
	levels [2]*Level
	terrains [2]*Terrain
	currentLevel int
}

func GetGame(ctx *sdl.Context) (*Game, error) {
	if gameInstance != nil {
		return gameInstance, nil
	}
	
	p, err := NewPlayer(ctx)
	if err != nil {
		return nil, err
	}
	
	m0, err := LoadLevel(level1AFile)
	if err != nil {
		return nil, err
	}
	m1, err := LoadLevel(level1BFile)
	if err != nil {
		return nil, err
	}
	
	t0, err := NewTerrain(ctx, m0)
	if err != nil {
		return nil, err
	}
	t1, err := NewTerrain(ctx, m1)
	if err != nil {
		return nil, err
	}

	gameInstance = &Game{
		ctx:    ctx,
		t0:     time.Now(),
		ticker: time.NewTicker(gameTickerDuration),
		player: p,
		levels: [2]*Level{m0, m1},
		terrains: [2]*Terrain{t0, t1},
		currentLevel: 0,
	}
	p.x, p.y = tileWidth * m0.StartX, tileHeight * m0.StartY
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

func (g *Game) Draw(r *sdl.Renderer) error {
	// Keep the player in view.
	// HAX HAX HAX
	if g.player.x + r.OffsetX > 768 {
		r.OffsetX = 768 - g.player.x
	}
	if g.player.x + r.OffsetX < 256 {
		r.OffsetX = 256 - g.player.x
	}
	if r.OffsetX > 0 {
		r.OffsetX = 0
	}
	// TODO: base on level size, pls
	if r.OffsetX < -96 * tileWidth {
		r.OffsetX = -96 * tileWidth
	}
	
	if g.player.y + r.OffsetY > 576 {
		r.OffsetY = 576 - g.player.y
	}
	if g.player.y + r.OffsetY < 192 {
		r.OffsetY = 192 - g.player.y
	}
	if r.OffsetY < 0 {
		r.OffsetY = 0
	}
	if r.OffsetY > 0 {
		r.OffsetY = 0
	}
	
	// Draw current terrain
	if err := g.terrains[g.currentLevel].Draw(r); err != nil {
		return err
	}
	
	// Draw everything else
	return g.Base.Draw(r)
}

func (g *Game) Destroy() {
	g.player.Controller <- Quit
	g.ticker.Stop()
	g.Base.Destroy()
}

func (g *Game) Level() *Level {
	return g.levels[g.currentLevel]
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
				if v.Type == sdl.KeyUp {
					g.currentLevel = (g.currentLevel+1)%2
					g.player.Controller <- Teleport
				}
			}
	}
	return nil
}
