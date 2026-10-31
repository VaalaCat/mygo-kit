package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

var diffState = editor.DiffState{Old: editor.NewBuffer("package main\n\nfunc message() string {\n    return \"before\"\n}\n"), New: editor.NewBuffer("package main\n\nfunc message() string {\n    return \"native workspace\"\n}\n"), Layout: editor.DiffUnified, ExpandUnchanged: true}

func init() {
	register(page{name: "Diff", order: 33, view: func(c *ui.Context) {
		kit.Diff(c, &diffState).Grow(1)
	}})
}
