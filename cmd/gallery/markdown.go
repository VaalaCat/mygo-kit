package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Markdown", order: 30, view: func(c *ui.Context) {
		ui.Scroll(c).Grow(1).Padding(24).Children(func() {
			kit.Markdown(c, "# Markdown\n\nSelect text across **paragraphs**, `inline code` and code blocks.\n\n| Component | Mode |\n|---|---|\n| Markdown | Native |\n| Selection | Shared |\n\n- [x] Task list\n- [ ] More content\n\n```go\nfunc hello() string {\n    return \"你好，native UI\"\n}\n```\n\n[MyGo documentation](https://mygo.egoist.dev)\n\n![Grid preview](demo://image)", kit.MarkdownOpt{Image: func(string) kit.ImageResource { return galleryImage }})
		})
	}})
}
