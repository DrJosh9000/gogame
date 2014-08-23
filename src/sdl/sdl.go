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

func Rect(x, y, w, h int) *C.SDL_Rect {
	return &C.SDL_Rect{C.int(x), C.int(y), C.int(w), C.int(h)}
}

type Context struct {
	*Window
	*Renderer
	*Surface
}

// NewContext creates a Context referring to a new window with a given title and
// size, and a sensible renderer. Don't forget to Close the context when done.
func NewContext(title string, width, height int) (*Context, error) {

	if errno := C.SDL_Init(C.kInitEverything); errno < 0 {
		return nil, fmt.Errorf("error no from SDL_init(everything): %d", errno)
	}
	one := C.CString("1")
	defer C.free(unsafe.Pointer(one))
	if errno := C.SDL_SetHint(C.kHintRenderScaleQuality, one); errno == 0 {
		return nil, fmt.Errorf("unable to set hint: %d %s", Err())
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
		Surface:  w.Surface(),
		Renderer: r,
	}
	ctx.Renderer.SetDrawColor(BlackColor)
	return ctx, nil
}

func (c *Context) Close() {
	c.Renderer.Destroy()
	c.Window.Destroy()
	C.SDL_Quit()
}
