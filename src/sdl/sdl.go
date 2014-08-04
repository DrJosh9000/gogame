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

const int kQuit = SDL_QUIT;
const int kKeyDown = SDL_KEYDOWN;
const int kKeyUp = SDL_KEYUP;

Uint32 getType(SDL_Event *ev) {
	return ev->type;
}

SDL_Keycode getKeyCode(SDL_Event *ev) {
	return ev->key.keysym.sym;
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type Context struct {
	window unsafe.Pointer
	renderer unsafe.Pointer
	primarySurface unsafe.Pointer
}

func Delay(delay uint32) {
	C.SDL_Delay(C.Uint32(delay))
}

// Err returns the last error as a string from an SDL library call.
func Err() string {
	return C.GoString(C.SDL_GetError())
}


type KeyEventType int
const (
	KeyDown KeyEventType = iota
	KeyUp
)
type QuitEvent int
type KeyEvent struct {
	Type KeyEventType
	KeyCode uint32
}

func HandleEvents(handler func(e interface{}) error) error {
	var ev C.SDL_Event
	for C.SDL_PollEvent(&ev) != 0 {
		var err error
		switch C.getType(&ev) {
		case C.kQuit:
			err = handler(QuitEvent(0))
		case C.kKeyDown:
			err = handler(KeyEvent{
				Type: KeyDown,
				KeyCode: uint32(C.getKeyCode(&ev)),
			})
		case C.kKeyUp:
			err = handler(KeyEvent{
				Type: KeyUp,
				KeyCode: uint32(C.getKeyCode(&ev)),
			})
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// NewContext creates a Context referring to a new window with a given title and
// size.
func NewContext(title string, width, height int) (*Context, error) {
	
	if errno := C.SDL_Init(C.kInitVideo); errno < 0 {
		return nil, fmt.Errorf("error no from SDL_init(video): %d", errno)
	}
	one := C.CString("1")
	defer C.free(unsafe.Pointer(one))
	if errno := C.SDL_SetHint(C.kHintRenderScaleQuality, one); errno == 0 {
		return nil, fmt.Errorf("unable to set hint: %d %s", Err());
	}
	winTitle := C.CString(title)
	defer C.free(unsafe.Pointer(winTitle))
	w := C.SDL_CreateWindow(winTitle, C.kWindowPosUndefined, C.kWindowPosUndefined, C.int(width), C.int(height), C.kWindowShown)
	if w == nil {
		return nil, fmt.Errorf("unable to create SDL Window: %s", Err());
	}
	s := C.SDL_GetWindowSurface(w)
	r := C.SDL_GetRenderer(w)
	C.SDL_SetRenderDrawColor(r, 0x00, 0x00, 0x00, 0xFF);
	
	ctx := &Context{
		window: unsafe.Pointer(w),
		primarySurface: unsafe.Pointer(s),
		renderer: unsafe.Pointer(r),
	}
	runtime.SetFinalizer(ctx, func(x interface{}) {
		c := x.(*Context)
		if c.renderer != nil {
			C.SDL_DestroyRenderer((*C.SDL_Renderer)(c.renderer))
		}
		if c.window != nil {
			C.SDL_DestroyWindow((*C.SDL_Window)(c.window))
		}
		C.SDL_Quit()
	})
	return ctx, nil
}

func (c *Context) Render() {
	r := (*C.SDL_Renderer)(c.renderer)
	C.SDL_RenderClear(r)
	C.SDL_RenderPresent(r)
}
