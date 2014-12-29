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
		b := &button{
			Label:  mi.text,
			sprite: &sprite{TemplateKey: "button"},
			action: mi.action,
			parent: m,
		}
		go b.life()
		b.Y = y
		y += 96
		m.addChild(b)
	}
	return m, nil
}
