package game

import (
	"fmt"
	"math"
	"time"

	"sdl"
)

var playerTemplate = &spriteTemplate{
	name:        "player",
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

	kmp("global", p.inbox)
	kmp("input.event", p.inbox)
	go p.life()
	return p, nil
}

func (p *player) destroy() {
	fmt.Println("player.destroy")
}

func (p *player) String() string {
	return "player"
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

func (p *player) handleMessage(msg message) bool {
	// TODO: Replace with configurable control map.
	var ctl playerControl
	switch v := msg.v.(type) {
	case *sdl.KeyDownEvent:
		switch v.KeyCode {
		case ' ':
			ctl = StartJump
		case 'w':
			fmt.Println("w down")
			return false
		case 'a':
			ctl = StartWalkLeft
		case 's':
			fmt.Println("s down")
			return false
		case 'd':
			ctl = StartWalkRight
		default:
			return false
		}
	case *sdl.KeyUpEvent:
		switch v.KeyCode {
		case ' ':
			ctl = StopJump
		case 'w':
			fmt.Println("w up")
			return false
		case 'a':
			ctl = StopWalkLeft
		case 's':
			fmt.Println("s up")
			return false
		case 'd':
			ctl = StopWalkRight
		case 'e':
			ctl = Teleport
		default:
			return false
		}
	case basicMsg:
		if v == quitMsg {
			return true
		}
	default:
		return false
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
	return false
}

func (p *player) life() {
	updater := time.NewTicker(playerUpdateInterval)
	defer func() {
		updater.Stop()
		close(p.inbox)
		fmt.Println("player.end of life")
	}()
	t0 := time.Now()
	for {
		select {
		case c := <-p.inbox:
			//fmt.Printf("player.inbox got %+v\n", c)
			if p.handleMessage(c) {
				return
			}
		case t := <-updater.C:
			//fmt.Printf("player.updater got %+v\n", t)
			p.update(t.Sub(t0))
		}
	}
}
