package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var patchRows ui.ListState

func init() {
	register(page{name: "Patch", order: 36, view: func(c *ui.Context) {
		number := 42
		kit.PatchLines(c, &patchRows, []kit.PatchLine{{Kind: kit.PatchHeader, Text: "@@ component.go @@"}, {Kind: kit.PatchDeleted, Text: "return oldValue"}, {Kind: kit.PatchAdded, Text: "return newValue", NewLine: &number}}, kit.PatchOpt{Source: "@@ component.go @@\n-return oldValue\n+return newValue\n"}).Grow(1)
	}})
}
