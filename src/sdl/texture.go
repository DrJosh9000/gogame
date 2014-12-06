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

func (t *Texture) Destroy() {
	if t.t() != nil {
		C.SDL_DestroyTexture(t.t())
	}
	t.texture = nil
}

func (t *Texture) SetBlendMode(b BlendMode) error {
	if errno := C.SDL_SetTextureBlendMode(t.t(), C.SDL_BlendMode(b)); errno != 0 {
		return Err()
	}
	return nil
}

type TextureManager struct {
	assets map[string]*Texture
	r      *Renderer
}

func NewTextureManager(r *Renderer) *TextureManager {
	return &TextureManager{
		assets: make(map[string]*Texture),
		r:      r,
	}
}

func (a *TextureManager) GetTexture(name string) (*Texture, error) {
	if t, ok := a.assets[name]; ok {
		return t, nil
	}
	t, err := a.r.LoadImage(name)
	if err != nil {
		return nil, err
	}
	a.assets[name] = t
	return t, nil
}

func (a *TextureManager) Destroy() {
	if a.assets != nil {
		for _, x := range a.assets {
			x.Destroy()
		}
	}
	a.assets = nil
}
