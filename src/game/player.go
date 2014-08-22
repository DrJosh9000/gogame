package game

import (
	"time"
)

type Player struct {
	*Base
}

func NewPlayer() *Player {
	return nil
}

func (p *Player) Update(t time.Duration) {
	// TODO: implement
}

func (p *Player) Destroy() {
	// TODO: implement
}