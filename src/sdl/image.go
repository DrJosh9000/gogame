package sdl

/*
#cgo LDFLAGS: -framework SDL2 -framework SDL2_image

#include <stdlib.h>
#include "SDL2_image/SDL_image.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func init() {
	if e := C.IMG_Init(C.IMG_INIT_PNG); e&C.IMG_INIT_PNG == 0 {
		panic(fmt.Sprintf("Unable to init PNG support with SDL_image (got init mask %x)", e))
	}
}

func ImgErr() string {
	return C.GoString(C.IMG_GetError())
}

func LoadImage(file string) (*Surface, error) {
	cf := C.CString(file)
	defer C.free(unsafe.Pointer(cf))
	img := C.IMG_Load(cf)
	if img == nil {
		return nil, fmt.Errorf("unable to load image at %q: %v", file, ImgErr())
	}
	return NewSurface(img), nil
}
