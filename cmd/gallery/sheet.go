package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var sheetLeft, sheetRight bool

func init() {
	register(page{name: "Sheet", order: 26, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(8).Children(func() {
				if kit.Button(c.Key("sheet-left"), t(c, "Open left sheet")).Clicked() {
					sheetLeft = true
				}
				if kit.Button(c.Key("sheet-right"), t(c, "Open right sheet")).Clicked() {
					sheetRight = true
				}
			})
			kit.Sheet(c.Key("left"), &sheetLeft, kit.SheetLeft, 280, func() {
				kit.Header(c, t(c, "Conversations"), nil)
				ui.Text(c, t(c, "Click the backdrop or press Escape.")).Padding(12).TextColor(kit.GetTheme(c).Colors.Muted)
			})
			kit.Sheet(c.Key("right"), &sheetRight, kit.SheetRight, 360, func() {
				kit.Header(c, t(c, "Details"), nil)
				ui.Text(c, t(c, "A drawer from the right edge.")).Padding(12).TextColor(kit.GetTheme(c).Colors.Muted)
			})
		})
	}})
}
