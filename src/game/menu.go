package game

type menuItem struct {
	action func() error
	text   string
}

var mainMenu = []menuItem{
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

func newMenu(menus []menuItem, wm *windowManager) (*menu, error) {
	m := &menu{}
	w := 384
	x, y := (1024-w)/2, 200
	for _, mi := range menus {
		b := &button{
			Label:  mi.text,
			W:      w,
			X:      x,
			Y:      y,
			Action: mi.action,
			parent: m,
		}
		if err := b.load(); err != nil {
			return nil, err
		}
		y += 96
		m.addChild(b)
		wm.add(b)
	}
	return m, nil
}
