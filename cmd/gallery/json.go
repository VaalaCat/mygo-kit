package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var jsonState, _ = kit.NewJSONState([]byte(`{"name":"mygo-kit","enabled":true,"features":["native","keyboard","streaming"],"limit":null}`))

func init() {
	register(page{name: "JSON", order: 31, view: func(c *ui.Context) {
		kit.JSONTree(c, jsonState, kit.JSONOpt{}).Grow(1)
	}})
}
