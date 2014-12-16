package game

import (
	"sdl"
)

type drawer interface {
	draw(*sdl.Renderer) error
}

type destroyer interface {
	destroy()
}

type complexObject interface {
	drawer
	addChild(drawer)
	children() []drawer
}

// complexBase is a starting point for implementing complexObject.
type complexBase struct {
	kids      []drawer
	invisible bool
	x, y      int
}

func (b *complexBase) addChild(c drawer) {
	b.kids = append(b.kids, c)
}

func (b *complexBase) children() []drawer {
	return b.kids
}

func (b *complexBase) draw(r *sdl.Renderer) error {
	if b.invisible {
		return nil
	}
	r.PushOffset(b.x, b.y)
	defer r.PopOffset()
	for _, c := range b.kids {
		if c != nil {
			if err := c.draw(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *complexBase) destroy() {
	for _, c := range b.kids {
		if d, ok := c.(destroyer); ok {
			d.destroy()
		}
	}
}

// unionObject is like a complex object, but only one subobject is ever drawn
// (kind of like a C union - only one element is useful at a time).
type unionObject struct {
	complexBase
	active int
}

func (u *unionObject) draw(r *sdl.Renderer) error {
	if u.invisible {
		return nil
	}
	r.PushOffset(u.x, u.y)
	defer r.PopOffset()
	return u.kids[u.active].draw(r)
}
