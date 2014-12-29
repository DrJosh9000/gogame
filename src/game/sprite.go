package game

import (
	"sdl"
)

type spriteTemplate struct {
	baseX, baseY            int // point within the sprite frame of the "centre"
	framesX, framesY        int // number of frames along X, Y axes
	frameWidth, frameHeight int // size of a frame
	sheetFile               string
}

func (s *spriteTemplate) new(ctx *sdl.Context) (*sprite, error) {
	tex, err := ctx.GetTexture(s.sheetFile)
	if err != nil {
		return nil, err
	}
	return &sprite{
		template: s,
		tex:      tex,
	}, nil
}

type sprite struct {
	Invisible      bool
	TemplateKey    string
	X, Y, Z, Frame int

	template *spriteTemplate
	tex      *sdl.Texture
}

func (s *sprite) load() {
	// Ensure the template & texture are loaded.
	if s.template == nil {
		s.template = templateLibrary[s.TemplateKey]
	}
}

func (s *sprite) draw(r *sdl.Renderer) error {
	if s.Invisible {
		return nil
	}
	s.load()
	if s.tex == nil {
		tex, err := ctx().GetTexture(s.template.sheetFile)
		if err != nil {
			return err
		}
		s.tex = tex
	}

	// Compute the frame bounds and draw.
	srcX := (s.Frame % s.template.framesX) * s.template.frameWidth
	srcY := ((s.Frame / s.template.framesX) % s.template.framesY) * s.template.frameHeight
	return r.Copy(s.tex,
		sdl.Rect{srcX, srcY, s.template.frameWidth, s.template.frameHeight},
		sdl.Rect{s.X - s.template.baseX, s.Y + s.Z - s.template.baseY, s.template.frameWidth, s.template.frameHeight})
}

func (s *sprite) z() int {
	return s.Z
}
