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
	Text         string
	X, Y, Z      int
	Invisible    bool
	Draw, Shadow sdl.Colour
	Align        sdl.Alignment

	tex  *sdl.Texture
	w, h int
}

// load will re-render the text. Useful for changing the text.
func (t *text) load() error {
	if defaultFont == nil {
		f, err := sdl.LoadFont(defaultFontFile, defaultFontSize)
		if err != nil {
			return err
		}
		defaultFont = f
	}
	var texts, shadows []*sdl.Surface
	b := bufio.NewScanner(strings.NewReader(t.Text))
	for b.Scan() {
		surf, err := defaultFont.RenderSolid(b.Text(), t.Draw)
		if err != nil {
			return err
		}
		texts = append(texts, surf)
		if t.Shadow != sdl.TransparentColour {
			surf, err := defaultFont.RenderSolid(b.Text(), t.Shadow)
			if err != nil {
				return err
			}
			shadows = append(shadows, surf)
		}
	}
	if err := b.Err(); err != nil {
		return err
	}
	w, h := sdl.StackSize(texts)
	if shadows != nil {
		w, h = w+2, h+2
	}
	surf, err := sdl.CreateRGBASurface(w, h)
	if err != nil {
		return err
	}
	if err := surf.SetBlendMode(sdl.BlendModeBlend); err != nil {
		return err
	}
	if err := sdl.Stack(surf, shadows, 2, 2, w, t.Align); err != nil {
		return err
	}
	if err := sdl.Stack(surf, texts, 0, 0, w, t.Align); err != nil {
		return err
	}
	tex, err := ctx().Renderer.TextureFromSurface(surf)
	if err != nil {
		return err
	}
	t.tex, t.w, t.h = tex, w, h
	return nil
}

func (t *text) draw(r *sdl.Renderer) error {
	if t == nil || t.Invisible {
		return nil
	}
	if t.tex == nil {
		if err := t.load(); err != nil {
			return err
		}
	}
	return r.Copy(t.tex,
		sdl.Rect{0, 0, t.w, t.h},
		sdl.Rect{t.X, t.Y, t.w, t.h})
}

func (t *text) invisible() bool {
	return t == nil || t.Invisible
}

func (t *text) z() int {
	return t.Z
}
