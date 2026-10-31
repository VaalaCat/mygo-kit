package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	treeState    ui.OutlineState[string]
	treeChildren = map[string][]string{
		"workspace": {"cmd", "internal", "README.md"},
		"cmd":       {"gallery", "server"},
		"internal":  {"desktop", "runtime"},
	}
)

func init() {
	register(page{name: "Tree", order: 37, view: func(c *ui.Context) {
		kit.Tree(c, &treeState, []string{"workspace"}, func(id string) []string { return treeChildren[id] }, func(id string) {
			ui.Row(c).Height(28).Padding(0, 8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, id).SingleLine()
			})
		}).Grow(1)
	}})
}
