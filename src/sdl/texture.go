package sdl

/*
#cgo CFLAGS: -I/Library/Frameworks/SDL2.framework/Headers
#cgo LDFLAGS: -F/Library/Frameworks -framework SDL2

#include <stdlib.h>
#include <SDL.h>
*/
import "C"
import (
	"unsafe"
)

type Texture struct {
	texture unsafe.Pointer
}

func (t *Texture) t() *C.SDL_Texture {
	return (*C.SDL_Texture)(t.texture)
}

func (t *Texture) Destroy() {
	C.SDL_DestroyTexture(t.t())
}