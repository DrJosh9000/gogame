package sdl

/*
#cgo CFLAGS: -I/Users/josh/Code/gogame/include
#cgo LDFLAGS: -F/Library/Frameworks -framework SDL2

#include <SDL2/SDL.h>

const int kQuit = SDL_QUIT;
const int kKeyDown = SDL_KEYDOWN;
const int kKeyUp = SDL_KEYUP;

Uint32 getType(SDL_Event *ev) {
	return ev->type;
}

SDL_Keycode getKeyCode(SDL_Event *ev) {
	return ev->key.keysym.sym;
}
*/
import "C"

type KeyEventType int

const (
	KeyDown KeyEventType = iota
	KeyUp
)

type QuitEvent int
type KeyEvent struct {
	Type    KeyEventType
	KeyCode uint32
}

func HandleEvents(handler func(e interface{}) error) error {
	var ev C.SDL_Event
	for C.SDL_PollEvent(&ev) != 0 {
		var err error
		switch C.getType(&ev) {
		case C.kQuit:
			err = handler(QuitEvent(0))
		case C.kKeyDown:
			err = handler(KeyEvent{
				Type:    KeyDown,
				KeyCode: uint32(C.getKeyCode(&ev)),
			})
		case C.kKeyUp:
			err = handler(KeyEvent{
				Type:    KeyUp,
				KeyCode: uint32(C.getKeyCode(&ev)),
			})
		}
		if err != nil {
			return err
		}
	}
	return nil
}
