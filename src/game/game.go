// Package game contains the game logic and stuff.
package game

import (
	"time"
)

const (
	gameTickerDuration = 10 * time.Millisecond
)
type Game struct {
	*Base
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
		g.Update(t.Sub(g.t0))
	}
}

func (g *Game) Update(t time.Duration) {
	// TODO: implement
	//fmt.Printf("The game time is %v\n", t)
}

func (g *Game) Destroy() {
	g.ticker.Stop()
}

func (g *Game) HandleKey(k uint32) error {
	switch k {
	default:
		return nil
	}
	return nil
}