// Package game contains the game logic and stuff.
package game

import (
	"fmt"
	"time"

	"sdl"
)

const (
	gameTickerDuration = 10 * time.Millisecond

	level0File  = "assets/level0.txt"
	level1AFile = "assets/level1a.txt"
	level1BFile = "assets/level1b.txt"
)

var gameInstance *Game

type Game struct {
	ctx    *sdl.Context
	t0     time.Time
	ticker *time.Ticker
	inbox  chan message

	sr    *sdl.Renderer
	wr    *worldRenderer
	world ComplexBase
	hud   ComplexBase

	// Special game objects.
	// TODO: OHDOG
	// cursor: keep this around to send mouse events to directly.
	cursor       *cursor
	player       *player
	exit         *exit
	levels       [2]*level
	terrains     [2]*terrain
	currentLevel int
}

func NewGame(ctx *sdl.Context) (*Game, error) {
	if gameInstance != nil {
		return gameInstance, nil
	}

	c, err := newCursor(ctx)
	if err != nil {
		return nil, err
	}

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

	g := &Game{
		ctx: ctx,
		t0:  time.Now(),
		//ticker: time.NewTicker(gameTickerDuration),
		sr: ctx.Renderer,
		wr: &worldRenderer{
			r:     ctx.Renderer,
			view:  sdl.Rect{0, 0, 1024, 768},
			world: sdl.Rect{0, 0, 3072, 768}, // TODO: derive from terrain
		},
		inbox:        make(chan message, 10),
		cursor:       c,
		player:       p,
		levels:       [2]*level{m0, m1},
		terrains:     [2]*terrain{t0, t1},
		currentLevel: 0,
	}
	gameInstance = g
	p.x, p.y = tileTemplate.frameWidth*m0.startX, tileTemplate.frameHeight*m0.startY
	g.wr.focus(p.x, p.y)

	g.world.AddChild(p)
	g.hud.AddChild(c)

	kmp("player.location", g.inbox)
	go g.messageLoop()

	return g, nil
}

func (g *Game) messageLoop() {
	for msg := range g.inbox {
		switch m := msg.v.(type) {
		case locationMsg:
			if msg.k == "player.location" {
				g.wr.focus(m.x, m.y)
			}
		} // switch msg.(type)
	} // for range g.inbox
}

func (g *Game) Draw() error {
	// Draw current terrain in world coordinates.
	if err := g.terrains[g.currentLevel].Draw(g.wr); err != nil {
		return err
	}

	// Draw everything in the world in world coordinates.
	if err := g.world.Draw(g.wr); err != nil {
		return err
	}

	// Draw the HUD in screen coordinates.
	return g.hud.Draw(g.sr)
}

func (g *Game) Destroy() {
	notify("game", quitMsg)
	g.player.controller <- Quit
	//	g.ticker.Stop()
	g.world.Destroy()
	g.hud.Destroy()
}

func (g *Game) level() *level {
	return g.levels[g.currentLevel]
}

func (g *Game) HandleEvent(ev interface{}) error {
	switch v := ev.(type) {
	case sdl.KeyDownEvent:
		switch v.KeyCode {
		case ' ':
			g.player.controller <- StartJump
		case 'w':
			fmt.Println("w down")
		case 'a':
			g.player.controller <- StartWalkLeft
		case 's':
			fmt.Println("s down")
		case 'd':
			g.player.controller <- StartWalkRight
		}
	case sdl.KeyUpEvent:
		switch v.KeyCode {
		case ' ':
			g.player.controller <- StopJump
		case 'w':
			fmt.Println("w up")
		case 'a':
			g.player.controller <- StopWalkLeft
		case 's':
			fmt.Println("s up")
		case 'd':
			g.player.controller <- StopWalkRight
		case 'e':
			g.currentLevel = (g.currentLevel + 1) % 2
			g.player.controller <- Teleport
		}
	case sdl.MouseMotionEvent:
		g.cursor.controller <- v
	}
	return nil
}

type worldRenderer struct {
	r           Renderer
	view, world sdl.Rect
}

// focus moves the world renderer viewport to include the point. Generally this
// would be used to focus on the player. It snaps immediately, no smoothing.
func (r *worldRenderer) focus(x, y int) {
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

func (r *worldRenderer) Copy(t *sdl.Texture, src, dst sdl.Rect) error {
	dst.X -= r.view.X
	dst.Y -= r.view.Y
	return r.r.Copy(t, src, dst)
}

func (r *worldRenderer) CopyEx(t *sdl.Texture, src, dst sdl.Rect, angle float64, center sdl.Point, flip sdl.RendererFlip) error {
	dst.X -= r.view.X
	dst.Y -= r.view.Y
	center.X -= r.view.X
	center.Y -= r.view.Y
	return r.r.CopyEx(t, src, dst, angle, center, flip)
}
