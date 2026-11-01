package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var activityOpen = [5]bool{false, true, false, false, true}

func init() {
	register(page{name: "ActivityCard", order: 41, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		code := func(s string, color ui.Color) {
			theme.TextRole(ui.Text(c, s), kit.TextCode).TextColor(color)
		}
		inspect := func() { kit.IconButton(c, kit.IconBug, t(c, "Inspect")).TextColor(theme.Colors.Muted) }
		pageBody(c, func() {
			ui.Column(c).MaxWidth(640).Gap(0).Children(func() {
				kit.ActivityCard(c, &activityOpen[0], kit.ActivityCardOpt{Icon: kit.IconTerminal, Title: "bash", Summary: "go test -race ./... -count=1", Status: kit.ActivityRunning}, inspect, func() {
					code("$ go test -race ./... -count=1", theme.Colors.Text)
					code("ok  github.com/VaalaCat/mygo-kit  3.2s", theme.Colors.Muted)
				})
				kit.ActivityCard(c, &activityOpen[1], kit.ActivityCardOpt{Icon: kit.IconFileText, Title: "read_file", Summary: "internal/desktop/app.go" + t(c, " · lines 1-120")}, inspect, func() {
					code("package desktop", theme.Colors.Text)
					code("// App is the desktop shell.", theme.Colors.Muted)
					code("type App struct{ ... }", theme.Colors.Text)
				})
				kit.ActivityCard(c, &activityOpen[2], kit.ActivityCardOpt{Icon: kit.IconPencil, Title: "apply_patch", Summary: "internal/desktop/chat.go", Status: kit.ActivityFailed}, inspect, func() {
					code("patch does not apply: context mismatch at line 42", theme.Colors.Danger)
				})
				kit.ActivityCard(c, &activityOpen[3], kit.ActivityCardOpt{Icon: kit.IconTriangleAlert, Title: t(c, "Model call failed"), Tone: kit.ActivityDanger}, nil, func() {
					code("rate_limit: too many requests", theme.Colors.Danger)
				})
				kit.ActivityCard(c, new(bool), kit.ActivityCardOpt{Icon: kit.IconBot, Title: "spawn_agent", Summary: t(c, "review the diff of the last commit"), Status: kit.ActivityInterrupted}, nil, nil)
				kit.ActivityCard(c, new(bool), kit.ActivityCardOpt{Icon: kit.IconBan, Title: t(c, "Turn interrupted"), Tone: kit.ActivityWarning}, nil, nil)
				kit.ActivityCard(c, &activityOpen[4], kit.ActivityCardOpt{Icon: kit.IconBrain, Title: t(c, "Thinking"), Shimmer: true}, nil, func() { ui.Text(c, t(c, "Reasoning in progress.")) })
				kit.ActivityCard(c, new(bool), kit.ActivityCardOpt{Icon: kit.IconReplace, Title: t(c, "Model changed from A to B")}, nil, nil)
			})
		})
	}})
}
