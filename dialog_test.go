package kit

import (
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

// dialogFixture 是 DialogWith 测试夹具：记录确认次数和正文里的两个输入。
type dialogFixture struct {
	open       bool
	confirmed  int
	opt        DialogOpt
	title      string
	notes      string
	withInput  bool
	withNotes  bool
	tt         *ui.Tester
	lastLocale string
}

// newDialogFixture 创建已打开对话框的夹具；withInput 加单行输入，withNotes 加多行输入。
func newDialogFixture(opt DialogOpt, withInput, withNotes bool) *dialogFixture {
	f := &dialogFixture{open: true, opt: opt, withInput: withInput, withNotes: withNotes}
	f.tt = ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		if DialogWith(c, &f.open, f.opt, func() {
			if f.withInput {
				TextInput(c.Key("title"), &f.title).Label("Title")
			}
			if f.withNotes {
				TextArea(c.Key("notes"), &f.notes).Label("Notes")
			}
		}) {
			f.confirmed++
		}
	}, 700, 500)
	f.tt.Frame()
	return f
}

// TestDialogWith 验证对话框页脚、Enter 确认、Esc 取消、初始焦点和文案。
func TestDialogWith(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"footer-labels-and-description", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Compact", Description: "Older turns are summarized."}, false, false)
			for _, want := range []string{"Compact", "Older turns are summarized.", "Cancel", "Confirm"} {
				if !f.tt.HasText(want) {
					t.Fatalf("missing %q in %q", want, f.tt.Texts())
				}
			}
			cancel, _ := f.tt.Find("Cancel")
			confirm, _ := f.tt.Find("Confirm")
			if confirm.X <= cancel.X {
				t.Fatalf("confirm %v should be right of cancel %v", confirm, cancel)
			}
			f.tt.SetLocale("zh-Hans")
			if !f.tt.HasText("取消") || !f.tt.HasText("确定") {
				t.Fatal(f.tt.Texts())
			}
		}},
		{"click-confirm-and-cancel", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Delete", ConfirmLabel: "Delete it", Danger: true}, false, false)
			f.tt.Click("Delete it")
			if f.confirmed != 1 || f.open {
				t.Fatalf("confirmed=%d open=%v", f.confirmed, f.open)
			}
			f = newDialogFixture(DialogOpt{Title: "Delete"}, false, false)
			f.tt.Click("Cancel")
			if f.confirmed != 0 || f.open {
				t.Fatalf("confirmed=%d open=%v", f.confirmed, f.open)
			}
		}},
		{"enter-confirms-from-single-line-input", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Rename", ConfirmLabel: "Save"}, true, false)
			if !f.tt.Focused("Title") {
				t.Fatal("first input did not take the initial focus")
			}
			f.tt.Type("New name")
			f.tt.Key(0, ui.KeyEnter)
			if f.confirmed != 1 || f.open || f.title != "New name" {
				t.Fatalf("confirmed=%d open=%v title=%q", f.confirmed, f.open, f.title)
			}
		}},
		{"enter-in-text-area-inserts-newline", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Notes"}, false, true)
			f.tt.Type("a")
			f.tt.Key(0, ui.KeyEnter)
			f.tt.Type("b")
			if f.confirmed != 0 || !f.open || f.notes != "a\nb" {
				t.Fatalf("confirmed=%d open=%v notes=%q", f.confirmed, f.open, f.notes)
			}
		}},
		{"escape-cancels", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Rename"}, true, false)
			f.tt.Key(0, ui.KeyEscape)
			if f.confirmed != 0 || f.open {
				t.Fatalf("confirmed=%d open=%v", f.confirmed, f.open)
			}
		}},
		{"no-input-focuses-confirm-and-enter-confirms", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Delete", ConfirmLabel: "Delete it", Danger: true}, false, false)
			if !f.tt.Focused("Delete it") {
				t.Fatal("confirm button did not take the initial focus")
			}
			f.tt.Key(0, ui.KeyEnter)
			if f.confirmed != 1 || f.open {
				t.Fatalf("confirmed=%d open=%v", f.confirmed, f.open)
			}
		}},
		{"enter-on-cancel-button-cancels", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Delete"}, false, false)
			f.tt.Key(0, ui.KeyTab)
			if !f.tt.Focused("Cancel") {
				t.Fatal("Tab did not reach Cancel")
			}
			f.tt.Key(0, ui.KeyEnter)
			if f.confirmed != 0 || f.open {
				t.Fatalf("confirmed=%d open=%v", f.confirmed, f.open)
			}
		}},
		{"confirm-disabled-blocks-click-and-enter", func(t *testing.T) {
			f := newDialogFixture(DialogOpt{Title: "Rename", ConfirmDisabled: true}, true, false)
			f.tt.Click("Confirm")
			f.tt.Key(0, ui.KeyEnter)
			if f.confirmed != 0 || !f.open {
				t.Fatalf("confirmed=%d open=%v", f.confirmed, f.open)
			}
			f.opt.ConfirmDisabled = false
			f.tt.Frame()
			f.tt.Key(0, ui.KeyEnter)
			if f.confirmed != 1 || f.open {
				t.Fatalf("after enabling: confirmed=%d open=%v", f.confirmed, f.open)
			}
		}},
		{"initial-focus-skips-tabs-and-takes-first-editable-input", func(t *testing.T) {
			f := newDialogTabsFixture(false)
			if f.tt.Focused("General") || f.tt.Focused("Advanced") {
				t.Fatal("initial focus is on a tab")
			}
			if !f.tt.Focused("Title") {
				t.Fatal("initial focus is not on the editable input")
			}
			f.tt.Type("ab")
			if f.title != "providerab" {
				t.Fatalf("title=%q", f.title)
			}
		}},
		{"initial-focus-without-editable-input-is-confirm-without-ring", func(t *testing.T) {
			f := newDialogTabsFixture(true)
			if f.tt.Focused("General") || f.tt.Focused("Advanced") {
				t.Fatal("initial focus is on a tab")
			}
			if !f.tt.Focused("Confirm") {
				t.Fatal("initial focus is not on the confirm button")
			}
			if n := countColor(f.tt.Image(), OneLight().Colors.Focus); n != 0 {
				t.Fatalf("%d focus ring pixels on open", n)
			}
			// 焦点离开确认按钮后环照常显示。
			f.tt.Key(0, ui.KeyTab)
			if n := countColor(f.tt.Image(), OneLight().Colors.Focus); n == 0 {
				t.Fatal("focus ring missing after Tab")
			}
		}},
		{"enter-confirms-without-editable-input", func(t *testing.T) {
			f := newDialogTabsFixture(true)
			f.tt.Key(0, ui.KeyEnter)
			if f.confirmed != 1 || f.open {
				t.Fatalf("confirmed=%d open=%v", f.confirmed, f.open)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

// TestDialogWithSearchableSelect 验证对话框里带搜索的 Select：搜索框按 Enter 只选中并关闭下拉，Esc 只关闭下拉，都不确认或关闭对话框。
func TestDialogWithSearchableSelect(t *testing.T) {
	// 1. 每个用例打开对话框、展开下拉、输入过滤词后按给定的键。
	tests := []struct {
		name         string
		key          ui.Key
		wantValue    string
		wantDropdown bool
	}{
		{"enter-selects-match-keeps-dialog", ui.KeyEnter, "b", false},
		{"escape-closes-dropdown-only", ui.KeyEscape, "", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			open, confirmed, value := true, 0, ""
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if DialogWith(c, &open, DialogOpt{Title: "Agent"}, func() {
					Select(c.Key("agent"), &value, SelectOpt{Label: "Agent", Searchable: true, SearchPlaceholder: "Find agent", Bordered: true, Options: []SelectOption{{Value: "a", Label: "Alpha"}, {Value: "b", Label: "Beta"}}})
				}) {
					confirmed++
				}
			}, 700, 500)
			tt.Frame()
			tt.Click("Agent")
			if !tt.HasText("Find agent") {
				t.Fatal("dropdown did not open")
			}
			tt.Type("be")
			tt.Key(0, test.key)
			tt.Frame()
			dropdownOpen := tt.HasText("Find agent")
			if value != test.wantValue || dropdownOpen != test.wantDropdown || confirmed != 0 || !open {
				t.Fatalf("value=%q dropdown=%v confirmed=%d dialogOpen=%v", value, dropdownOpen, confirmed, open)
			}
		})
	}
}

// dialogTabsFixture 是正文以页签开头的 DialogWith 夹具，在键盘聚焦可见时打开。
type dialogTabsFixture struct {
	open      bool
	confirmed int
	title     string
	tab       int
	tt        *ui.Tester
}

// newDialogTabsFixture 先按一次 Tab 让键盘焦点可见，再打开对话框；readOnly 为真时正文只有只读输入。
func newDialogTabsFixture(readOnly bool) *dialogTabsFixture {
	f := &dialogTabsFixture{title: "provider"}
	f.tt = ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		ui.Column(c).Padding(16).Children(func() { Button(c, "Open") })
		if DialogWith(c, &f.open, DialogOpt{Title: "Settings"}, func() {
			Tabs(c, &f.tab, []string{"General", "Advanced"})
			TextInputWith(c.Key("title"), &f.title, InputOpt{ReadOnly: readOnly}).Label("Title")
		}) {
			f.confirmed++
		}
	}, 700, 500)
	f.tt.Key(0, ui.KeyTab)
	f.open = true
	f.tt.Frame()
	return f
}

// countColor 返回图像中恰好等于 c 的像素个数。
func countColor(img *image.RGBA, c ui.Color) int {
	want, n := color.RGBA{c.R, c.G, c.B, c.A}, 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.RGBAAt(x, y) == want {
				n++
			}
		}
	}
	return n
}
