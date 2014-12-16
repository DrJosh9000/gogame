package game

import (
	"sdl"
)

var exitTemplate = &spriteTemplate{
	sheetFile:   "assets/door.png",
	framesX:     4,
	framesY:     1,
	frameWidth:  64,
	frameHeight: 64,
}

type exit struct {
	*sprite
	state, wantState int
	inbox            chan message
}

func newExit(ctx *sdl.Context) (*exit, error) {
	s, err := exitTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	e := &exit{
		sprite: s,
		inbox:  make(chan message, 10),
	}
	kmp("quit", e.inbox)
	kmp("clock", e.inbox)
	kmp("player.location", e.inbox)
	go e.life()
	return e, nil
}

func (e *exit) life() {
	for msg := range e.inbox {
		switch msg.k {
		case "quit":
			return
		case "clock":
			if e.wantState > e.state {
				e.state++
			}
			if e.wantState < e.state {
				e.state--
			}
			e.frame = e.state / 10
		case "player.location":
			m := msg.v.(locationMsg)
			// If the player is near the door, open it;
			// If the player is not near the door, close it.
			if m.x > e.x-200 && m.x < e.x+200 &&
				m.y > e.y-200 && m.y < e.y+200 {
				e.wantState = 30
			}
			if m.x <= e.x-200 || m.x >= e.x+200 ||
				m.y <= e.y-200 || m.y >= e.y+200 {
				e.wantState = 0
			}
		}
	}
}
