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
	tex        *sdl.Texture
	x, y, w, h int
	invisible  bool
}

func newText(s string, draw, shadow sdl.Colour, al sdl.Alignment) (*text, error) {
	if defaultFont == nil {
		f, err := sdl.LoadFont(defaultFontFile, defaultFontSize)
		if err != nil {
			return nil, err
		}
		defaultFont = f
	}
	if draw == sdl.TransparentColour {
		return nil, nil
	}
	var texts, shadows []*sdl.Surface
	b := bufio.NewScanner(strings.NewReader(s))
	for b.Scan() {
		surf, err := defaultFont.RenderSolid(b.Text(), draw)
		if err != nil {
			return nil, err
		}
		texts = append(texts, surf)
		if shadow != sdl.TransparentColour {
			surf, err := defaultFont.RenderSolid(b.Text(), shadow)
			if err != nil {
				return nil, err
			}
			shadows = append(shadows, surf)
		}
	}
	if err := b.Err(); err != nil {
		return nil, err
	}

	w, h := sdl.StackSize(texts)
	if shadows != nil {
		w, h = w+2, h+2
	}
	surf, err := sdl.CreateRGBASurface(w, h)
	if err != nil {
		return nil, err
	}
	if err := surf.SetBlendMode(sdl.BlendModeBlend); err != nil {
		return nil, err
	}
	if err := sdl.Stack(surf, shadows, 2, 2, w, al); err != nil {
		return nil, err
	}
	if err := sdl.Stack(surf, texts, 0, 0, w, al); err != nil {
		return nil, err
	}
	tex, err := ctx().Renderer.TextureFromSurface(surf)
	if err != nil {
		return nil, err
	}
	return &text{tex: tex, w: w, h: h}, nil
}

func (t *text) draw(r *sdl.Renderer) error {
	if t == nil || t.invisible {
		return nil
	}
	return r.Copy(t.tex,
		sdl.Rect{0, 0, t.w, t.h},
		sdl.Rect{t.x, t.y, t.w, t.h})
}
