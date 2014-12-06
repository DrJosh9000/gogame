package sdl

/*
#cgo LDFLAGS: -framework SDL2

#include <stdlib.h>
#include <SDL2/SDL.h>

// Because SDL_LoadBMP is a macro.
SDL_Surface* loadbmp(const char *file) {
	return SDL_LoadBMP(file);
}
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type Surface struct {
	surface unsafe.Pointer
}

func (s *Surface) s() *C.SDL_Surface {
	return (*C.SDL_Surface)(s.surface)
}

func LoadBMP(path string) (*Surface, error) {
	cp := C.CString(path)
	defer C.free(unsafe.Pointer(cp))
	s := C.loadbmp(cp)
	if s == nil {
		return nil, fmt.Errorf("unable to load BMP %q: %s", path, Err())
	}
	return NewSurface(s), nil
}

func NewSurface(s *C.SDL_Surface) *Surface {
	r := &Surface{unsafe.Pointer(s)}
	runtime.SetFinalizer(r, func(x interface{}) {
		C.SDL_FreeSurface(x.(*Surface).s())
	})
	return r
}

func CreateSurface(w, h, bitDepth, rMask, gMask, bMask, aMask uint32) (*Surface, error) {
	surf := C.SDL_CreateRGBSurface(0, C.int(w), C.int(h), C.int(bitDepth), C.Uint32(rMask), C.Uint32(gMask), C.Uint32(bMask), C.Uint32(aMask))
	if surf == nil {
		return nil, Err()
	}
	return NewSurface(surf), nil
}

func CreateRGBASurface(w, h int) (*Surface, error) {
	surf := C.SDL_CreateRGBSurface(0, C.int(w), C.int(h), 32, 0xff000000, 0xff0000, 0xff00, 0xff)
	if surf == nil {
		return nil, Err()
	}
	return NewSurface(surf), nil
}

func (s *Surface) BlitOnto(dest *Surface, srcRect, dstRect *Rect) error {
	if errno := C.SDL_BlitSurface(s.s(), srcRect.r(), dest.s(), dstRect.r()); errno != 0 {
		return Err()
	}
	return nil
}

func (s *Surface) Fill(fill Colour) error {
	if errno := C.SDL_FillRect(s.s(), nil, fill.MapRGBA(s.s())); errno != 0 {
		return Err()
	}
	return nil
}

func (s *Surface) Size() (w, h int) {
	return int(s.s().w), int(s.s().h)
}

func (s *Surface) SetBlendMode(b BlendMode) error {
	if errno := C.SDL_SetSurfaceBlendMode(s.s(), C.SDL_BlendMode(b)); errno != 0 {
		return Err()
	}
	return nil
}

type Alignment int

const (
	LeftAlign Alignment = iota
	CentreAlign
	RightAlign
)

// StackSize computes the combined size of a slice of surfaces.
func StackSize(surfs []*Surface) (width, height int) {
	sumH, maxW := 0, 0
	for _, s := range surfs {
		w, h := s.Size()
		if w > maxW {
			maxW = w
		}
		sumH += h
	}
	return maxW, sumH
}

func Stack(dest *Surface, surfs []*Surface, x, y, maxW int, al Alignment) error {
	dstRect := Rect{Y: y}
	for _, s := range surfs {
		w, h := s.Size()
		dstRect.W = w
		switch al {
		case CentreAlign:
			dstRect.X = x + (maxW-w)/2
		case RightAlign:
			dstRect.X = x + maxW - w
		}
		if err := s.BlitOnto(dest, nil, &dstRect); err != nil {
			return err
		}
		dstRect.Y += h
	}
	return nil
}
