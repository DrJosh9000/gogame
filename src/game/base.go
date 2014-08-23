package game

import (
	"time"

	"sdl"
)

type Object interface {
	AddChild(Object)
	Children() []Object
	Destroy()
	Draw(*sdl.Renderer) error
	Update(t time.Duration)
}

type Base struct {
	children []Object
}

func (b *Base) AddChild(c Object) {
	b.children = append(b.children, c)
}

func (b *Base) Children() []Object {
	return b.children
}

func (b *Base) Draw(r *sdl.Renderer) error {
	for _, c := range b.children {
		if c != nil {
			if err := c.Draw(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Base) Update(t time.Duration) {
	for _, c := range b.children {
		if c != nil {
			c.Update(t)
		}
	}
}

func (b *Base) Destroy() {
	for _, c := range b.children {
		if c != nil {
			c.Destroy()
		}
	}
}
