package sdl

/*
#cgo LDFLAGS: -framework SDL2

#include <SDL2/SDL.h>

const int kEnable  = SDL_ENABLE;
const int kDisable = SDL_DISABLE;
const int kQuery   = SDL_QUERY;
*/
import "C"

func ShowCursor() error {
	if errno := C.SDL_ShowCursor(C.kEnable); errno != 0 {
		return Err()
	}
	return nil
}

func HideCursor() error {
	if errno := C.SDL_ShowCursor(C.kDisable); errno != 0 {
		return Err()
	}
	return nil
}

func IsCursorShown() bool {
	return C.SDL_ShowCursor(C.kQuery) == C.kEnable
}
