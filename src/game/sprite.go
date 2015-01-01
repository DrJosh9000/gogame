package game

import (
	"sdl"
	"sync"
)

type spriteTemplate struct {
	baseX, baseY            int // point within the sprite frame of the "centre"
	framesX, framesY        int // number of frames along X, Y axes
	frameWidth, frameHeight int // size of a frame
	sheetFile               string

	tex    *sdl.Texture
	loadMu sync.Once
}

func (t *spriteTemplate) load() (err error) {
	t.loadMu.Do(func() {
		t.tex, err = ctx().LoadImage(t.sheetFile)
	})
	return
}

type sprite struct {
	Invisible      bool
	TemplateKey    string
	X, Y, Z, Frame int

	template *spriteTemplate
}

func (s *sprite) load() error {
	// Ensure the template & texture are loaded.
	if s.template == nil {
		s.template = templateLibrary[s.TemplateKey]
	}
	return s.template.load()
}

func (s *sprite) draw(r *sdl.Renderer) error {
	if s == nil || s.Invisible {
		return nil
	}
	if err := s.load(); err != nil {
		return err
	}

	// Compute the frame bounds and draw.
	srcX := (s.Frame % s.template.framesX) * s.template.frameWidth
	srcY := ((s.Frame / s.template.framesX) % s.template.framesY) * s.template.frameHeight
	return r.Copy(s.template.tex,
		sdl.Rect{srcX, srcY, s.template.frameWidth, s.template.frameHeight},
		sdl.Rect{s.X - s.template.baseX, s.Y + s.Z - s.template.baseY, s.template.frameWidth, s.template.frameHeight})
}

func (s *sprite) z() int {
	return s.Z
}
