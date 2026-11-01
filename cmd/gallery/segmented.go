package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var segmentedAppearance, segmentedView, segmentedSections, segmentedFull, segmentedDisabled = 0, 1, 0, 1, 0

func init() {
	register(page{name: "Segmented", order: 51, view: func(c *ui.Context) {
		pageBody(c, func() {
			sectionTitle(c, t(c, "Text only, hugging content (arrows move the choice)"))
			kit.Segmented(c.Key("appearance"), &segmentedAppearance, []string{t(c, "System"), t(c, "Light"), t(c, "Dark")}).Label(t(c, "Appearance")).AlignSelf(ui.Start)
			sectionTitle(c, t(c, "Two segments, second selected"))
			kit.Segmented(c.Key("view"), &segmentedView, []string{t(c, "Conversations"), t(c, "Archive")}).Label(t(c, "View")).AlignSelf(ui.Start)
			sectionTitle(c, t(c, "With icons and tooltips"))
			sections := []kit.SegmentedItem{
				{Icon: kit.IconFolder, Label: t(c, "Workspace"), Tooltip: t(c, "Executor and directory")},
				{Icon: kit.IconServerCog, Label: "MCP", Tooltip: t(c, "Connected MCP servers")},
				{Icon: kit.IconBookOpen, Label: t(c, "Skill"), Tooltip: t(c, "Skills this conversation can use")},
			}
			kit.SegmentedWith(c.Key("sections"), &segmentedSections, sections, kit.SegmentedOpt{}).Label(t(c, "Sections")).AlignSelf(ui.Start)
			sectionTitle(c, t(c, "Full width, equal segments"))
			kit.SegmentedWith(c.Key("full"), &segmentedFull, sections, kit.SegmentedOpt{FullWidth: true}).Label(t(c, "Full width"))
			sectionTitle(c, t(c, "Disabled"))
			kit.Segmented(c.Key("disabled"), &segmentedDisabled, []string{t(c, "One"), t(c, "Two"), t(c, "Three")}).Label(t(c, "Disabled")).Disabled(true).AlignSelf(ui.Start)
		})
	}})
}
