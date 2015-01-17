// Package game contains the game logic and stuff.
package game

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sdl"
)

const (
	clockDuration = 10 * time.Millisecond
)

type gameState int

const (
	gameStateRunning = iota
	gameStateMenu
	gameStateQuitting
)

var gameInstance *Game

type Game struct {
	state gameState
	ctx   *sdl.Context
	clock *time.Ticker

	renderer   *sdl.Renderer
	wv         *worldView
	world, hud complexBase
	windows    windowManager

	cursor *sprite
	menu   *menu
}

func ctx() *sdl.Context {
	return gameInstance.ctx
}

func NewGame(ctx *sdl.Context, width, height int) (*Game, error) {
	if gameInstance != nil {
		return gameInstance, nil
	}

	rand.Seed(time.Now().Unix())
	gameInstance = &Game{
		state:    gameStateMenu,
		ctx:      ctx,
		clock:    time.NewTicker(clockDuration),
		renderer: ctx.Renderer,
		wv: &worldView{
			view:  sdl.Rect{0, 0, width, height},
			world: sdl.Rect{0, 0, 4096, 768}, // TODO: derive from terrain
		},
		cursor: &sprite{TemplateKey: "cursor"},
	}

	go gameInstance.windows.life()
	go cursorLife(gameInstance.cursor)

	menu, err := newMenu(mainMenu, &gameInstance.windows)
	if err != nil {
		return nil, err
	}
	gameInstance.menu = menu

	go gameInstance.life()
	go gameInstance.pulse()
	return gameInstance, nil
}

// life listens to and handles events.
func (g *Game) life() {
	inbox := make(chan message, 10)
	kmp("quit", inbox)
	kmp("player.location", inbox)
	kmp("input.event", inbox)
	kmp("menuAction", inbox)
	defer func() {
		g.state = gameStateQuitting
	}()
	for msg := range inbox {
		//log.Printf("game.inbox got %+v\n", msg)
		switch msg.k {
		case "quit":
			return
		case "menuAction":
			switch msg.v.(string) {
			case "start":
				g.menu.Invisible = true
				g.state = gameStateRunning
			case "levelEdit":
				g.menu.Invisible = true
			}
		}
		switch m := msg.v.(type) {
		case locationMsg:
			if msg.k == "player.location" {
				g.wv.focus(m.x, m.y)
			}
		case sdl.QuitEvent:
			quit()
		case *sdl.KeyUpEvent:
			switch m.KeyCode {
			case 'q':
				quit()
			case 'e':
				// Do teleport
				//g.lev.active = (g.lev.active + 1) % 2
			}
		}
	}
}

// loadWorld loads g.world from a gob-encoded file.
func (g *Game) loadWorld(worldFile string) error {
	f, err := os.Open(worldFile)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	return dec.Decode(&g.world)
}

// saveWorld writes g.world to a temporary file, and then moves it
// onto worldFile.
func (g *Game) saveWorld(worldFile string) error {
	f, err := ioutil.TempFile(filepath.Split(worldFile))
	if err != nil {
		return err
	}
	tmpFile := f.Name()
	defer func() {
		f.Close()
		os.Remove(tmpFile)
	}()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(&g.world); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmpFile, worldFile)
}

// pulse notifies the "clock" key with events from g.clock when the game is
// playing.
func (g *Game) pulse() {
	t0 := time.Now()
	for t := range g.clock.C {
		if g.state == gameStateRunning {
			notify("clock", t.Sub(t0))
		}
	}
}

// Draw renders the game scene & any menus/HUD.
func (g *Game) Draw() error {
	// Draw everything in the world in world coordinates.
	g.renderer.PushOffset(-g.wv.view.X, -g.wv.view.Y)
	defer g.renderer.ResetOffset()
	if err := g.world.draw(g.renderer); err != nil {
		return err
	}
	g.renderer.PopOffset()

	if err := g.hud.draw(g.renderer); err != nil {
		return err
	}
	if err := g.menu.draw(g.renderer); err != nil {
		return err
	}
	return g.cursor.draw(g.renderer)
}

func (g *Game) Destroy() {
	//log.Print("game.destroy")
	g.clock.Stop()
	g.world.destroy()
	g.hud.destroy()
}

// Exec runs a command string. Commands are usually entered on the terminal.
func (g *Game) Exec(cmd string) {
	//log.Printf("game.Exec(%q)\n", cmd)
	argv := strings.Split(cmd, " ")
	switch argv[0] {
	case "quit":
		quit()
	case "help":
		if len(argv) == 1 {
			fmt.Println("help: Usage: help <command>")
		} else {
			fmt.Println("help: Not yet implemented")
		}
	case "load":
		if len(argv) != 2 {
			fmt.Println("load: Wrong number of arguments")
		}
		if err := g.loadWorld(argv[1]); err != nil {
			fmt.Println("load:", err)
			return
		}
		fmt.Println("load: Success")
	case "save":
		if len(argv) != 2 {
			fmt.Println("save: Wrong number of arguments")
		}
		if err := g.saveWorld(argv[1]); err != nil {
			fmt.Println("save:", err)
			return
		}
		fmt.Println("save: Success")
	case "testworld":
		if len(argv) != 1 {
			fmt.Println("instaworld: Wrong number of arguments")
		}
		g.testWorld()
	case "n", "notify":
		if len(argv) != 3 {
			fmt.Printf("notify: Wrong number of arguments [%d != 3]\n", len(argv))
		}
		notify(argv[1], argv[2])
	case "t", "twoot":
		g.hud.addChild(easyTwoot(strings.Join(argv[1:], " ")))
	case "":
		return
	default:
		fmt.Println("Bad command or file name")
	}
}

// testWorld replaces the world with a single screen full of hexes and an orb.
func (g *Game) testWorld() {
	g.world.Kids = nil
	for i := 0; i < 25; i++ {
		for j := 0; j < 8; j++ {
			g.world.addChild(&sprite{
				TemplateKey: "hex",
				X:           192*j + 96*(i%2) - 32,
				//Y:           -rand.Intn(5) * 2,
				Z: 32 * (i - 1),
			})
		}
	}
	g.world.addChild(&orb{X: 150, Y: 22, Z: 100, Selected: true})
}

func (g *Game) Quitting() bool {
	return g.state == gameStateQuitting
}

func (g *Game) HandleEvent(ev sdl.Event) error {
	notify("input.event", ev)
	return nil
}

type worldView struct {
	view, world sdl.Rect
}

// focus moves the world viewport to include the point. Generally this
// would be used to focus on the player. It snaps immediately, no smoothing.
func (r *worldView) focus(x, y int) {
	// Keep the point in view.
	left, right := r.view.W/4, 3*r.view.W/4
	if x-r.view.X > right {
		r.view.X = x - right
	}
	if x-r.view.X < left {
		r.view.X = x - left
	}
	// Clamp to world bounds.
	if r.view.X < r.world.X {
		r.view.X = r.world.X
	}
	if r.view.X+r.view.W > r.world.X+r.world.W {
		r.view.X = r.world.X + r.world.W - r.view.W
	}

	top, bottom := r.view.H/4, 3*r.view.H/4
	if y-r.view.Y < top {
		r.view.Y = y - top
	}
	if y-r.view.Y > bottom {
		r.view.Y = y - bottom
	}
	if r.view.Y < r.world.Y {
		r.view.Y = r.world.Y
	}
	if r.view.Y+r.view.H > r.world.Y+r.world.H {
		r.view.Y = r.world.Y + r.world.H - r.view.H
	}
}
