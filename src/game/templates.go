package game

var templateLibrary = map[string]*spriteTemplate{
	"button": {
		sheetFile:   "assets/button.png",
		framesX:     1,
		framesY:     2,
		frameWidth:  256,
		frameHeight: 64,
	},
	"cursor": {
		sheetFile:   "assets/cursor.png",
		baseX:       32,
		baseY:       32,
		framesX:     2,
		framesY:     1,
		frameWidth:  64,
		frameHeight: 64,
	},
	"hex": {
		sheetFile:   "assets/hex.png",
		framesX:     1,
		framesY:     1,
		frameWidth:  128,
		frameHeight: 128,
	},
}
