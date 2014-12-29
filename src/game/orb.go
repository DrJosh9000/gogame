package game

import (
	"math"
	"sdl"
	"time"
)

type orb struct {
	X, Y, Z             int
	Selected, Invisible bool
	py                  int
	orb, shadow         *sprite
	selection           *ellipse
}

func (o *orb) load() {
	if o.orb == nil {
		o.orb = &sprite{TemplateKey: "orb"}
		o.shadow = &sprite{TemplateKey: "orbShadow"}
		o.selection = &ellipse{
			w:      20,
			h:      12,
			colour: sdl.Colour{R: 0x00, G: 0xAA, B: 0xEE, A: 0xFF},
		}
		go o.life()
	}
}

func (o *orb) draw(r *sdl.Renderer) error {
	if o == nil || o.Invisible {
		return nil
	}
	o.load()
	r.PushOffset(o.X, o.Z)
	defer r.PopOffset()
	y := o.Y + o.py
	r.PushOffset(int(2*float64(y)/sqrt3), y/2)
	if err := o.shadow.draw(r); err != nil {
		return err
	}
	r.PopOffset()
	if o.Selected {
		if err := o.selection.draw(r); err != nil {
			return err
		}
	}
	r.PushOffset(0, -y)
	defer r.PopOffset()
	return o.orb.draw(r)
}

func (o *orb) life() {
	inbox := make(chan message, 10)
	kmp("clock", inbox)
	for m := range inbox {
		if d, ok := m.v.(time.Duration); ok {
			o.py = int(3.0 * math.Sin(3.0*d.Seconds()))
		}
	}
}

func (o *orb) z() int {
	return o.Z
}
