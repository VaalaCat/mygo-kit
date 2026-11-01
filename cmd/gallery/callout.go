package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Callout", order: 43, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		pageBody(c, func() {
			ui.Column(c).MaxWidth(640).Gap(theme.Space(2)).Children(func() {
				kit.Callout(c, kit.CalloutError, t(c, "Model request failed"), t(c, "The provider returned 429: rate limit exceeded. Retry in a few seconds."), func() {
					kit.Button(c, t(c, "Resend"))
				})
				kit.Callout(c, kit.CalloutWarning, t(c, "Context is nearly full"), t(c, "86% of the window is used; older turns may be compacted."), nil)
				kit.Callout(c, kit.CalloutInfo, t(c, "Conversation forked"), "", func() {
					kit.Button(c, t(c, "Open"))
					kit.IconButton(c, kit.IconX, t(c, "Dismiss"))
				})
			})
		})
	}})
}
