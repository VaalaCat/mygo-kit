package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	dialogOpen, popoverOpen             bool
	confirmOpen, dangerOpen, renameOpen bool
	tabsDialogOpen, readOnlyDialogOpen  bool
	dialogTab                           int
	renameText                          = "Refactor the parser"
	dialogResult                        = "(none)"
)

func init() {
	register(page{name: "Dialog", order: 17, view: func(c *ui.Context) {
		pageBody(c, func() {
			sectionTitle(c, "Dialog")
			if kit.Button(c, "Open dialog").AlignSelf(ui.Start).Clicked() {
				dialogOpen = true
			}
			kit.Dialog(c, &dialogOpen, "Native dialog", func() {
				ui.Text(c, "Focus stays inside this dialog. Escape returns it to the trigger.")
				if kit.PrimaryButton(c, "Done").Clicked() {
					dialogOpen = false
				}
			})
			kit.Separator(c)
			sectionTitle(c, "DialogWith (cancel + confirm, Enter confirms, Esc cancels)")
			ui.Row(c).Gap(8).Children(func() {
				if kit.Button(c.Key("open-confirm"), "Confirm").Clicked() {
					confirmOpen = true
				}
				if kit.Button(c.Key("open-danger"), "Danger").Clicked() {
					dangerOpen = true
				}
				if kit.Button(c.Key("open-rename"), "With input").Clicked() {
					renameOpen = true
				}
				if kit.Button(c.Key("open-tabs"), "With tabs").Clicked() {
					tabsDialogOpen = true
				}
				if kit.Button(c.Key("open-readonly"), "Read-only values").Clicked() {
					readOnlyDialogOpen = true
				}
			})
			ui.Text(c, "Last result: "+dialogResult).TextColor(kit.GetTheme(c).Colors.Muted)
			if kit.DialogWith(c.Key("confirm"), &confirmOpen, kit.DialogOpt{Title: "Compact context", Description: "Older turns are summarized to free up the context window.", ConfirmLabel: "Compact"}, func() {}) {
				dialogResult = "compact"
			}
			if kit.DialogWith(c.Key("danger"), &dangerOpen, kit.DialogOpt{Title: "Delete conversation", Description: "This cannot be undone.", ConfirmLabel: "Delete", Danger: true}, func() {}) {
				dialogResult = "delete"
			}
			renameError := ""
			if renameText == "" {
				renameError = "A title is required."
			}
			if kit.DialogWith(c.Key("rename"), &renameOpen, kit.DialogOpt{Title: "Rename", ConfirmLabel: "Save", ConfirmDisabled: renameText == ""}, func() {
				kit.Field(c, kit.FieldOpt{Label: "Title", Error: renameError}, func() ui.Element {
					return kit.TextInputWith(c.Key("rename-input"), &renameText, kit.InputOpt{Invalid: renameText == ""})
				})
			}) {
				dialogResult = "rename: " + renameText
			}
			// 正文是页签加一个输入框：打开时焦点落在输入框上，页签不显示焦点环。
			kit.DialogWith(c.Key("tabs"), &tabsDialogOpen, kit.DialogOpt{Title: "Settings", ConfirmLabel: "Save"}, func() {
				kit.Tabs(c, &dialogTab, []string{"General", "Advanced"})
				kit.Field(c, kit.FieldOpt{Label: "Title"}, func() ui.Element {
					return kit.TextInput(c.Key("tabs-input"), &renameText)
				})
			})
			// 正文没有可编辑输入：打开时焦点不落在任何可见控件上，Enter 仍然确认。
			kit.DialogWith(c.Key("readonly"), &readOnlyDialogOpen, kit.DialogOpt{Title: "Provider details", ConfirmLabel: "Close"}, func() {
				kit.Field(c, kit.FieldOpt{Label: "Provider ID"}, func() ui.Element {
					return kit.TextInputWith(c.Key("provider-id"), &readOnlyValue, kit.InputOpt{ReadOnly: true})
				})
			})
			kit.Separator(c)
			sectionTitle(c, "Popover and tooltip")
			ui.Row(c).Gap(8).Children(func() {
				anchor := kit.Button(c.Key("popover"), "Popover")
				if anchor.Clicked() {
					popoverOpen = !popoverOpen
				}
				kit.Popover(c, anchor, &popoverOpen, func() { ui.Text(c, "Anchored content") })
				kit.Tooltip(c, kit.Button(c.Key("tooltip"), "Hover for tooltip"), "A short hint")
			})
		})
	}})
}
