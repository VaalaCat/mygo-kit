package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Spinner", order: 22, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				kit.Spinner(c, 12)
				kit.Spinner(c, 16)
				kit.Spinner(c, 24)
				ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
					kit.Spinner(c, 14)
					ui.Text(c, "Thinking…").TextColor(kit.GetTheme(c).Colors.Muted)
				})
			})
		})
	}})
}
