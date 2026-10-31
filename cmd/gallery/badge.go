package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Badge", order: 21, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				kit.Badge(c, "Neutral", kit.BadgeNeutral)
				kit.Badge(c, "Accent", kit.BadgeAccent)
				kit.Badge(c, "Success", kit.BadgeSuccess)
				kit.Badge(c, "Warning", kit.BadgeWarning)
				kit.Badge(c, "Danger", kit.BadgeDanger)
			})
			kit.Separator(c)
			t := kit.GetTheme(c)
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				for _, state := range []struct {
					label string
					color ui.Color
				}{{"Connected", t.Colors.Success}, {"Connecting", t.Colors.Warning}, {"Failed", t.Colors.Danger}, {"Unread", t.Colors.Accent}} {
					ui.Row(c.Key(state.label)).Gap(8).AlignItems(ui.Center).Children(func() {
						kit.StatusDot(c, state.color).Label(state.label)
						ui.Text(c, state.label)
					})
				}
			})
		})
	}})
}
