package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var switchOn, switchOff, switchBare = true, false, true

func init() {
	register(page{name: "Switch", order: 50, view: func(c *ui.Context) {
		pageBody(c, func() {
			sectionTitle(c, "With label (the whole row toggles; Space toggles while focused)")
			ui.Row(c).Gap(24).Children(func() {
				kit.Switch(c.Key("on"), &switchOn, "Reduce motion")
				kit.Switch(c.Key("off"), &switchOff, "Notify when done")
			})
			kit.Separator(c)
			sectionTitle(c, "Track only (name it with Label)")
			kit.Switch(c.Key("bare"), &switchBare, "").Label("Auto-save").AlignSelf(ui.Start)
			kit.Separator(c)
			sectionTitle(c, "Disabled (framework opacity)")
			ui.Row(c).Gap(24).Children(func() {
				kit.Switch(c.Key("disabled-on"), &switchOn, "On").Disabled(true)
				kit.Switch(c.Key("disabled-off"), &switchOff, "Off").Disabled(true)
			})
		})
	}})
}
