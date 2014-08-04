package main

import (
	"errors"
	"sdl"
)

const (
	defaultWidth, defaultHeight = 1024, 768
	gameName = "gogame"
)

func main() {
	ctx, err := sdl.NewContext(gameName, defaultWidth, defaultHeight)
	if err != nil {
		panic(err)
	}

	quit := 	errors.New("quitting")
	for {
		err = sdl.HandleEvents(func(e interface{}) error {
			switch v := e.(type) {
			case sdl.QuitEvent:
				return quit
			case sdl.KeyEvent:
				if v.Type == sdl.KeyUp && v.KeyCode == 'q' {
					return quit
				}
			}
			return nil
		})
		if err == quit {
			return
		}
		ctx.Render()
		sdl.Delay(1)
	}
}