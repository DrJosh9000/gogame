package main

/*
#cgo CFLAGS: -I/Library/Frameworks/SDL2.framework/Headers
#cgo LDFLAGS: -framework SDL2

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
	"unsafe"
)

const (
	defaultWidth, defaultHeight = 1024, 768	
	gameName = "My SDL Game"
)

var (
	window *C.SDL_Window
	renderer *C.SDL_Renderer
	primarySurface *C.SDL_Surface
	running bool
)

func sdlErr() string {
	return C.GoString(C.SDL_GetError())
}

func setup() error {
	if errno := C.SDL_Init(C.kInitVideo); errno < 0 {
		return fmt.Errorf("error no from SDL_init(video): %d", errno)
	}
	
	one := C.CString("1")
	defer C.free(unsafe.Pointer(one))
	if errno := C.SDL_SetHint(C.kHintRenderScaleQuality, one); errno == 0 {
		return fmt.Errorf("unable to Init hinting: %d %s", sdlErr());
	}
	
	winTitle := C.CString(gameName)
	defer C.free(unsafe.Pointer(winTitle))
	if window = C.SDL_CreateWindow(winTitle, C.kWindowPosUndefined, C.kWindowPosUndefined, defaultWidth, defaultHeight, C.kWindowShown); window == nil {
		return fmt.Errorf("unable to create SDL Window: %s", sdlErr());
	}
	
	primarySurface = C.SDL_GetWindowSurface(window);
	
	renderer = C.SDL_GetRenderer(window)
	C.SDL_SetRenderDrawColor(renderer, 0x00, 0x00, 0x00, 0xFF);
	return nil
}

func render() {
	C.SDL_RenderClear(renderer)
	C.SDL_RenderPresent(renderer)
}

func cleanup() {
	if renderer != nil {
		C.SDL_DestroyRenderer(renderer)
		renderer = nil
	}
	if window != nil {
		C.SDL_DestroyWindow(window)
		window = nil
	}
	C.SDL_Quit()
}

func main() {	
	if err := setup(); err != nil {
		panic(err)
	}
	defer cleanup()

	for {
		// Poll events
		var ev C.SDL_Event
		for C.SDL_PollEvent(&ev) != 0 {
			switch C.getType(&ev) {
			case C.kQuit:
				return
			case C.kKeyUp:
				if C.getKeyCode(&ev) == C.SDLK_q {
					return
				}
			}
		}
		
		render()
		C.SDL_Delay(1)
	}
}