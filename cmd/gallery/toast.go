package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

func init() {
	register(page{name: "Toast", order: 27, view: func(c *ui.Context) {
		pageBody(c, func() {
			ui.Row(c).Gap(8).Children(func() {
				if kit.Button(c.Key("toast-info"), "Info").Clicked() {
					kit.ShowToast(c, kit.ToastOpt{Title: "Conversation exported", Kind: kit.ToastInfo})
				}
				if kit.Button(c.Key("toast-success"), "Success with action").Clicked() {
					kit.ShowToast(c, kit.ToastOpt{Title: "Conversation archived", Kind: kit.ToastSuccess, Action: "Undo", OnAction: func() {}})
				}
				if kit.Button(c.Key("toast-error"), "Error").Clicked() {
					kit.ShowToast(c, kit.ToastOpt{Title: "Connection lost", Description: "The runtime stopped responding.", Kind: kit.ToastError})
				}
			})
		})
	}})
}
