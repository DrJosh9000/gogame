package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type TileProps struct {
	index int
	solid, deadly bool
}

// Row, Column
type LevelLayer [][]TileProps
type Level []LevelLayer

var (
	tileMap = map[byte]TileProps {
		' ': {index: 0}, // transparent, don't even paint a tile.
		'.': {index: 1}, // space panel
		'a': {index: 2}, // space platform
		's': {index: 3},
		'd': {index: 4},
		'f': {index: 5},
		'g': {index: 6}, // end space platform
		'[': {index: 7},
		'_': {index: 8},
		'Z': {index: 9},
		'z': {index: 10, solid:false},
		'x': {index: 11, solid:true},
		'c': {index: 12, solid:true},
		'v': {index: 13, solid:true},
		'b': {index: 14, solid:false},
		']': {index: 15},
		',': {index: 16},
		'1': {index: 18},
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
	}
	outOfBounds = TileProps{solid: true}
	
	loadedMaps = make(map[string]Level)
)

func LoadLevel(name string) (Level, error) {
	if m, ok := loadedMaps[name]; ok {
		return m, nil
	}
	
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	width := 32
	var l Level
	var m LevelLayer
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
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
			l = append(l, m)
			continue
		}
	
		var r []TileProps
		for _, c := range []byte(line) {
			t, ok := tileMap[c]
			if !ok {
				return nil, fmt.Errorf("unknown map tile '%c'", c)
			}
			r = append(r, t)
		}
		// Ensure the row is the right length.
		if len(r) < width {
			for i:=len(r); i<width; i++ {
				r = append(r, tileMap[' '])
			}
		}
		m = append(m, r)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return l, nil
}

func (l LevelLayer) QueryPoint(x, y int) TileProps {
	tx, ty := x/tileWidth, y/tileHeight
	if ty < 0 || ty >= len(l) || tx < 0 || tx >= len(l[ty]) {
		return outOfBounds
	}
	return l[ty][tx]
}

func (l Level) IsPointSolid(x, y int) bool {
	for _, m := range l {
		if m.QueryPoint(x, y).solid {
			return true
		}
	}
	return false
}