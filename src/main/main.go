// Press E To Teleport (working title) - game(ish) made for Ludum Dare 30
// Programmer, etc: @DrJosh9000
package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"game"
	"sdl"
)

const (
	defaultWidth, defaultHeight = 1024, 768
	gameName                    = "Gate Simulator"
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

	go console(g)
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

// console processes text entered on the tty.
func console(g *game.Game) {
	fmt.Printf("> ")
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		g.Exec(in.Text())
		fmt.Printf("> ")
	}
	if err := in.Err(); err != nil {
		panic(err)
	}
}
