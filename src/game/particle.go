package game

import (
	"math/rand"
	"time"
)

type Particle struct {
	Base
	x, y int32
}

func (p *Particle) Update(t time.Duration) {
	// TODO: implement

	p.x += (rand.Int31() % 5) - 2
	p.y += rand.Int31() % 5
}
