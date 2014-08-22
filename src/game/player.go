package game

import (
	"time"
)

type Player struct {
	
}

func NewPlayer() *Player {
	return nil
}

func (p *Player) Children() []Object {
	return nil
}

func (p *Player) Parent() Object {
	return nil
}

func (p *Player) Update(t time.Duration) {
	// TODO: implement
}

func (p *Player) Destroy() {
	// TODO: implement
}