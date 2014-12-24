package game

import (
	"sdl"
)

type rect struct {
	bounds sdl.Rect
	colour sdl.Colour
	z      int
}

func (r *rect) draw(ren *sdl.Renderer) error {
	ren.SetDrawColour(r.colour)
	m1 := r.bounds
	m1.X, m1.Y = m1.X+1, m1.Y+1
	m1.W, m1.H = m1.W-2, m1.H-2
	if err := ren.DrawRect(m1); err != nil {
		return err
	}
	return ren.DrawRect(r.bounds)
}

func (r *rect) Z() int {
	return r.z
}

type circle struct {
	x, y, r int // centre, radius
	colour  sdl.Colour
	z       int
}

func (e *circle) draw(r *sdl.Renderer) error {
	r.SetDrawColour(e.colour)
	rx := float64(e.r / 2)
	x, y := e.r/2, 0
	for x >= y {
		if err := r.DrawFatPoint(e.x+x*2, e.y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-x*2, e.y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x+x*2, e.y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-x*2, e.y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x+y*2, e.y+x*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-y*2, e.y+x*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x+y*2, e.y-x*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-y*2, e.y-x*2, 2); err != nil {
			return err
		}
		y++
		rx -= float64(y) / rx
		x = int(rx + 0.5)
	}
	return nil
}

func (e *circle) Z() int {
	return e.z
}

type ellipse struct {
	x, y, w, h int
	colour     sdl.Colour
	z          int
}

func (e *ellipse) draw(r *sdl.Renderer) error {
	r.SetDrawColour(e.colour)
	rx := float64(e.w / 2)
	x, y, step := e.w/2, 0, 0.0
	fr := float64(e.w*e.w) / float64(e.h*e.h)
	for step <= 1.0 {
		if err := r.DrawFatPoint(e.x+x*2, e.y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-x*2, e.y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x+x*2, e.y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-x*2, e.y-y*2, 2); err != nil {
			return err
		}
		y++
		step = float64(y) * fr / rx
		rx -= step
		x = int(rx + 0.5)
	}
	ry := float64(e.h / 2)
	x, y, step = 0, e.h/2, 0.0
	for step <= 1.0 {
		if err := r.DrawFatPoint(e.x+x*2, e.y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-x*2, e.y+y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x+x*2, e.y-y*2, 2); err != nil {
			return err
		}
		if err := r.DrawFatPoint(e.x-x*2, e.y-y*2, 2); err != nil {
			return err
		}
		x++
		step = float64(x) / ry / fr
		ry -= step
		y = int(ry + 0.5)
	}
	return nil
}

func (e *ellipse) Z() int {
	return e.z
}
