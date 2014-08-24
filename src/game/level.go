package game

import (
	"bufio"
	"fmt"
	"os"
)

type TileProps struct {
	index int
	solid, deadly bool
}

type LevelMap [][]TileProps

var (
	tileMap = map[byte]TileProps {
		' ': {index: 0},
		'.': {index: 1},
		'a': {index: 2},
		's': {index: 3},
		'd': {index: 4},
		'f': {index: 5},
		'g': {index: 6},
		'[': {index: 7},
		'_': {index: 8},
		'Z': {index: 9},
		'z': {index: 10, solid:false},
		'x': {index: 11, solid:true},
		'c': {index: 12, solid:true},
		'v': {index: 13, solid:true},
		'b': {index: 14, solid:false},
		']': {index: 15},
		'+': {index: 22, solid: true},
		'^': {index: 30, solid: true, deadly: true},
	}
	outOfBounds = TileProps{solid: true}
	
	loadedMaps = make(map[string]LevelMap)
)

func LoadMap(name string) (LevelMap, error) {
	if m, ok := loadedMaps[name]; ok {
		return m, nil
	}
	
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	var m LevelMap
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		var r []TileProps
		for _, c := range []byte(line) {
			t, ok := tileMap[c]
			if !ok {
				return nil, fmt.Errorf("unknown map tile '%c'", c)
			}
			r = append(r, t)
		}
		m = append(m, r)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func (l LevelMap) QueryPoint(x, y int) TileProps {
	tx, ty := x/tileWidth, y/tileHeight
	if ty < 0 || ty >= len(l) || tx < 0 || tx >= len(l[ty]) {
		return outOfBounds
	}
	return l[ty][tx]
}