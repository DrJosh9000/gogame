package sdl

/*
#cgo LDFLAGS: -framework SDL2

#include <SDL2/SDL.h>

const int kQuit 			= SDL_QUIT;
const int kKeyDown 		= SDL_KEYDOWN;
const int kKeyUp 		= SDL_KEYUP;
const int kMouseMotion 	= SDL_MOUSEMOTION;
const int kMouseDown 	= SDL_MOUSEBUTTONDOWN;
const int kMouseUp 		= SDL_MOUSEBUTTONUP;
const int kMouseWheel 	= SDL_MOUSEWHEEL;

const unsigned int kMouseLeftMask   = SDL_BUTTON_LMASK;
const unsigned int kMouseMiddleMask = SDL_BUTTON_MMASK;
const unsigned int kMouseRightMask  = SDL_BUTTON_RMASK;
const unsigned int kMouseX1Mask     = SDL_BUTTON_X1MASK;
const unsigned int kMouseX2Mask     = SDL_BUTTON_X2MASK;

Uint32 getType(SDL_Event *ev) {
	return ev->type;
}

SDL_Keycode getKeyCode(SDL_Event *ev) {
	return ev->key.keysym.sym;
}
*/
import "C"

import (
	"unsafe"
)

type QuitEvent int
type KeyEvent struct {
	Timestamp, WindowID, KeyCode uint32
	State, Repeat                uint8
}

type KeyDownEvent KeyEvent
type KeyUpEvent KeyEvent

type MouseMotionEvent struct {
	Timestamp, WindowID, MouseID, ButtonState uint32
	X, Y, XRel, YRel                          int
}

type MouseButtonEvent struct {
	Timestamp, WindowID, MouseID uint32
	Button, State, Clicks        uint8
	X, Y                         int
}

type MouseButtonDownEvent MouseButtonEvent
type MouseButtonUpEvent MouseButtonEvent

type MouseWheelEvent struct {
	Timestamp, WindowID, MouseID uint32
	X, Y                         int
}

const (
	MouseLeftMask   = uint32(C.kMouseLeftMask)
	MouseMiddleMask = uint32(C.kMouseMiddleMask)
	MouseRightMask  = uint32(C.kMouseRightMask)
	MouseX1Mask     = uint32(C.kMouseX1Mask)
	MouseX2Mask     = uint32(C.kMouseX2Mask)
)

func HandleEvents(handler func(e interface{}) error) error {
	var ev C.SDL_Event
	for C.SDL_PollEvent(&ev) != 0 {
		var gev interface{}
		switch C.getType(&ev) {
		case C.kQuit:
			gev = QuitEvent(0)
		case C.kKeyDown:
			kev := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))
			gev = KeyDownEvent{
				Timestamp: uint32(kev.timestamp),
				WindowID:  uint32(kev.windowID),
				State:     uint8(kev.state),
				Repeat:    uint8(kev.repeat),
				KeyCode:   uint32(C.getKeyCode(&ev)),
			}
		case C.kKeyUp:
			kev := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))
			gev = KeyUpEvent{
				Timestamp: uint32(kev.timestamp),
				WindowID:  uint32(kev.windowID),
				State:     uint8(kev.state),
				Repeat:    uint8(kev.repeat),
				KeyCode:   uint32(C.getKeyCode(&ev)),
			}
		case C.kMouseMotion:
			mmev := (*C.SDL_MouseMotionEvent)(unsafe.Pointer(&ev))
			gev = MouseMotionEvent{
				Timestamp:   uint32(mmev.timestamp),
				WindowID:    uint32(mmev.windowID),
				MouseID:     uint32(mmev.which),
				ButtonState: uint32(mmev.state),
				X:           int(mmev.x),
				Y:           int(mmev.y),
				XRel:        int(mmev.xrel),
				YRel:        int(mmev.yrel),
			}
		case C.kMouseDown:
			mbev := (*C.SDL_MouseButtonEvent)(unsafe.Pointer(&ev))
			gev = MouseButtonDownEvent{
				Timestamp: uint32(mbev.timestamp),
				WindowID:  uint32(mbev.windowID),
				MouseID:   uint32(mbev.which),
				Button:    uint8(mbev.button),
				State:     uint8(mbev.state),
				Clicks:    uint8(mbev.clicks),
				X:         int(mbev.x),
				Y:         int(mbev.y),
			}
		case C.kMouseUp:
			mbev := (*C.SDL_MouseButtonEvent)(unsafe.Pointer(&ev))
			gev = MouseButtonUpEvent{
				Timestamp: uint32(mbev.timestamp),
				WindowID:  uint32(mbev.windowID),
				MouseID:   uint32(mbev.which),
				Button:    uint8(mbev.button),
				State:     uint8(mbev.state),
				Clicks:    uint8(mbev.clicks),
				X:         int(mbev.x),
				Y:         int(mbev.y),
			}
		case C.kMouseWheel:
			mwev := (*C.SDL_MouseWheelEvent)(unsafe.Pointer(&ev))
			gev = MouseWheelEvent{
				Timestamp: uint32(mwev.timestamp),
				WindowID:  uint32(mwev.windowID),
				MouseID:   uint32(mwev.which),
				X:         int(mwev.x),
				Y:         int(mwev.y),
			}
		}
		if err := handler(gev); err != nil {
			return err
		}
	}
	return nil
}
