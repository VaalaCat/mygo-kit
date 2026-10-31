package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var progressValue float32 = .42

func init() {
	register(page{name: "Progress", order: 23, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				for _, value := range []float32{0, .14, .25, .5, .84, .85, .94, .95, 1} {
					kit.ProgressRing(c, value, 20)
				}
			})
			sectionTitle(c, "ProgressBar (4 high)")
			ui.Column(c).Gap(8).Width(320).Children(func() {
				for _, value := range []float32{0, .25, progressValue, 1} {
					kit.ProgressBar(c, value)
				}
			})
			sectionTitle(c, "Context usage")
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				kit.ProgressRing(c, progressValue, 16)
				ui.Text(c, "Context window").TextColor(kit.GetTheme(c).Colors.Muted)
				if kit.Button(c.Key("less"), "−10%").Clicked() {
					progressValue = max(progressValue-.1, 0)
				}
				if kit.Button(c.Key("more"), "+10%").Clicked() {
					progressValue = min(progressValue+.1, 1)
				}
			})
		})
	}})
}
