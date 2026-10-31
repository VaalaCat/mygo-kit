package kit

import (
	"github.com/egoist/mygo/ui"
	"testing"
	"time"
)

// TestButton 验证按钮的公开交互。
func TestButton(t *testing.T) {
	// 1. 使用真实离屏输入逐项驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"button-keyboard-disabled", func(t *testing.T) {
			count := 0
			disabled := false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if Button(c, "Run").Disabled(disabled).Clicked() {
					count++
				}
			}, 300, 120)
			if err := tt.Click("Run"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyEnter)
			if count != 2 {
				t.Fatal(count)
			}
			disabled = true
			tt.Frame()
			tt.Click("Run")
			tt.Key(0, ui.KeyEnter)
			if count != 2 {
				t.Fatal(count)
			}
		}},
		{"small-icon-button-size-and-click", func(t *testing.T) {
			count := 0
			var size ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				b := SmallIconButton(c, IconCopy, "Copy")
				size = b.Bounds()
				if b.Clicked() {
					count++
				}
			}, 200, 80)
			tt.Frame()
			if size.W != 20 || size.H != 20 {
				t.Fatalf("size=%v", size)
			}
			if err := tt.Click("Copy"); err != nil || count != 1 {
				t.Fatalf("count=%d err=%v", count, err)
			}
		}},
		{"flush-button-aligns-text-and-hover-keeps-background", func(t *testing.T) {
			var flush, normal ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				ui.Column(c).Padding(0, 16).AlignItems(ui.Start).Children(func() {
					flush = FlushButton(c, "Refresh").Bounds()
					normal = Button(c, "Refresh").Bounds()
				})
			}, 300, 120)
			tt.Frame()
			if flush.X != 16 || flush.H != 28 || normal.W-flush.W != 16 {
				t.Fatalf("flush=%v normal=%v", flush, normal)
			}
			before := tt.Image().RGBAAt(int(flush.X), int(flush.Y)+2)
			tt.Move(flush.X+flush.W/2, flush.Y+flush.H/2)
			if after := tt.Image().RGBAAt(int(flush.X), int(flush.Y)+2); after != before {
				t.Fatalf("hover repainted background: %v -> %v", before, after)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestIconButtonTip 验证带提示的图标按钮：点击、可访问名称和悬停 600ms 后的提示。
func TestIconButtonTip(t *testing.T) {
	// 1. 一个调用同时给出按钮、名称和提示。
	count := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		ui.Row(c).Padding(40).Children(func() {
			if IconButtonTip(c, IconCopy, "Copy").Clicked() {
				count++
			}
		})
	}, 300, 160)
	button, ok := tt.Find("Copy")
	if !ok || button.W != 28 || button.H != 28 {
		t.Fatalf("button %v ok=%v", button, ok)
	}
	// 2. 提示在悬停满 600ms 后出现，点击仍触发按钮。
	tt.Move(button.X+button.W/2, button.Y+button.H/2)
	tt.Advance(500 * time.Millisecond)
	before := len(tt.Texts())
	tt.Advance(200 * time.Millisecond)
	if after := len(tt.Texts()); after != before+1 {
		t.Fatalf("tooltip did not appear after 700ms: %v", tt.Texts())
	}
	tt.Click("Copy")
	if count != 1 {
		t.Fatalf("count=%d", count)
	}
}

// firstInkX 返回控件内第一个与背景不同的像素 x（背景取控件上沿内侧的像素）；找不到返回负数。
// 按钮的文字节点和按钮同名，Find 只能拿到整个按钮，所以靠像素量文字起点。
func firstInkX(tt *ui.Tester, control ui.Rect, from, to float32) float32 {
	bg := tt.Image().RGBAAt(int(from)+1, int(control.Y)+3)
	for x := int(from); x < int(to); x++ {
		for y := int(control.Y) + 6; y < int(control.Y+control.H)-6; y++ {
			got := tt.Image().RGBAAt(x, y)
			if abs(int(got.R)-int(bg.R))+abs(int(got.G)-int(bg.G))+abs(int(got.B)-int(bg.B)) > 60 {
				return float32(x)
			}
		}
	}
	return -1
}

// TestStretchedControlsKeepTextAtStart 验证被拉宽的按钮、工具按钮、下拉触发器和链接文字不居中：
// 在 480 宽的纵向容器里（控件拉到 448 宽），文字起点紧跟控件左缘加它自己的左内边距，而不是在控件中央。
func TestStretchedControlsKeepTextAtStart(t *testing.T) {
	// 1. 每个用例画一个被拉宽的控件，padding 是控件的左内边距加边线。
	theme := OneLight()
	space := theme.Space(2)
	selected := "a"
	options := []SelectOption{{Value: "a", Label: "Alpha"}}
	tests := []struct {
		name    string
		padding float32
		build   func(c *ui.Context) ui.Element
	}{
		{"button", space, func(c *ui.Context) ui.Element { return Button(c, "Plain") }},
		{"primary-button", space, func(c *ui.Context) ui.Element { return PrimaryButton(c, "Primary") }},
		{"danger-button", space, func(c *ui.Context) ui.Element { return DangerButton(c, "Danger") }},
		{"flush-button", 0, func(c *ui.Context) ui.Element { return FlushButton(c, "Flush") }},
		{"tool-button", space, func(c *ui.Context) ui.Element { return ToolButton(c, nil, "Tool", false) }},
		{"tool-button-chevron", space, func(c *ui.Context) ui.Element { return ToolButton(c, nil, "Menu", true) }},
		{"select-toolbar", space, func(c *ui.Context) ui.Element { return Select(c, &selected, SelectOpt{Options: options}) }},
		{"select-bordered", space + 1, func(c *ui.Context) ui.Element {
			return Select(c, &selected, SelectOpt{Options: options, Bordered: true})
		}},
		{"link", 0, func(c *ui.Context) ui.Element { return Link(c, "Docs") }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var control ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Fill().Padding(theme.Space(4)).Children(func() { control = test.build(c).Bounds() })
			}, 480, 200)
			tt.Frame()
			// 2. 文字笔画起点在左内边距之后几像素内（字形有左侧间距）。
			start := control.X + test.padding
			ink := firstInkX(tt, control, control.X+1+test.padding-1, control.X+control.W/2)
			if control.W < 400 || ink < start || ink > start+4 {
				t.Fatalf("control=%v first ink x=%v, want within [%v, %v]", control, ink, start, start+4)
			}
		})
	}
}

// TestStretchedToolButtonChevronAtEnd 验证被拉宽的带箭头工具按钮：下拉箭头贴右缘内边距，行中部没有内容。
func TestStretchedToolButtonChevronAtEnd(t *testing.T) {
	theme := OneLight()
	var control ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, theme)
		ui.Column(c).Fill().Padding(theme.Space(4)).Children(func() { control = ToolButton(c, IconFolder, "Menu", true).Bounds() })
	}, 480, 200)
	tt.Frame()
	// 1. 右侧内边距之内的 12 宽区域有箭头笔画，控件右半部其余位置是空的。
	right := control.X + control.W - theme.Space(2)
	if ink := firstInkX(tt, control, right-IconSizeSmall-1, right+1); ink < 0 {
		t.Fatalf("no chevron at the right edge of %v", control)
	}
	if ink := firstInkX(tt, control, control.X+control.W/2, right-IconSizeSmall-1); ink >= 0 {
		t.Fatalf("content in the middle of the stretched button at x=%v", ink)
	}
}
