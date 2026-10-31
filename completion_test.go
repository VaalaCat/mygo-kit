package kit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// completionFixture 是补全测试的输入区：底部一个带边框的锚点，里面是多行输入。
type completionFixture struct {
	text     string
	opt      CompletionOpt
	selected []int
	sent     int
	anchor   ui.Rect
	tt       *ui.Tester
}

// newCompletionFixture 创建测试夹具；opt.Open 为空时按输入以 / 开头决定是否打开，Enter 不被补全接管时计为发送。
func newCompletionFixture(text string, opt CompletionOpt) *completionFixture {
	f := &completionFixture{text: text, opt: opt}
	autoOpen := !opt.Open
	f.tt = ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		ui.Column(c).Fill().Justify(ui.End).Padding(20).Children(func() {
			anchor := ui.Column(c.Key("anchor")).Padding(8).Border(1, OneLight().Colors.Border)
			var input ui.Element
			anchor.Children(func() {
				input = TextAreaWith(c.Key("input"), &f.text, InputOpt{Borderless: true}).Label("Message").MinHeight(40)
				input.OnTextKey(func(k ui.TextKey) bool {
					if k.Key == ui.KeyEnter && k.Mods == 0 {
						f.sent++
						return true
					}
					return false
				})
			})
			o := f.opt
			if autoOpen {
				o.Open = strings.HasPrefix(f.text, "/")
			}
			if i := CompletionMenu(c, anchor, input, o); i >= 0 {
				f.selected = append(f.selected, i)
			}
			f.anchor = anchor.Bounds()
		})
	}, 500, 500)
	f.tt.Click("Message")
	return f
}

// completionItems 返回 n 个候选 item-0 … item-(n-1)。
func completionItems(n int) []CompletionItem {
	items := make([]CompletionItem, n)
	for i := range items {
		items[i] = CompletionItem{Icon: IconBookOpen, Title: fmt.Sprintf("item-%d", i), Description: fmt.Sprintf("description %d", i)}
	}
	return items
}

// TestCompletionMenu 验证补全浮层的焦点、按键、选择、状态行和布局。
func TestCompletionMenu(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"open-keeps-input-focus-and-typing", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(3)})
			f.tt.Frame()
			if !f.tt.HasText("item-0") || !f.tt.HasText("description 2") {
				t.Fatal(f.tt.Texts())
			}
			f.tt.Type("re")
			if f.text != "/re" || !f.tt.Focused("Message") {
				t.Fatalf("text=%q focused=%v", f.text, f.tt.Focused("Message"))
			}
		}},
		{"closed-shows-nothing", func(t *testing.T) {
			f := newCompletionFixture("hello", CompletionOpt{Items: completionItems(3)})
			if f.tt.HasText("item-0") {
				t.Fatal("menu visible while closed")
			}
		}},
		{"arrows-wrap-and-enter-selects-without-sending", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(3)})
			f.tt.Key(0, ui.KeyDown)
			f.tt.Key(0, ui.KeyEnter)
			f.tt.Frame()
			f.text = ""
			f.tt.Frame()
			f.text = "/"
			f.tt.Frame()
			f.tt.Key(0, ui.KeyUp)
			f.tt.Key(0, ui.KeyEnter)
			if len(f.selected) != 2 || f.selected[0] != 1 || f.selected[1] != 2 || f.sent != 0 {
				t.Fatalf("selected=%v sent=%d", f.selected, f.sent)
			}
		}},
		{"tab-selects", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(3)})
			f.tt.Key(0, ui.KeyTab)
			if len(f.selected) != 1 || f.selected[0] != 0 || !f.tt.Focused("Message") {
				t.Fatalf("selected=%v focused=%v", f.selected, f.tt.Focused("Message"))
			}
		}},
		{"click-selects-and-input-keeps-focus", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(3)})
			f.tt.Click("item-2")
			if len(f.selected) != 1 || f.selected[0] != 2 || !f.tt.Focused("Message") {
				t.Fatalf("selected=%v focused=%v", f.selected, f.tt.Focused("Message"))
			}
		}},
		{"escape-closes-until-reopened", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(3)})
			f.tt.Key(0, ui.KeyEscape)
			if f.tt.HasText("item-0") {
				t.Fatal("Escape did not close")
			}
			f.tt.Type("a")
			if f.tt.HasText("item-0") || f.text != "/a" {
				t.Fatalf("menu came back while still typing the same word: text=%q", f.text)
			}
			f.text = "plain"
			f.tt.Frame()
			f.text = "/"
			f.tt.Frame()
			if !f.tt.HasText("item-0") {
				t.Fatal("menu did not reopen after Open went false")
			}
		}},
		{"no-items-leaves-keys-to-the-editor", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Empty: "No matching skills"})
			if !f.tt.HasText("No matching skills") {
				t.Fatal(f.tt.Texts())
			}
			f.tt.Key(0, ui.KeyEnter)
			if f.sent != 1 || len(f.selected) != 0 {
				t.Fatalf("sent=%d selected=%v", f.sent, f.selected)
			}
		}},
		{"loading-error-and-default-empty-rows-are-visible", func(t *testing.T) {
			// 候选之后的加载行、错误行同样要在可视区里：高度大于 0，并在最后一个候选的下方。
			for _, test := range []struct {
				opt  CompletionOpt
				want string
			}{
				{CompletionOpt{Loading: true}, "Loading…"},
				{CompletionOpt{Error: "Could not load skills"}, "Could not load skills"},
				{CompletionOpt{}, "No results"},
				{CompletionOpt{Items: completionItems(2), Loading: true}, "Loading…"},
				{CompletionOpt{Items: completionItems(2), Error: "Could not load skills"}, "Could not load skills"},
				{CompletionOpt{Items: completionItems(6), Loading: true, Error: "Could not load skills"}, "Could not load skills"},
			} {
				f := newCompletionFixture("/", test.opt)
				f.tt.Frame()
				row, ok := f.tt.Find(test.want)
				if !ok || row.H <= 0 {
					t.Fatalf("%+v: row %v ok=%v texts %v", test.opt, row, ok, f.tt.Texts())
				}
				if len(test.opt.Items) > 0 {
					last, _ := f.tt.Find(fmt.Sprintf("item-%d", len(test.opt.Items)-1))
					if row.Y < last.Y+last.H {
						t.Fatalf("%+v: status row %v is not below the last item %v", test.opt, row, last)
					}
				}
			}
		}},
		{"composing-keeps-enter-for-the-input-method", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(3)})
			f.tt.Frame()
			f.tt.Compose("ni", 2)
			f.tt.Key(0, ui.KeyEnter)
			if len(f.selected) != 0 {
				t.Fatalf("Enter during composition selected %v", f.selected)
			}
		}},
		{"floats-above-without-moving-the-input", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(3)})
			open := f.anchor
			f.text = "plain"
			f.tt.Frame()
			if closed := f.anchor; closed != open {
				t.Fatalf("anchor moved %v => %v", open, closed)
			}
			f.text = "/"
			f.tt.Frame()
			first, _ := f.tt.Find("item-0")
			if first.Y+first.H > open.Y {
				t.Fatalf("row %v is not above the input %v", first, open)
			}
		}},
		{"eight-rows-then-scroll-keeps-highlight-visible", func(t *testing.T) {
			f := newCompletionFixture("/", CompletionOpt{Items: completionItems(12)})
			top, _ := f.tt.Find("item-0")
			if seventh, _ := f.tt.Find("item-7"); seventh.H == 0 {
				t.Fatal("8th row not visible")
			}
			if ninth, _ := f.tt.Find("item-8"); ninth.H != 0 {
				t.Fatalf("9th row visible without scrolling: %v", ninth)
			}
			for i := 0; i < 10; i++ {
				f.tt.Key(0, ui.KeyDown)
			}
			f.tt.Frame()
			row, ok := f.tt.Find("item-10")
			if !ok || row.Y < top.Y-1 || row.Y+row.H > f.anchor.Y {
				t.Fatalf("highlighted row %v outside the list [%v, %v)", row, top.Y, f.anchor.Y)
			}
			// 面板最高 8 行，滚动后第一行不再显示（被裁掉的行没有区域）。
			if first, _ := f.tt.Find("item-0"); first.H != 0 && first.Y >= top.Y {
				t.Fatalf("list did not scroll: item-0 %v => %v", top, first)
			}
		}},
		{"row-text-aligns-with-input-text", func(t *testing.T) {
			// 输入框和候选行的文字都是 "Alpha"，两者最左的笔画像素应在同一列（笔画抗锯齿允许差 1 像素，未对齐时相差 3）。
			f := newCompletionFixture("Alpha", CompletionOpt{Open: true, Items: []CompletionItem{{Title: "Alpha"}, {Title: "Alphabet"}}})
			f.tt.Frame()
			row, rowOK := f.tt.Find("Alphabet")
			input, inputOK := f.tt.Find("Message")
			if !rowOK || !inputOK {
				t.Fatalf("row=%v input=%v", rowOK, inputOK)
			}
			first := func(r ui.Rect, y0, y1 float32) int {
				img, bg := f.tt.Image(), f.tt.Image().RGBAAt(int(r.X)+1, int(y0))
				for x := int(r.X) + 1; x < int(r.X)+40; x++ {
					for y := int(y0); y < int(y1); y++ {
						if img.RGBAAt(x, y) != bg {
							return x
						}
					}
				}
				return -1
			}
			got, want := first(ui.Rect{X: row.X}, row.Y+6, row.Y+row.H-6), first(ui.Rect{X: input.X}, input.Y+4, input.Y+28)
			if got < 0 || want < 0 || got-want < -1 || got-want > 1 {
				t.Fatalf("row ink x=%d, input ink x=%d", got, want)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
