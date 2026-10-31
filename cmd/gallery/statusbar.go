package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var statusActive bool

func init() {
	register(page{name: "StatusBar", order: 24, view: func(c *ui.Context) {
		ui.Column(c).Grow(1).Children(func() {
			ui.Spacer(c)
			kit.StatusBar(c, func() {
				button := kit.StatusButton(c.Key("status-left"), kit.IconPanelLeft, "Toggle sidebar", statusActive)
				kit.Tooltip(c, button, "Toggle sidebar")
				if button.Clicked() {
					statusActive = !statusActive
				}
				ui.Text(c, "Connected")
				ui.Spacer(c)
				kit.StatusButton(c.Key("status-bottom"), kit.IconPanelBottom, "Terminal", false)
				kit.StatusButton(c.Key("status-right"), kit.IconPanelRight, "Outline", true)
			})
		})
	}})
}
