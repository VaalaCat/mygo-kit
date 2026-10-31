package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var segmentedAppearance, segmentedView, segmentedSections, segmentedFull, segmentedDisabled = 0, 1, 0, 1, 0

func init() {
	register(page{name: "Segmented", order: 51, view: func(c *ui.Context) {
		pageBody(c, func() {
			sectionTitle(c, "Text only, hugging content (arrows move the choice)")
			kit.Segmented(c.Key("appearance"), &segmentedAppearance, []string{"System", "Light", "Dark"}).Label("Appearance").AlignSelf(ui.Start)
			sectionTitle(c, "Two segments, second selected")
			kit.Segmented(c.Key("view"), &segmentedView, []string{"Conversations", "Archive"}).Label("View").AlignSelf(ui.Start)
			sectionTitle(c, "With icons and tooltips")
			sections := []kit.SegmentedItem{
				{Icon: kit.IconFolder, Label: "Workspace", Tooltip: "Executor and directory"},
				{Icon: kit.IconServerCog, Label: "MCP", Tooltip: "Connected MCP servers"},
				{Icon: kit.IconBookOpen, Label: "Skill", Tooltip: "Skills this conversation can use"},
			}
			kit.SegmentedWith(c.Key("sections"), &segmentedSections, sections, kit.SegmentedOpt{}).Label("Sections").AlignSelf(ui.Start)
			sectionTitle(c, "Full width, equal segments")
			kit.SegmentedWith(c.Key("full"), &segmentedFull, sections, kit.SegmentedOpt{FullWidth: true}).Label("Full width")
			sectionTitle(c, "Disabled")
			kit.Segmented(c.Key("disabled"), &segmentedDisabled, []string{"One", "Two", "Three"}).Label("Disabled").Disabled(true).AlignSelf(ui.Start)
		})
	}})
}
