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
	String() string
}

type complexObject interface {
	object
	addChild(object)
	children() []object
}

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

func (b *complexBase) String() string {
	return "complexBase"
}
