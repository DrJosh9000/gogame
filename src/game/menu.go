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
			// TODO: start a game
			return nil
		},
	},
	{
		text: "Level editor",
		action: func() error {
			// TODO: launch level editor
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

func newMenu(ctx *sdl.Context) (*complexBase, error) {
	base := new(complexBase)
	y := 200
	for _, mi := range menus {
		b, err := newButton(ctx, buttonTemplate, mi.text, mi.action)
		if err != nil {
			return nil, err
		}
		b.x, b.y = 512-128, y
		y += 128
		base.addChild(b)
	}
	return base, nil
}
