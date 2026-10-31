package main

import (
	"fmt"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var listRows = ui.ListState{Key: func(i int) any { return i }}

func init() {
	register(page{name: "Virtual list", order: 35, view: func(c *ui.Context) {
		kit.List(c, &listRows, 10000, func(i int) {
			ui.Row(c).Height(float32(28+i%3*8)).Padding(0, 16).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf("Item %05d", i)).Font(kit.GetTheme(c).Typography.MonoFont)
				ui.Spacer(c)
				ui.Text(c, "Variable height").TextColor(kit.GetTheme(c).Colors.Muted)
			})
		}).Grow(1)
	}})
}
