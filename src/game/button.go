package game

import (
	"log"

	"sdl"
)

// button is a clickable region on screen with optional text.
type button struct {
	*sprite
	Label  string
	text   *text
	action func() error
	parent *menu
}

func (b *button) draw(r *sdl.Renderer) error {
	if b == nil || b.Invisible {
		return nil
	}
	if err := b.sprite.draw(r); err != nil {
		return err
	}

	r.PushOffset(b.X, b.Y)
	defer r.PopOffset()
	if b.Frame == 1 {
		r.PushOffset(0, 8)
		defer r.PopOffset()
	}
	if b.text == nil {
		t, err := newText(b.Label, sdl.BlackColour, sdl.TransparentColour, sdl.CentreAlign)
		if err != nil {
			return err
		}
		t.x = (b.template.frameWidth - t.w) / 2
		t.y = (b.template.frameHeight - t.h) / 2
		b.text = t
	}
	if err := b.text.draw(r); err != nil {
		return err
	}
	return nil
}

// hitTest tests screen coordinates against the button bounds.
func (b *button) hitTest(x, y int) bool {
	if b == nil || b.Invisible {
		return false
	}
	if b.parent == nil {
		return x >= b.X && x <= b.X+b.template.frameWidth &&
			y >= b.Y && y <= b.Y+b.template.frameHeight
	}
	return !b.parent.Invisible &&
		x >= b.X+b.parent.X &&
		x <= b.X+b.parent.X+b.template.frameWidth &&
		y >= b.Y+b.parent.Y &&
		y <= b.Y+b.parent.Y+b.template.frameHeight
}

func (b *button) life() {
	inbox := make(chan message, 10)
	kmp("quit", inbox)
	kmp("input.event", inbox)
	for msg := range inbox {
		if msg.k == "quit" {
			return
		}
		switch m := msg.v.(type) {
		case *sdl.MouseButtonDownEvent:
			if b.hitTest(m.X, m.Y) {
				b.Frame = 1
			}
		case *sdl.MouseButtonUpEvent:
			if b.hitTest(m.X, m.Y) {
				if b.Frame == 1 {
					if err := b.action(); err != nil {
						log.Printf("error running menu item action: %v\n", err)
					}
				}
				b.Frame = 0
			}
		case *sdl.MouseMotionEvent:
			if !b.hitTest(m.X, m.Y) {
				b.Frame = 0
			}
		}
	}
}
