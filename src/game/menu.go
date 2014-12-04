package game

// TODO: menu system goes here

type menuAction func() error
var menus = map[string]menuAction {
	"Start game": func() error {
		// TODO: start a game
		return nil
	},
	"Quit": func() error {
		notify("global", quitMsg)
		return nil
	},
}