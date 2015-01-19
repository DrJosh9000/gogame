package game

import (
	"log"
	"sdl"
)

type buttonState int

const (
	buttonUp buttonState = iota
	buttonDown
)

// button is a clickable region on screen with optional text.
type button struct {
	W, X, Y, Z int
	Label      string
	Action     func() error
	State      buttonState
	Invisible  bool

	img    []*sprite
	text   text
	parent *menu
}

func (b *button) bounds() sdl.Rect {
	return sdl.Rect{b.X, b.Y, b.W, 64}
}

func (b *button) load() error {
	if b.img != nil {
		return nil
	}

	// Button image setup.
	for i := 0; i < 3; i++ {
		s := &sprite{TemplateKey: "button"}
		if err := s.load(); err != nil {
			return err
		}
		b.img = append(b.img, s)
	}

	// Button image layout.
	w := b.img[0].template.frameWidth
	b.img[1].X = w
	b.img[1].w = b.W - 2*w
	b.img[2].X = b.W - w

	// Label setup.
	b.text.Text = b.Label
	b.text.Draw = sdl.BlackColour
	b.text.Align = sdl.CentreAlign
	if err := b.text.load(); err != nil {
		return err
	}
	return nil
}

func (b *button) draw(r *sdl.Renderer) error {
	if b == nil || b.Invisible {
		return nil
	}
	if err := b.load(); err != nil {
		return err
	}
	r.PushOffset(b.X, b.Y)
	defer r.PopOffset()
	for i, s := range b.img {
		s.Frame = (int(b.State) * 3) + i
		if err := s.draw(r); err != nil {
			return err
		}
	}
	if b.State == buttonDown {
		r.PushOffset(0, 8)
		defer r.PopOffset()
	}
	b.text.X = (b.W - b.text.w) / 2
	b.text.Y = (64-b.text.h)/2 - 4
	if err := b.text.draw(r); err != nil {
		return err
	}
	return nil
}

func (b *button) invisible() bool {
	return b == nil || b.parent == nil || b.parent.invisible() || b.Invisible
}

func (b *button) setDown(down bool) {
	b.State = buttonUp
	if down {
		b.State = buttonDown
	}
}

func (b *button) click() {
	if err := b.Action(); err != nil {
		log.Print(err)
	}
}

func (b *button) z() int {
	return b.Z
}
