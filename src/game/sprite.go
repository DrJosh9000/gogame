package game

import (
	"fmt"
	"sdl"
)

type spriteTemplate struct {
	name, sheetFile                           string
	framesX, framesY, frameWidth, frameHeight int
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

func (s *spriteTemplate) String() string {
	return s.name
}

type sprite struct {
	template    *spriteTemplate
	x, y, frame int
	tex         *sdl.Texture
	invisible   bool
}

func (s *sprite) destroy() {}

func (s *sprite) draw(r renderer) error {
	if s.invisible {
		return nil
	}
	srcX := (s.frame % s.template.framesX) * s.template.frameWidth
	srcY := ((s.frame / s.template.framesX) % s.template.framesY) * s.template.frameHeight
	return r.Copy(s.tex,
		sdl.Rect{srcX, srcY, s.template.frameWidth, s.template.frameHeight},
		sdl.Rect{s.x, s.y, s.template.frameWidth, s.template.frameHeight})
}

func (s *sprite) String() string {
	return fmt.Sprintf("sprite(%s)", s.template)
}
