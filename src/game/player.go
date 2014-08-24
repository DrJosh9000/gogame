package game

import (
	//"fmt"
	"math"
	"time"

	"sdl"
)

const (
	texturePlayerFile = "assets/spacepsn.png"
	playerFramesWidth  = 4
	playerWidth, playerHeight = 32, 32
	
	playerWalkSpeed = 384 // pixels per second?
	playerJumpSpeed = -512
	playerGravity = 2048
	
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
	Base
	tex  *sdl.Texture
	lastUpdate time.Duration
	
	Controller chan Control
	Updater chan time.Duration
	
	// All struct elements below should be controlled only by the life goroutine.
	facing Facing
	anim Animation
	wx, wy float64
	x, y, frame int
	fx, fy, dx, dy, ddx, ddy float64
}

func NewPlayer(ctx *sdl.Context) (*Player, error) {
	tex, err := ctx.GetTexture(texturePlayerFile)
	if err != nil {
		return nil, err
	}
	p := &Player{
		fx:   64,
		fy:   768-64,
		facing: Right,
		anim: Standing,
		tex: tex,
		Controller: make(chan Control),
		Updater: make(chan time.Duration),
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
		// TODO: separate animations for facing right
		return r.CopyEx(p.tex,
			sdl.Rect(fx, 0, playerWidth, playerHeight),
			sdl.Rect(p.x, p.y, playerWidth, playerHeight), 0, nil, sdl.FlipHorizontal)
	}
	return nil
}

func (p *Player) Destroy() {
	close(p.Controller)
	close(p.Updater)
}

func (p *Player) update(t time.Duration) {
	if p.lastUpdate == 0 {
		p.lastUpdate = t
		return
	}
	delta := float64(t - p.lastUpdate) / float64(time.Second)
	
	switch p.anim {
	case Walking:
		p.frame = (int(2 * t / time.Millisecond) % 1000) / 250
		fallthrough
	default:
		if p.anim != Walking {
			p.frame = 0
		}
		tau := playerTau * math.Exp(delta)
		p.dx = tau * p.wx + (1-tau) * p.dx
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
		if !gameInstance.level.QueryPoint(nx, ny+32).solid && !gameInstance.level.QueryPoint(nx+31, ny+32).solid {
				p.anim = Falling
				p.ddy = playerGravity
		}
	case Falling:
		if tile := gameInstance.level.QueryPoint(nx, ny+32); tile.solid {
				p.anim = Standing
				ny = (ny/tileHeight)*tileHeight
				p.fy = float64(ny)
				p.dy = 0
				p.ddy = 0
		}
		if tile := gameInstance.level.QueryPoint(nx+31, ny+32); tile.solid {
				p.anim = Standing
				ny = (ny/tileHeight)*tileHeight
				p.fy = float64(ny)
				p.dy = 0
				p.ddy = 0
		}
	}
	if tile := gameInstance.level.QueryPoint(nx, ny+31); tile.solid {
			nx = ((nx/tileWidth)+1)*tileWidth
			p.fx, p.fy = float64(nx), float64(ny)
			p.dx = 0
	}
	if tile := gameInstance.level.QueryPoint(nx+31, ny+31); tile.solid {
			nx = (nx/tileWidth)*tileWidth
			p.fx, p.fy = float64(nx), float64(ny)
			p.dx = 0
	}
	
	p.x, p.y = nx, ny
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
	default:
		// TODO: more actions
	}
	return false
}

func (p *Player) life() {
	for {
		select {
		case c := <-p.Controller:
			if p.control(c) {
				return
			}
		case t := <-p.Updater:
			p.update(t)
		}
	}
}