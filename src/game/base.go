package game

import (
	"sdl"
	"sort"
)

type clicker interface {
	object
	click()
	setDown(bool)
}

type destroyer interface {
	destroy()
}

type loader interface {
	load() error
}

type object interface {
	bounds() sdl.Rect
	draw(*sdl.Renderer) error
	invisible() bool
	z() int
}

// objectSlices are z-orderable.
type objectSlice []object

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

func (b *complexBase) invisible() bool {
	return b == nil || b.Invisible
}

func (b *complexBase) load() error {
	for _, c := range b.Kids {
		if l, ok := c.(loader); ok {
			if err := l.load(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *complexBase) z() int {
	return b.Z
}

// unionObject is like a complex object, but only one subobject is ever drawn
// (kind of like a C union - only one element is useful at a time).
type unionObject struct {
	Active int
	Kids   []object
}

func (u *unionObject) addChild(c object) {
	u.Kids = append(u.Kids, c)
}

func (u *unionObject) children() objectSlice {
	return objectSlice(u.Kids)
}

func (u *unionObject) bounds() sdl.Rect {
	return u.Kids[u.Active].bounds()
}

func (u *unionObject) draw(r *sdl.Renderer) error {
	return u.Kids[u.Active].draw(r)
}

func (u *unionObject) invisible() bool {
	return u.Kids[u.Active].invisible()
}

func (u *unionObject) z() int {
	return u.Kids[u.Active].z()
}
