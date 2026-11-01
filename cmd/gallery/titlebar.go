package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "TitleBar", order: 25, view: func(c *ui.Context) {
		ui.Column(c).Grow(1).Children(func() {
			kit.TitleBar(c, func() {
				kit.IconButton(c.Key("title-folder"), kit.IconFolder, t(c, "Workspace"))
				ui.Text(c, "~/code/project")
			})
			kit.EmptyState(c, t(c, "Window content"), t(c, "The title bar drags the window and keeps clear of system buttons."))
		})
	}})
}
