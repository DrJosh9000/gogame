package game

import (
	"time"
	
	"sdl"
)

type Object interface {
	AddChild(Object)
	Children() []Object
	Destroy()
	Draw(*sdl.Renderer)
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

func (b *Base) Draw(r *sdl.Renderer) {
	for _, c := range b.children {
		if c != nil {
			c.Draw(r)
		}
	}
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
