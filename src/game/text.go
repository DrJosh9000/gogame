package game

import (
	"bufio"
	"sdl"
	"sort"
	"strings"
)

const (
	defaultFontFile = "assets/munro.ttf"
	defaultFontSize = 20
)

var defaultFont *sdl.Font

type text struct {
	Text         string
	X, Y, Z      int // Top left corner, Z order.
	Invisible    bool
	Draw, Shadow sdl.Colour
	Align        sdl.Alignment
	MaxWidth     int // At which wrapping occurs. 0 == no limit

	tex  *sdl.Texture
	w, h int
}

// load will render the text. Useful for changing the text.
func (t *text) load() error {
	if defaultFont == nil {
		f, err := sdl.LoadFont(defaultFontFile, defaultFontSize)
		if err != nil {
			return err
		}
		defaultFont = f
	}
	var texts, shadows []*sdl.Surface
	addLine := func(l string) error {
		surf, err := defaultFont.RenderSolid(l, t.Draw)
		if err != nil {
			return err
		}
		texts = append(texts, surf)
		if t.Shadow != sdl.TransparentColour {
			surf, err := defaultFont.RenderSolid(l, t.Shadow)
			if err != nil {
				return err
			}
			shadows = append(shadows, surf)
		}
		return nil
	}

	b := bufio.NewScanner(strings.NewReader(t.Text))
	for b.Scan() {
		if t.MaxWidth <= 0 {
			if err := addLine(b.Text()); err != nil {
				return err
			}
			continue
		}
		words := strings.Split(b.Text(), " ")
		for len(words) > 0 {
			wc := sort.Search(len(words), func(i int) bool {
				w, _, err := defaultFont.SizeText(strings.Join(words[:i+1], " "))
				if err != nil {
					return false
				}
				return w > t.MaxWidth
			})
			if err := addLine(strings.Join(words[:wc], " ")); err != nil {
				return err
			}
			words = words[wc:]
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
