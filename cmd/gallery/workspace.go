package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var workspaceState = kit.WorkspaceState{SidebarOpen: true, Right: kit.DockState{Open: true, ActivePanel: "outline"}, Bottom: kit.DockState{Open: true, ActivePanel: "terminal"}}

func init() {
	register(page{name: "Workspace", order: 29, full: true, view: func(c *ui.Context) {
		t := kit.GetTheme(c)
		muted := func(text string) { ui.Text(c, text).Padding(12).TextColor(t.Colors.Muted) }
		kit.Workspace(c, &workspaceState, kit.WorkspaceOpt{
			TitleBar: func() {
				if kit.IconButton(c.Key("back"), kit.IconChevronLeft, "Back to gallery").Clicked() {
					showPage("Overview")
				}
				ui.Text(c, "~/code/project")
			},
			Sidebar: func() {
				kit.Header(c, "CONVERSATIONS", nil)
				for _, name := range []string{"Refactor the parser", "Fix flaky test", "Explain this error"} {
					kit.ListRow(c.Key(name), kit.ListRowOpt{Title: name, Selected: name == "Fix flaky test"})
				}
			},
			Center: func() {
				kit.EmptyState(c, "Conversation", "Drag the dividers; narrow the window below 900 to get sheets.")
			},
			RightPanels: []kit.Panel{
				{ID: "outline", Title: "Outline", Icon: kit.IconFiles, Content: func() { muted("Files touched by the agent") }},
				{ID: "info", Title: "Details", Icon: kit.IconInfo, Content: func() { muted("Model, usage and context") }},
			},
			BottomPanels: []kit.Panel{
				{ID: "terminal", Title: "Terminal", Icon: kit.IconTerminal, Content: func() { muted("$ go test ./...") }},
			},
			StatusLeft: func() {
				ui.Text(c, "Connected")
				kit.Badge(c, "local", kit.BadgeNeutral)
			},
			StatusRight: func() {
				kit.ProgressRing(c, .42, 14)
			},
		})
	}})
}
