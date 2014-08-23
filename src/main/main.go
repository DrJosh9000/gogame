// "gogame"
// Something that might turn into a game.
// Is currently my own custom set of bindings for SDL2.
// On Mac. Yay.
package main

import (
	"errors"
	
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
	ctx, err := sdl.NewContext(gameName, defaultWidth, defaultHeight)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()
	r := ctx.Renderer

	g, err := game.NewGame(ctx)
	if err != nil {
		panic(err)
	}
	defer g.Destroy()

	for {
		err := sdl.HandleEvents(func (e interface{}) error {
			switch v := e.(type) {
			case sdl.QuitEvent:
				return quitting
			case sdl.KeyEvent:
				if v.Type == sdl.KeyUp {
					switch v.KeyCode {
					case 'q':
						return quitting
					default:
						// Get the game to handle all other events
						if g != nil {
							return g.HandleKey(v.KeyCode)
						}
					}
				}
			}
			return nil
		})
		if err == quitting {
			return
		}
		r.Clear()
		g.Draw()
		r.Present()
		sdl.Delay(1)
	}
}
