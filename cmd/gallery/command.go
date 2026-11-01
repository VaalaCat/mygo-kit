package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	commandOpen   bool
	commandResult = "(none)"
)

func init() {
	register(page{name: "CommandPalette", order: 28, view: func(c *ui.Context) {
		commandResultLabel := commandResult
		if commandResultLabel == "(none)" {
			commandResultLabel = t(c, "(none)")
		}
		pageBody(c, func() {
			if kit.Button(c.Key("open-palette"), t(c, "Open command palette")).AlignSelf(ui.Start).Clicked() {
				commandOpen = true
			}
			ui.Text(c, t(c, "Last command: ")+commandResultLabel).TextColor(kit.GetTheme(c).Colors.Muted)
			if id := kit.CommandPalette(c, &commandOpen, []kit.Command{
				{ID: "conversation.new", Title: t(c, "New conversation"), Group: t(c, "Conversation"), Keys: "⌘N", Icon: kit.IconPlus},
				{ID: "conversation.archive", Title: t(c, "Archive conversation"), Group: t(c, "Conversation"), Keys: "⌘E"},
				{ID: "conversation.delete", Title: t(c, "Delete conversation"), Group: t(c, "Conversation"), Disabled: true},
				{ID: "view.sidebar", Title: t(c, "Toggle sidebar"), Group: t(c, "View"), Keys: "⌘B", Icon: kit.IconPanelLeft},
				{ID: "view.terminal", Title: t(c, "Toggle terminal"), Group: t(c, "View"), Keys: "⌘J", Icon: kit.IconTerminal},
				{ID: "app.settings", Title: t(c, "Open settings"), Group: t(c, "Application"), Keys: "⌘,", Icon: kit.IconSettings},
			}); id != "" {
				commandResult = id
			}
		})
	}})
}
