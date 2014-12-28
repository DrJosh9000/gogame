package game

import (
	"sdl"
)

type orb struct {
	X, Y, Z             int
	orb, shadow         *sprite
	selection           *ellipse
	Selected, Invisible bool
}

func (o *orb) load() {
	if o.orb == nil {
		o.orb = &sprite{TemplateKey: "orb"}
	}
	if o.shadow == nil {
		o.shadow = &sprite{TemplateKey: "orbShadow"}
	}
	if o.selection == nil {
		o.selection = &ellipse{
			w:      16,
			h:      9,
			colour: sdl.Colour{R: 0x00, G: 0xCC, B: 0x11, A: 0xFF},
		}
	}
}

func (o *orb) draw(r *sdl.Renderer) error {
	if o == nil || o.Invisible {
		return nil
	}
	o.load()
	r.PushOffset(o.X, o.Y+o.Z)
	defer r.PopOffset()
	if err := o.shadow.draw(r); err != nil {
		return err
	}
	if o.Selected {
		if err := o.selection.draw(r); err != nil {
			return err
		}
	}
	return o.orb.draw(r)
}

func (o *orb) z() int {
	return o.Z
}
