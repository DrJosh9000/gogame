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
			return nil
		},
	},
	{
		text: "Level editor",
		action: func() error {
			// TODO: launch level editor
			log.Print("level editor button")
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

func newMenu(ctx *sdl.Context) (*complexBase, error) {
	base := &complexBase{
		x: (1024 - buttonTemplate.frameWidth) / 2,
		y: 200,
	}
	y := 0
	for _, mi := range menus {
		b, err := newButton(ctx, base, buttonTemplate, mi.text, mi.action)
		if err != nil {
			return nil, err
		}
		b.y = y
		y += 96
		base.addChild(b)
	}
	return base, nil
}
