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
	// Z returns the Z order index. Lesser numbers are drawn before greater.
	Z() int
}

type object interface {
	drawer
	zer
}

type objectSlice []object

// The following implement sort.Interface.

func (o objectSlice) Len() int           { return len(o) }
func (o objectSlice) Less(i, j int) bool { return o[i].Z() < o[j].Z() }
func (o objectSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

type complexObject interface {
	object
	addChild(object)
	children() objectSlice
}

// complexBase is a starting point for implementing complexObject.
type complexBase struct {
	kids      []object
	invisible bool
	x, y, z   int
}

func (b *complexBase) addChild(c object) {
	b.kids = append(b.kids, c)
}

func (b *complexBase) children() objectSlice {
	return objectSlice(b.kids)
}

func (b *complexBase) draw(r *sdl.Renderer) error {
	if b.invisible {
		return nil
	}
	r.PushOffset(b.x, b.y)
	defer r.PopOffset()
	// Do not rely on the z order of children remaining static...
	sort.Sort(objectSlice(b.kids))
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

func (b *complexBase) Z() int {
	return b.z
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
