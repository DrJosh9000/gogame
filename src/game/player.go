package game

import (
	//"fmt"
	"math"
	"time"

	"sdl"
)

var playerTemplate = &spriteTemplate{
	name: "player",
	sheetFile: "assets/spacepsn.png",
	framesX: 4,
	framesY: 2,
	frameWidth: 32,
	frameHeight: 32,
}

const (
	playerUpdateInterval = 10 * time.Millisecond

	playerWalkSpeed = 384 // pixels per second?
	playerJumpSpeed = -512
	playerGravity   = 2048

	playerTau = 0.1
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
	Falling
	Floating
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
	Land
	Teleport
)

type Player struct {
	*sprite
	lastUpdate time.Duration

	Controller chan Control
	Updater    *time.Ticker

	// All struct elements below should be controlled only by the life goroutine.
	facing                   Facing
	anim                     Animation
	wx, wy                   float64
	fx, fy, dx, dy, ddx, ddy float64
}

func NewPlayer(ctx *sdl.Context) (*Player, error) {
	s, err := playerTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	p := &Player{
		sprite:	    s,
		fx:         64,
		fy:         768 - 64,
		facing:     Right,
		anim:       Standing,
		Controller: make(chan Control),
		Updater:    time.NewTicker(playerUpdateInterval),
	}
	go p.life()
	return p, nil
}

func (p *Player) Destroy() {
	close(p.Controller)
	p.Updater.Stop()
}

func (p *Player) String() string {
	return "player"
}

func (p *Player) update(t time.Duration) {
	if p.lastUpdate == 0 {
		p.lastUpdate = t
		return
	}
	delta := float64(t-p.lastUpdate) / float64(time.Second)

	switch p.anim {
	case Walking:
		p.frame = (int(2*t/time.Millisecond) % 1000) / 250
		fallthrough
	default:
		if p.anim != Walking {
			p.frame = 0
		}
		if p.facing == Right {
			p.frame += 4
		}
		tau := playerTau * math.Exp(delta)
		p.dx = tau*p.wx + (1-tau)*p.dx
		/*
				p.dy = tau * p.wy + (1-tau) * p.dy
			case Falling:
				p.frame = 0
				p.dx += p.ddx * delta
		*/
		p.dy += p.ddy * delta
	}

	// FISIXX
	p.fx += p.dx * delta
	p.fy += p.dy * delta
	nx, ny := int(p.fx), int(p.fy)
	p.lastUpdate = t

	switch p.anim {
	case Standing, Walking:
		if !gameInstance.Level().IsPointSolid(nx, ny+32) && !gameInstance.Level().IsPointSolid(nx+31, ny+32) {
			p.anim = Falling
			p.ddy = playerGravity
		}
	case Falling:
		if gameInstance.Level().IsPointSolid(nx, ny+32) {
			p.anim = Standing
			ny = (ny / tileHeight) * tileHeight
			p.fy = float64(ny)
			p.dy = 0
			p.ddy = 0
		}
		if gameInstance.Level().IsPointSolid(nx+31, ny+32) {
			p.anim = Standing
			ny = (ny / tileHeight) * tileHeight
			p.fy = float64(ny)
			p.dy = 0
			p.ddy = 0
		}
	}
	if gameInstance.Level().IsPointSolid(nx, ny+31) {
		nx = ((nx / tileWidth) + 1) * tileWidth
		p.fx, p.fy = float64(nx), float64(ny)
		p.dx = 0
	}
	if gameInstance.Level().IsPointSolid(nx+31, ny+31) {
		nx = (nx / tileWidth) * tileWidth
		p.fx, p.fy = float64(nx), float64(ny)
		p.dx = 0
	}

	p.x, p.y = nx, ny

	notify("player.location", locationMsg{o: p, x: p.x, y: p.y})
}

func (p *Player) control(ctl Control) bool {
	switch ctl {
	case Quit:
		return true
	case StartWalkLeft:
		switch p.anim {
		case Standing, Walking:
			p.anim = Walking
			p.facing = Left
			p.wx = -playerWalkSpeed
		}
	case StopWalkLeft:
		if p.anim == Walking {
			p.anim = Standing
		}
		p.wx = 0
	case StartWalkRight:
		switch p.anim {
		case Standing, Walking:
			p.anim = Walking
			p.facing = Right
			p.wx = playerWalkSpeed
		}
	case StopWalkRight:
		if p.anim == Walking {
			p.anim = Standing
		}
		p.wx = 0
	case StartJump:
		switch p.anim {
		case Standing, Walking:
			/*
						p.anim = Jumping
					}
				case StopJump:
					if p.anim == Jumping {
			*/
			p.anim = Falling
			p.dy = playerJumpSpeed
			p.ddy = playerGravity
		}
	case Land:
		if p.anim == Falling {
			p.anim = Standing
			p.dy = 0
			p.ddy = 0
		}
	case Teleport:
		p.anim = Falling
		p.ddy = playerGravity
	default:
		// TODO: more actions
	}
	return false
}

func (p *Player) life() {
	t0 := time.Now()
	for {
		select {
		case c := <-p.Controller:
			if p.control(c) {
				return
			}
		case t := <-p.Updater.C:
			p.update(t.Sub(t0))
		}
	}
}
