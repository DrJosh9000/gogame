package game

import (
	"sdl"
)

var cursorTemplate = &spriteTemplate{
	name:        "cursor",
	sheetFile:   "assets/cursor.png",
	framesX:     2,
	framesY:     1,
	frameWidth:  64,
	frameHeight: 64,
}

type cursor struct {
	*sprite
	controller chan interface{}
}

func newCursor(ctx *sdl.Context) (*cursor, error) {
	s, err := cursorTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	c := &cursor{
		sprite:     s,
		controller: make(chan interface{}, 3),
	}
	go c.life()
	return c, nil
}

func (c *cursor) Destroy() {
	close(c.controller)
}

func (c *cursor) life() {
	for ev := range c.controller {
		switch m := ev.(type) {
		case sdl.MouseButtonDownEvent:
			c.x = m.X - cursorTemplate.frameWidth/2
			c.y = m.Y - cursorTemplate.frameHeight/2
			c.frame = 1
		case sdl.MouseButtonUpEvent:
			c.x = m.X - cursorTemplate.frameWidth/2
			c.y = m.Y - cursorTemplate.frameHeight/2
			c.frame = 0
		case sdl.MouseMotionEvent:
			c.x = m.X - cursorTemplate.frameWidth/2
			c.y = m.Y - cursorTemplate.frameHeight/2
			c.frame = 0
			if m.ButtonState&sdl.MouseLeftMask != 0 {
				c.frame = 1
			}
			if c.invisible {
				c.invisible = false
				sdl.HideCursor()
			}
		case sdl.WindowEvent:
			c.invisible = true
			switch m.EventID {
			case sdl.WindowEnter:
				c.invisible = false
			}
			if c.invisible {
				sdl.ShowCursor()
			} else {
				sdl.HideCursor()
			}
		}

	}
}
