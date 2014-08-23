package game

import (
	"fmt"
	"math/rand"
	"time"

	"sdl"
)

const (
	texturePlayerFile = "assets/spacepsn.png"
	playerWalkFrames  = 4
	playerWidth, playerHeight = 32, 32
	
	playerWalkSpeed = 256 // pixels per second?
)

type Facing int
const (
	Left Facing = iota
	Right
)

type Control int
const (
	Quit Control = iota
	StartWalkLeft
	StopWalkLeft
	StartWalkRight
	StopWalkRight
	StartJump
	StopJump
	StartFire
	StopFire
)

type Player struct {
	Base
	x, y, frame int
	fx, fy float64
	tex  *sdl.Texture
	lastUpdate time.Duration
	
	Controller chan Control
	
	// All struct elements below should be controlled only by the life goroutine.
	facing Facing
	walking bool
}

func NewPlayer(ctx *sdl.Context) (*Player, error) {
	tex, err := ctx.GetTexture(texturePlayerFile)
	if err != nil {
		return nil, err
	}
	p := &Player{
		fx:   float64(rand.Int31() % (1024 - playerWidth)),
		fy:   float64(rand.Int31() % (768 - playerHeight)),
		tex: tex,
		Controller: make(chan Control),
	}
	go p.life()
	return p, nil
}

func (p *Player) Draw(r *sdl.Renderer) error {
	fx := (p.frame%playerWalkFrames)*playerWidth
	switch p.facing {
	case Left:
		return r.Copy(p.tex,
			sdl.Rect(fx, 0, playerWidth, playerHeight),
			sdl.Rect(p.x, p.y, playerWidth, playerHeight))
	case Right:
		return r.CopyEx(p.tex,
			sdl.Rect(fx, 0, playerWidth, playerHeight),
			sdl.Rect(p.x, p.y, playerWidth, playerHeight), 0, nil, sdl.FlipHorizontal)
	}
	return nil
}

func (p *Player) Update(t time.Duration) {
	if p.lastUpdate == 0 {
		p.lastUpdate = t
		return
	}
	delta := t - p.lastUpdate
	//fmt.Printf("delta %d\n", delta)
	if p.walking {
		dist := float64(playerWalkSpeed * delta) / float64(time.Second)
		fmt.Printf("dist %06f\n", dist)
		switch p.facing {
		case Left:
			p.fx -= dist
		case Right:
			p.fx += dist
		}
		p.frame = (int(2 * t / time.Millisecond) % 1000) / 250
	} else {
		p.frame = 0
	}
	p.x = int(p.fx)
	p.y = int(p.fy)
	p.lastUpdate = t
}

func (p *Player) life() {
	for {
		select {
		case ctl := <-p.Controller:
			switch ctl {
			case Quit:
				return
			case StartWalkLeft:
				p.walking = true
				p.facing = Left
			case StopWalkLeft:
				p.walking = false
			case StartWalkRight:
				p.walking = true
				p.facing = Right
			case StopWalkRight:
				p.walking = false
			default:
				// TODO: more actions
			}
		}
	}
}