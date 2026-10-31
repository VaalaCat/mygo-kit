package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	inputValue     string
	invalidValue   = "not a valid name"
	readOnlyValue  = "provider-id"
	searchValue    string
	longLeading    = "/Users/example/projects/workspace/packages/application/internal/desktop/very/long/nested/directory/name"
	borderless     string
	areaBorderless string
)

func init() {
	register(page{name: "Input", order: 11, view: func(c *ui.Context) {
		pageBody(c, func() {
			sectionTitle(c, "TextInput")
			kit.TextInput(c.Key("gallery-input"), &inputValue).Label("Name").Placeholder("A project name")
			sectionTitle(c, "TextArea")
			kit.TextArea(c.Key("gallery-area"), &inputValue).Label("Description").Placeholder("Multi-line text and IME").Height(100)
			kit.Separator(c)
			sectionTitle(c, "Invalid")
			kit.TextInputWith(c.Key("gallery-invalid"), &invalidValue, kit.InputOpt{Invalid: true}).Label("Invalid name")
			// 只读值没有边线和底色，看起来就是一段可选中复制的文字；文字起点、28 高与普通输入一致，放进 Field 时与可编辑字段对齐。
			sectionTitle(c, "ReadOnly (a plain selectable value, aligned with editable fields)")
			kit.Field(c, kit.FieldOpt{Label: "Provider ID"}, func() ui.Element {
				return kit.TextInputWith(c.Key("gallery-readonly"), &readOnlyValue, kit.InputOpt{ReadOnly: true})
			})
			kit.Field(c, kit.FieldOpt{Label: "Name (editable, for comparison)"}, func() ui.Element {
				return kit.TextInput(c.Key("gallery-readonly-peer"), &inputValue).Placeholder("A project name")
			})
			sectionTitle(c, "Leading icon")
			kit.TextInputWith(c.Key("gallery-search"), &searchValue, kit.InputOpt{Leading: kit.IconSearch, Placeholder: "Search"}).Label("Search")
			sectionTitle(c, "Leading icon with long text (the text is clipped beside the icon)")
			kit.TextInputWith(c.Key("gallery-long"), &longLeading, kit.InputOpt{Leading: kit.IconFolder}).Label("Directory")
			sectionTitle(c, "Borderless (the surrounding box decides the look)")
			box := ui.Column(c).Gap(0).Padding(kit.GetTheme(c).Space(3)).Radius(kit.GetTheme(c).Metrics.RadiusLarge).Border(1, kit.GetTheme(c).Colors.Border)
			box.Children(func() {
				kit.TextAreaWith(c.Key("gallery-borderless-area"), &areaBorderless, kit.InputOpt{Borderless: true, Placeholder: "Message"}).Label("Message").MinHeight(kit.GetTheme(c).Space(12))
				kit.TextInputWith(c.Key("gallery-borderless"), &borderless, kit.InputOpt{Borderless: true, Leading: kit.IconSearch, Placeholder: "Filter"}).Label("Filter")
			})
		})
	}})
}
