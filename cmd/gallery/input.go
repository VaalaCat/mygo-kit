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
			sectionTitle(c, t(c, "TextInput"))
			kit.TextInput(c.Key("gallery-input"), &inputValue).Label(t(c, "Name")).Placeholder(t(c, "A project name"))
			sectionTitle(c, t(c, "TextArea"))
			kit.TextArea(c.Key("gallery-area"), &inputValue).Label(t(c, "Description")).Placeholder(t(c, "Multi-line text and IME")).Height(100)
			kit.Separator(c)
			sectionTitle(c, t(c, "Invalid"))
			kit.TextInputWith(c.Key("gallery-invalid"), &invalidValue, kit.InputOpt{Invalid: true}).Label(t(c, "Invalid name"))
			// 只读值没有边线和底色，看起来就是一段可选中复制的文字；文字起点、28 高与普通输入一致，放进 Field 时与可编辑字段对齐。
			sectionTitle(c, t(c, "ReadOnly (a plain selectable value, aligned with editable fields)"))
			kit.Field(c, kit.FieldOpt{Label: t(c, "Provider ID")}, func() ui.Element {
				return kit.TextInputWith(c.Key("gallery-readonly"), &readOnlyValue, kit.InputOpt{ReadOnly: true})
			})
			kit.Field(c, kit.FieldOpt{Label: t(c, "Name (editable, for comparison)")}, func() ui.Element {
				return kit.TextInput(c.Key("gallery-readonly-peer"), &inputValue).Placeholder(t(c, "A project name"))
			})
			sectionTitle(c, t(c, "Leading icon"))
			kit.TextInputWith(c.Key("gallery-search"), &searchValue, kit.InputOpt{Leading: kit.IconSearch, Placeholder: t(c, "Search")}).Label(t(c, "Search"))
			sectionTitle(c, t(c, "Leading icon with long text (the text is clipped beside the icon)"))
			kit.TextInputWith(c.Key("gallery-long"), &longLeading, kit.InputOpt{Leading: kit.IconFolder}).Label(t(c, "Directory"))
			sectionTitle(c, t(c, "Borderless (the surrounding box decides the look)"))
			box := ui.Column(c).Gap(0).Padding(kit.GetTheme(c).Space(3)).Radius(kit.GetTheme(c).Metrics.RadiusLarge).Border(1, kit.GetTheme(c).Colors.Border)
			box.Children(func() {
				kit.TextAreaWith(c.Key("gallery-borderless-area"), &areaBorderless, kit.InputOpt{Borderless: true, Placeholder: t(c, "Message")}).Label(t(c, "Message")).MinHeight(kit.GetTheme(c).Space(12))
				kit.TextInputWith(c.Key("gallery-borderless"), &borderless, kit.InputOpt{Borderless: true, Leading: kit.IconSearch, Placeholder: t(c, "Filter")}).Label(t(c, "Filter"))
			})
		})
	}})
}
