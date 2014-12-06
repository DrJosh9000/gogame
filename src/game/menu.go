package game

// TODO: menu system goes here

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
