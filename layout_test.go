package kit

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLayouts 验证模态焦点、虚拟列表和分页锚点。
func TestLayouts(t *testing.T) {
	// 1. 对话框使用原生焦点陷阱并响应 Escape。
	t.Run("dialog-focus", func(t *testing.T) {
		open := false
		tt := ui.NewTester(func(c *ui.Context) {
			ApplyTheme(c, OneLight())
			if Button(c, "Open").Clicked() {
				open = true
			}
			Dialog(c, &open, "Confirm", func() { Button(c, "Accept").AutoFocus(); Button(c, "Cancel") })
		}, 600, 400)
		tt.Click("Open")
		if !tt.HasText("Accept") {
			t.Fatal("dialog missing")
		}
		tt.Key(0, ui.KeyEscape)
		if open {
			t.Fatal("Escape did not close")
		}
	})
	// 2. 变高列表插入历史后保持可见项位置。
	t.Run("list-anchor", func(t *testing.T) {
		items := make([]int, 100)
		for i := range items {
			items[i] = i + 100
		}
		state := ui.ListState{FollowEnd: true}
		state.Key = func(i int) any { return items[i] }
		tt := ui.NewTester(func(c *ui.Context) {
			ApplyTheme(c, OneLight())
			List(c, &state, len(items), func(i int) { ui.Text(c, fmt.Sprintf("row-%d", items[i])).Height(float32(24 + items[i]%3*8)) }).Fill()
		}, 400, 300)
		tt.Scroll(100, 100, 0, -180)
		anchor := ""
		var before ui.Rect
		for _, label := range tt.Texts() {
			if r, ok := tt.Find(label); ok && r.Y >= 0 && r.Y+r.H < 300 {
				anchor, before = label, r
				break
			}
		}
		if anchor == "" {
			t.Fatal("anchor missing")
		}
		history := make([]int, 50)
		for i := range history {
			history[i] = i + 50
		}
		items = append(history, items...)
		tt.Frame()
		after, ok := tt.Find(anchor)
		if !ok || math.Abs(float64(after.Y-before.Y)) > .1 {
			t.Fatalf("anchor moved %v => %v", before, after)
		}
	})
}

// TestTooltipDelay 验证不同锚点均等待 600ms，不继承上一提示的热启动。
func TestTooltipDelay(t *testing.T) {
	// 1. 第一个提示等待完整延迟。
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		ui.Row(c).Padding(40).Gap(8).Children(func() {
			a := Button(c, "First")
			Tooltip(c, a, "First hint")
			b := Button(c, "Second")
			Tooltip(c, b, "Second hint")
		})
	}, 400, 200)
	a, _ := tt.Find("First")
	tt.Move(a.X+a.W/2, a.Y+a.H/2)
	tt.Advance(599 * time.Millisecond)
	if tt.HasText("First hint") {
		t.Fatal("early first tooltip")
	}
	tt.Advance(2 * time.Millisecond)
	if !tt.HasText("First hint") {
		t.Fatal("first tooltip missing")
	}
	// 2. 第二个锚点独立计时。
	b, _ := tt.Find("Second")
	tt.Move(b.X+b.W/2, b.Y+b.H/2)
	tt.Advance(200 * time.Millisecond)
	if tt.HasText("Second hint") {
		t.Fatal("warm tooltip shown early")
	}
	tt.Advance(401 * time.Millisecond)
	if !tt.HasText("Second hint") {
		t.Fatal("second tooltip missing")
	}
}

// TestPhysicalDividerPixels 验证不同缩放下共享分隔线连续、整像素且没有双线。
func TestPhysicalDividerPixels(t *testing.T) {
	// 1. 使用高对比试验色隔离线条，正文区域无额外边框。
	for _, scale := range []float32{1, 1.25, 1.5, 2} {
		t.Run(fmt.Sprint(scale), func(t *testing.T) {
			theme := OneLight()
			theme.Colors.Panel = ui.RGB(255, 255, 255)
			theme.Colors.Editor = ui.RGB(255, 255, 255)
			theme.Colors.Border = ui.RGB(0, 0, 0)
			theme.Colors.BorderSubtle = theme.Colors.Border
			width := float32(80)
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Fill().Gap(0).Children(func() {
					Header(c, "", nil)
					Split(c, &width, func() { ui.Box(c).Fill() }, func() { ui.Box(c).Fill() }).Grow(1)
				})
			}, 240, 180)
			tt.SetScale(scale)
			img := tt.Image()
			black := color.RGBA{0, 0, 0, 255}
			white := color.RGBA{255, 255, 255, 255}
			want := int(math.Round(float64(scale)))
			// 2. 水平和垂直截面只有一段实色边线，无透明或灰色过渡像素。
			for _, axis := range []string{"horizontal", "vertical"} {
				count := 0
				segments := 0
				inside := false
				limit := img.Bounds().Dx()
				if axis == "horizontal" {
					limit = int(60 * scale)
				}
				for i := 0; i < limit; i++ {
					var pixel color.RGBA
					if axis == "horizontal" {
						pixel = img.RGBAAt(int(180*scale), i)
					} else {
						pixel = img.RGBAAt(i, int(100*scale))
					}
					if pixel == black {
						count++
						if !inside {
							segments++
						}
						inside = true
					} else {
						inside = false
						if pixel != white {
							t.Fatalf("%s scale %v blurred pixel %v at %d", axis, scale, pixel, i)
						}
					}
				}
				if count != want || segments != 1 {
					t.Fatalf("%s scale=%v pixels=%d segments=%d want pixels=%d", axis, scale, count, segments, want)
				}
			}
			if dir := os.Getenv("KIT_CAPTURE_DIR"); dir != "" {
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
				file, err := os.Create(filepath.Join(dir, fmt.Sprintf("dividers-%g.png", scale)))
				if err != nil {
					t.Fatal(err)
				}
				if err := png.Encode(file, img); err != nil {
					t.Fatal(err)
				}
				file.Close()
			}
		})
	}
}

// TestHeader 验证标题栏的高度和左内边距，对话区可用 16 对齐消息列。
func TestHeader(t *testing.T) {
	// 1. 默认 12，HeaderOpt 可改成 16；高度始终是标题栏档 36。
	tests := []struct {
		name  string
		opt   *HeaderOpt
		wantX float32
	}{
		{"default-inset-12", nil, 12},
		{"sidebar-inset-12", &HeaderOpt{Inset: 12}, 12},
		{"conversation-inset-16", &HeaderOpt{Inset: 16}, 16},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if test.opt == nil {
					bounds = Header(c, "Title", nil).Bounds()
					return
				}
				bounds = HeaderWith(c, "Title", *test.opt, nil).Bounds()
			}, 300, 100)
			tt.Frame()
			title, ok := tt.Find("Title")
			if !ok || title.X != test.wantX || bounds.H != 36 {
				t.Fatalf("title=%v ok=%v header=%v", title, ok, bounds)
			}
		})
	}
}

// TestHeaderLeading 验证自定义左侧内容：替换标题文字，仍从 Inset 起、36 高、右侧操作保留。
func TestHeaderLeading(t *testing.T) {
	// 1. 有 Leading 时标题文字不显示；没有时显示标题。
	tests := []struct {
		name      string
		leading   func(c *ui.Context)
		wantTitle bool
		wantX     float32
		wantText  string
	}{
		{"title-when-no-leading", nil, true, 16, "Title"},
		{"leading-replaces-title", func(c *ui.Context) { ui.Text(c, "Leading") }, false, 16, "Leading"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				opt := HeaderOpt{Inset: 16}
				if test.leading != nil {
					opt.Leading = func() { test.leading(c) }
				}
				bounds = HeaderWith(c, "Title", opt, func() { ui.Text(c, "Action") }).Bounds()
			}, 300, 100)
			tt.Frame()
			text, ok := tt.Find(test.wantText)
			if !ok || text.X != test.wantX || bounds.H != 36 || tt.HasText("Title") != test.wantTitle || !tt.HasText("Action") {
				t.Fatalf("text=%v ok=%v header=%v texts=%v", text, ok, bounds, tt.Texts())
			}
		})
	}
}

// TestFrameHeaderCopy 验证内容框的复制按钮：复制按钮是 20 见方的纯图标按钮，复制后 2 秒内显示「已复制」，随后恢复，语言跟随窗口；禁用时不复制。
func TestFrameHeaderCopy(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动 JSON 阅读器的标题行。
	state, _ := NewJSONState([]byte("null"))
	tests := []struct {
		name string
		run  func(*testing.T, *ui.Tester)
	}{
		{"copy-button-is-icon-only-20-square", func(t *testing.T, tt *ui.Tester) {
			button, ok := tt.Find("Copy JSON")
			if !ok || button.W != 20 || button.H != 20 {
				t.Fatalf("button %v ok=%v texts=%v", button, ok, tt.Texts())
			}
			// 悬停停留后提示出现，文字就是复制标签；按钮本身不再显示文字。
			tt.Move(button.X+button.W/2, button.Y+button.H/2)
			tt.Advance(700 * time.Millisecond)
			shown := 0
			for _, text := range tt.Texts() {
				if text == "Copy JSON" {
					shown++
				}
			}
			if shown != 2 {
				t.Fatalf("want accessible label plus tooltip, texts=%v", tt.Texts())
			}
		}},
		{"copied-feedback-lasts-two-seconds", func(t *testing.T, tt *ui.Tester) {
			tt.Click("Copy JSON")
			if tt.Clipboard() != "null" || !tt.HasText("Copied") || tt.HasText("Copy JSON") {
				t.Fatalf("after click: clipboard=%q texts=%v", tt.Clipboard(), tt.Texts())
			}
			// 指针一直停在按钮上，悬停提示的帧不能吞掉到期唤醒。
			tt.Advance(1900 * time.Millisecond)
			if !tt.HasText("Copied") {
				t.Fatal("feedback ended before 2 seconds")
			}
			tt.Advance(200 * time.Millisecond)
			if tt.HasText("Copied") || !tt.HasText("Copy JSON") {
				t.Fatalf("feedback did not end: %v", tt.Texts())
			}
		}},
		{"copied-follows-locale", func(t *testing.T, tt *ui.Tester) {
			tt.SetLocale("zh-Hans")
			tt.Click("复制 JSON")
			if !tt.HasText("已复制") {
				t.Fatal(tt.Texts())
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); JSONTree(c, state, JSONOpt{}).Fill() }, 400, 240)
			test.run(t, tt)
		})
	}
}

// TestEmptyState 检查空状态在区域中央：标题和说明水平居中，带主操作时按钮也居中在文字下方，点击触发 OnAction。
func TestEmptyState(t *testing.T) {
	// 1. 无操作、有操作两种用例，区域宽 400。
	tests := []struct {
		name   string
		action bool
	}{
		{"text-only", false},
		{"with-action", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clicked := false
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				ui.Column(c).Fill().Children(func() {
					opt := EmptyStateOpt{}
					if test.action {
						opt = EmptyStateOpt{ActionLabel: "Add provider", OnAction: func() { clicked = true }}
					}
					EmptyStateWith(c, "Nothing here", "Detail", opt)
				})
			}, 400, 300)
			tt.Frame()
			title, ok := tt.Find("Nothing here")
			detail, okDetail := tt.Find("Detail")
			if !ok || !okDetail {
				t.Fatalf("title=%v detail=%v", title, detail)
			}
			centered := func(r ui.Rect) bool { return abs(int(r.X+r.W/2-200)) <= 1 }
			if !centered(title) || !centered(detail) || detail.Y < title.Y+title.H || title.Y < 60 {
				t.Fatalf("not centered: title=%v detail=%v", title, detail)
			}
			if !test.action {
				return
			}
			button, ok := tt.Find("Add provider")
			if !ok || button.Y < detail.Y+detail.H || !centered(button) {
				t.Fatalf("action misplaced: %v detail=%v", button, detail)
			}
			tt.ClickAt(button.X+button.W/2, button.Y+button.H/2)
			if !clicked {
				t.Fatal("OnAction was not called")
			}
		})
	}
}
