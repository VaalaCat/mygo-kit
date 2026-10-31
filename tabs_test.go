package kit

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestTabs 验证页签的公开交互。
func TestTabs(t *testing.T) {
	// 1. 使用真实离屏输入逐项驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"tabs-keyboard", func(t *testing.T) {
			selected := 0
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Tabs(c, &selected, []string{"First", "Second"}) }, 400, 200)
			tt.Click("First")
			tt.Key(0, ui.KeyRight)
			if selected != 1 {
				t.Fatal(selected)
			}
		}},
		{"tabs-focus-ring-stays-inside-focused-tab", func(t *testing.T) {
			theme := OneLight()
			selected := 0
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				Tabs(c, &selected, []string{"First", "Second", "Third"})
			}, 400, 120)
			tt.Click("First")
			tt.Key(0, ui.KeyRight)
			bounds, ok := tt.Find("Second")
			if !ok || bounds.H != 36 {
				t.Fatalf("tab bounds=%v ok=%v", bounds, ok)
			}
			ring, image := color.RGBA{theme.Colors.Focus.R, theme.Colors.Focus.G, theme.Colors.Focus.B, theme.Colors.Focus.A}, tt.Image()
			// 环的像素只出现在键盘聚焦的页签范围内，右缘一列没有被后画的下一个页签盖住。
			inside, outside := 0, 0
			for y := 0; y < image.Bounds().Dy(); y++ {
				for x := 0; x < image.Bounds().Dx(); x++ {
					if image.RGBAAt(x, y) != ring {
						continue
					}
					if float32(x) >= bounds.X && float32(x) < bounds.X+bounds.W && float32(y) >= bounds.Y && float32(y) < bounds.Y+bounds.H {
						inside++
					} else {
						outside++
					}
				}
			}
			if inside == 0 || outside != 0 {
				t.Fatalf("ring pixels inside=%d outside=%d", inside, outside)
			}
			right, bottom := int(bounds.X+bounds.W)-1, int(bounds.Y+bounds.H)-1
			if image.RGBAAt(right, int(bounds.Y+bounds.H/2)) != ring || image.RGBAAt(int(bounds.X+bounds.W/2), bottom) != ring {
				t.Fatal("ring is cut on the right or bottom edge")
			}
		}},
		{"tabs-focus-follows-external-selection", func(t *testing.T) {
			selected := 0
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				Tabs(c, &selected, []string{"First", "Second", "Third"})
			}, 400, 120)
			tt.Click("First")
			tt.Key(0, ui.KeyRight)
			selected = 0
			tt.Frame()
			// 选中项被外部改回第一项后，焦点跟到第一项：再按右键应选第二项，而不是第三项。
			tt.Key(0, ui.KeyRight)
			if selected != 1 {
				t.Fatalf("selected=%d", selected)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestTabsDisabled 验证已有焦点或展开状态的页签禁用后不再修改值。
func TestTabsDisabled(t *testing.T) {
	// 1. 先建立交互状态，再用下一帧禁用，最后发送输入。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"tabs", func(t *testing.T) {
			selected := 0
			disabled := false
			changes := 0
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if Tabs(c, &selected, []string{"First", "Second"}).Disabled(disabled).Changed() {
					changes++
				}
			}, 400, 200)
			tt.Click("First")
			disabled = true
			tt.Frame()
			before := changes
			tt.Key(0, ui.KeyRight)
			tt.Click("Second")
			if selected != 0 || changes != before {
				t.Fatalf("disabled change: value=%v changes=%d before=%d", selected, changes, before)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
