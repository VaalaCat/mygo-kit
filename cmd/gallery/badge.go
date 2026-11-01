package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Badge", order: 21, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				kit.Badge(c, t(c, "Neutral"), kit.BadgeNeutral)
				kit.Badge(c, t(c, "Accent"), kit.BadgeAccent)
				kit.Badge(c, t(c, "Success"), kit.BadgeSuccess)
				kit.Badge(c, t(c, "Warning"), kit.BadgeWarning)
				kit.Badge(c, t(c, "Danger"), kit.BadgeDanger)
			})
			kit.Separator(c)
			theme := kit.GetTheme(c)
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				for _, state := range []struct {
					label string
					color ui.Color
				}{{"Connected", theme.Colors.Success}, {"Connecting", theme.Colors.Warning}, {"Failed", theme.Colors.Danger}, {"Unread", theme.Colors.Accent}} {
					ui.Row(c.Key(state.label)).Gap(8).AlignItems(ui.Center).Children(func() {
						kit.StatusDot(c, state.color).Label(t(c, state.label))
						ui.Text(c, t(c, state.label))
					})
				}
			})
		})
	}})
}
