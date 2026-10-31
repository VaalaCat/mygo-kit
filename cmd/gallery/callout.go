package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Callout", order: 43, view: func(c *ui.Context) {
		t := kit.GetTheme(c)
		pageBody(c, func() {
			ui.Column(c).MaxWidth(640).Gap(t.Space(2)).Children(func() {
				kit.Callout(c, kit.CalloutError, "Model request failed", "The provider returned 429: rate limit exceeded. Retry in a few seconds.", func() {
					kit.Button(c, "Resend")
				})
				kit.Callout(c, kit.CalloutWarning, "Context is nearly full", "86% of the window is used; older turns may be compacted.", nil)
				kit.Callout(c, kit.CalloutInfo, "Conversation forked", "", func() {
					kit.Button(c, "Open")
					kit.IconButton(c, kit.IconX, "Dismiss")
				})
			})
		})
	}})
}
