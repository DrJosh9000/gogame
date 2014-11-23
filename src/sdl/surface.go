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

func (s *Surface) Size() (w, h int) {
	return int(s.s().w), int(s.s().h)
}

type Alignment int

const (
	LeftAlign Alignment = iota
	CentreAlign
	RightAlign
)

func Stack(surfs []*Surface, fill Colour, al Alignment) (*Surface, error) {
	sumH, maxW := 0, 0
	for _, s := range surfs {
		w, h := s.Size()
		if w > maxW {
			maxW = w
		}
		sumH += h
	}

	surf := C.SDL_CreateRGBSurface(0, C.int(maxW), C.int(sumH), 32, 0, 0, 0, 0)
	if surf == nil {
		return nil, Err()
	}

	if errno := C.SDL_FillRect(surf, nil, fill.MapRGBA(surf)); errno != 0 {
		return nil, Err()
	}

	dstRect := Rect{}
	for _, s := range surfs {
		w, h := s.Size()
		dstRect.W = w
		switch al {
		case CentreAlign:
			dstRect.X = (maxW - w) / 2
		case RightAlign:
			dstRect.X = maxW - w
		}
		if errno := C.SDL_BlitSurface(s.s(), nil, surf, dstRect.r()); errno != 0 {
			return nil, Err()
		}
		dstRect.Y += h
	}
	return NewSurface(surf), nil
}
