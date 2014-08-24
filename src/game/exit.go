package game

import (
	"time"
	
	"sdl"
)

const (
	exitFile = "assets/door.png"
	exitFrames = 4
	exitFrameWidth, exitFrameHeight = 64, 64
)

type DoorState int
const (
	DoorStateClosed DoorState = iota
	DoorStateOpen
	DoorStateQuit
)

type Exit struct {
	tex *sdl.Texture
	DoorState
	x, y, frame int
	
	Updater *time.Ticker
	Controller chan DoorState
}

func NewExit(ctx *sdl.Context) (*Exit, error) {
	tex, err := ctx.GetTexture(exitFile)
	if err != nil {
		return nil, err
	}
	e := &Exit{
		tex: tex,
		Controller: make(chan DoorState),
		Updater: time.NewTicker(100 * time.Millisecond),
	}
	go e.life()
	return e, nil
}

func (e *Exit) AddChild(Object) {}
func (e *Exit) Children() []Object {return nil}

func (e *Exit) Destroy() {
	e.Controller <- DoorStateQuit
	e.Updater.Stop()
}

func (e *Exit) Draw(r *sdl.Renderer) error {
	return r.Copy(e.tex, sdl.Rect(e.frame*exitFrameWidth, 0, exitFrameWidth, exitFrameHeight), sdl.Rect(e.x, e.y, exitFrameWidth, exitFrameHeight))
}

func (e* Exit) life() {
	for {
		select {
		case c := <-e.Controller:
			if c == DoorStateQuit {
				return
			}
			e.DoorState = c
		case <-e.Updater.C:
			switch  {
			case e.DoorState == DoorStateClosed && e.frame > 0:
				e.frame--
			case e.DoorState == DoorStateOpen && e.frame < 3:
				e.frame++
			}
		}
	}
}