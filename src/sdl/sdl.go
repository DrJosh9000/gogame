// Package sdl wraps the SDL2 framework. hack hack hack hack
package sdl

/*
#cgo CFLAGS: -I/Library/Frameworks/SDL2.framework/Headers
#cgo LDFLAGS: -F/Library/Frameworks -framework SDL2

#include <stdlib.h>
#include <SDL.h>

const int kInitEverything = SDL_INIT_EVERYTHING;
const int kInitVideo = SDL_INIT_VIDEO;
const char* kHintRenderScaleQuality = SDL_HINT_RENDER_SCALE_QUALITY;
const int kRendererAccelerated = SDL_RENDERER_ACCELERATED;
const int kWindowPosUndefined = SDL_WINDOWPOS_UNDEFINED;
const int kWindowShown = SDL_WINDOW_SHOWN;
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

type Renderer struct {
	renderer unsafe.Pointer
}

func (r *Renderer) r() *C.SDL_Renderer {
	return (*C.SDL_Renderer)(r.renderer)
}

func (r *Renderer) Render() {
	C.SDL_RenderClear(r.r())
	C.SDL_RenderPresent(r.r())
}

func (r *Renderer) Destroy() {
	if r.r() != nil {
		C.SDL_DestroyRenderer(r.r())
	}
}

func (r *Renderer) SetDrawColor(red, green, blue, alpha byte) {
	C.SDL_SetRenderDrawColor(r.r(), C.Uint8(red), C.Uint8(green), C.Uint8(blue), C.Uint8(alpha));
}

type Surface struct {
	surface unsafe.Pointer
}

func (s *Surface) s() *C.SDL_Surface {
	return (*C.SDL_Surface)(s.surface)
}

type Window struct {
	window unsafe.Pointer
}

func NewWindow(title string, width, height int) (*Window, error) {
	winTitle := C.CString(title)
	defer C.free(unsafe.Pointer(winTitle))
	w := C.SDL_CreateWindow(winTitle, C.kWindowPosUndefined, C.kWindowPosUndefined, C.int(width), C.int(height), C.kWindowShown)
	if w == nil {
		return nil, fmt.Errorf("unable to create SDL Window: %s", Err());
	}
	return &Window{ unsafe.Pointer(w) }, nil
}

func (w *Window) w() *C.SDL_Window {
	return (*C.SDL_Window)(w.window)
}

func (w *Window) Renderer() *Renderer {
	r := C.SDL_GetRenderer(w.w())
	return &Renderer{ unsafe.Pointer(r) }
}

func (w *Window) Surface() *Surface {
	s := C.SDL_GetWindowSurface(w.w())
	return &Surface{ unsafe.Pointer(s) }
}

func (w *Window) Destroy() {
	if w.w() != nil {
		C.SDL_DestroyWindow(w.w())
	}
}

type Context struct {
	*Window
	*Renderer
	*Surface
}

// NewContext creates a Context referring to a new window with a given title and
// size. Don't forget to Close the context when done.
func NewContext(title string, width, height int) (*Context, error) {
	
	if errno := C.SDL_Init(C.kInitVideo); errno < 0 {
		return nil, fmt.Errorf("error no from SDL_init(video): %d", errno)
	}
	one := C.CString("1")
	defer C.free(unsafe.Pointer(one))
	if errno := C.SDL_SetHint(C.kHintRenderScaleQuality, one); errno == 0 {
		return nil, fmt.Errorf("unable to set hint: %d %s", Err());
	}

	w, err := NewWindow(title, width, height)
	if err != nil {
		return nil, err
	}
	ctx := &Context{
		Window: w,
		Surface: w.Surface(),
		Renderer: w.Renderer(),
	}
	ctx.Renderer.SetDrawColor(0x00, 0x00, 0x00, 0xff)
	return ctx, nil
}

func (c *Context) Close() {
	c.Renderer.Destroy()
	c.Window.Destroy()
	C.SDL_Quit()
}

