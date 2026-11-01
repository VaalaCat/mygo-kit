package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var listRowSelected = "Refactor the parser"

func init() {
	register(page{name: "ListRow", order: 16, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		list := func(title string, rows func()) {
			ui.Column(c.Key(title)).Gap(0).Width(300).Children(func() {
				kit.GetTheme(c).TextRole(ui.Text(c, t(c, title)), kit.TextCaption).TextColor(theme.Colors.Muted).Padding(0, 0, 4, 0)
				ui.Column(c).Gap(0).Background(theme.Colors.Panel).Border(1, theme.Colors.Border).Children(rows)
			})
		}
		row := func(opt kit.ListRowOpt) {
			renderOpt := opt
			if opt.Icon != kit.IconFolder && opt.Icon != kit.IconFileText && opt.Title != "local · 001" {
				renderOpt.Title = t(c, opt.Title)
				if opt.Description != "" {
					renderOpt.Description = t(c, opt.Description)
				}
			}
			if kit.ListRow(c.Key(opt.Title+opt.Description), renderOpt).Clicked() {
				listRowSelected = opt.Title
			}
		}
		pageBody(c, func() {
			ui.Row(c).Gap(24).AlignItems(ui.Start).Children(func() {
				list("Single line · icon · selection", func() {
					for _, name := range []string{"Refactor the parser", "Fix flaky test", "Explain this error"} {
						row(kit.ListRowOpt{Icon: kit.IconMessageSquare, Title: name, Selected: name == listRowSelected})
					}
				})
				list("Two lines (44) · trailing", func() {
					row(kit.ListRowOpt{Title: "Add retry to the client", Description: "openvaala · 3 minutes ago", Trailing: func() { kit.Spinner(c, kit.IconSizeDefault) }})
					row(kit.ListRowOpt{Title: "Update the changelog", Description: "docs · yesterday", Trailing: func() { kit.StatusDot(c, theme.Colors.Accent).Label(t(c, "Unread")) }})
					row(kit.ListRowOpt{Title: "Draft release notes", Description: "docs · 2 days ago"})
				})
			})
			ui.Row(c).Gap(24).AlignItems(ui.Start).Children(func() {
				list("Depth (16 per level)", func() {
					row(kit.ListRowOpt{Icon: kit.IconFolder, Title: "internal"})
					row(kit.ListRowOpt{Icon: kit.IconFolder, Title: "desktop", Depth: 1})
					row(kit.ListRowOpt{Icon: kit.IconFileText, Title: "composer.go", Depth: 2})
					row(kit.ListRowOpt{Icon: kit.IconFileText, Title: "settings.go", Depth: 2})
				})
				list("Actions on hover / selected / focus", func() {
					for _, name := range []string{"Hover this row", "Or select it", "Tab into it"} {
						row(kit.ListRowOpt{Icon: kit.IconBot, Title: name, Selected: name == listRowSelected, Trailing: func() { kit.Badge(c, t(c, "local"), kit.BadgeNeutral) }, Actions: func() {
							kit.IconButtonTip(c.Key("more-"+name), kit.IconEllipsis, t(c, "More"))
						}})
					}
				})
			})
			list("Trailing text on the title baseline · disabled", func() {
				row(kit.ListRowOpt{Icon: kit.IconTerminal, Title: "local · 001", Trailing: func() {
					theme.TextRole(ui.Text(c, t(c, "running")), kit.TextCaption).TextColor(theme.Colors.Muted)
				}})
				kit.ListRow(c.Key("disabled"), kit.ListRowOpt{Icon: kit.IconBan, Title: t(c, "Unavailable")}).Disabled(true)
			})
		})
	}})
}
