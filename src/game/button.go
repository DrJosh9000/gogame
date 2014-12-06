package game

import (
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

type buttonState int

const (
	buttonUp = buttonState(iota)
	buttonDown
)

type button struct {
	sprite *sprite
	text   *text
}

func newButton(ctx *sdl.Context, label string) (*button, error) {
	s, err := buttonTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	t, err := newText(ctx, label, sdl.BlackColour, sdl.TransparentColour, sdl.CentreAlign)
	if err != nil {
		return nil, err
	}
	return &button{sprite: s, text: t}, nil
}

func (b *button) draw(r renderer) error {
	if b == nil || b.sprite.invisible {
		return nil
	}
	if err := b.sprite.draw(r); err != nil {
		return err
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

func (b *button) setPosition(x, y int) {
	if b == nil {
		return
	}
	b.sprite.x, b.sprite.y = x, y
	b.text.x = x + (b.sprite.template.frameWidth-b.text.w)/2
	b.text.y = y + (b.sprite.template.frameHeight-b.text.h)/2
}
