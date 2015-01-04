package game

import (
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
			notify("menuAction", "start")
			return nil
		},
	},
	{
		text: "Level editor",
		action: func() error {
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
			text: text{
				Text:   mi.text,
				Draw:   sdl.BlackColour,
				Shadow: sdl.TransparentColour,
				Align:  sdl.CentreAlign,
			},
			sprite: &sprite{
				X:           x,
				Y:           y,
				TemplateKey: "button",
			},
			action: mi.action,
			parent: m,
		}
		if err := b.load(); err != nil {
			return nil, err
		}
		if err := b.text.load(); err != nil {
			return nil, err
		}
		y += 96
		m.addChild(b)
		wm.add(b)
	}
	return m, nil
}
