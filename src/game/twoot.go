package game

import (
	"sdl"
)

type twoot struct {
	Avatar    *sprite
	Bounds    sdl.Rect
	Invisible bool
	Text      text
	Z         int

	frame []sprite
}

func (t *twoot) bounds() sdl.Rect {
	return t.Bounds
}

func (t *twoot) draw(r *sdl.Renderer) error {
	if t == nil || t.Invisible {
		return nil
	}
	if err := t.load(); err != nil {
		return err
	}
	r.PushOffset(t.Bounds.X, t.Bounds.Y)
	defer r.PopOffset()
	// Place the blue frame around it.
	for _, s := range t.frame {
		if err := s.draw(r); err != nil {
			return err
		}
	}
	if err := t.Avatar.draw(r); err != nil {
		return err
	}
	if err := t.Text.draw(r); err != nil {
		return err
	}
	return nil
}

func (t *twoot) load() error {
	if t.frame != nil {
		return nil
	}
	t.frame = make([]sprite, 9)
	for i := 0; i < 9; i++ {
		t.frame[i] = sprite{
			TemplateKey: "twootFrame",
			Frame:       i,
		}
		if err := t.frame[i].load(); err != nil {
			return err
		}
	}
	ft := t.frame[0].template
	t.frame[1].X = ft.frameWidth
	t.frame[1].w = t.Bounds.W - 2*ft.frameWidth
	t.frame[2].X = t.Bounds.W - ft.frameWidth
	t.frame[3].Y = ft.frameHeight
	t.frame[3].h = t.Bounds.H - 2*ft.frameHeight
	t.frame[4].X = t.frame[1].X
	t.frame[4].Y = t.frame[3].Y
	t.frame[4].w = t.frame[1].w
	t.frame[4].h = t.frame[3].h
	t.frame[5].X = t.frame[2].X
	t.frame[5].Y = ft.frameHeight
	t.frame[5].h = t.frame[3].h
	t.frame[6].Y = t.Bounds.H - ft.frameHeight
	t.frame[7].X = t.frame[1].X
	t.frame[7].Y = t.frame[6].Y
	t.frame[7].w = t.frame[1].w
	t.frame[8].X = t.frame[2].X
	t.frame[8].Y = t.frame[6].Y

	t.Avatar.X, t.Avatar.Y = 32, 32
	t.Text.X, t.Text.Y = 112, 32
	return nil
}

func (t *twoot) invisible() bool {
	return t.Invisible
}

func (t *twoot) z() int {
	return t.Z
}
