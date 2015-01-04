package game

import (
	"sdl"
)

// button is a clickable region on screen with optional text.
type button struct {
	*sprite
	text   text
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
	b.text.X = (b.template.frameWidth - b.text.w) / 2
	b.text.Y = (b.template.frameHeight - b.text.h) / 2
	if err := b.text.draw(r); err != nil {
		return err
	}
	return nil
}

func (b *button) invisible() bool {
	return b == nil || b.parent == nil || b.parent.invisible() || b.sprite.invisible()
}

func (b *button) setDown(down bool) {
	b.Frame = 0
	if down {
		b.Frame = 1
	}
}

func (b *button) click() {
	b.action()
}
