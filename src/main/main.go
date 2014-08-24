// "gogame"
// Something that might turn into a game.
// Is currently my own custom set of bindings for SDL2.
// On Mac. Yay.
package main

import (
	"errors"
	"runtime"

	"game"
	"sdl"
)

const (
	defaultWidth, defaultHeight = 1024, 768
	gameName                    = "Connected Worlds"
)

var (
	quitting = errors.New("quitting")
)

func main() {
	// Must do rendering from the main thread, duh.
	runtime.LockOSThread()
	
	ctx, err := sdl.NewContext(gameName, defaultWidth, defaultHeight)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()
	r := ctx.Renderer

	g, err := game.GetGame(ctx)
	if err != nil {
		panic(err)
	}
	defer g.Destroy()

	for {
		if err := sdl.HandleEvents(func(e interface{}) error {
			switch v := e.(type) {
			case sdl.QuitEvent:
				return quitting
			case sdl.KeyEvent:
				if v.Type == sdl.KeyUp && v.KeyCode == 'q' {
					return quitting
				}
			}
			// Get the game to handle all other keys
			if g != nil {
				return g.HandleEvent(e)
			}
			return nil
		}); err == quitting {
			return
		}
		if err := r.Clear(); err != nil {
			panic(err)
		}
		if err := g.Draw(r); err != nil {
			panic(err)
		}
		r.Present()
		sdl.Delay(1)
	}
}
