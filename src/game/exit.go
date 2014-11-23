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

type exitState int

const (
	exitStateClosed exitState = iota
	exitStateOpen
	exitStateQuit
)

type exit struct {
	*sprite
	exitState

	inbox   chan message
	updater *time.Ticker
}

func newExit(ctx *sdl.Context) (*exit, error) {
	s, err := exitTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	e := &exit{
		sprite:  s,
		inbox:   make(chan message, 10),
		updater: time.NewTicker(100 * time.Millisecond),
	}
	kmp("global", e.inbox)
	kmp("player.location", e.inbox)
	go e.life()
	return e, nil
}

func (e *exit) Destroy() {
	e.updater.Stop()
}

func (e *exit) life() {
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
					if e.exitState == exitStateClosed &&
						m.x > e.x-200 && m.x < e.x+200 &&
						m.y > e.y-200 && m.y < e.y+200 {
						e.exitState = exitStateOpen
					}
					if e.exitState == exitStateOpen && (m.x <= e.x-200 || m.x >= e.x+200 ||
						m.y <= e.y-200 || m.y >= e.y+200) {
						e.exitState = exitStateClosed
					}
				}
			}
		case <-e.updater.C:
			switch {
			case e.exitState == exitStateClosed && e.frame > 0:
				e.frame--
			case e.exitState == exitStateOpen && e.frame < 3:
				e.frame++
			}
		}
	}
}
