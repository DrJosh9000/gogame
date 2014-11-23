// Package sdl wraps the SDL2 framework. hack hack hack hack
package sdl

/*
#include <stdlib.h>
#include <SDL2/SDL.h>

const int kInitEverything = SDL_INIT_EVERYTHING;
const int kInitVideo = SDL_INIT_VIDEO;
const char* kHintRenderScaleQuality = SDL_HINT_RENDER_SCALE_QUALITY;
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func Delay(delay uint32) {
	C.SDL_Delay(C.Uint32(delay))
}

// Err returns the last error as a string from an SDL library call.
func Err() string {
	return C.GoString(C.SDL_GetError())
}

type Color struct {
	red, green, blue, alpha uint8
}

var BlackColor = Color{0x00, 0x00, 0x00, 0xff}

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
	return &C.SDL_Rect{x: C.int(r.X), y: C.int(r.Y), w: C.int(r.W), h: C.int(r.H)}
}

type Context struct {
	*Window
	*Renderer
	*TextureManager
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
		Window:         w,
		Renderer:       r,
		TextureManager: NewTextureManager(r),
	}
	ctx.Renderer.SetDrawColor(BlackColor)
	return ctx, nil
}

func (c *Context) Close() {
	c.TextureManager.Destroy()
	c.Renderer.Destroy()
	c.Window.Destroy()
	C.SDL_Quit()
}
