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
const int kWindow	    = SDL_WINDOWEVENT;

const unsigned int kMouseLeftMask   = SDL_BUTTON_LMASK;
const unsigned int kMouseMiddleMask = SDL_BUTTON_MMASK;
const unsigned int kMouseRightMask  = SDL_BUTTON_RMASK;
const unsigned int kMouseX1Mask     = SDL_BUTTON_X1MASK;
const unsigned int kMouseX2Mask     = SDL_BUTTON_X2MASK;

const unsigned char kWindowEventShown       = SDL_WINDOWEVENT_SHOWN;
const unsigned char kWindowEventHidden      = SDL_WINDOWEVENT_HIDDEN;
const unsigned char kWindowEventExposed     = SDL_WINDOWEVENT_EXPOSED;
const unsigned char kWindowEventMoved       = SDL_WINDOWEVENT_MOVED;
const unsigned char kWindowEventResized     = SDL_WINDOWEVENT_RESIZED;
const unsigned char kWindowEventSizeChanged = SDL_WINDOWEVENT_SIZE_CHANGED;
const unsigned char kWindowEventMinimized   = SDL_WINDOWEVENT_MINIMIZED;
const unsigned char kWindowEventMaximized   = SDL_WINDOWEVENT_MAXIMIZED;
const unsigned char kWindowEventRestored    = SDL_WINDOWEVENT_RESTORED;
const unsigned char kWindowEventEnter       = SDL_WINDOWEVENT_ENTER;
const unsigned char kWindowEventLeave       = SDL_WINDOWEVENT_LEAVE;
const unsigned char kWindowEventFocusGained = SDL_WINDOWEVENT_FOCUS_GAINED;
const unsigned char kWindowEventFocusLost   = SDL_WINDOWEVENT_FOCUS_LOST;
const unsigned char kWindowEventClose       = SDL_WINDOWEVENT_CLOSE;

Uint32 getType(SDL_Event *ev) {
	return ev->type;
}

SDL_Keycode getKeyCode(SDL_Event *ev) {
	return ev->key.keysym.sym;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Event interface {
	fmt.Stringer
}

type Timestamp uint32

func (t Timestamp) String() string {
	return fmt.Sprintf("@%d", uint32(t))
}

type QuitEvent struct {
	Timestamp
}
type KeyEvent struct {
	Timestamp
	WindowID, KeyCode uint32
	State, Repeat     uint8
}

type KeyDownEvent KeyEvent
type KeyUpEvent KeyEvent

type MouseMotionEvent struct {
	Timestamp
	WindowID, MouseID, ButtonState uint32
	X, Y, XRel, YRel               int
}

type MouseButtonEvent struct {
	Timestamp
	WindowID, MouseID     uint32
	Button, State, Clicks uint8
	X, Y                  int
}

type MouseButtonDownEvent MouseButtonEvent
type MouseButtonUpEvent MouseButtonEvent

const (
	MouseLeftMask   = uint32(C.kMouseLeftMask)
	MouseMiddleMask = uint32(C.kMouseMiddleMask)
	MouseRightMask  = uint32(C.kMouseRightMask)
	MouseX1Mask     = uint32(C.kMouseX1Mask)
	MouseX2Mask     = uint32(C.kMouseX2Mask)
)

type MouseWheelEvent struct {
	Timestamp
	WindowID, MouseID uint32
	X, Y              int
}

type WindowEvent struct {
	Timestamp
	WindowID     uint32
	EventID      WindowEventID
	Data1, Data2 int32
}

type WindowEventID uint8

const (
	WindowShown       = WindowEventID(C.kWindowEventShown)
	WindowHidden      = WindowEventID(C.kWindowEventHidden)
	WindowExposed     = WindowEventID(C.kWindowEventExposed)
	WindowMoved       = WindowEventID(C.kWindowEventMoved)
	WindowResized     = WindowEventID(C.kWindowEventResized)
	WindowSizeChanged = WindowEventID(C.kWindowEventSizeChanged)
	WindowMinimized   = WindowEventID(C.kWindowEventMinimized)
	WindowMaximized   = WindowEventID(C.kWindowEventMaximized)
	WindowRestored    = WindowEventID(C.kWindowEventRestored)
	WindowEnter       = WindowEventID(C.kWindowEventEnter)
	WindowLeave       = WindowEventID(C.kWindowEventLeave)
	WindowFocusGained = WindowEventID(C.kWindowEventFocusGained)
	WindowFocusLost   = WindowEventID(C.kWindowEventFocusLost)
	WindowClose       = WindowEventID(C.kWindowEventClose)
)

func HandleEvents(handler func(e Event) error) error {
	var ev C.SDL_Event
	for C.SDL_PollEvent(&ev) != 0 {
		var gev Event
		switch C.getType(&ev) {
		case C.kQuit:
			gev = QuitEvent{}
		case C.kKeyDown, C.kKeyUp:
			kev := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))
			e := &KeyEvent{
				Timestamp: Timestamp(kev.timestamp),
				WindowID:  uint32(kev.windowID),
				State:     uint8(kev.state),
				Repeat:    uint8(kev.repeat),
				KeyCode:   uint32(C.getKeyCode(&ev)),
			}
			if C.getType(&ev) == C.kKeyUp {
				gev = (*KeyUpEvent)(e)
			} else {
				gev = (*KeyDownEvent)(e)
			}
		case C.kMouseMotion:
			mmev := (*C.SDL_MouseMotionEvent)(unsafe.Pointer(&ev))
			gev = &MouseMotionEvent{
				Timestamp:   Timestamp(mmev.timestamp),
				WindowID:    uint32(mmev.windowID),
				MouseID:     uint32(mmev.which),
				ButtonState: uint32(mmev.state),
				X:           int(mmev.x),
				Y:           int(mmev.y),
				XRel:        int(mmev.xrel),
				YRel:        int(mmev.yrel),
			}
		case C.kMouseDown, C.kMouseUp:
			mbev := (*C.SDL_MouseButtonEvent)(unsafe.Pointer(&ev))
			e := &MouseButtonEvent{
				Timestamp: Timestamp(mbev.timestamp),
				WindowID:  uint32(mbev.windowID),
				MouseID:   uint32(mbev.which),
				Button:    uint8(mbev.button),
				State:     uint8(mbev.state),
				Clicks:    uint8(mbev.clicks),
				X:         int(mbev.x),
				Y:         int(mbev.y),
			}
			if C.getType(&ev) == C.kMouseUp {
				gev = (*MouseButtonUpEvent)(e)
			} else {
				gev = (*MouseButtonDownEvent)(e)
			}
		case C.kMouseWheel:
			mwev := (*C.SDL_MouseWheelEvent)(unsafe.Pointer(&ev))
			gev = &MouseWheelEvent{
				Timestamp: Timestamp(mwev.timestamp),
				WindowID:  uint32(mwev.windowID),
				MouseID:   uint32(mwev.which),
				X:         int(mwev.x),
				Y:         int(mwev.y),
			}
		case C.kWindow:
			wev := (*C.SDL_WindowEvent)(unsafe.Pointer(&ev))
			gev = &WindowEvent{
				Timestamp: Timestamp(wev.timestamp),
				WindowID:  uint32(wev.windowID),
				EventID:   WindowEventID(wev.event),
				Data1:     int32(wev.data1),
				Data2:     int32(wev.data2),
			}
		}
		if err := handler(gev); err != nil {
			return err
		}
	}
	return nil
}
