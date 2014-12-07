package game

import (
	"log"
	"math"
	"time"

	"sdl"
)

var playerTemplate = &spriteTemplate{
	sheetFile:   "assets/spacepsn.png",
	framesX:     4,
	framesY:     2,
	frameWidth:  32,
	frameHeight: 32,
}

const (
	playerUpdateInterval = 10 * time.Millisecond

	playerWalkSpeed = 384 // pixels per second?
	playerJumpSpeed = -512
	playerGravity   = 2048

	playerTau = 0.1
)

type playerFacing int

const (
	Left playerFacing = iota
	Right
)

type playerAnimation int

const (
	Standing playerAnimation = iota
	Walking
	Jumping
	Falling
	Floating
)

type playerControl int

const (
	StartWalkLeft playerControl = iota
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

type player struct {
	*sprite
	lastUpdate time.Duration

	inbox chan message

	// All struct elements below should be controlled only by the life goroutine.
	facing                   playerFacing
	anim                     playerAnimation
	wx, wy                   float64
	fx, fy, dx, dy, ddx, ddy float64
}

func newPlayer(ctx *sdl.Context) (*player, error) {
	s, err := playerTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	p := &player{
		sprite: s,
		fx:     64, // TODO: ohgod fix
		fy:     768 - 64,
		facing: Right,
		anim:   Standing,
		inbox:  make(chan message, 10),
	}

	kmp("quit", p.inbox)
	kmp("clock", p.inbox)
	kmp("input.event", p.inbox)
	go p.life()
	return p, nil
}

func (p *player) destroy() {
	log.Print("player.destroy")
}

func (p *player) update(t time.Duration) {
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

	if !gameInstance.level().isPointSolid(nx, ny+32) && !gameInstance.level().isPointSolid(nx+31, ny+32) {
		p.anim = Falling
		p.ddy = playerGravity
	} else {
		p.anim = Standing
		if math.Abs(p.wx) > 1.0 {
			p.anim = Walking
		}
		ny = (ny / tileTemplate.frameHeight) * tileTemplate.frameHeight
		p.fy = float64(ny)
		p.dy = 0
		p.ddy = 0
	}

	if gameInstance.level().isPointSolid(nx, ny+31) {
		nx = ((nx / tileTemplate.frameWidth) + 1) * tileTemplate.frameWidth
		p.fx, p.fy = float64(nx), float64(ny)
		p.dx = 0
	}
	if gameInstance.level().isPointSolid(nx+31, ny+31) {
		nx = (nx / tileTemplate.frameWidth) * tileTemplate.frameWidth
		p.fx, p.fy = float64(nx), float64(ny)
		p.dx = 0
	}

	p.x, p.y = nx, ny

	notify("player.location", locationMsg{o: p, x: p.x, y: p.y})
}

func (p *player) handleMessage(msg message) {
	// TODO: Replace with configurable control map.
	var ctl playerControl
	switch v := msg.v.(type) {
	case *sdl.KeyDownEvent:
		switch v.KeyCode {
		case ' ':
			ctl = StartJump
		case 'w':
			log.Print("w down")
			return
		case 'a':
			ctl = StartWalkLeft
		case 's':
			log.Print("s down")
			return
		case 'd':
			ctl = StartWalkRight
		default:
			return
		}
	case *sdl.KeyUpEvent:
		switch v.KeyCode {
		case ' ':
			ctl = StopJump
		case 'w':
			log.Print("w up")
			return
		case 'a':
			ctl = StopWalkLeft
		case 's':
			log.Print("s up")
			return
		case 'd':
			ctl = StopWalkRight
		case 'e':
			ctl = Teleport
		default:
			return
		}
	default:
		return
	}

	switch ctl {
	case StartWalkLeft:
		p.facing = Left
		p.wx = -playerWalkSpeed
	case StopWalkLeft:
		p.wx = 0
	case StartWalkRight:
		p.facing = Right
		p.wx = playerWalkSpeed
	case StopWalkRight:
		p.wx = 0
	case StartJump:
		p.dy = playerJumpSpeed
		p.ddy = playerGravity
	case Land:
		p.dy = 0
		p.ddy = 0
	case Teleport:
		p.ddy = playerGravity
	default:
		// TODO: more actions
	}
	return
}

func (p *player) life() {
	defer log.Print("player.end of life")
	for msg := range p.inbox {
		//fmt.Printf("player.inbox got %+v\n", msg)
		switch msg.k {
		case "quit":
			return
		case "clock":
			p.update(msg.v.(time.Duration))
		default:
			p.handleMessage(msg)
		}
	}
}
