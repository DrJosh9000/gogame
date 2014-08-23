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


type TextureManager struct {
	assets map[string]*Texture
	r *Renderer
}

func NewTextureManager(r *Renderer) *TextureManager {
	return &TextureManager{
		assets: make(map[string]*Texture),
		r: r,
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