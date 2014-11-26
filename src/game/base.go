package game

import (
	"sdl"
)

type renderer interface {
	Copy(t *sdl.Texture, src, dst sdl.Rect) error
	CopyEx(t *sdl.Texture, src, dst sdl.Rect, angle float64, center sdl.Point, flip sdl.RendererFlip) error
}

type object interface {
	destroy()
	draw(renderer) error
}

type complexObject interface {
	object
	addChild(object)
	children() []object
}

// complexBase is a starting point for implementing complexObject.
type complexBase struct {
	kids []object
}

func (b *complexBase) addChild(c object) {
	b.kids = append(b.kids, c)
}

func (b *complexBase) children() []object {
	return b.kids
}

func (b *complexBase) draw(r renderer) error {
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
		if c != nil {
			c.destroy()
		}
	}
}

// unionObject is like a complex object, but only one subobject is ever drawn
// (kind of like a C union - only one element is useful at a time).
type unionObject struct {
	complexBase
	active int
}

func (u *unionObject) draw(r renderer) error {
	return u.kids[u.active].draw(r)
}
