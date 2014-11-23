package game

import (
	"sdl"
)

var cursorTemplate = &spriteTemplate{
	name:        "cursor",
	sheetFile:   "assets/cursor.png",
	framesX:     1,
	framesY:     1,
	frameWidth:  64,
	frameHeight: 64,
}

type cursor struct {
	*sprite
	controller chan sdl.MouseMotionEvent
}

func newCursor(ctx *sdl.Context) (*cursor, error) {
	s, err := cursorTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	c := &cursor{
		sprite: s,
		controller: make(chan sdl.MouseMotionEvent, 3),
	}
	go c.life()
	return c, nil
}

func (c *cursor) Destroy() {
	close(c.controller)
}

func (c *cursor) life() {
	for ev := range c.controller {
		c.x = ev.X - cursorTemplate.frameWidth / 2
		c.y = ev.Y - cursorTemplate.frameHeight / 2
	}
}