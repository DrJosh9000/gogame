package sdl

/*
#include <stdlib.h>
#include <SDL2/SDL.h>
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
	if t.t() != nil {
		C.SDL_DestroyTexture(t.t())
	}
	t.texture = nil
}
