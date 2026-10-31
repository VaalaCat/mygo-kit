package kit

import (
	"github.com/egoist/mygo/ui"
	"testing"
	"time"
)

// TestLocale 验证组件文案随窗口语言变化，应用仍拥有语言选择。
func TestLocale(t *testing.T) {
	// 1. 两种语言与未知语言回退使用同一组件实例。
	state, _ := NewJSONState([]byte("null"))
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); JSONTree(c, state, JSONOpt{}).Fill() }, 400, 240)
	for _, test := range []struct{ locale, label string }{{"en", "Copy JSON"}, {"zh-Hans", "复制 JSON"}, {"fr", "Copy JSON"}} {
		t.Run(test.locale, func(t *testing.T) {
			tt.SetLocale(test.locale)
			if !tt.HasText(test.label) {
				t.Fatal(tt.Texts())
			}
			if err := tt.Click(test.label); err != nil {
				t.Fatal(err)
			}
			if tt.Clipboard() != "null" {
				t.Fatal(tt.Clipboard())
			}
			// 复制后按钮会显示 2 秒「已复制」，等它恢复再测下一种语言。
			tt.Advance(3 * time.Second)
		})
	}
	// 2. 调用方可为相同阅读器设置不同可访问复制名称。
	tt = ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		JSONTree(c, state, JSONOpt{Title: "Result", CopyLabel: "Copy result"}).Fill()
	}, 400, 240)
	tt.Click("Copy result")
	if tt.Clipboard() != "null" {
		t.Fatal(tt.Clipboard())
	}
}
