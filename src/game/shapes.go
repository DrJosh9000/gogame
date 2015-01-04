package game

import (
	"encoding/gob"
	"sdl"
)

const sqrt3 = 1.73205080757

func init() {
	gob.Register(&rect{})
	gob.Register(&circle{})
	gob.Register(&ellipse{})
}

type rect struct {
	sdl.Rect
	Colour sdl.Colour
	Z      int
}

func (r *rect) bounds() sdl.Rect {
	if r == nil {
		return sdl.Rect{}
	}
	return r.Rect
}

func (r *rect) draw(ren *sdl.Renderer) error {
	ren.SetDrawColour(r.Colour)
	// Make it 2 thick by drawing a second inset rect.
	if err := ren.DrawRect(r.Rect.Inset(1, 1)); err != nil {
		return err
	}
	return ren.DrawRect(r.Rect)
}

func (r *rect) z() int {
	return r.Z
}

type circle struct {
	X, Y, R int // centre, radius
	Colour  sdl.Colour
	Z       int
}

func (e *circle) bounds() sdl.Rect {
	if e == nil {
		return sdl.Rect{}
	}
	return sdl.Rect{X: e.X - e.R, Y: e.Y - e.R, W: 2 * e.R, H: 2 * e.R}
}

func (e *circle) draw(r *sdl.Renderer) error {
	r.SetDrawColour(e.Colour)
	rx := float64(e.R / 2)
	x, y := e.R/2, 0
	for x >= y {
		if err := r.DrawFatPoint(e.X+x*2, e.Y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-x*2, e.Y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X+x*2, e.Y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-x*2, e.Y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X+y*2, e.Y+x*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-y*2, e.Y+x*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X+y*2, e.Y-x*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-y*2, e.Y-x*2, 2); err != nil {
			return err
		}
		y++
		rx -= float64(y) / rx
		x = int(rx + 0.5)
	}
	return nil
}

func (e *circle) z() int {
	return e.Z
}

type ellipse struct {
	X, Y, W, H int
	Colour     sdl.Colour
	Z          int
}

func (e *ellipse) bounds() sdl.Rect {
	if e == nil {
		return sdl.Rect{}
	}
	return sdl.Rect{X: e.X - e.W, Y: e.Y - e.H, W: 2 * e.W, H: 2 * e.H}
}

func (e *ellipse) draw(r *sdl.Renderer) error {
	r.SetDrawColour(e.Colour)
	rx := float64(e.W / 2)
	x, y, step := e.W/2, 0, 0.0
	fr := float64(e.W*e.W) / float64(e.H*e.H)
	for step <= 1.0 {
		if err := r.DrawFatPoint(e.X+x*2, e.Y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-x*2, e.Y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X+x*2, e.Y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-x*2, e.Y-y*2, 2); err != nil {
			return err
		}
		y++
		step = float64(y) * fr / rx
		rx -= step
		x = int(rx + 0.5)
	}
	ry := float64(e.H / 2)
	x, y, step = 0, e.H/2, 0.0
	for step <= 1.0 {
		if err := r.DrawFatPoint(e.X+x*2, e.Y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-x*2, e.Y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X+x*2, e.Y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.X-x*2, e.Y-y*2, 2); err != nil {
			return err
		}
		x++
		step = float64(x) / (ry * fr)
		ry -= step
		y = int(ry + 0.5)
	}
	return nil
}

func (e *ellipse) z() int {
	return e.Z
}
