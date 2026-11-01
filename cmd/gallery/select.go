package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	selectChoice    = "One Light"
	selectModel     = "claude-sonnet-5-5"
	selectReasoning = "auto"
)

func init() {
	register(page{name: "Select", order: 12, view: func(c *ui.Context) {
		pageBody(c, func() {
			options := []kit.SelectOption{{Value: "One Light", Label: "One Light"}, {Value: "One Dark", Label: "One Dark"}, {Value: "System", Label: t(c, "System")}}
			sectionTitle(c, t(c, "Toolbar select"))
			kit.Select(c, &selectChoice, kit.SelectOpt{Searchable: true, Options: options}).Label(t(c, "Appearance")).AlignSelf(ui.Start)
			sectionTitle(c, t(c, "Bordered select"))
			kit.Select(c.Key("bordered"), &selectChoice, kit.SelectOpt{Bordered: true, Options: options}).Label(t(c, "Theme")).Width(240)
			sectionTitle(c, t(c, "Descriptions and group headers"))
			models := []kit.SelectOption{
				{Label: "Anthropic", Header: true},
				{Value: "claude-sonnet-5-5", Label: "Claude Sonnet 5.5", Description: "Anthropic"},
				{Value: "claude-opus-5", Label: "Claude Opus 5", Description: "Anthropic"},
				{Label: "OpenAI", Header: true},
				{Value: "gpt-5", Label: "GPT-5", Description: "OpenAI"},
				{Value: "gpt-5-mini", Label: "GPT-5 mini", Description: "OpenAI", Disabled: true},
				{Label: t(c, "Local"), Header: true},
				{Value: "qwen3", Label: "Qwen3 32B", Description: "Ollama"},
			}
			kit.Select(c.Key("models"), &selectModel, kit.SelectOpt{Searchable: true, Options: models}).Label(t(c, "Model")).AlignSelf(ui.Start)
			sectionTitle(c, t(c, "Custom trigger text (options read just the level)"))
			levels := []kit.SelectOption{{Value: "auto", Label: t(c, "Auto")}, {Value: "low", Label: t(c, "Low")}, {Value: "medium", Label: t(c, "Medium")}, {Value: "high", Label: t(c, "High")}}
			kit.Select(c.Key("reasoning"), &selectReasoning, kit.SelectOpt{Options: levels, TriggerLabel: func(selected kit.SelectOption) string { return t(c, "Reasoning: ") + selected.Label }}).Label(t(c, "Reasoning")).AlignSelf(ui.Start)
		})
	}})
}
