package game

import (
	"sdl"
)

var (
	orbTemplate = &spriteTemplate{
		sheetFile:   "assets/orb.png",
		baseX:       15,
		baseY:       31,
		framesX:     1,
		framesY:     1,
		frameWidth:  32,
		frameHeight: 32,
	}

	orbShadowTemplate = &spriteTemplate{
		sheetFile:   "assets/orbshadow.png",
		baseX:       0,
		baseY:       12,
		framesX:     1,
		framesY:     1,
		frameWidth:  32,
		frameHeight: 32,
	}
)

type orb struct {
	X, Y, Z             int
	orb, shadow         *sprite
	selection           *ellipse
	Selected, Invisible bool
}

func newOrb(ctx *sdl.Context) (*orb, error) {
	o, err := orbTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	shad, err := orbShadowTemplate.new(ctx)
	if err != nil {
		return nil, err
	}
	sel := &ellipse{
		w:      16,
		h:      9,
		colour: sdl.Colour{R: 0x00, G: 0xCC, B: 0x11, A: 0xFF},
	}
	return &orb{
		orb:       o,
		shadow:    shad,
		selection: sel,
	}, nil
}

func (o *orb) draw(r *sdl.Renderer) error {
	if o.Invisible {
		return nil
	}
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
