package sdl

/*
#cgo LDFLAGS: -framework SDL2 -framework SDL2_ttf

#include <stdlib.h>
#include "SDL2_ttf/SDL_ttf.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

func init() {
	if errno := C.TTF_Init(); errno != 0 {
		panic(Err())
	}
}

type Font struct {
	font unsafe.Pointer
}

func (f *Font) f() *C.TTF_Font {
	return (*C.TTF_Font)(f.font)
}

func LoadFont(path string, size uint) (*Font, error) {
	cp := C.CString(path)
	defer C.free(unsafe.Pointer(cp))
	f := C.TTF_OpenFont(cp, C.int(size))
	if f == nil {
		return nil, Err()
	}
	r := &Font{unsafe.Pointer(f)}
	runtime.SetFinalizer(r, func(x interface{}) {
		C.TTF_CloseFont(x.(*Font).f())
	})
	return r, nil
}

func (f *Font) RenderSolid(text string, c Colour) (*Surface, error) {
	cp := C.CString(text)
	defer C.free(unsafe.Pointer(cp))
	
	s := C.TTF_RenderUTF8_Solid(f.f(), cp, c.c())
	if s == nil {
		return nil, Err()
	}
	return NewSurface(s), nil
}