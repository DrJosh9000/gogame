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
		'-': {index: 2},
		'^': {index: 3, solid: true, deadly: true},
		'+': {index: 4, solid: true},
		'[': {index:5},
		'_': {index: 8},
		'#': {index: 10, solid: true},
		']': {index: 13},
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