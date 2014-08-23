package sdl

/*
#include <stdlib.h>
#include <SDL2/SDL.h>

const int kWindowPosUndefined = SDL_WINDOWPOS_UNDEFINED;
const int kWindowShown = SDL_WINDOW_SHOWN;
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Window struct {
	window unsafe.Pointer
}

func NewWindow(title string, width, height int) (*Window, error) {
	winTitle := C.CString(title)
	defer C.free(unsafe.Pointer(winTitle))
	w := C.SDL_CreateWindow(winTitle, C.kWindowPosUndefined, C.kWindowPosUndefined, C.int(width), C.int(height), C.kWindowShown)
	if w == nil {
		return nil, fmt.Errorf("unable to create SDL Window: %s", Err())
	}
	return &Window{unsafe.Pointer(w)}, nil
}

func (w *Window) w() *C.SDL_Window {
	return (*C.SDL_Window)(w.window)
}

func (w *Window) CreateRenderer(opts RendererOption) (*Renderer, error) {
	r := C.SDL_CreateRenderer(w.w(), -1, C.Uint32(opts))
	if r == nil {
		return nil, fmt.Errorf("unable to create SDL Renderer with opts %d: %s", opts, Err())
	}
	return &Renderer{unsafe.Pointer(r)}, nil
}

func (w *Window) Renderer() *Renderer {
	r := C.SDL_GetRenderer(w.w())
	return &Renderer{unsafe.Pointer(r)}
}

func (w *Window) Surface() *Surface {
	s := C.SDL_GetWindowSurface(w.w())
	return &Surface{unsafe.Pointer(s)}
}

func (w *Window) Destroy() {
	if w.w() != nil {
		C.SDL_DestroyWindow(w.w())
	}
	w.window = nil
}
