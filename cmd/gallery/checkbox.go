package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var checkboxChecked bool

func init() {
	register(page{name: "Checkbox", order: 13, view: func(c *ui.Context) {
		pageBody(c, func() {
			kit.Checkbox(c, &checkboxChecked, "Auto-save")
			kit.Checkbox(c.Key("disabled"), &checkboxChecked, "Disabled").Disabled(true)
		})
	}})
}
