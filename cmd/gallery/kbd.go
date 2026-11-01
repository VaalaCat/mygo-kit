package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Kbd", order: 20, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, t(c, "Open the command palette"))
				kit.Kbd(c, "⌘K")
			})
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, t(c, "New conversation"))
				kit.Kbd(c, "Ctrl+N")
			})
		})
	}})
}
