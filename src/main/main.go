// Press E To Teleport (working title) - game(ish) made for Ludum Dare 30
// Programmer, etc: @DrJosh9000
package main

import (
	"errors"
	"runtime"

	"game"
	"sdl"
)

const (
	defaultWidth, defaultHeight = 1024, 768
	gameName                    = "Press E to Teleport"
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

	g, err := game.NewGame(ctx)
	if err != nil {
		panic(err)
	}
	defer g.Destroy()

	for !g.Quitting() {
		if err := sdl.HandleEvents(g.HandleEvent); err != nil {
			panic(err)
		}
		r.SetDrawColour(sdl.BlackColour)
		if err := r.Clear(); err != nil {
			panic(err)
		}
		if err := g.Draw(); err != nil {
			panic(err)
		}
		r.Present()
		sdl.Delay(1)
	}
}
