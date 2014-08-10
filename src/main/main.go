package main

import (
	"errors"
	"sdl"
)

const (
	defaultWidth, defaultHeight = 1024, 768
	gameName                    = "gogame"
)

var (
	quitting = errors.New("quitting")
)

func eventHandler(e interface{}) error {
	switch v := e.(type) {
	case sdl.QuitEvent:
		return quitting
	case sdl.KeyEvent:
		if v.Type == sdl.KeyUp {
			switch v.KeyCode {
			case 'q':
				return quitting
			}
		}
	}
	return nil
}

func main() {
	ctx, err := sdl.NewContext(gameName, defaultWidth, defaultHeight)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	for {
		err = sdl.HandleEvents(eventHandler)
		if err == quitting {
			return
		}
		ctx.Render()
		sdl.Delay(1)
	}
}
