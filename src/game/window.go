package game

// Sigh.

import (
	"sdl"
)

type window interface {
	bounds() sdl.Rect
	zer
}
