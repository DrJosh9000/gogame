package game

import (
	"time"

	"sdl"
)

var exitTemplate = &spriteTemplate{
	name:        "exit",
	sheetFile:   "assets/door.png",
	framesX:     4,
	framesY:     1,
	frameWidth:  64,
	frameHeight: 64,
}

type DoorState int

const (
	DoorStateClosed DoorState = iota
	DoorStateOpen
	DoorStateQuit
)

type Exit struct {
	*sprite
	DoorState

	inbox   chan message
	updater *time.Ticker
}

func NewExit(ctx *sdl.Context) (*Exit, error) {
	s, err := exitTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	e := &Exit{
		sprite:  s,
		inbox:   make(chan message, 10),
		updater: time.NewTicker(100 * time.Millisecond),
	}
	kmp("global", e.inbox)
	kmp("player.location", e.inbox)
	go e.life()
	return e, nil
}

func (e *Exit) Destroy() {
	e.updater.Stop()
}

func (e *Exit) life() {
	for {
		select {
		case msg := <-e.inbox:
			switch m := msg.v.(type) {
			case basicMsg:
				if msg.k == "game" && m == quitMsg {
					return
				}
			case locationMsg:
				if msg.k == "player.location" {
					// If the player is near the door, open it;
					// If the player is not near the door, close it.
					if e.DoorState == DoorStateClosed &&
						m.x > e.x-200 && m.x < e.x+200 &&
						m.y > e.y-200 && m.y < e.y+200 {
						e.DoorState = DoorStateOpen
					}
					if e.DoorState == DoorStateOpen && (m.x <= e.x-200 || m.x >= e.x+200 ||
						m.y <= e.y-200 || m.y >= e.y+200) {
						e.DoorState = DoorStateClosed
					}
				}
			}
		case <-e.updater.C:
			switch {
			case e.DoorState == DoorStateClosed && e.frame > 0:
				e.frame--
			case e.DoorState == DoorStateOpen && e.frame < 3:
				e.frame++
			}
		}
	}
}
