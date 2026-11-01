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
				button := kit.StatusButton(c.Key("status-left"), kit.IconPanelLeft, t(c, "Toggle sidebar"), statusActive)
				kit.Tooltip(c, button, t(c, "Toggle sidebar"))
				if button.Clicked() {
					statusActive = !statusActive
				}
				ui.Text(c, t(c, "Connected"))
				ui.Spacer(c)
				kit.StatusButton(c.Key("status-bottom"), kit.IconPanelBottom, t(c, "Terminal"), false)
				kit.StatusButton(c.Key("status-right"), kit.IconPanelRight, t(c, "Outline"), true)
			})
		})
	}})
}
