package kit

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestListRow 验证列表行的高度、缩进、选中底色、点击和操作按钮的显示时机。
func TestListRow(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	theme := OneLight()
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"height-28-and-44-with-description", func(t *testing.T) {
			var single, double ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Width(300).Gap(0).Children(func() {
					single = ListRow(c.Key("single"), ListRowOpt{Icon: IconBot, Title: "One"}).Bounds()
					double = ListRow(c.Key("double"), ListRowOpt{Icon: IconBot, Title: "Two", Description: "second line"}).Bounds()
				})
			}, 400, 200)
			tt.Frame()
			if single.H != 28 || double.H != 44 {
				t.Fatalf("heights %v %v", single.H, double.H)
			}
		}},
		{"two-line-text-keeps-4-padding-top-and-bottom", func(t *testing.T) {
			var row ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Width(300).Gap(0).Children(func() {
					row = ListRow(c, ListRowOpt{Title: "Two", Description: "second line"}).Bounds()
				})
			}, 400, 200)
			tt.Frame()
			description, ok := tt.Find("second line")
			if !ok || description.Y+description.H != row.Y+row.H-4 {
				t.Fatalf("description %v row %v ok=%v", description, row, ok)
			}
		}},
		{"depth-indents-16-per-level", func(t *testing.T) {
			// 行按钮的 Label 改成别的名字，Find 才会返回标题文字本身的位置。
			names := []string{"zero", "one", "two"}
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Width(300).Gap(0).Children(func() {
					for depth, name := range names {
						ListRow(c.Key(depth), ListRowOpt{Icon: IconFolder, Title: name, Depth: depth}).Label("row " + name)
					}
				})
			}, 400, 200)
			tt.Frame()
			base, _ := tt.Find("zero")
			for depth, name := range names[1:] {
				title, ok := tt.Find(name)
				if want := float32(depth+1) * 16; !ok || title.X-base.X != want {
					t.Fatalf("%s indent %v want %v (ok=%v)", name, title.X-base.X, want, ok)
				}
			}
		}},
		{"disabled-row-has-no-hover-background", func(t *testing.T) {
			// 同一个悬停位置：可用的行画悬停底色，禁用的行不画。
			for _, disabled := range []bool{false, true} {
				tt := ui.NewTester(func(c *ui.Context) {
					ApplyTheme(c, theme)
					ui.Column(c).Width(300).Gap(0).Children(func() {
						ListRow(c.Key("row"), ListRowOpt{Title: "row"}).Disabled(disabled)
					})
				}, 400, 200)
				tt.Frame()
				row, _ := tt.Find("row")
				tt.Move(row.X+row.W-8, row.Y+row.H/2)
				if got := pixelNear(tt, row.X+row.W-8, row.Y+row.H/2, theme.Colors.Hover, 1); got == disabled {
					t.Fatalf("disabled=%v hover background drawn=%v", disabled, got)
				}
			}
		}},
		{"selected-background-and-click", func(t *testing.T) {
			selected, clicked := "b", ""
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Width(300).Gap(0).Children(func() {
					for _, name := range []string{"a", "b"} {
						if ListRow(c.Key(name), ListRowOpt{Title: name, Selected: name == selected}).Clicked() {
							clicked = name
						}
					}
				})
			}, 400, 200)
			a, _ := tt.Find("a")
			b, _ := tt.Find("b")
			if !pixelNear(tt, b.X+b.W-8, b.Y+b.H/2, theme.Colors.Selected, 1) || pixelNear(tt, a.X+a.W-8, a.Y+a.H/2, theme.Colors.Selected, 1) {
				t.Fatal("selected row background")
			}
			tt.Click("a")
			if clicked != "a" {
				t.Fatalf("clicked %q", clicked)
			}
		}},
		{"actions-show-on-hover-selected-and-focus", func(t *testing.T) {
			selected, actions := false, 0
			var row ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Width(300).Gap(0).Children(func() {
					ListRow(c.Key("other"), ListRowOpt{Title: "other"})
					row = ListRow(c.Key("row"), ListRowOpt{Title: "row", Selected: selected, Actions: func() {
						if IconButton(c, IconEllipsis, "More").Clicked() {
							actions++
						}
					}}).Bounds()
				})
			}, 400, 200)
			tt.Frame()
			// 按钮可见时图标笔画让按钮范围内出现与角落底色不同的像素；不可见时整个范围都是行底色，与主题色和图标形状无关。
			visible := func() bool {
				button, _ := tt.Find("More")
				img := tt.Image()
				bg := img.RGBAAt(int(button.X)+1, int(button.Y)+1)
				for y := int(button.Y); y < int(button.Y+button.H); y++ {
					for x := int(button.X); x < int(button.X+button.W); x++ {
						if img.RGBAAt(x, y) != bg {
							return true
						}
					}
				}
				return false
			}
			if visible() {
				t.Fatal("actions visible with no hover, selection or focus")
			}
			tt.Move(row.X+20, row.Y+row.H/2)
			if !visible() {
				t.Fatal("actions hidden while hovering")
			}
			tt.Move(row.X+row.W+100, row.Y+row.H/2)
			if visible() {
				t.Fatal("actions still visible after the pointer left")
			}
			selected = true
			tt.Frame()
			if !visible() {
				t.Fatal("actions hidden while selected")
			}
			selected = false
			tt.Frame()
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("row") || !visible() {
				t.Fatal("actions hidden while the row has the focus")
			}
			if err := tt.Click("More"); err != nil || actions != 1 {
				t.Fatalf("action click err=%v actions=%d", err, actions)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
