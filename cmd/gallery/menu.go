package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var menuLast = "(none)"

func init() {
	register(page{name: "Menu", order: 54, view: func(c *ui.Context) {
		t := kit.GetTheme(c)
		pick := func(name string) func() { return func() { menuLast = name } }
		items := []kit.MenuItem{
			{Icon: kit.IconPencil, Title: "Rename", Shortcut: "F2", OnSelect: pick("Rename")},
			{Icon: kit.IconGitFork, Title: "Fork", OnSelect: pick("Fork")},
			{Icon: kit.IconCopy, Title: "Export", Disabled: true, OnSelect: pick("Export")},
			{Separator: true},
			{Icon: kit.IconArchive, Title: "Archive", Shortcut: "⌘E", OnSelect: pick("Archive")},
			{Icon: kit.IconX, Title: "Delete", Danger: true, OnSelect: pick("Delete")},
		}
		pageBody(c, func() {
			sectionTitle(c, "MenuButton (Up/Down, Enter, Esc, first letter)")
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				kit.MenuButton(c, kit.IconButton(c.Key("menu-more"), kit.IconEllipsis, "More"), items)
				kit.MenuButton(c, kit.ToolButton(c.Key("menu-tool"), kit.IconFolder, "Workspace", true), []kit.MenuItem{
					{Title: "Local", OnSelect: pick("Local")}, {Title: "Remote", OnSelect: pick("Remote")}, {Title: "Sandbox", OnSelect: pick("Sandbox")},
				})
			})
			ui.Text(c, "Last chosen: "+menuLast).TextColor(t.Colors.Muted)
			kit.Separator(c)
			sectionTitle(c, "ContextMenu (right click, Shift+F10 or the Menu key while focused)")
			area := ui.Column(c.Key("context-area")).Height(120).Width(360).Center().Radius(t.Metrics.RadiusLarge).Border(1, t.Colors.Border).Focusable().Children(func() {
				ui.Text(c, "Right click here").TextColor(t.Colors.Muted)
			})
			kit.ContextMenu(c, area, items)
			kit.Separator(c)
			sectionTitle(c, "ContextMenu on Selectable text (drag to select, then right click: Copy is built in)")
			text := ui.Column(c.Key("context-text")).Width(360).Gap(4).Padding(12).Radius(t.Metrics.RadiusLarge).Border(1, t.Colors.Border).Selectable().Children(func() {
				ui.Text(c, "Select any part of this text, then right click.")
				ui.Text(c, "The menu is kit's own, not the system's.").TextColor(t.Colors.Muted)
			})
			kit.ContextMenu(c, text, items)
		})
	}})
}
