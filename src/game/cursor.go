package game

import (
	"sdl"
)

var cursorTemplate = &spriteTemplate{
	sheetFile:   "assets/cursor.png",
	baseX:       32,
	baseY:       32,
	framesX:     2,
	framesY:     1,
	frameWidth:  64,
	frameHeight: 64,
}

type cursor struct {
	*sprite
	inbox chan message
}

func newCursor(ctx *sdl.Context) (*cursor, error) {
	s, err := cursorTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	c := &cursor{
		sprite: s,
		inbox:  make(chan message, 10),
	}
	kmp("input.event", c.inbox)
	go c.life()
	return c, nil
}

func (c *cursor) life() {
	for msg := range c.inbox {
		//log.Printf("cursor.inbox got %+v\n", msg)
		switch m := msg.v.(type) {
		case *sdl.MouseButtonDownEvent:
			c.x, c.y = m.X, m.Y
			c.frame = 1
		case *sdl.MouseButtonUpEvent:
			c.x, c.y = m.X, m.Y
			c.frame = 0
		case *sdl.MouseMotionEvent:
			c.x, c.y = m.X, m.Y
			c.frame = 0
			if m.ButtonState&sdl.MouseLeftMask != 0 {
				c.frame = 1
			}
			if c.invisible {
				c.invisible = false
				sdl.HideCursor()
			}
		case *sdl.WindowEvent:
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
