// Package game contains the game logic and stuff.
package game

import (
	"fmt"
	"time"
)

const (
	gameTickerDuration = 10 * time.Millisecond
)

type Game struct {
	t0 time.Time
	ticker *time.Ticker
}

func NewGame() *Game {
	g := &Game {
		t0: time.Now(),
		ticker: time.NewTicker(gameTickerDuration),
	}
	go g.tickLoop()
	return g
}

func (g *Game) tickLoop() {
	for t := range g.ticker.C {
		g.Update(t)
	}
}

func (g *Game) Update(t time.Time) {
	// TODO: implement
	//fmt.Printf("The time is %v\n", t.Sub(g.t0))
}

func (g *Game) Destroy() {
	g.ticker.Stop()
}