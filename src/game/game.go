// Package game contains the game logic and stuff.
package game

import (
	"fmt"
	"time"

	"sdl"
)

const (
	gameTickerDuration = 10 * time.Millisecond

	level0File  = "assets/level0.txt"
	level1AFile = "assets/level1a.txt"
	level1BFile = "assets/level1b.txt"
)

var gameInstance *Game

type Game struct {
	ComplexBase
	ctx              *sdl.Context
	t0               time.Time
	ticker           *time.Ticker
	inbox            chan message
	offsetX, offsetY int

	player       *Player
	exit         *Exit
	levels       [2]*Level
	terrains     [2]*Terrain
	currentLevel int
}

func NewGame(ctx *sdl.Context) (*Game, error) {
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

	g := &Game{
		ctx: ctx,
		t0:  time.Now(),
		//ticker: time.NewTicker(gameTickerDuration),
		inbox:        make(chan message, 10),
		player:       p,
		levels:       [2]*Level{m0, m1},
		terrains:     [2]*Terrain{t0, t1},
		currentLevel: 0,
	}
	gameInstance = g
	p.x, p.y = tileTemplate.frameWidth*m0.StartX, tileTemplate.frameHeight*m0.StartY
	g.AddChild(p)

	kmp("player.location", g.inbox)
	go g.messageLoop()

	return g, nil
}

func (g *Game) messageLoop() {
	for msg := range g.inbox {
		switch m := msg.v.(type) {
		case locationMsg:
			if msg.k == "player.location" {
				// Keep the player in view.
				if m.x+g.offsetX > 768 {
					g.offsetX = 768 - m.x
				}
				if m.x+g.offsetX < 256 {
					g.offsetX = 256 - m.x
				}
				if g.offsetX > 0 {
					g.offsetX = 0
				}
				// TODO: base on level size, pls
				if g.offsetX < -96*tileTemplate.frameWidth {
					g.offsetX = -96 * tileTemplate.frameWidth
				}

				if m.y+g.offsetY > 576 {
					g.offsetY = 576 - m.y
				}
				if m.y+g.offsetY < 192 {
					g.offsetY = 192 - m.y
				}
				// TODO: when levels have height, adapt here:
				g.offsetY = 0
				/*
					if g.offsetY < 0 {
						g.offsetY = 0
					}
					if g.offsetY > 0 {
						g.offsetY = 0
					}
				*/
			}
		} // switch msg.(type)
	} // for range g.inbox
}

func (g *Game) Draw(r *sdl.Renderer) error {
	r.OffsetX, r.OffsetY = g.offsetX, g.offsetY

	// Draw current terrain
	if err := g.terrains[g.currentLevel].Draw(r); err != nil {
		return err
	}

	// Draw everything else
	return g.ComplexBase.Draw(r)
}

func (g *Game) Destroy() {
	notify("game", quitMsg)
	g.player.Controller <- Quit
	//	g.ticker.Stop()
	g.ComplexBase.Destroy()
}

func (g *Game) Level() *Level {
	return g.levels[g.currentLevel]
}

func (g *Game) HandleEvent(ev interface{}) error {
	switch v := ev.(type) {
	case sdl.KeyDownEvent:
		switch v.KeyCode {
		case ' ':
			g.player.Controller <- StartJump
		case 'w':
			fmt.Println("w down")
		case 'a':
			g.player.Controller <- StartWalkLeft
		case 's':
			fmt.Println("s down")
		case 'd':
			g.player.Controller <- StartWalkRight
		}
	case sdl.KeyUpEvent:
		switch v.KeyCode {
		case ' ':
			g.player.Controller <- StopJump
		case 'w':
			fmt.Println("w up")
		case 'a':
			g.player.Controller <- StopWalkLeft
		case 's':
			fmt.Println("s up")
		case 'd':
			g.player.Controller <- StopWalkRight
		case 'e':
			g.currentLevel = (g.currentLevel + 1) % 2
			g.player.Controller <- Teleport

		}
	}
	return nil
}
