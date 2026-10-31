package kit

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

// TestDisclosure 验证折叠标题的公开交互。
func TestDisclosure(t *testing.T) {
	// 1. 使用真实离屏输入逐项驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"disclosure", func(t *testing.T) {
			open := false
			tt := ui.NewTester(func(c *ui.Context) {
				theme := OneDark()
				theme.Motion.Reduced = true
				ApplyTheme(c, theme)
				Disclosure(c, "Details", &open, func() { ui.Text(c, "Visible detail") })
			}, 400, 200)
			tt.Click("Details")
			if !open || !tt.HasText("Visible detail") {
				t.Fatal("not expanded")
			}
			tt.Key(0, ui.KeyEnter)
			if open || tt.HasText("Visible detail") {
				t.Fatal("not collapsed")
			}
		}},
		{"disabled-has-no-hover-background", func(t *testing.T) {
			// 同一个悬停位置：可用的标题画悬停底色，禁用的不画。
			theme := OneLight()
			for _, disabled := range []bool{false, true} {
				open := false
				tt := ui.NewTester(func(c *ui.Context) {
					ApplyTheme(c, theme)
					Disclosure(c, "Details", &open, func() {}).Width(300).Disabled(disabled)
				}, 400, 200)
				tt.Frame()
				title, _ := tt.Find("Details")
				tt.Move(title.X+title.W-8, title.Y+title.H/2)
				if got := pixelNear(tt, title.X+title.W-8, title.Y+title.H/2, theme.Colors.Hover, 1); got == disabled {
					t.Fatalf("disabled=%v hover background drawn=%v", disabled, got)
				}
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
