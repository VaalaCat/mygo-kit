package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var disclosureOpen = true

func init() {
	register(page{name: "Disclosure", order: 15, view: func(c *ui.Context) {
		pageBody(c, func() {
			kit.Disclosure(c, t(c, "More settings"), &disclosureOpen, func() { ui.Text(c, t(c, "Native keyboard behavior, independently styled.")) })
		})
	}})
}
