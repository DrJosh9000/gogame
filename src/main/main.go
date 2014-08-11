package main

import (
	"errors"
	
	"game"
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
	r := ctx.Renderer
	
	hello, err := r.LoadBMP("assets/hello.bmp")
	if err != nil {
		panic(err)
	}

	g := game.NewGame()
	defer g.Destroy()

	for {
		err = sdl.HandleEvents(eventHandler)
		if err == quitting {
			return
		}
		r.Clear()
		r.Copy(hello, nil, sdl.Rect(500, 500, 200, 200))
		r.Present()
		sdl.Delay(1)
	}
}
