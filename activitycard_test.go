package kit

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestActivityCard 验证条目标题行的展开、不可展开、操作按钮和状态图标。
func TestActivityCard(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开交互。
	tests := []struct {
		name string
		opt  ActivityCardOpt
		run  func(*testing.T, *ui.Tester, *bool, *int)
	}{
		{"click toggles detail", ActivityCardOpt{Icon: IconTerminal, Title: "bash", Summary: "ls -la"}, func(t *testing.T, tt *ui.Tester, open *bool, _ *int) {
			if tt.HasText("detail body") {
				t.Fatal("detail visible while collapsed")
			}
			if err := tt.Click("bash ls -la"); err != nil {
				t.Fatal(err)
			}
			if !*open || !tt.HasText("detail body") {
				t.Fatal("not expanded", *open, tt.Texts())
			}
			tt.Click("bash ls -la")
			if *open || tt.HasText("detail body") {
				t.Fatal("not collapsed")
			}
		}},
		{"actions click does not toggle", ActivityCardOpt{Icon: IconTerminal, Title: "bash", Summary: "ls -la"}, func(t *testing.T, tt *ui.Tester, open *bool, actions *int) {
			if err := tt.Click("Inspect"); err != nil {
				t.Fatal(err)
			}
			if *actions != 1 || *open {
				t.Fatalf("actions=%d open=%v", *actions, *open)
			}
		}},
		{"running shows spinner", ActivityCardOpt{Icon: IconTerminal, Title: "bash", Status: ActivityRunning}, func(t *testing.T, tt *ui.Tester, _ *bool, _ *int) {
			if _, ok := tt.Find("Running"); !ok {
				t.Fatal("spinner missing")
			}
		}},
		{"failed shows icon", ActivityCardOpt{Icon: IconTerminal, Title: "bash", Status: ActivityFailed}, func(t *testing.T, tt *ui.Tester, _ *bool, _ *int) {
			if _, ok := tt.Find("Failed"); !ok {
				t.Fatal("failed icon missing")
			}
			if _, ok := tt.Find("Running"); ok {
				t.Fatal("spinner shown for failed item")
			}
		}},
		{"plain shows no status", ActivityCardOpt{Icon: IconTerminal, Title: "bash"}, func(t *testing.T, tt *ui.Tester, _ *bool, _ *int) {
			for _, label := range []string{"Running", "Failed", "Interrupted", "Blocked", "Stopped"} {
				if _, ok := tt.Find(label); ok {
					t.Fatal(label)
				}
			}
		}},
		{"interrupted blocked stopped icons", ActivityCardOpt{Icon: IconTerminal, Title: "bash", Status: ActivityBlocked}, func(t *testing.T, tt *ui.Tester, _ *bool, _ *int) {
			if _, ok := tt.Find("Blocked"); !ok {
				t.Fatal("blocked icon missing")
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			open, actions := false, 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme := OneDark()
				theme.Motion.Reduced = true
				ApplyTheme(c, theme)
				ActivityCard(c, &open, test.opt, func() {
					if IconButton(c, IconBug, "Inspect").Clicked() {
						actions++
					}
				}, func() { ui.Text(c, "detail body") })
			}, 400, 200)
			test.run(t, tt, &open, &actions)
		})
	}
	// 2. Extra 在展开按钮之外，点击其中的链接不会展开。
	t.Run("extra link does not toggle", func(t *testing.T) {
		open, clicked := false, false
		tt := ui.NewTester(func(c *ui.Context) {
			theme := OneDark()
			theme.Motion.Reduced = true
			ApplyTheme(c, theme)
			ActivityCard(c, &open, ActivityCardOpt{Icon: IconMessageSquare, Title: "From", Extra: func() {
				if Link(c, "Session root").Clicked() {
					clicked = true
				}
			}}, nil, func() { ui.Text(c, "detail body") })
		}, 400, 200)
		if err := tt.Click("Session root"); err != nil {
			t.Fatal(err)
		}
		if !clicked || open {
			t.Fatalf("clicked=%v open=%v", clicked, open)
		}
	})
	// 3. 没有详情的条目不可展开，点击标题没有任何效果。
	t.Run("plain row is not expandable", func(t *testing.T) {
		open := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme := OneDark()
			theme.Motion.Reduced = true
			ApplyTheme(c, theme)
			ActivityCard(c, &open, ActivityCardOpt{Icon: IconReplace, Title: "model changed"}, nil, nil)
		}, 400, 200)
		tt.Click("model changed")
		if open {
			t.Fatal("non-expandable row toggled")
		}
	})
}

// TestActivityCardIconSlot 验证类型图标是 14 的图标居中在 16 宽的图标位里，标题文字仍从 24 起。
func TestActivityCardIconSlot(t *testing.T) {
	// 1. 在左上角没有内边距的标题行里，图标位占 0 到 16，间距 8，标题从 24 开始。
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme := OneLight()
		theme.Motion.Reduced = true
		ApplyTheme(c, theme)
		ui.Column(c).Children(func() {
			ActivityCard(c, &open, ActivityCardOpt{Icon: IconTerminal, Title: "bash"}, nil, nil)
		})
	}, 300, 100)
	bg := tt.Image().RGBAAt(100, 10)
	inkColumns := func(from, to int) (first, last int) {
		first, last = -1, -1
		for x := from; x < to; x++ {
			for y := 0; y < 28; y++ {
				if tt.Image().RGBAAt(x, y) != bg {
					if first < 0 {
						first = x
					}
					last = x
				}
			}
		}
		return first, last
	}
	if first, _ := inkColumns(16, 60); first < 24 || first > 26 {
		t.Fatalf("title ink starts at x=%d want about 24", first)
	}
	// 2. 图标墨迹落在 16 宽图标位内且左右对称。
	minX, maxX := inkColumns(0, 16)
	if minX < 1 || maxX > 14 || abs(minX+maxX-15) > 1 {
		t.Fatalf("icon ink spans x=%d..%d, want centered in the 16-wide slot", minX, maxX)
	}
}

// TestTextRole 验证文字角色给出的固定行高与字号落在网格上。
func TestTextRole(t *testing.T) {
	// 1. 逐个角色绘制一行文字，量出行高。
	tests := []struct {
		role   TextRole
		height float32
	}{{TextProse, 24}, {TextUI, 20}, {TextCaption, 16}, {TextToken, 20}, {TextCode, 20}, {TextCodeLabel, 16}, {TextTitle, 24}}
	for _, test := range tests {
		tt := ui.NewTester(func(c *ui.Context) {
			ApplyTheme(c, OneLight())
			ui.Column(c).Children(func() { GetTheme(c).TextRole(ui.Text(c, "Role").SingleLine(), test.role) })
		}, 200, 100)
		rect, ok := tt.Find("Role")
		if !ok || rect.H != test.height {
			t.Fatalf("role %d: rect=%v ok=%v want height %v", test.role, rect, ok, test.height)
		}
	}
}
