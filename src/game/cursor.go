package game

import (
	"sdl"
)

var cursorTemplate = &spriteTemplate{
	name:        "cursor",
	sheetFile:   "assets/cursor.png",
	framesX:     1,
	framesY:     1,
	frameWidth:  64,
	frameHeight: 64,
}

type cursor struct {
	*sprite

	controller chan sdl.MouseMotionEvent
}
