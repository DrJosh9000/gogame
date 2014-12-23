// Package game contains the game logic and stuff.
package game

import (
	"log"
	"time"

	"sdl"
)

const (
	clockDuration = 10 * time.Millisecond

	level0File  = "assets/level0.txt"
	level1AFile = "assets/level1a.txt"
	level1BFile = "assets/level1b.txt"
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
	inbox chan message

	renderer   *sdl.Renderer
	wv         *worldView
	world, hud complexBase

	// Special game objects.
	// TODO: OHDOG
	cursor *cursor
	menu   *menu
	player *player
	exit   *exit
	lev    unionObject
	levels [2]*level
}

func NewGame(ctx *sdl.Context) (*Game, error) {
	if gameInstance != nil {
		return gameInstance, nil
	}

	c, err := newCursor(ctx)
	if err != nil {
		return nil, err
	}

	/*
		p, err := newPlayer(ctx)
		if err != nil {
			return nil, err
		}

		m0, err := loadLevel(level1AFile)
		if err != nil {
			return nil, err
		}
		m1, err := loadLevel(level1BFile)
		if err != nil {
			return nil, err
		}

		t0, err := newTerrain(ctx, m0)
		if err != nil {
			return nil, err
		}
		t1, err := newTerrain(ctx, m1)
		if err != nil {
			return nil, err
		}
	*/

	menu, err := newMenu(ctx)
	if err != nil {
		return nil, err
	}

	g := &Game{
		state:    gameStateMenu,
		ctx:      ctx,
		clock:    time.NewTicker(clockDuration),
		renderer: ctx.Renderer,
		wv: &worldView{
			view:  sdl.Rect{0, 0, 1024, 768},
			world: sdl.Rect{0, 0, 4096, 768}, // TODO: derive from terrain
		},
		inbox:  make(chan message, 10),
		cursor: c,
		menu:   menu,
		//player: p,
		//levels: [2]*level{m0, m1},
	}
	gameInstance = g
	//p.x, p.y = tileTemplate.frameWidth*m0.startX, tileTemplate.frameHeight*m0.startY
	//g.wv.focus(p.x, p.y)

	//g.lev.addChild(t0)
	//g.lev.addChild(t1)
	//g.world.addChild(&g.lev)
	for i := 0; i < 4; i++ {
		h, err := newHex(ctx)
		if err != nil {
			return nil, err
		}
		h.x, h.y = i*96, i*32
		g.world.addChild(h)
	}
	//g.world.addChild(p)

	kmp("quit", g.inbox)
	kmp("player.location", g.inbox)
	kmp("input.event", g.inbox)
	kmp("menuAction", g.inbox)
	go g.life()
	go g.pulse()
	return g, nil
}

func (g *Game) life() {
	defer func() {
		g.state = gameStateQuitting
	}()
	for msg := range g.inbox {
		//log.Printf("game.inbox got %+v\n", msg)
		switch msg.k {
		case "quit":
			return
		case "menuAction":
			switch msg.v.(string) {
			case "start":
				g.menu.invisible = true
				g.state = gameStateRunning
			case "levelEdit":
				g.menu.invisible = true
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
				g.lev.active = (g.lev.active + 1) % 2
			}
		}
	}
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

func (g *Game) Draw() error {
	// Draw everything in the world in world coordinates.
	g.renderer.PushOffset(-g.wv.view.X, -g.wv.view.Y)
	defer g.renderer.ResetOffset()
	if err := g.world.draw(g.renderer); err != nil {
		return err
	}
	g.renderer.PopOffset()

	// Draw the HUD in screen coordinates.
	if err := g.hud.draw(g.renderer); err != nil {
		return err
	}

	// Draw the menu in screen coordinates.
	if err := g.menu.draw(g.renderer); err != nil {
		return err
	}

	// Draw the cursor, always, in screen coordinates.
	return g.cursor.draw(g.renderer)
}

func (g *Game) Destroy() {
	log.Print("game.destroy")
	g.clock.Stop()
	g.world.destroy()
	g.hud.destroy()
}

func (g *Game) Quitting() bool {
	return g.state == gameStateQuitting
}

func (g *Game) level() *level {
	return g.levels[g.lev.active]
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
