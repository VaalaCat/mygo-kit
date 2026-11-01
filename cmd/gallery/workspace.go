package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var workspaceState = kit.WorkspaceState{SidebarOpen: true, Right: kit.DockState{Open: true, ActivePanel: "outline"}, Bottom: kit.DockState{Open: true, ActivePanel: "terminal"}}

func init() {
	register(page{name: "Workspace", order: 29, full: true, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		muted := func(text string) { ui.Text(c, text).Padding(12).TextColor(theme.Colors.Muted) }
		kit.Workspace(c, &workspaceState, kit.WorkspaceOpt{
			TitleBar: func() {
				if kit.IconButton(c.Key("back"), kit.IconChevronLeft, t(c, "Back to gallery")).Clicked() {
					showPage("Overview")
				}
				ui.Text(c, "~/code/project")
			},
			Sidebar: func() {
				kit.Header(c, t(c, "CONVERSATIONS"), nil)
				for _, name := range []string{"Refactor the parser", "Fix flaky test", "Explain this error"} {
					kit.ListRow(c.Key(name), kit.ListRowOpt{Title: t(c, name), Selected: name == "Fix flaky test"})
				}
			},
			Center: func() {
				kit.EmptyState(c, t(c, "Conversation"), t(c, "Drag the dividers; narrow the window below 900 to get sheets."))
			},
			RightPanels: []kit.Panel{
				{ID: "outline", Title: t(c, "Outline"), Icon: kit.IconFiles, Content: func() { muted(t(c, "Files touched by the agent")) }},
				{ID: "info", Title: t(c, "Details"), Icon: kit.IconInfo, Content: func() { muted(t(c, "Model, usage and context")) }},
			},
			BottomPanels: []kit.Panel{
				{ID: "terminal", Title: t(c, "Terminal"), Icon: kit.IconTerminal, Content: func() { muted("$ go test ./...") }},
			},
			StatusLeft: func() {
				ui.Text(c, t(c, "Connected"))
				kit.Badge(c, t(c, "local"), kit.BadgeNeutral)
			},
			StatusRight: func() {
				kit.ProgressRing(c, .42, 14)
			},
		})
	}})
}
