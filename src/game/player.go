package game

import (
	//"fmt"
	"math"
	"math/rand"
	"time"

	"sdl"
)

const (
	texturePlayerFile = "assets/spacepsn.png"
	playerFramesWidth  = 4
	playerWidth, playerHeight = 32, 32
	
	playerWalkSpeed = 256 // pixels per second?
	
	playerTau = 0.2
)

type Facing int
const (
	Left Facing = iota
	Right
)

type Animation int
const (
	Standing Animation = iota
	Walking
	Jumping
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
	fx, fy, dx, dy float64
	tex  *sdl.Texture
	lastUpdate time.Duration
	
	Controller chan Control
	
	// All struct elements below should be controlled only by the life goroutine.
	facing Facing
	anim Animation
	wx, wy float64
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
	fx := (p.frame%playerFramesWidth)*playerWidth
	switch p.facing {
	case Left:
		return r.Copy(p.tex,
			sdl.Rect(fx, 0, playerWidth, playerHeight),
			sdl.Rect(p.x, p.y, playerWidth, playerHeight))
	case Right:
		// TODO: separate animation for walking right
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
	delta := float64(t - p.lastUpdate) / float64(time.Second)
	// TODO: more frame calcs
	if p.anim == Walking {
		p.frame = (int(2 * t / time.Millisecond) % 1000) / 250
	} else {
		p.frame = 0
	}
	tau := playerTau * math.Exp(delta)
	p.dx = tau * p.wx + (1-tau) * p.dx
	p.dy = tau * p.wy + (1-tau) * p.dy
	
	// FISIXX
	p.fx += p.dx * delta
	p.fy += p.dy * delta
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
				p.anim = Walking
				p.facing = Left
				p.wx = -playerWalkSpeed
			case StopWalkLeft:
				p.anim = Standing
				p.wx = 0
			case StartWalkRight:
				p.anim = Walking
				p.facing = Right
				p.wx = playerWalkSpeed
			case StopWalkRight:
				p.anim = Standing
				p.wx = 0
			case StartJump:
				p.anim = Jumping
				p.wy = -1000
			case StopJump:
				p.anim = Standing
				p.wy = 0
			default:
				// TODO: more actions
			}
		}
	}
}