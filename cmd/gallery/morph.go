package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	morphIndex  int
	morphValues = []string{"Reading files", "Reading config", "Running tests", "Running 3 tests", "正在读取文件", "正在运行测试"}
)

func init() {
	register(page{name: "MorphText", order: 46, view: func(c *ui.Context) {
		t := kit.GetTheme(c)
		pageBody(c, func() {
			ui.Column(c).Gap(t.Space(2)).Children(func() {
				kit.MorphText(c, morphValues[morphIndex])
				if kit.Button(c.Key("morph-next"), "Next status").Clicked() {
					morphIndex = (morphIndex + 1) % len(morphValues)
				}
			})
		})
	}})
}
