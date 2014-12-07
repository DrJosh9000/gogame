package game

import (
	"log"
	"sdl"
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

func newMenu(ctx *sdl.Context) (*menu, error) {
	m := &menu{
		complexBase: complexBase{
			x: (1024 - buttonTemplate.frameWidth) / 2,
			y: 200,
		},
	}
	y := 0
	for _, mi := range menus {
		b, err := newButton(ctx, m, buttonTemplate, mi.text, mi.action)
		if err != nil {
			return nil, err
		}
		b.y = y
		y += 96
		m.addChild(b)
	}
	return m, nil
}
