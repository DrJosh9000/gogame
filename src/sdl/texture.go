package sdl

/*
#include <stdlib.h>
#include <SDL2/SDL.h>

const int kBlendModeNone  = SDL_BLENDMODE_NONE;
const int kBlendModeBlend = SDL_BLENDMODE_BLEND;
const int kBlendModeAdd   = SDL_BLENDMODE_ADD;
const int kBlendModeMod   = SDL_BLENDMODE_MOD;
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type BlendMode int

const (
	BlendModeNone  = BlendMode(C.kBlendModeNone)
	BlendModeBlend = BlendMode(C.kBlendModeBlend)
	BlendModeAdd   = BlendMode(C.kBlendModeAdd)
	BlendModeMod   = BlendMode(C.kBlendModeMod)
)

type Texture struct {
	texture unsafe.Pointer
}

func (t *Texture) t() *C.SDL_Texture {
	return (*C.SDL_Texture)(t.texture)
}

func NewTexture(t *C.SDL_Texture) *Texture {
	r := &Texture{unsafe.Pointer(t)}
	runtime.SetFinalizer(r, func(x interface{}) {
		tex := x.(*Texture)
		C.SDL_DestroyTexture(tex.t())
		tex.texture = nil
	})
	return r
}

func (t *Texture) SetBlendMode(b BlendMode) error {
	if errno := C.SDL_SetTextureBlendMode(t.t(), C.SDL_BlendMode(b)); errno != 0 {
		return Err()
	}
	return nil
}
