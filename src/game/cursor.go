package game

import (
	"sdl"
)

func cursorLife(c *sprite) {
	inbox := make(chan message, 10)
	kmp("input.event", inbox)
	for msg := range inbox {
		//log.Printf("cursor.inbox got %+v\n", msg)
		switch m := msg.v.(type) {
		case *sdl.MouseButtonDownEvent:
			c.X, c.Y = m.X, m.Y
			c.Frame = 1
		case *sdl.MouseButtonUpEvent:
			c.X, c.Y = m.X, m.Y
			c.Frame = 0
		case *sdl.MouseMotionEvent:
			c.X, c.Y = m.X, m.Y
			c.Frame = 0
			if m.ButtonState&sdl.MouseLeftMask != 0 {
				c.Frame = 1
			}
			if c.Invisible {
				c.Invisible = false
				sdl.HideCursor()
			}
		case *sdl.WindowEvent:
			c.Invisible = true
			switch m.EventID {
			case sdl.WindowEnter:
				c.Invisible = false
			}
			if c.Invisible {
				sdl.ShowCursor()
			} else {
				sdl.HideCursor()
			}
		}
	}
}
