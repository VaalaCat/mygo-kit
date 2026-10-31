package kit

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// menuTestItems 返回测试用菜单：第二项禁用，第三项前有分隔线。
func menuTestItems(chosen *[]string) []MenuItem {
	pick := func(name string) func() { return func() { *chosen = append(*chosen, name) } }
	return []MenuItem{
		{Icon: IconPencil, Title: "Rename", Shortcut: "F2", OnSelect: pick("Rename")},
		{Title: "Export", Disabled: true, OnSelect: pick("Export")},
		{Separator: true},
		{Title: "Archive", OnSelect: pick("Archive")},
		{Title: "Delete", Danger: true, OnSelect: pick("Delete")},
	}
}

// TestMenuButton 验证菜单按钮的点击、键盘、禁用项、首字母跳转和焦点归还。
func TestMenuButton(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T, *ui.Tester, *[]string)
	}{
		{"click-opens-and-click-item-selects", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			if tt.HasText("Rename") {
				t.Fatal("menu visible while closed")
			}
			tt.Click("More")
			if !tt.HasText("Rename") || !tt.HasText("F2") {
				t.Fatal(tt.Texts())
			}
			tt.Click("Archive")
			if len(*chosen) != 1 || (*chosen)[0] != "Archive" || tt.HasText("Rename") {
				t.Fatalf("chosen=%v texts=%v", *chosen, tt.Texts())
			}
		}},
		{"click-disabled-item-does-nothing", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			tt.Click("More")
			tt.Click("Export")
			if len(*chosen) != 0 || !tt.HasText("Rename") {
				t.Fatalf("chosen=%v texts=%v", *chosen, tt.Texts())
			}
		}},
		{"tab-and-shift-tab-close-the-menu", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			for _, mods := range []ui.Modifiers{0, ui.Shift} {
				tt.Key(0, ui.KeyTab)
				tt.Key(0, ui.KeyEnter)
				if !tt.HasText("Rename") {
					t.Fatalf("mods %v: Enter did not open the menu", mods)
				}
				tt.Key(mods, ui.KeyTab)
				if tt.HasText("Rename") || len(*chosen) != 0 {
					t.Fatalf("mods %v: menu still open or chose %v", mods, *chosen)
				}
			}
		}},
		{"keyboard-open-highlights-first-and-enter-selects", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyEnter)
			if !tt.HasText("Rename") {
				t.Fatal("Enter did not open the menu")
			}
			tt.Key(0, ui.KeyEnter)
			if len(*chosen) != 1 || (*chosen)[0] != "Rename" || tt.HasText("Archive") {
				t.Fatalf("chosen=%v", *chosen)
			}
		}},
		{"arrows-skip-disabled-and-separators-and-wrap", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			tt.Click("More")
			// 鼠标点开时没有高亮：Down 到第一项，再 Down 跳过 Export 和分隔线到 Archive。
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			tt.Click("More")
			// Up 从没有高亮开始取最后一项，再 Down 绕回第一项。
			tt.Key(0, ui.KeyUp)
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			tt.Click("More")
			tt.Key(0, ui.KeyEnd)
			tt.Key(0, ui.KeyEnter)
			tt.Click("More")
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyHome)
			tt.Key(0, ui.KeyEnter)
			want := []string{"Archive", "Rename", "Delete", "Rename"}
			if len(*chosen) != len(want) {
				t.Fatalf("chosen=%v want %v", *chosen, want)
			}
			for i := range want {
				if (*chosen)[i] != want[i] {
					t.Fatalf("chosen=%v want %v", *chosen, want)
				}
			}
		}},
		{"first-letter-jumps", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			tt.Click("More")
			tt.Type("d")
			tt.Key(0, ui.KeyEnter)
			if len(*chosen) != 1 || (*chosen)[0] != "Delete" {
				t.Fatalf("chosen=%v", *chosen)
			}
		}},
		{"escape-closes-and-returns-focus-to-trigger", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			tt.Click("More")
			if tt.Focused("More") {
				t.Fatal("focus stayed on the trigger while the menu is open")
			}
			tt.Key(0, ui.KeyEscape)
			if tt.HasText("Rename") || !tt.Focused("More") || len(*chosen) != 0 {
				t.Fatalf("texts=%v focused=%v chosen=%v", tt.Texts(), tt.Focused("More"), *chosen)
			}
		}},
		{"items-are-28-high-and-panel-is-8-radius-padded", func(t *testing.T, tt *ui.Tester, _ *[]string) {
			tt.Click("More")
			rename, _ := tt.Find("Rename")
			archive, _ := tt.Find("Archive")
			if rename.H != 28 || archive.H != 28 {
				t.Fatalf("item heights %v %v", rename.H, archive.H)
			}
			// 两项之间隔着禁用项（28）和分隔线（1 加上下各 4 的外边距）。
			if got := archive.Y - rename.Y; got != 28+28+9 {
				t.Fatalf("Rename to Archive distance %v", got)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var chosen []string
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				ui.Row(c).Padding(40).Children(func() {
					MenuButton(c, IconButton(c.Key("trigger"), IconEllipsis, "More"), menuTestItems(&chosen))
				})
			}, 400, 300)
			test.run(t, tt, &chosen)
		})
	}
}

// TestContextMenu 验证右键在指针处弹出，以及键盘唤出、选中和关闭。
func TestContextMenu(t *testing.T) {
	// 1. 每个场景用真实离屏输入驱动公开组件。
	tests := []struct {
		name string
		run  func(*testing.T, *ui.Tester, *[]string)
	}{
		{"right-click-opens-at-pointer-and-enter-selects", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			tt.RightClickAt(150, 100)
			rename, ok := tt.Find("Rename")
			if !ok || rename.X < 150 || rename.Y < 100 {
				t.Fatalf("menu at %v, pointer at (150,100)", rename)
			}
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if len(*chosen) != 1 || (*chosen)[0] != "Rename" || tt.HasText("Archive") {
				t.Fatalf("chosen=%v", *chosen)
			}
		}},
		{"shift-f10-and-menu-key-open-below-target", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			for _, key := range []struct {
				mods ui.Modifiers
				key  ui.Key
			}{{ui.Shift, ui.KeyF10}, {0, ui.KeyContextMenu}} {
				tt.Click("Area")
				tt.Key(key.mods, key.key)
				if !tt.HasText("Rename") {
					t.Fatalf("key %v did not open the menu", key)
				}
				tt.Key(0, ui.KeyEscape)
				if tt.HasText("Rename") || !tt.Focused("Area") {
					t.Fatalf("Escape: menu visible=%v, target focused=%v", tt.HasText("Rename"), tt.Focused("Area"))
				}
			}
		}},
		{"click-outside-closes", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			tt.RightClickAt(150, 100)
			tt.ClickAt(390, 290)
			if tt.HasText("Rename") || len(*chosen) != 0 {
				t.Fatalf("texts=%v chosen=%v", tt.Texts(), *chosen)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var chosen []string
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				area := ui.Column(c.Key("area")).Size(360, 160).Margin(20).Focusable().Label("Area")
				ContextMenu(c, area, menuTestItems(&chosen))
			}, 400, 300)
			test.run(t, tt, &chosen)
		})
	}
}

// TestContextMenuOnSelectableText 验证 Selectable 文字容器上的右键菜单：弹出 kit 自己的菜单而不是系统菜单，
// 有选中文字时带「复制」，拖选照常工作。
func TestContextMenuOnSelectableText(t *testing.T) {
	// 1. 每个场景在容器里拖选一段文字后用真实右键驱动。
	const body = "alpha beta gamma"
	tests := []struct {
		name string
		run  func(*testing.T, *ui.Tester, *[]string)
	}{
		{"right-click-on-selection-shows-kit-menu-with-copy", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			x, y := dragSelectText(tt, body, 0, 1)
			tt.RightClickAt(x, y)
			if len(tt.Menu()) != 0 {
				t.Fatalf("native text menu opened: %v", tt.Menu())
			}
			if !tt.HasText("Rename") || !tt.HasText("Copy") {
				t.Fatalf("kit menu missing items: %v", tt.Texts())
			}
			if err := tt.Click("Copy"); err != nil {
				t.Fatal(err)
			}
			if tt.Clipboard() != body || tt.HasText("Rename") {
				t.Fatalf("clipboard=%q texts=%v", tt.Clipboard(), tt.Texts())
			}
		}},
		{"partial-selection-copies-only-selected-text", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			x, y := dragSelectText(tt, body, 0, 0.3)
			tt.RightClickAt(x, y)
			tt.Click("Copy")
			if got := tt.Clipboard(); got == "" || got == body || !strings.HasPrefix(body, got) {
				t.Fatalf("clipboard=%q", got)
			}
		}},
		{"caller-items-still-work-after-copy-item", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			x, y := dragSelectText(tt, body, 0, 1)
			tt.RightClickAt(x, y)
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if len(*chosen) != 1 || (*chosen)[0] != "Rename" || tt.Clipboard() != "" {
				t.Fatalf("chosen=%v clipboard=%q", *chosen, tt.Clipboard())
			}
		}},
		{"right-click-without-selection-has-no-copy", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			text, _ := tt.Find(body)
			tt.RightClickAt(text.X+text.W/2, text.Y+text.H/2)
			if len(tt.Menu()) != 0 || !tt.HasText("Rename") || tt.HasText("Copy") {
				t.Fatalf("native=%v texts=%v", tt.Menu(), tt.Texts())
			}
		}},
		{"dragging-still-selects-text", func(t *testing.T, tt *ui.Tester, chosen *[]string) {
			x, y := dragSelectText(tt, body, 0, 1)
			tt.RightClickAt(x, y)
			if !tt.HasText("Copy") {
				t.Fatal("drag did not select text")
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var chosen []string
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				area := ui.Column(c.Key("area")).Size(360, 160).Margin(20).Selectable().Children(func() {
					ui.Text(c, body).SingleLine()
				})
				ContextMenu(c, area, menuTestItems(&chosen))
			}, 400, 300)
			test.run(t, tt, &chosen)
		})
	}
}

// dragSelectText 在文字 body 上从 from 拖到 to（占文字宽度的比例），返回拖选范围内的一个点。
func dragSelectText(tt *ui.Tester, body string, from, to float32) (x, y float32) {
	text, _ := tt.Find(body)
	y = text.Y + text.H/2
	start := text.X + 1 + (text.W-2)*from
	end := text.X + 1 + (text.W-2)*to
	tt.Press(start, y)
	tt.Move((start+end)/2, y)
	tt.Move(end, y)
	tt.Release(end, y)
	return (start + end) / 2, y
}

// TestMenuClosed 验证菜单关闭的那一帧 MenuButton 和 ContextMenu 都返回 true，调用方可以在这一帧把焦点放到别处。
func TestMenuClosed(t *testing.T) {
	// 1. 每个场景记录返回 true 的帧数。
	tests := []struct {
		name string
		run  func(*testing.T, *ui.Tester)
	}{
		{"menu-button-escape", func(t *testing.T, tt *ui.Tester) {
			tt.Click("More")
			tt.Key(0, ui.KeyEscape)
		}},
		{"menu-button-click-outside", func(t *testing.T, tt *ui.Tester) {
			tt.Click("More")
			tt.ClickAt(390, 290)
		}},
		{"menu-button-item-selected", func(t *testing.T, tt *ui.Tester) {
			tt.Click("More")
			tt.Click("Archive")
		}},
		{"context-menu-escape", func(t *testing.T, tt *ui.Tester) {
			tt.RightClickAt(200, 200)
			tt.Key(0, ui.KeyEscape)
		}},
		{"context-menu-click-outside", func(t *testing.T, tt *ui.Tester) {
			tt.RightClickAt(200, 200)
			tt.ClickAt(5, 5)
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var chosen []string
			closedFrames := 0
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				// 关闭的那一帧把焦点放到旁边的元素上，框架归还给触发元素的焦点不能盖过它。
				other := ui.Box(c.Key("other")).Absolute().Left(300).Top(10).Size(40, 40).Focusable().Label("Other")
				ui.Row(c).Padding(40).Children(func() {
					if MenuButton(c, IconButton(c.Key("trigger"), IconEllipsis, "More"), menuTestItems(&chosen)) {
						closedFrames++
						other.Focus()
					}
				})
				area := ui.Column(c.Key("area")).Absolute().Left(100).Top(150).Size(250, 120).Focusable().Label("Area")
				if ContextMenu(c, area, menuTestItems(&chosen)) {
					closedFrames++
					other.Focus()
				}
			}, 400, 300)
			test.run(t, tt)
			if closedFrames != 1 || !tt.Focused("Other") {
				t.Fatalf("closed reported in %d frames, other focused=%v", closedFrames, tt.Focused("Other"))
			}
		})
	}
}
