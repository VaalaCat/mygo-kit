package kit

import (
	"github.com/egoist/mygo/ui"
	"image/color"
	"testing"
)

// TestInput 验证输入的公开交互。
func TestInput(t *testing.T) {
	// 1. 使用真实离屏输入逐项驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"controlled-local-input", func(t *testing.T) {
			value := ""
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				local := value
				e := TextArea(c.Key("input"), &local).Label("Message")
				if e.Changed() {
					value = local
				}
			}, 400, 200)
			if err := tt.Click("Message"); err != nil {
				t.Fatal(err)
			}
			tt.Compose("中文", 2)
			if value != "" {
				t.Fatal(value)
			}
			tt.Type("中文 input")
			tt.Frame()
			if value != "中文 input" {
				t.Fatal(value)
			}
			tt.Command("selectAll")
			tt.Command("copy")
			if tt.Clipboard() != value {
				t.Fatal(tt.Clipboard())
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestInputVariants 验证输入框的高度、边线、只读和无边框变体。
func TestInputVariants(t *testing.T) {
	// 1. 每一行用一组选项渲染单行输入，再检查公开表现。
	theme := OneLight()
	rgba := func(c ui.Color) color.RGBA { return color.RGBA{c.R, c.G, c.B, c.A} }
	tests := []struct {
		name string
		opt  InputOpt
		run  func(t *testing.T, tt *ui.Tester, bounds ui.Rect, value *string)
	}{
		{"default-is-28-high-with-border", InputOpt{}, func(t *testing.T, tt *ui.Tester, b ui.Rect, _ *string) {
			if b.H != 28 {
				t.Fatalf("height=%v", b.H)
			}
			if got := tt.Image().RGBAAt(int(b.X), int(b.Y+b.H/2)); got != rgba(theme.Colors.Border) {
				t.Fatalf("border=%v", got)
			}
		}},
		{"invalid-uses-danger-border", InputOpt{Invalid: true}, func(t *testing.T, tt *ui.Tester, b ui.Rect, _ *string) {
			if got := tt.Image().RGBAAt(int(b.X), int(b.Y+b.H/2)); got != rgba(theme.Colors.Danger) {
				t.Fatalf("border=%v", got)
			}
		}},
		{"focused-uses-focus-border", InputOpt{}, func(t *testing.T, tt *ui.Tester, b ui.Rect, _ *string) {
			if err := tt.Click("Field"); err != nil {
				t.Fatal(err)
			}
			if got := tt.Image().RGBAAt(int(b.X), int(b.Y+b.H/2)); got != rgba(theme.Colors.Focus) {
				t.Fatalf("border=%v", got)
			}
		}},
		{"read-only-ignores-typing", InputOpt{ReadOnly: true}, func(t *testing.T, tt *ui.Tester, _ ui.Rect, value *string) {
			if err := tt.Click("Field"); err != nil {
				t.Fatal(err)
			}
			tt.Type("more")
			tt.Frame()
			if *value != "fixed" {
				t.Fatalf("value=%q", *value)
			}
		}},
		{"read-only-has-no-border-or-fill", InputOpt{ReadOnly: true}, func(t *testing.T, tt *ui.Tester, b ui.Rect, _ *string) {
			outside := tt.Image().RGBAAt(int(b.X)-4, int(b.Y+b.H/2))
			if b.H != 28 {
				t.Fatalf("height=%v", b.H)
			}
			for _, p := range [][2]int{{int(b.X), int(b.Y + b.H/2)}, {int(b.X) + 3, int(b.Y) + 3}, {int(b.X+b.W) - 1, int(b.Y + b.H/2)}} {
				if got := tt.Image().RGBAAt(p[0], p[1]); got != outside {
					t.Fatalf("pixel %v=%v differs from surrounding %v", p, got, outside)
				}
			}
		}},
		{"borderless-draws-no-border", InputOpt{Borderless: true}, func(t *testing.T, tt *ui.Tester, b ui.Rect, _ *string) {
			inside := tt.Image().RGBAAt(int(b.X), int(b.Y+b.H/2))
			outside := tt.Image().RGBAAt(int(b.X)-4, int(b.Y+b.H/2))
			if inside != outside {
				t.Fatalf("inside=%v outside=%v", inside, outside)
			}
		}},
		{"leading-icon-is-drawn-in-left-padding", InputOpt{Leading: IconSearch}, func(t *testing.T, tt *ui.Tester, b ui.Rect, _ *string) {
			background := rgba(theme.Colors.Input)
			drawn := false
			for y := int(b.Y) + 7; y < int(b.Y+b.H)-7 && !drawn; y++ {
				for x := int(b.X) + 9; x < int(b.X)+23; x++ {
					drawn = drawn || tt.Image().RGBAAt(x, y) != background
				}
			}
			if !drawn {
				t.Fatal("icon missing")
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := "fixed"
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Padding(16).Children(func() {
					local := value
					e := TextInputWith(c.Key("input"), &local, test.opt).Label("Field")
					bounds = e.Bounds()
					if e.Changed() {
						value = local
					}
				})
			}, 300, 120)
			tt.Frame()
			test.run(t, tt, bounds, &value)
		})
	}
}

// TestInputLeadingClipsText 验证带前缀图标的输入框里，滚动的长文字不会画到图标和图标右侧的间距下面。
func TestInputLeadingClipsText(t *testing.T) {
	theme := OneLight()
	tests := []struct {
		name string
		opt  InputOpt
	}{
		{"bordered", InputOpt{Leading: IconFolder}},
		{"invalid", InputOpt{Leading: IconFolder, Invalid: true}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := "/a/very/long/path/that/does/not/fit/inside/the/small/input/box/at/all/so/it/scrolls"
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Padding(16).Children(func() {
					local := value
					e := TextInputWith(c.Key("input"), &local, test.opt).Label("Field").Width(160)
					bounds = e.Bounds()
				})
			}, 300, 120)
			tt.Click("Field")
			tt.Key(0, ui.KeyEnd)
			tt.Frame()
			// 图标占 [9, 23)，其后 8 是间距，文字区从 31 开始：间距内不能有任何笔画。
			background := color.RGBA{theme.Colors.Input.R, theme.Colors.Input.G, theme.Colors.Input.B, theme.Colors.Input.A}
			for y := int(bounds.Y) + 4; y < int(bounds.Y+bounds.H)-4; y++ {
				for x := int(bounds.X) + 24; x < int(bounds.X)+31; x++ {
					if got := tt.Image().RGBAAt(x, y); got != background {
						t.Fatalf("text drawn beside the icon at (%d,%d): %v", x, y, got)
					}
				}
			}
		})
	}
}
