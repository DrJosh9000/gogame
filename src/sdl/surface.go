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
	return &Surface{unsafe.Pointer(s)}, nil
}

func (s *Surface) Free() {
	if s.s() != nil {
		C.SDL_FreeSurface(s.s())
	}
	s.surface = nil
}
