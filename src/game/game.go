package game

import (
	"fmt"
	"time"
)

const (
	gameTickerDuration = 10 * time.Millisecond
)

type Game struct {
	T0 time.Time
	*time.Ticker
}

func NewGame() *Game {
	g := &Game {
		T0: time.Now(),
		Ticker: time.NewTicker(gameTickerDuration),
	}
	go g.tickLoop()
	return g
}

func (g *Game) tickLoop() {
	for t := range g.Ticker.C {
		g.Update(t)
	}
}

func (g *Game) Update(t time.Time) {
	// TODO: implement
	fmt.Printf("The time is %v\n", t)
}

func (g *Game) Destroy() {
	g.Ticker.Stop()
}