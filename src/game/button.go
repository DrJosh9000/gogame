package game

import (
	"fmt"
	"log"

	"sdl"
)

// button is a clickable region on screen with optional text.
type button struct {
	*sprite
	text   *text
	action func() error
	parent *menu
}

func newButton(parent *menu, template string, label string, action func() error) (*button, error) {
	b := &button{
		sprite: &sprite{TemplateKey: template},
		action: action,
		parent: parent,
	}
	if label != "" {
		tmpl, ok := templateLibrary[template]
		if !ok {
			return nil, fmt.Errorf("no template named %q", template)
		}

		t, err := newText(ctx(), label, sdl.BlackColour, sdl.TransparentColour, sdl.CentreAlign)
		if err != nil {
			return nil, err
		}
		t.x = (tmpl.frameWidth - t.w) / 2
		t.y = (tmpl.frameHeight - t.h) / 2
		b.text = t
	}
	go b.life()
	return b, nil
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
	if err := b.text.draw(r); err != nil {
		return err
	}
	return nil
}

func (b *button) destroy() {
	if b == nil {
		return
	}
	b.text.destroy()
}

// hitTest tests screen coordinates against the button bounds.
func (b *button) hitTest(x, y int) bool {
	return !b.Invisible &&
		!b.parent.Invisible &&
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
