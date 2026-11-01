package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Shimmer", order: 45, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		pageBody(c, func() {
			ui.Column(c).Gap(theme.Space(2)).Children(func() {
				kit.Shimmer(c.Key("a"), t(c, "Thinking…"))
				kit.Shimmer(c.Key("b"), t(c, "Generating a title for this conversation"))
				kit.Shimmer(c.Key("c"), "正在思考…")
			})
		})
	}})
}
