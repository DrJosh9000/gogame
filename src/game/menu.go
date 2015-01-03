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

func newMenu(wm *windowManager) (*menu, error) {
	m := &menu{}
	x, y := (1024-256)/2, 200
	for _, mi := range menus {
		b := &button{
			Label: mi.text,
			sprite: &sprite{
				X:           x,
				Y:           y,
				TemplateKey: "button",
			},
			action: mi.action,
			parent: m,
		}
		y += 96
		m.addChild(b)
		wm.add(b)
	}
	return m, nil
}
