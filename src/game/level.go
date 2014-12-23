package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type tileProps struct {
	index         int
	solid, deadly bool
}

// Row, Column
type levelLayer [][]tileProps
type level struct {
	levelMap       []levelLayer
	startX, startY int
	hasExit        bool
	exitX, exitY   int
}

var (
	transparentTile = tileProps{index: 0}
	tileMap         = map[byte]tileProps{
		' ': transparentTile,
		'.': {index: 1}, // space panel
		// 2: currently blank
		's': {index: 3}, // space platform
		'd': {index: 4},
		'f': {index: 5},
		'g': {index: 6}, // end space platform
		'[': {index: 7},
		'_': {index: 8},
		// 9: currently blank
		'z': {index: 10, solid: false},
		'x': {index: 11, solid: true},
		'c': {index: 12, solid: true},
		'v': {index: 13, solid: true},
		'b': {index: 14, solid: false},
		']': {index: 15},
		',': {index: 16},
		// 17, 18: currently blank
		'2': {index: 19},
		'3': {index: 20},
		'4': {index: 21},
		'5': {index: 22},
		'+': {index: 23, solid: true},
		'@': {index: 24, solid: true},
		'%': {index: 25, solid: true},
		'q': {index: 26, solid: false},
		'w': {index: 27, solid: true},
		'e': {index: 28, solid: true},
		'r': {index: 29, solid: true},
		't': {index: 30, solid: false},
		'^': {index: 31, solid: true, deadly: true},
		'n': {index: 32}, // Exit sign
		'm': {index: 33},
		'(': {index: 34},
		'R': {index: 35, solid: false},
		'Y': {index: 36, solid: false},
		'T': {index: 37, solid: false},
		')': {index: 38},

		'*': transparentTile, // Start position
		'X': transparentTile, // Exit
	}
	outOfBounds = tileProps{solid: true}

	loadedMaps = make(map[string]*level)
)

func loadLevel(name string) (*level, error) {
	if m, ok := loadedMaps[name]; ok {
		return m, nil
	}

	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	width := 32
	l := &level{}
	var m levelLayer
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if len(line) > 0 {
			// Filter comments.
			if line[0] == '#' {
				continue
			}
			// ! = directive.
			if line[0] == '!' {
				w, err := strconv.ParseInt(line[1:], 10, 32)
				if err != nil {
					return nil, err
				}
				width = int(w)
				continue
			}
			// Braces delineate layers.
			if line == "{" {
				m = nil
				continue
			}
			if line == "}" {
				l.levelMap = append(l.levelMap, m)
				continue
			}
		}

		var r []tileProps
		for j, c := range []byte(line) {
			switch c {
			case '*':
				l.startX = j
				l.startY = len(m)
			case 'X':
				l.hasExit = true
				l.exitX = j
				l.exitY = len(m)
			}
			t, ok := tileMap[c]
			if !ok {
				return nil, fmt.Errorf("unknown map tile '%c'", c)
			}
			r = append(r, t)
		}
		// Ensure the row is the right length.
		if len(r) < width {
			for i := len(r); i < width; i++ {
				r = append(r, transparentTile)
			}
		}
		m = append(m, r)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	loadedMaps[name] = l
	return l, nil
}

func floorDiv(n, d int) int {
	if n < 0 {
		return (n+1)/d - 1
	}
	return n / d
}

func (l levelLayer) queryPoint(x, y int) tileProps {
	if x < 0 || y < 0 {
		return outOfBounds
	}
	tx, ty := x/tileTemplate.frameWidth, y/tileTemplate.frameHeight
	if ty >= len(l) || tx >= len(l[ty]) {
		return outOfBounds
	}
	return l[ty][tx]
}

func (l *level) isPointSolid(x, y int) bool {
	for _, m := range l.levelMap {
		if m.queryPoint(x, y).solid {
			return true
		}
	}
	return false
}
