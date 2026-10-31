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
		t := kit.GetTheme(c)
		pageBody(c, func() {
			sectionTitle(c, "Field: label above, description or error below (16 between fields)")
			ui.Column(c).Width(360).Gap(t.Space(4)).Children(func() {
				kit.Field(c.Key("name"), kit.FieldOpt{Label: "Provider ID", Description: "Used in config.yaml and as the model prefix."}, func() ui.Element {
					return kit.TextInputWith(c.Key("name-input"), &fieldName, kit.InputOpt{Placeholder: "openai"})
				})
				kit.Field(c.Key("bad"), kit.FieldOpt{Label: "Name", Description: "Hidden while there is an error.", Error: "Only letters, digits and dashes."}, func() ui.Element {
					return kit.TextInputWith(c.Key("bad-input"), &fieldBadName, kit.InputOpt{Invalid: true})
				})
				kit.Field(c.Key("url"), kit.FieldOpt{Label: "Base URL"}, func() ui.Element {
					return kit.TextInputWith(c.Key("url-input"), &fieldBaseURL, kit.InputOpt{ReadOnly: true})
				})
				kit.Field(c.Key("protocol"), kit.FieldOpt{Label: "Protocol", Description: "Wire format of the provider."}, func() ui.Element {
					return kit.Select(c.Key("protocol-select"), &fieldProtocol, kit.SelectOpt{Bordered: true, Options: []kit.SelectOption{{Value: "openai-chat", Label: "OpenAI Chat"}, {Value: "anthropic", Label: "Anthropic"}}})
				})
				kit.Field(c.Key("inherit"), kit.FieldOpt{Label: "Inherit settings", Description: "Use the default provider's headers."}, func() ui.Element {
					return kit.Switch(c.Key("inherit-switch"), &fieldInherit, "")
				})
			})
			kit.Separator(c)
			sectionTitle(c, "FieldMessage")
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				kit.FieldMessage(c, kit.MessageDescription, "A description in muted text.")
				kit.FieldMessage(c, kit.MessageWarning, "A warning in the warning color.")
				kit.FieldMessage(c, kit.MessageError, "An error in the danger color.")
			})
			kit.Separator(c)
			sectionTitle(c, "InlineHint")
			ui.Column(c).Width(280).Background(t.Colors.Panel).Children(func() {
				kit.InlineHint(c.Key("loading"), "Loading…", true)
				kit.InlineHint(c.Key("empty"), "No conversations yet", false)
			})
		})
	}})
	register(page{name: "Settings form", order: 53, view: func(c *ui.Context) {
		t := kit.GetTheme(c)
		options := []kit.SelectOption{{Value: "One Light", Label: "One Light"}, {Value: "One Dark", Label: "One Dark"}}
		kit.SettingPage(c, "Settings", func() {
			kit.Button(c.Key("revert"), "Revert")
			kit.PrimaryButton(c.Key("apply"), "Apply")
		}, func() {
			kit.SettingSection(c.Key("general"), "General", func() {
				kit.SettingRow(c.Key("appearance"), "Appearance", "System follows the operating system.", func() {
					kit.Segmented(c.Key("appearance-segmented"), &settingsAppearance, []string{"System", "Light", "Dark"}).Label("Appearance")
				})
				kit.SettingRow(c.Key("theme"), "Color theme", "", func() {
					kit.Select(c.Key("theme-select"), &settingsTheme, kit.SelectOpt{Bordered: true, Options: options}).Label("Color theme").Width(t.Space(40))
				})
				kit.SettingRow(c.Key("motion"), "Reduce motion", "Turns off transitions and spinners.", func() {
					kit.Switch(c.Key("motion-switch"), &settingsReduceMotion, "").Label("Reduce motion")
				})
				kit.SettingRow(c.Key("notify"), "Notify when a turn finishes", "", func() {
					kit.Switch(c.Key("notify-switch"), &settingsNotify, "").Label("Notify")
				})
			})
			kit.SettingSection(c.Key("provider"), "Model provider", func() {
				ui.Column(c).Padding(t.Space(3)).Gap(t.Space(4)).Children(func() {
					kit.Field(c.Key("provider-id"), kit.FieldOpt{Label: "Provider ID"}, func() ui.Element {
						return kit.TextInputWith(c.Key("provider-id-input"), &settingsProviderID, kit.InputOpt{ReadOnly: true})
					})
					kit.Field(c.Key("provider-url"), kit.FieldOpt{Label: "Base URL", Description: "Include the /v1 suffix."}, func() ui.Element {
						return kit.TextInputWith(c.Key("provider-url-input"), &settingsProviderURL, kit.InputOpt{})
					})
					kit.Field(c.Key("provider-protocol"), kit.FieldOpt{Label: "Protocol"}, func() ui.Element {
						return kit.Select(c.Key("provider-protocol-select"), &settingsProvider, kit.SelectOpt{Bordered: true, Options: []kit.SelectOption{{Value: "openai-chat", Label: "OpenAI Chat"}, {Value: "anthropic", Label: "Anthropic"}}})
					})
				})
			})
		})
	}})
}
