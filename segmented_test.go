package kit

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// TestSegmented 验证分段控件的点击、方向键、禁用、高度、选中段浅底、外框、段间竖线和键盘焦点环。
func TestSegmented(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	labels := []string{"System", "Light", "Dark"}
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"click-chooses-segment", func(t *testing.T) {
			selected := 0
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Segmented(c, &selected, labels) }, 300, 80)
			tt.Click("Dark")
			if selected != 2 {
				t.Fatalf("selected=%d", selected)
			}
		}},
		{"arrow-keys-move-choice", func(t *testing.T) {
			selected := 0
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Segmented(c, &selected, labels).Label("Appearance") }, 300, 80)
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			tt.Key(0, ui.KeyRight)
			if selected != 2 {
				t.Fatalf("after Right Right selected=%d", selected)
			}
			tt.Key(0, ui.KeyLeft)
			if selected != 1 {
				t.Fatalf("after Left selected=%d", selected)
			}
		}},
		{"disabled-ignores-click", func(t *testing.T) {
			selected := 0
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Segmented(c, &selected, labels).Disabled(true) }, 300, 80)
			tt.Click("Light")
			if selected != 0 {
				t.Fatalf("disabled segmented chose %d", selected)
			}
		}},
		{"height-and-selected-tint", func(t *testing.T) {
			theme := OneLight()
			selected := 1
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Row(c).Padding(theme.Space(4)).Children(func() { bounds = Segmented(c, &selected, labels).Bounds() })
			}, 300, 80)
			tt.Frame()
			if bounds.H != 28 {
				t.Fatalf("height %v", bounds.H)
			}
			// Find 返回按标签命名的段本身；段左内边距 8，边缘内 4 处只有段底色：选中段是强调色浅底，其余是透明（露出窗口底色）。
			tint := theme.Colors.AccentSubtle.Over(theme.Colors.Editor)
			for i, want := range []ui.Color{theme.Colors.Editor, tint, theme.Colors.Editor} {
				segment, _ := tt.Find(labels[i])
				if !pixelNear(tt, segment.X+4, segment.Y+segment.H/2, want, 2) {
					t.Fatalf("segment %d background %v want %v", i, tt.Image().RGBAAt(int(segment.X+4), int(segment.Y+segment.H/2)), want)
				}
			}
		}},
		{"outer-border-and-dividers", func(t *testing.T) {
			theme := OneLight()
			selected := 0
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Row(c).Padding(theme.Space(4)).Children(func() { bounds = Segmented(c, &selected, labels).Bounds() })
			}, 300, 80)
			tt.Frame()
			// 外框和竖线都是 60% 透明度的边线色；竖线在两段之间的 1 像素，外框在轨道最外一圈。
			line := theme.Colors.Border.Alpha(0.6).Over(theme.Colors.Editor)
			middle := bounds.Y + bounds.H/2
			if !pixelNear(tt, bounds.X, middle, line, 2) || !pixelNear(tt, bounds.X+bounds.W-1, middle, line, 2) {
				t.Fatalf("outer border missing: left %v right %v want %v", tt.Image().RGBAAt(int(bounds.X), int(middle)), tt.Image().RGBAAt(int(bounds.X+bounds.W-1), int(middle)), line)
			}
			for i := range labels[:len(labels)-1] {
				segment, _ := tt.Find(labels[i])
				x := segment.X + segment.W
				if !pixelNear(tt, x, middle, line, 2) || !pixelNear(tt, x, bounds.Y+2, line, 2) {
					t.Fatalf("divider after segment %d missing at x=%v: %v want %v", i, x, tt.Image().RGBAAt(int(x), int(middle)), line)
				}
			}
		}},
		{"focus-ring-stays-inside-focused-segment", func(t *testing.T) {
			theme := OneLight()
			selected := 0
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Row(c).Padding(theme.Space(4)).Children(func() { Segmented(c, &selected, labels).Label("Appearance") })
			}, 300, 80)
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			bounds, _ := tt.Find("Light")
			ring, image := color.RGBA{theme.Colors.Focus.R, theme.Colors.Focus.G, theme.Colors.Focus.B, theme.Colors.Focus.A}, tt.Image()
			// 环的像素只出现在聚焦的段范围内，且贴着段的上下缘。
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
			if image.RGBAAt(int(bounds.X+bounds.W/2), int(bounds.Y)) != ring || image.RGBAAt(int(bounds.X+bounds.W/2), int(bounds.Y+bounds.H)-1) != ring {
				t.Fatal("ring is cut on the top or bottom edge")
			}
		}},
		{"tooltip-shows-after-hover", func(t *testing.T) {
			selected := 0
			items := []SegmentedItem{{Label: "One", Tooltip: "First section"}, {Label: "Two"}}
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); SegmentedWith(c, &selected, items, SegmentedOpt{}) }, 300, 120)
			bounds, _ := tt.Find("One")
			tt.Move(bounds.X+bounds.W/2, bounds.Y+bounds.H/2)
			tt.Advance(700 * time.Millisecond)
			if !tt.HasText("First section") {
				t.Fatalf("tooltip missing: %v", tt.Texts())
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestSegmentedWidth 验证分段控件的宽度：默认各段贴合自己的内容、靠左；FullWidth 占满容器、各段等宽；图标让段变宽。
func TestSegmentedWidth(t *testing.T) {
	// 1. 同一组文字长度不同的段，放进纵向容器，按选项和图标分别测量轨道与各段。
	theme := OneLight()
	labels := []string{"local", "sandbox-long"}
	tests := []struct {
		name  string
		opt   SegmentedOpt
		icon  bool
		check func(t *testing.T, track, first, second ui.Rect)
	}{
		{"auto-width-hugs-each-label", SegmentedOpt{}, false, func(t *testing.T, track, first, second ui.Rect) {
			if track.W > 300 || first.W >= second.W {
				t.Fatalf("track=%v first=%v second=%v: not hugging content", track, first, second)
			}
			if track.X != 16 {
				t.Fatalf("track %v not on the left content line", track)
			}
		}},
		{"full-width-fills-container-with-equal-segments", SegmentedOpt{FullWidth: true}, false, func(t *testing.T, track, first, second ui.Rect) {
			if track.W < 440 || abs(int(first.W-second.W)) > 1 {
				t.Fatalf("track=%v first=%v second=%v: not full width with equal segments", track, first, second)
			}
			// 文字居中在各自的一半里。
			inner := track.W - 2
			if abs(int(first.X+first.W/2-(track.X+1+inner/4))) > 1 || abs(int(second.X+second.W/2-(track.X+1+inner*3/4))) > 1 {
				t.Fatalf("track=%v first=%v second=%v: segments are not centered", track, first, second)
			}
		}},
		{"icon-widens-segment", SegmentedOpt{}, true, func(t *testing.T, track, first, second ui.Rect) {
			// 图标 12 加图标与文字的间隙 4；这里只比较第一段，它比不带图标时多出 16。
			plain := hugWidth(theme, labels[0], false)
			if got := first.W - plain; got != IconSizeSmall+theme.Space(1) {
				t.Fatalf("icon segment is %v wide, plain %v", first.W, plain)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selected := 0
			var track ui.Rect
			items := []SegmentedItem{{Label: labels[0]}, {Label: labels[1]}}
			if test.icon {
				items[0].Icon = IconFolder
			}
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Fill().Padding(theme.Space(4)).Children(func() {
					track = SegmentedWith(c, &selected, items, test.opt).AlignSelf(ui.Start).Bounds()
				})
			}, 480, 200)
			tt.Frame()
			first, _ := tt.Find(labels[0])
			second, _ := tt.Find(labels[1])
			test.check(t, track, first, second)
		})
	}
}

// hugWidth 返回只有一段的分段控件里这一段贴合内容的宽度。
func hugWidth(theme Theme, label string, icon bool) float32 {
	item := SegmentedItem{Label: label}
	if icon {
		item.Icon = IconFolder
	}
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, theme)
		ui.Column(c).Fill().Padding(theme.Space(4)).Children(func() {
			SegmentedWith(c, &selected, []SegmentedItem{item}, SegmentedOpt{}).AlignSelf(ui.Start)
		})
	}, 480, 200)
	tt.Frame()
	bounds, _ := tt.Find(label)
	return bounds.W
}
