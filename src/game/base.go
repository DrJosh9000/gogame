package game

import (
	"sdl"
	"sort"
)

type drawer interface {
	draw(*sdl.Renderer) error
}

type destroyer interface {
	destroy()
}

type zer interface {
	// z returns the Z order index. Lesser numbers are drawn before greater.
	z() int
}

type object interface {
	drawer
	zer
}

type objectSlice []object

// The following implement sort.Interface.

func (o objectSlice) Len() int           { return len(o) }
func (o objectSlice) Less(i, j int) bool { return o[i].z() < o[j].z() }
func (o objectSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

// complexBase is a starting point for implementing complexObject.
type complexBase struct {
	Kids      []object
	Invisible bool
	X, Y, Z   int
}

func (b *complexBase) addChild(c object) {
	b.Kids = append(b.Kids, c)
}

func (b *complexBase) children() objectSlice {
	return objectSlice(b.Kids)
}

func (b *complexBase) draw(r *sdl.Renderer) error {
	if b.Invisible {
		return nil
	}
	r.PushOffset(b.X, b.Y)
	defer r.PopOffset()
	// Do not rely on the z order of children remaining static...
	sort.Sort(objectSlice(b.Kids))
	for _, c := range b.Kids {
		if c != nil {
			if err := c.draw(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *complexBase) destroy() {
	for _, c := range b.Kids {
		if d, ok := c.(destroyer); ok {
			d.destroy()
		}
	}
}

func (b *complexBase) z() int {
	return b.Z
}

// unionObject is like a complex object, but only one subobject is ever drawn
// (kind of like a C union - only one element is useful at a time).
type unionObject struct {
	complexBase
	Active int
}

func (u *unionObject) draw(r *sdl.Renderer) error {
	if u.Invisible {
		return nil
	}
	r.PushOffset(u.X, u.Y)
	defer r.PopOffset()
	return u.Kids[u.Active].draw(r)
}
