package game

import (
	"fmt"

	"sdl"
)

type Object interface {
	Destroy()
	Draw(*sdl.Renderer) error
	String() string
}

type ComplexObject interface {
	Object
	AddChild(Object)
	Children() []Object
}

type ComplexBase struct {
	children []Object
}

func (b *ComplexBase) AddChild(c Object) {
	b.children = append(b.children, c)
}

func (b *ComplexBase) Children() []Object {
	return b.children
}

func (b *ComplexBase) Draw(r *sdl.Renderer) error {
	for _, c := range b.children {
		if c != nil {
			if err := c.Draw(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *ComplexBase) Destroy() {
	for _, c := range b.children {
		if c != nil {
			c.Destroy()
		}
	}
}

func (b *ComplexBase) String() string {
	return fmt.Sprintf("%T", b)
}
