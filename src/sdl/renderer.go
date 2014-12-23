package sdl

/*
#include <stdlib.h>
#include <SDL2/SDL.h>

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

type offset struct {
	x, y int
}

type Renderer struct {
	renderer    unsafe.Pointer
	offsetStack []offset
	offset      offset
}

func (r *Renderer) r() *C.SDL_Renderer {
	return (*C.SDL_Renderer)(r.renderer)
}

func (r *Renderer) Clear() error {
	if errno := C.SDL_RenderClear(r.r()); errno != 0 {
		return fmt.Errorf("error in SDL_RenderClear: %d %s", errno, Err())
	}
	return nil
}

func (r *Renderer) Present() {
	C.SDL_RenderPresent(r.r())
}

func (r *Renderer) Copy(t *Texture, src, dst Rect) error {
	dst.X += r.offset.x
	dst.Y += r.offset.y
	if errno := C.SDL_RenderCopy(r.r(), t.t(), src.r(), dst.r()); errno != 0 {
		return fmt.Errorf("error in SDL_RenderCopy: %d %s", errno, Err())
	}
	return nil
}

func (r *Renderer) CopyEx(t *Texture, src, dst Rect, angle float64, center Point, flip RendererFlip) error {
	dst.X += r.offset.x
	dst.Y += r.offset.y
	if errno := C.SDL_RenderCopyEx(r.r(), t.t(), src.r(), dst.r(), C.double(angle), center.p(), C.SDL_RendererFlip(flip)); errno != 0 {
		return fmt.Errorf("error in SDL_RenderCopyEx: %d %s", errno, Err())
	}
	return nil
}

func (r *Renderer) DrawPoint(x, y int) error {
	if errno := C.SDL_RenderDrawPoint(r.r(), C.int(x+r.offset.x), C.int(y+r.offset.y)); errno != 0 {
		return fmt.Errorf("error in SDL_RenderDrawPoint: %d %s", errno, Err())
	}
	return nil
}

func (r *Renderer) DrawRect(dst Rect) error {
	dst.X += r.offset.x
	dst.Y += r.offset.y
	if errno := C.SDL_RenderDrawRect(r.r(), dst.r()); errno != 0 {
		return fmt.Errorf("error in SDL_RenderDrawRect: %d %s", errno, Err())
	}
	return nil
}

func (r *Renderer) FillRect(dst Rect) error {
	dst.X += r.offset.x
	dst.Y += r.offset.y
	if errno := C.SDL_RenderFillRect(r.r(), dst.r()); errno != 0 {
		return fmt.Errorf("error in SDL_RenderFillRect: %d %s", errno, Err())
	}
	return nil
}

func (r *Renderer) PushOffset(x, y int) {
	r.offsetStack = append(r.offsetStack, offset{x, y})
	r.offset.x += x
	r.offset.y += y
}

func (r *Renderer) PopOffset() {
	if len(r.offsetStack) == 0 {
		return
	}
	l := len(r.offsetStack) - 1
	r.offset.x -= r.offsetStack[l].x
	r.offset.y -= r.offsetStack[l].y
	r.offsetStack = r.offsetStack[:l]
}

func (r *Renderer) ResetOffset() {
	r.offsetStack = nil
	r.offset = offset{}
}

func (r *Renderer) Destroy() {
	if r.r() != nil {
		C.SDL_DestroyRenderer(r.r())
	}
	r.renderer = nil
}

func (r *Renderer) SetDrawBlendMode(b BlendMode) error {
	if errno := C.SDL_SetRenderDrawBlendMode(r.r(), C.SDL_BlendMode(b)); errno != 0 {
		return Err()
	}
	return nil
}

func (r *Renderer) SetDrawColour(c Colour) {
	C.SDL_SetRenderDrawColor(r.r(), C.Uint8(c.Red), C.Uint8(c.Green), C.Uint8(c.Blue), C.Uint8(c.Alpha))
}

func (r *Renderer) TextureFromSurface(s *Surface) (*Texture, error) {
	t := C.SDL_CreateTextureFromSurface(r.r(), s.s())
	if t == nil {
		return nil, fmt.Errorf("unable to create texture from surface: %s", Err())
	}
	return &Texture{unsafe.Pointer(t)}, nil
}

// Convenience function that loads a BMP into a surface,
// adds it to the renderer as a texture.
func (r *Renderer) LoadBMP(path string) (*Texture, error) {
	s, err := LoadBMP(path)
	if err != nil {
		return nil, err
	}
	return r.TextureFromSurface(s)
}

// Similar but uses the SDL_image library.
func (r *Renderer) LoadImage(path string) (*Texture, error) {
	s, err := LoadImage(path)
	if err != nil {
		return nil, err
	}
	return r.TextureFromSurface(s)
}
