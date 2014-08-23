package game

import (
	"time"
)

type Player struct {
	Base
}

func NewPlayer() *Player {
	return &Player {	}
}

func (p *Player) Update(t time.Duration) {
	// TODO: implement
	
	p.Base.Update(t)
}
