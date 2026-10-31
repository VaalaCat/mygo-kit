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
		pageBody(c, func() {
			if kit.Button(c.Key("open-palette"), "Open command palette").AlignSelf(ui.Start).Clicked() {
				commandOpen = true
			}
			ui.Text(c, "Last command: "+commandResult).TextColor(kit.GetTheme(c).Colors.Muted)
			if id := kit.CommandPalette(c, &commandOpen, []kit.Command{
				{ID: "conversation.new", Title: "New conversation", Group: "Conversation", Keys: "⌘N", Icon: kit.IconPlus},
				{ID: "conversation.archive", Title: "Archive conversation", Group: "Conversation", Keys: "⌘E"},
				{ID: "conversation.delete", Title: "Delete conversation", Group: "Conversation", Disabled: true},
				{ID: "view.sidebar", Title: "Toggle sidebar", Group: "View", Keys: "⌘B", Icon: kit.IconPanelLeft},
				{ID: "view.terminal", Title: "Toggle terminal", Group: "View", Keys: "⌘J", Icon: kit.IconTerminal},
				{ID: "app.settings", Title: "Open settings", Group: "Application", Keys: "⌘,", Icon: kit.IconSettings},
			}); id != "" {
				commandResult = id
			}
		})
	}})
}
