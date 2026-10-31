package kit

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestField 验证字段的标签、说明、错误和标签点击聚焦。
func TestField(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	tests := []struct {
		name string
		opt  FieldOpt
		run  func(*testing.T, *ui.Tester, *string)
	}{
		{"label-click-focuses-control", FieldOpt{Label: "Name"}, func(t *testing.T, tt *ui.Tester, value *string) {
			label, _ := tt.Find("Name")
			tt.ClickAt(label.X+1, label.Y+label.H/2)
			tt.Type("abc")
			if *value != "abc" {
				t.Fatalf("typed value %q: label click did not focus the control", *value)
			}
		}},
		{"description-below-control", FieldOpt{Label: "Name", Description: "Shown in the sidebar."}, func(t *testing.T, tt *ui.Tester, _ *string) {
			input, _ := tt.Find("Name")
			description, ok := tt.Find("Shown in the sidebar.")
			if !ok {
				t.Fatal(tt.Texts())
			}
			if description.Y < input.Y+input.H {
				t.Fatalf("description %v not below control %v", description, input)
			}
		}},
		{"error-replaces-description", FieldOpt{Label: "Name", Description: "Shown in the sidebar.", Error: "Required."}, func(t *testing.T, tt *ui.Tester, _ *string) {
			if !tt.HasText("Required.") || tt.HasText("Shown in the sidebar.") {
				t.Fatal(tt.Texts())
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := ""
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				Field(c, test.opt, func() ui.Element { return TextInput(c.Key("input"), &value) })
			}, 400, 200)
			test.run(t, tt, &value)
		})
	}
}

// TestFieldMessageAndInlineHint 验证提示文字的三种语义和加载提示的内置文案。
func TestFieldMessageAndInlineHint(t *testing.T) {
	// 1. 三种语义取不同文字颜色，加载提示在没有文字时用内置「加载中」。
	theme := OneLight()
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, theme)
		ui.Column(c).Children(func() {
			FieldMessage(c.Key("d"), MessageDescription, "plain")
			FieldMessage(c.Key("w"), MessageWarning, "careful")
			FieldMessage(c.Key("e"), MessageError, "broken")
			InlineHint(c.Key("empty"), "Nothing here", false)
			InlineHint(c.Key("loading"), "", true)
		})
	}, 300, 200)
	for _, want := range []string{"plain", "careful", "broken", "Nothing here", "Loading…"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	tt.SetLocale("zh-Hans")
	if !tt.HasText("加载中…") {
		t.Fatal(tt.Texts())
	}
}

// TestInlineHintInset 验证提示的左对齐线：默认 12，InlineHintOpt 可改成 16。
func TestInlineHintInset(t *testing.T) {
	tests := []struct {
		name  string
		opt   InlineHintOpt
		wantX float32
	}{{"default-12", InlineHintOpt{}, 12}, {"conversation-16", InlineHintOpt{Inset: 16}, 16}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				ui.Column(c).Fill().Children(func() { InlineHintWith(c, "Nothing here", false, test.opt) })
			}, 300, 100)
			tt.Frame()
			hint, ok := tt.Find("Nothing here")
			if !ok || hint.X != test.wantX {
				t.Fatalf("hint=%v ok=%v", hint, ok)
			}
		})
	}
}

// TestSettingPage 验证设置页的边距、最大宽度、行高和页头在滚动时保持可见。
func TestSettingPage(t *testing.T) {
	// 1. 内容足够多以便滚动，页头右侧有一个操作按钮。
	theme := OneLight()
	var rowBounds ui.Rect
	var sectionBounds ui.Rect
	view := func(c *ui.Context) {
		ApplyTheme(c, theme)
		SettingPage(c, "Settings", func() { Button(c, "Apply") }, func() {
			for s := 0; s < 4; s++ {
				section := SettingSection(c.Key(s), fmt.Sprintf("Section %d", s), func() {
					for r := 0; r < 4; r++ {
						row := SettingRow(c.Key(r), fmt.Sprintf("Row %d-%d", s, r), "", func() { Button(c, "Edit").Label(fmt.Sprintf("Edit %d-%d", s, r)) })
						if s == 0 && r == 0 {
							rowBounds = row.Bounds()
						}
					}
				})
				if s == 0 {
					sectionBounds = section.Bounds()
				}
			}
		}).Fill()
	}
	tt := ui.NewTester(view, 1200, 400)
	tt.Frame()
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"left-padding-and-max-width", func(t *testing.T) {
			title, _ := tt.Find("Section 0")
			if title.X != 24 {
				t.Fatalf("section title x=%v want 24", title.X)
			}
			header, _ := tt.Find("Settings")
			if header.X != 24 {
				t.Fatalf("page title x=%v want 24", header.X)
			}
			if got := sectionBounds.X + sectionBounds.W; got > 720 {
				t.Fatalf("section right edge %v exceeds the 720 max width", got)
			}
			if sectionBounds.W < 600 {
				t.Fatalf("section width %v: content should fill up to the max width", sectionBounds.W)
			}
		}},
		{"row-height-at-least-36", func(t *testing.T) {
			if rowBounds.H < 36 {
				t.Fatalf("row height %v", rowBounds.H)
			}
		}},
		{"header-stays-while-scrolling", func(t *testing.T) {
			before, _ := tt.Find("Apply")
			last, _ := tt.Find("Section 3")
			tt.Scroll(300, 200, 0, 400)
			after, ok := tt.Find("Apply")
			if !ok || after != before {
				t.Fatalf("actions moved %v => %v", before, after)
			}
			if scrolled, _ := tt.Find("Section 3"); scrolled.Y >= last.Y {
				t.Fatalf("content did not scroll: Section 3 %v => %v", last, scrolled)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
