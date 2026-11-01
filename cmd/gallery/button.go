package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Button", order: 10, view: func(c *ui.Context) {
		pageBody(c, func() {
			sectionTitle(c, t(c, "Actions"))
			ui.Row(c).Gap(8).Children(func() {
				kit.Button(c, t(c, "Ghost"))
				kit.PrimaryButton(c, t(c, "Primary"))
				kit.DangerButton(c, t(c, "Delete"))
				kit.Button(c, t(c, "Disabled")).Disabled(true)
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "Icon and tool buttons"))
			ui.Row(c).Gap(8).Children(func() {
				kit.IconButton(c.Key("icon-settings"), kit.IconSettings, t(c, "Settings"))
				kit.IconButton(c.Key("icon-plus"), kit.IconPlus, t(c, "Add"))
				kit.ToolButton(c.Key("tool-branch"), kit.IconBranch, "main", true)
				kit.ToolButton(c.Key("tool-folder"), kit.IconFolder, t(c, "Workspace"), false)
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "IconButtonTip (hover for the tooltip)"))
			ui.Row(c).Gap(8).Children(func() {
				kit.IconButtonTip(c.Key("tip-copy"), kit.IconCopy, t(c, "Copy"))
				kit.IconButtonTip(c.Key("tip-refresh"), kit.IconRefresh, t(c, "Refresh"))
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "Compact icon button (20)"))
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				kit.SmallIconButton(c.Key("small-copy"), kit.IconCopy, t(c, "Copy"))
				kit.SmallIconButton(c.Key("small-close"), kit.IconX, t(c, "Close"))
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "Flush buttons align with the text line above"))
			ui.Column(c).Gap(0).AlignItems(ui.Start).Children(func() {
				ui.Text(c, t(c, "Provider list")).TextColor(kit.GetTheme(c).Colors.Muted)
				kit.FlushButton(c.Key("flush-refresh"), t(c, "Refresh"))
				kit.FlushButton(c.Key("flush-disabled"), t(c, "Disabled")).Disabled(true)
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "Link"))
			kit.Link(c, t(c, "MyGo documentation"))
		})
	}})
}
