package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var tabsSelected int

func init() {
	register(page{name: "Tabs", order: 14, view: func(c *ui.Context) {
		ui.Column(c).Grow(1).Gap(0).Children(func() {
			kit.Tabs(c, &tabsSelected, []string{"General", "Editor", "Advanced"})
			ui.Text(c, []string{"General settings", "Editor settings", "Advanced settings"}[tabsSelected]).Padding(16)
		})
	}})
}
