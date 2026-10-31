package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Button", order: 10, view: func(c *ui.Context) {
		pageBody(c, func() {
			sectionTitle(c, "Actions")
			ui.Row(c).Gap(8).Children(func() {
				kit.Button(c, "Ghost")
				kit.PrimaryButton(c, "Primary")
				kit.DangerButton(c, "Delete")
				kit.Button(c, "Disabled").Disabled(true)
			})
			kit.Separator(c)
			sectionTitle(c, "Icon and tool buttons")
			ui.Row(c).Gap(8).Children(func() {
				kit.IconButton(c.Key("icon-settings"), kit.IconSettings, "Settings")
				kit.IconButton(c.Key("icon-plus"), kit.IconPlus, "Add")
				kit.ToolButton(c.Key("tool-branch"), kit.IconBranch, "main", true)
				kit.ToolButton(c.Key("tool-folder"), kit.IconFolder, "Workspace", false)
			})
			kit.Separator(c)
			sectionTitle(c, "IconButtonTip (hover for the tooltip)")
			ui.Row(c).Gap(8).Children(func() {
				kit.IconButtonTip(c.Key("tip-copy"), kit.IconCopy, "Copy")
				kit.IconButtonTip(c.Key("tip-refresh"), kit.IconRefresh, "Refresh")
			})
			kit.Separator(c)
			sectionTitle(c, "Compact icon button (20)")
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				kit.SmallIconButton(c.Key("small-copy"), kit.IconCopy, "Copy")
				kit.SmallIconButton(c.Key("small-close"), kit.IconX, "Close")
			})
			kit.Separator(c)
			sectionTitle(c, "Flush buttons align with the text line above")
			ui.Column(c).Gap(0).AlignItems(ui.Start).Children(func() {
				ui.Text(c, "Provider list").TextColor(kit.GetTheme(c).Colors.Muted)
				kit.FlushButton(c.Key("flush-refresh"), "Refresh")
				kit.FlushButton(c.Key("flush-disabled"), "Disabled").Disabled(true)
			})
			kit.Separator(c)
			sectionTitle(c, "Link")
			kit.Link(c, "MyGo documentation")
		})
	}})
}
