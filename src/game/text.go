package game

import (
	"bufio"
	"sdl"
	"strings"
)

const (
	defaultFontFile = "assets/munro.ttf"
	defaultFontSize = 20
)

var defaultFont *sdl.Font

type text struct {
	t          string
	tex        *sdl.Texture
	x, y, w, h int
	invisible  bool
}

func newText(ctx *sdl.Context, s string, c, fill sdl.Colour, al sdl.Alignment) (*text, error) {
	if defaultFont == nil {
		f, err := sdl.LoadFont(defaultFontFile, defaultFontSize)
		if err != nil {
			return nil, err
		}
		defaultFont = f
	}
	var surfs []*sdl.Surface
	b := bufio.NewScanner(strings.NewReader(s))
	for b.Scan() {
		surf, err := defaultFont.RenderSolid(b.Text(), c)
		if err != nil {
			return nil, err
		}
		surfs = append(surfs, surf)
	}
	if err := b.Err(); err != nil {
		return nil, err
	}
	surf, err := sdl.Stack(surfs, fill, al)
	if err != nil {
		return nil, err
	}
	w, h := surf.Size()
	tex, err := ctx.Renderer.TextureFromSurface(surf)
	if err != nil {
		return nil, err
	}
	return &text{
		t:   s,
		tex: tex,
		w:   w,
		h:   h,
	}, nil
}

func (t *text) draw(r renderer) error {
	if t.invisible {
		return nil
	}
	return r.Copy(t.tex,
		sdl.Rect{0, 0, t.w, t.h},
		sdl.Rect{t.x, t.y, t.w, t.h})
}

func (t *text) destroy() {
	if t.tex != nil {
		t.tex.Destroy()
	}
}
