package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	fieldName, fieldBadName, fieldBaseURL = "my-provider", "not a valid id", "https://api.example.com/v1"
	fieldProtocol                         = "openai-chat"
	fieldInherit                          = true

	settingsAppearance                                        = 0
	settingsReduceMotion, settingsNotify                      bool
	settingsTheme                                             = "One Light"
	settingsProviderID, settingsProviderURL, settingsProvider = "openai", "https://api.openai.com/v1", "openai-chat"
)

func init() {
	register(page{name: "Field", order: 52, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		pageBody(c, func() {
			sectionTitle(c, t(c, "Field: label above, description or error below (16 between fields)"))
			ui.Column(c).Width(360).Gap(theme.Space(4)).Children(func() {
				kit.Field(c.Key("name"), kit.FieldOpt{Label: t(c, "Provider ID"), Description: t(c, "Used in config.yaml and as the model prefix.")}, func() ui.Element {
					return kit.TextInputWith(c.Key("name-input"), &fieldName, kit.InputOpt{Placeholder: "openai"})
				})
				kit.Field(c.Key("bad"), kit.FieldOpt{Label: t(c, "Name"), Description: t(c, "Hidden while there is an error."), Error: t(c, "Only letters, digits and dashes.")}, func() ui.Element {
					return kit.TextInputWith(c.Key("bad-input"), &fieldBadName, kit.InputOpt{Invalid: true})
				})
				kit.Field(c.Key("url"), kit.FieldOpt{Label: t(c, "Base URL")}, func() ui.Element {
					return kit.TextInputWith(c.Key("url-input"), &fieldBaseURL, kit.InputOpt{ReadOnly: true})
				})
				kit.Field(c.Key("protocol"), kit.FieldOpt{Label: t(c, "Protocol"), Description: t(c, "Wire format of the provider.")}, func() ui.Element {
					return kit.Select(c.Key("protocol-select"), &fieldProtocol, kit.SelectOpt{Bordered: true, Options: []kit.SelectOption{{Value: "openai-chat", Label: "OpenAI Chat"}, {Value: "anthropic", Label: "Anthropic"}}})
				})
				kit.Field(c.Key("inherit"), kit.FieldOpt{Label: t(c, "Inherit settings"), Description: t(c, "Use the default provider's headers.")}, func() ui.Element {
					return kit.Switch(c.Key("inherit-switch"), &fieldInherit, "")
				})
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "FieldMessage"))
			ui.Column(c).Gap(theme.Space(1)).Children(func() {
				kit.FieldMessage(c, kit.MessageDescription, t(c, "A description in muted text."))
				kit.FieldMessage(c, kit.MessageWarning, t(c, "A warning in the warning color."))
				kit.FieldMessage(c, kit.MessageError, t(c, "An error in the danger color."))
			})
			kit.Separator(c)
			sectionTitle(c, t(c, "InlineHint"))
			ui.Column(c).Width(280).Background(theme.Colors.Panel).Children(func() {
				kit.InlineHint(c.Key("loading"), t(c, "Loading…"), true)
				kit.InlineHint(c.Key("empty"), t(c, "No conversations yet"), false)
			})
		})
	}})
	register(page{name: "Settings form", order: 53, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		options := []kit.SelectOption{{Value: "One Light", Label: "One Light"}, {Value: "One Dark", Label: "One Dark"}}
		kit.SettingPage(c, t(c, "Settings"), func() {
			kit.Button(c.Key("revert"), t(c, "Revert"))
			kit.PrimaryButton(c.Key("apply"), t(c, "Apply"))
		}, func() {
			kit.SettingSection(c.Key("general"), t(c, "General"), func() {
				kit.SettingRow(c.Key("appearance"), t(c, "Appearance"), t(c, "System follows the operating system."), func() {
					kit.Segmented(c.Key("appearance-segmented"), &settingsAppearance, []string{t(c, "System"), t(c, "Light"), t(c, "Dark")}).Label(t(c, "Appearance"))
				})
				kit.SettingRow(c.Key("theme"), t(c, "Color theme"), "", func() {
					kit.Select(c.Key("theme-select"), &settingsTheme, kit.SelectOpt{Bordered: true, Options: options}).Label(t(c, "Color theme")).Width(theme.Space(40))
				})
				kit.SettingRow(c.Key("motion"), t(c, "Reduce motion"), t(c, "Turns off transitions and spinners."), func() {
					kit.Switch(c.Key("motion-switch"), &settingsReduceMotion, "").Label(t(c, "Reduce motion"))
				})
				kit.SettingRow(c.Key("notify"), t(c, "Notify when a turn finishes"), "", func() {
					kit.Switch(c.Key("notify-switch"), &settingsNotify, "").Label(t(c, "Notify"))
				})
			})
			kit.SettingSection(c.Key("provider"), t(c, "Model provider"), func() {
				ui.Column(c).Padding(theme.Space(3)).Gap(theme.Space(4)).Children(func() {
					kit.Field(c.Key("provider-id"), kit.FieldOpt{Label: t(c, "Provider ID")}, func() ui.Element {
						return kit.TextInputWith(c.Key("provider-id-input"), &settingsProviderID, kit.InputOpt{ReadOnly: true})
					})
					kit.Field(c.Key("provider-url"), kit.FieldOpt{Label: t(c, "Base URL"), Description: t(c, "Include the /v1 suffix.")}, func() ui.Element {
						return kit.TextInputWith(c.Key("provider-url-input"), &settingsProviderURL, kit.InputOpt{})
					})
					kit.Field(c.Key("provider-protocol"), kit.FieldOpt{Label: t(c, "Protocol")}, func() ui.Element {
						return kit.Select(c.Key("provider-protocol-select"), &settingsProvider, kit.SelectOpt{Bordered: true, Options: []kit.SelectOption{{Value: "openai-chat", Label: "OpenAI Chat"}, {Value: "anthropic", Label: "Anthropic"}}})
					})
				})
			})
		})
	}})
}
