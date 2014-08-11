package sdl

/*
#cgo CFLAGS: -I/Library/Frameworks/SDL2.framework/Headers
#cgo LDFLAGS: -F/Library/Frameworks -framework SDL2

#include <stdlib.h>
#include <SDL.h>

const int kRendererSoftware = SDL_RENDERER_SOFTWARE;
const int kRendererAccelerated = SDL_RENDERER_ACCELERATED;
const int kRendererPresentVSync = SDL_RENDERER_PRESENTVSYNC;
const int kRendererTargetTexture = SDL_RENDERER_TARGETTEXTURE;

const int kFlipNone = SDL_FLIP_NONE;
const int kFlipHorizontal = SDL_FLIP_HORIZONTAL;
const int kFlipVertical = SDL_FLIP_VERTICAL;
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type RendererOption uint32
type RendererFlip uint32

const (
	RendererSoftware      RendererOption = C.kRendererSoftware
	RendererAccelerated   RendererOption = C.kRendererAccelerated
	RendererPresentVSync  RendererOption = C.kRendererPresentVSync
	RendererTargetTexture RendererOption = C.kRendererTargetTexture

	FlipNone       RendererFlip = C.kFlipNone
	FlipHorizontal RendererFlip = C.kFlipHorizontal
	FlipVertical   RendererFlip = C.kFlipVertical
)

type Renderer struct {
	renderer unsafe.Pointer
}

func (r *Renderer) r() *C.SDL_Renderer {
	return (*C.SDL_Renderer)(r.renderer)
}

func (r *Renderer) Clear() {
	C.SDL_RenderClear(r.r())
}

func (r *Renderer) Present() {
	C.SDL_RenderPresent(r.r())
}

func (r *Renderer) Copy(t *Texture, src, dst *C.SDL_Rect) {
	C.SDL_RenderCopy(r.r(), t.t(), src, dst)
}

func (r *Renderer) CopyEx(t *Texture, src, dst *C.SDL_Rect, angle float64, center *C.SDL_Point, flip RendererFlip) {
	C.SDL_RenderCopyEx(r.r(), t.t(), src, dst, C.double(angle), center, C.SDL_RendererFlip(flip))
}

func (r *Renderer) Destroy() {
	if r.r() != nil {
		C.SDL_DestroyRenderer(r.r())
	}
}

func (r *Renderer) SetDrawColor(c Color) {
	C.SDL_SetRenderDrawColor(r.r(), C.Uint8(c.red), C.Uint8(c.green), C.Uint8(c.blue), C.Uint8(c.alpha))
}

func (r *Renderer) TextureFromSurface(s *Surface) (*Texture, error) {
	t := C.SDL_CreateTextureFromSurface(r.r(), s.s())
	if t == nil {
		return nil, fmt.Errorf("unable to create texture from surface: %s", Err())
	}
	return &Texture{unsafe.Pointer(t)}, nil
}

// Convenience function that loads a BMP into a surface,
// adds it to the renderer as a texture, and frees the surface.
func (r *Renderer) LoadBMP(path string) (*Texture, error) {
	s, err := LoadBMP(path)
	if err != nil {
		return nil, err
	}
	defer s.Free()
	return r.TextureFromSurface(s)
}
