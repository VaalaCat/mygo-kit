package kit

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestSwitch 验证开关的点击、键盘、禁用、尺寸和底色。
func TestSwitch(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"click-track-and-label-toggle", func(t *testing.T) {
			on := false
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Switch(c, &on, "Notify") }, 300, 80)
			tt.Click("Notify")
			if !on {
				t.Fatal("label click did not turn on")
			}
			tt.Click("Notify")
			if on {
				t.Fatal("second click did not turn off")
			}
		}},
		{"space-toggles-focused", func(t *testing.T) {
			on := false
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Switch(c, &on, "").Label("Auto-save") }, 300, 80)
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("Auto-save") {
				t.Fatal("switch did not take the focus")
			}
			tt.Key(0, ui.KeySpace)
			if !on {
				t.Fatal("Space did not turn on")
			}
			tt.Key(0, ui.KeySpace)
			if on {
				t.Fatal("Space did not turn off")
			}
		}},
		{"disabled-ignores-input", func(t *testing.T) {
			on := false
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Switch(c, &on, "Locked").Disabled(true) }, 300, 80)
			tt.Click("Locked")
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeySpace)
			if on {
				t.Fatal("disabled switch toggled")
			}
		}},
		{"track-size-and-accent-when-on", func(t *testing.T) {
			theme := OneLight()
			for _, on := range []bool{true, false} {
				var bounds ui.Rect
				value := on
				tt := ui.NewTester(func(c *ui.Context) {
					ApplyTheme(c, theme)
					ui.Row(c).Padding(theme.Space(4)).Children(func() { bounds = Switch(c, &value, "").Label("Track").Bounds() })
				}, 200, 80)
				tt.Frame()
				if bounds.W != 36 || bounds.H != 20 {
					t.Fatalf("track %vx%v", bounds.W, bounds.H)
				}
				// 开时滑块在右侧、左侧露出强调色；关时滑块在左侧、右侧露出边线色混入文字色。
				x, y := bounds.X+bounds.W-6, bounds.Y+bounds.H/2
				want := theme.Colors.Border.Mix(theme.Colors.Text, .15)
				if on {
					x, want = bounds.X+6, theme.Colors.Accent
				}
				if !pixelNear(tt, x, y, want, 4) {
					t.Fatalf("on=%v track color %v want %v", on, tt.Image().RGBAAt(int(x), int(y)), want)
				}
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
