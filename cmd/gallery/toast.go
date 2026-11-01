package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Toast", order: 27, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(8).Children(func() {
				if kit.Button(c.Key("toast-info"), t(c, "Info")).Clicked() {
					kit.ShowToast(c, kit.ToastOpt{Title: t(c, "Conversation exported"), Kind: kit.ToastInfo})
				}
				if kit.Button(c.Key("toast-success"), t(c, "Success with action")).Clicked() {
					kit.ShowToast(c, kit.ToastOpt{Title: t(c, "Conversation archived"), Kind: kit.ToastSuccess, Action: t(c, "Undo"), OnAction: func() {}})
				}
				if kit.Button(c.Key("toast-error"), t(c, "Error")).Clicked() {
					kit.ShowToast(c, kit.ToastOpt{Title: t(c, "Connection lost"), Description: t(c, "The runtime stopped responding."), Kind: kit.ToastError})
				}
			})
		})
	}})
}
