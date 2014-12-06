package game

import (
	"log"

	"sdl"
)

// TODO: more button styles
var buttonTemplate = &spriteTemplate{
	sheetFile:   "assets/button.png",
	framesX:     1,
	framesY:     2,
	frameWidth:  256,
	frameHeight: 64,
}

type button struct {
	*sprite
	text   *text
	inbox  chan message
	action func() error
}

func newButton(ctx *sdl.Context, label string, action func() error) (*button, error) {
	s, err := buttonTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	t, err := newText(ctx, label, sdl.BlackColour, sdl.TransparentColour, sdl.CentreAlign)
	if err != nil {
		return nil, err
	}
	b := &button{
		sprite: s,
		text:   t,
		inbox:  make(chan message, 10),
		action: action,
	}
	kmp("global", b.inbox)
	kmp("input.event", b.inbox)
	go b.life()
	return b, nil
}

func (b *button) draw(r renderer) error {
	if b == nil || b.invisible {
		return nil
	}
	if err := b.sprite.draw(r); err != nil {
		return err
	}
	b.text.x = b.x + (b.template.frameWidth-b.text.w)/2
	b.text.y = b.y + (b.template.frameHeight-b.text.h)/2
	if b.frame == 1 {
		b.text.y += 8
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

func (b *button) hitTest(x, y int) bool {
	return !b.invisible && x >= b.x && x <= b.x+b.template.frameWidth && y >= b.y && y <= b.y+b.template.frameHeight
}

func (b *button) life() {
	for msg := range b.inbox {
		switch m := msg.v.(type) {
		case basicMsg:
			if m == quitMsg {
				return
			}
		case *sdl.MouseButtonDownEvent:
			if b.hitTest(m.X, m.Y) {
				b.frame = 1
			}
		case *sdl.MouseButtonUpEvent:
			if b.hitTest(m.X, m.Y) {
				if b.frame == 1 {
					if err := b.action(); err != nil {
						log.Printf("error running menu item action: %v\n", err)
					}
				}
				b.frame = 0
			}
		case *sdl.MouseMotionEvent:
			if !b.hitTest(m.X, m.Y) {
				b.frame = 0
			}
		}
	}
}
