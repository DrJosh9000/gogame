package game

import (
	"time"
)

type Object interface {
	AddChild(Object)
	Children() []Object
	Destroy()
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