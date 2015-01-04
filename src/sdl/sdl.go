// Package sdl wraps the SDL2 framework. hack hack hack hack
package sdl

/*
#include <stdlib.h>
#include <SDL2/SDL.h>

const int kInitEverything = SDL_INIT_EVERYTHING;
const int kInitVideo      = SDL_INIT_VIDEO;

const char* kHintRenderScaleQuality = SDL_HINT_RENDER_SCALE_QUALITY;
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

func Delay(delay uint32) {
	C.SDL_Delay(C.Uint32(delay))
}

// Err returns the last error as a string from an SDL library call.
func Err() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}

type Colour struct {
	R, G, B, A uint8
}

func (c *Colour) c() C.SDL_Color {
	return C.SDL_Color{
		r: C.Uint8(c.R),
		g: C.Uint8(c.G),
		b: C.Uint8(c.B),
		a: C.Uint8(c.A),
	}
}

func (c *Colour) MapRGBA(s *C.SDL_Surface) C.Uint32 {
	return C.SDL_MapRGBA(s.format, C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A))
}

var TransparentColour = Colour{}
var BlackColour = Colour{0x00, 0x00, 0x00, 0xff}
var WhiteColour = Colour{0xff, 0xff, 0xff, 0xff}

type Point struct {
	X, Y int
}

func (p *Point) p() *C.SDL_Point {
	return &C.SDL_Point{x: C.int(p.X), y: C.int(p.Y)}
}

type Rect struct {
	X, Y, W, H int
}

func (r *Rect) r() *C.SDL_Rect {
	if r == nil {
		return nil
	}
	return &C.SDL_Rect{x: C.int(r.X), y: C.int(r.Y), w: C.int(r.W), h: C.int(r.H)}
}

func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x <= r.X+r.W && y >= r.Y && y <= r.Y+r.H
}

// Inset returns a smaller rect inside the outside rect.
func (r Rect) Inset(x, y int) Rect {
	return Rect{X: r.X + x, Y: r.Y + y, W: r.W - 2*x, H: r.H - 2*y}
}

type Context struct {
	*Window
	*Renderer
}

// NewContext creates a Context referring to a new window with a given title and
// size, and a sensible renderer. Don't forget to Close the context when done.
func NewContext(title string, width, height int) (*Context, error) {
	if errno := C.SDL_Init(C.kInitEverything); errno < 0 {
		return nil, fmt.Errorf("unable to SDL_init(everything): %d %s", errno, Err())
	}

	one := C.CString("1")
	defer C.free(unsafe.Pointer(one))
	if errno := C.SDL_SetHint(C.kHintRenderScaleQuality, one); errno == 0 {
		return nil, fmt.Errorf("unable to set hint: %d %s", errno, Err())
	}

	w, err := NewWindow(title, width, height)
	if err != nil {
		return nil, err
	}
	r, err := w.CreateRenderer(RendererAccelerated | RendererPresentVSync)
	if err != nil {
		return nil, err
	}
	ctx := &Context{
		Window:   w,
		Renderer: r,
	}
	r.SetDrawColour(BlackColour)
	r.SetDrawBlendMode(BlendModeBlend)
	return ctx, nil
}

func (c *Context) Close() {
	c.Renderer.Destroy()
	c.Window.Destroy()
	C.SDL_Quit()
}
