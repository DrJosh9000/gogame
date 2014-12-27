package game

import (
	"log"
)

type menuItem struct {
	action func() error
	text   string
}

var menus = []menuItem{
	{
		text: "Start game",
		action: func() error {
			// TODO: start a game
			log.Print("start game button")
			notify("menuAction", "start")
			return nil
		},
	},
	{
		text: "Level editor",
		action: func() error {
			// TODO: launch level editor
			log.Print("level editor button")
			notify("menuAction", "levelEdit")
			return nil
		},
	},
	{
		text: "Quit",
		action: func() error {
			log.Print("quit button")
			quit()
			return nil
		},
	},
}

type menu struct {
	complexBase
}

func newMenu() (*menu, error) {
	m := &menu{
		complexBase: complexBase{
			// TODO: Less magic numbers.
			X: (1024 - 256) / 2,
			Y: 200,
		},
	}
	y := 0
	for _, mi := range menus {
		b, err := newButton(m, "button", mi.text, mi.action)
		if err != nil {
			return nil, err
		}
		b.Y = y
		y += 96
		m.addChild(b)
	}
	return m, nil
}
