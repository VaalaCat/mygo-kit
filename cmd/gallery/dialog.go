package main

import (
	"strings"

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
			sectionTitle(c, t(c, "Dialog"))
			if kit.Button(c, t(c, "Open dialog")).AlignSelf(ui.Start).Clicked() {
				dialogOpen = true
			}
			kit.Dialog(c, &dialogOpen, t(c, "Native dialog"), func() {
				ui.Text(c, t(c, "Focus stays inside this dialog. Escape returns it to the trigger."))
				if kit.PrimaryButton(c, t(c, "Done")).Clicked() {
					dialogOpen = false
				}
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "DialogWith (cancel + confirm, Enter confirms, Esc cancels)"))
			ui.Row(c).Gap(8).Children(func() {
				if kit.Button(c.Key("open-confirm"), t(c, "Confirm")).Clicked() {
					confirmOpen = true
				}
				if kit.Button(c.Key("open-danger"), t(c, "Danger")).Clicked() {
					dangerOpen = true
				}
				if kit.Button(c.Key("open-rename"), t(c, "With input")).Clicked() {
					renameOpen = true
				}
				if kit.Button(c.Key("open-tabs"), t(c, "With tabs")).Clicked() {
					tabsDialogOpen = true
				}
				if kit.Button(c.Key("open-readonly"), t(c, "Read-only values")).Clicked() {
					readOnlyDialogOpen = true
				}
			})
			dialogResultLabel := dialogResult
			switch dialogResult {
			case "(none)":
				dialogResultLabel = t(c, "(none)")
			case "compact":
				dialogResultLabel = t(c, "compact")
			case "delete":
				dialogResultLabel = t(c, "delete")
			default:
				if title, ok := strings.CutPrefix(dialogResult, "rename: "); ok {
					dialogResultLabel = t(c, "rename: ") + title
				}
			}
			ui.Text(c, t(c, "Last result: ")+dialogResultLabel).TextColor(kit.GetTheme(c).Colors.Muted)
			if kit.DialogWith(c.Key("confirm"), &confirmOpen, kit.DialogOpt{Title: t(c, "Compact context"), Description: t(c, "Older turns are summarized to free up the context window."), ConfirmLabel: t(c, "Compact")}, func() {}) {
				dialogResult = "compact"
			}
			if kit.DialogWith(c.Key("danger"), &dangerOpen, kit.DialogOpt{Title: t(c, "Delete conversation"), Description: t(c, "This cannot be undone."), ConfirmLabel: t(c, "Delete"), Danger: true}, func() {}) {
				dialogResult = "delete"
			}
			renameError := ""
			if renameText == "" {
				renameError = t(c, "A title is required.")
			}
			if kit.DialogWith(c.Key("rename"), &renameOpen, kit.DialogOpt{Title: t(c, "Rename"), ConfirmLabel: t(c, "Save"), ConfirmDisabled: renameText == ""}, func() {
				kit.Field(c, kit.FieldOpt{Label: t(c, "Title"), Error: renameError}, func() ui.Element {
					return kit.TextInputWith(c.Key("rename-input"), &renameText, kit.InputOpt{Invalid: renameText == ""})
				})
			}) {
				dialogResult = "rename: " + renameText
			}
			// 正文是页签加一个输入框：打开时焦点落在输入框上，页签不显示焦点环。
			kit.DialogWith(c.Key("tabs"), &tabsDialogOpen, kit.DialogOpt{Title: t(c, "Settings"), ConfirmLabel: t(c, "Save")}, func() {
				kit.Tabs(c, &dialogTab, []string{t(c, "General"), t(c, "Advanced")})
				kit.Field(c, kit.FieldOpt{Label: t(c, "Title")}, func() ui.Element {
					return kit.TextInput(c.Key("tabs-input"), &renameText)
				})
			})
			// 正文没有可编辑输入：打开时焦点不落在任何可见控件上，Enter 仍然确认。
			kit.DialogWith(c.Key("readonly"), &readOnlyDialogOpen, kit.DialogOpt{Title: t(c, "Provider details"), ConfirmLabel: t(c, "Close")}, func() {
				kit.Field(c, kit.FieldOpt{Label: t(c, "Provider ID")}, func() ui.Element {
					return kit.TextInputWith(c.Key("provider-id"), &readOnlyValue, kit.InputOpt{ReadOnly: true})
				})
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "Popover and tooltip"))
			ui.Row(c).Gap(8).Children(func() {
				anchor := kit.Button(c.Key("popover"), t(c, "Popover"))
				if anchor.Clicked() {
					popoverOpen = !popoverOpen
				}
				kit.Popover(c, anchor, &popoverOpen, func() { ui.Text(c, t(c, "Anchored content")) })
				kit.Tooltip(c, kit.Button(c.Key("tooltip"), t(c, "Hover for tooltip")), t(c, "A short hint"))
			})
		})
	}})
}
