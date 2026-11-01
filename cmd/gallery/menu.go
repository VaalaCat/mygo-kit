package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var menuLast = "(none)"

func init() {
	register(page{name: "Menu", order: 54, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		pick := func(name string) func() { return func() { menuLast = name } }
		items := []kit.MenuItem{
			{Icon: kit.IconPencil, Title: t(c, "Rename"), Shortcut: "F2", OnSelect: pick("Rename")},
			{Icon: kit.IconGitFork, Title: t(c, "Fork"), OnSelect: pick("Fork")},
			{Icon: kit.IconCopy, Title: t(c, "Export"), Disabled: true, OnSelect: pick("Export")},
			{Separator: true},
			{Icon: kit.IconArchive, Title: t(c, "Archive"), Shortcut: "⌘E", OnSelect: pick("Archive")},
			{Icon: kit.IconX, Title: t(c, "Delete"), Danger: true, OnSelect: pick("Delete")},
		}
		pageBody(c, func() {
			sectionTitle(c, t(c, "MenuButton (Up/Down, Enter, Esc, first letter)"))
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				kit.MenuButton(c, kit.IconButton(c.Key("menu-more"), kit.IconEllipsis, t(c, "More")), items)
				kit.MenuButton(c, kit.ToolButton(c.Key("menu-tool"), kit.IconFolder, t(c, "Workspace"), true), []kit.MenuItem{
					{Title: t(c, "Local"), OnSelect: pick("Local")}, {Title: t(c, "Remote"), OnSelect: pick("Remote")}, {Title: t(c, "Sandbox"), OnSelect: pick("Sandbox")},
				})
			})
			ui.Text(c, t(c, "Last chosen: ")+t(c, menuLast)).TextColor(theme.Colors.Muted)
			kit.Separator(c)
			sectionTitle(c, t(c, "ContextMenu (right click, Shift+F10 or the Menu key while focused)"))
			area := ui.Column(c.Key("context-area")).Height(120).Width(360).Center().Radius(theme.Metrics.RadiusLarge).Border(1, theme.Colors.Border).Focusable().Children(func() {
				ui.Text(c, t(c, "Right click here")).TextColor(theme.Colors.Muted)
			})
			kit.ContextMenu(c, area, items)
			kit.Separator(c)
			sectionTitle(c, t(c, "ContextMenu on Selectable text (drag to select, then right click: Copy is built in)"))
			text := ui.Column(c.Key("context-text")).Width(360).Gap(4).Padding(12).Radius(theme.Metrics.RadiusLarge).Border(1, theme.Colors.Border).Selectable().Children(func() {
				ui.Text(c, t(c, "Select any part of this text, then right click."))
				ui.Text(c, t(c, "The menu is kit's own, not the system's.")).TextColor(theme.Colors.Muted)
			})
			kit.ContextMenu(c, text, items)
		})
	}})
}
