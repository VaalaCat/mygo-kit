package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var overviewTab int

func init() {
	register(page{name: "Overview", order: 0, view: func(c *ui.Context) {
		ui.Scroll(c).Grow(1).Children(func() {
			ui.Column(c).Padding(24).Gap(24).Children(func() {
				kit.GetTheme(c).TextRole(ui.Text(c, "A focused native workspace"), kit.TextTitle)
				ui.Text(c, "连续面板、紧凑工具栏与清晰的文字层级。\nReusable Go controls, without another application runtime.").TextColor(kit.GetTheme(c).Colors.Muted)
				kit.Separator(c)
				kit.Tabs(c, &overviewTab, []string{"Preview", "Source", "Details"})
				kit.Markdown(c, "## Reading, without a webview\n\nThe component library adds **one visual language** over MyGo's native behavior. Keyboard, selection and IME stay with the framework.\n\n- [x] Semantic One Light / One Dark themes\n- [x] Compact controls and shared borders\n- [ ] Your application's next feature\n\n```go\nkit.ApplyTheme(c, kit.OneLight())\nif kit.Button(c, \"Open\").Clicked() {\n    // application action\n}\n```", kit.MarkdownOpt{})
				ui.Row(c).Gap(8).Children(func() {
					if kit.PrimaryButton(c, "Open the workspace").Clicked() {
						showPage("Workspace")
					}
					kit.Button(c, "Documentation")
				})
			})
		})
	}})
}
