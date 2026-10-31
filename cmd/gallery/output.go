package main

import (
	"strings"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var outputRows = ui.ListState{FollowEnd: true}

func init() {
	register(page{name: "Output", order: 32, view: func(c *ui.Context) {
		kit.ANSIOutput(c, "$ go test ./...\n\x1b[32mPASS\x1b[0m  controls\n\x1b[36mINFO\x1b[0m  keyboard and selection ready\n\x1b[33mWARN\x1b[0m  software-rendered preview\n"+strings.Repeat("a long line  ", 30)+"\n", &outputRows, kit.ANSIOutputOpt{}).Grow(1)
	}})
}
