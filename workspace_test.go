package kit

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// workspaceFixture 返回带侧栏、两个右面板和一个下面板的工作台；state 由测试持有。
func workspaceFixture(state *WorkspaceState) func(*ui.Context) {
	return func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		Workspace(c, state, WorkspaceOpt{
			TitleBar: func() { ui.Text(c, "Title") },
			Sidebar:  func() { ui.Text(c, "Sidebar body") },
			Center:   func() { ui.Text(c, "Center body") },
			RightPanels: []Panel{
				{ID: "outline", Title: "Outline", Icon: IconFiles, Content: func() { ui.Text(c, "Outline body") }},
				{ID: "info", Title: "Info", Icon: IconInfo, Content: func() { ui.Text(c, "Info body") }},
			},
			BottomPanels: []Panel{{ID: "terminal", Title: "Terminal", Icon: IconTerminal, Content: func() { ui.Text(c, "Terminal body") }}},
		})
	}
}

// TestWorkspace 验证状态栏按钮开关 Dock、切换面板，以及窄窗口把侧栏改为抽屉。
func TestWorkspace(t *testing.T) {
	// 1. 每个用例使用独立状态和离屏窗口。
	tests := []struct {
		name  string
		width int
		state WorkspaceState
		run   func(*testing.T, *ui.Tester, *WorkspaceState)
	}{
		{"status-button-toggles-right-dock", 1200, WorkspaceState{}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			if tt.HasText("Outline body") {
				t.Fatal("dock open by default")
			}
			// Dock 标题行与状态栏按钮同名，先记下关闭状态时按钮的位置。
			outline, _ := tt.Find("Outline")
			info, _ := tt.Find("Info")
			click := func(r ui.Rect) { tt.ClickAt(r.X+r.W/2, r.Y+r.H/2) }
			click(outline)
			if !s.Right.Open || s.Right.ActivePanel != "outline" || !tt.HasText("Outline body") {
				t.Fatalf("not opened: %+v %v", s.Right, tt.Texts())
			}
			click(info)
			if !s.Right.Open || s.Right.ActivePanel != "info" || tt.HasText("Outline body") || !tt.HasText("Info body") {
				t.Fatalf("not switched: %+v", s.Right)
			}
			click(info)
			if s.Right.Open || tt.HasText("Info body") {
				t.Fatalf("not closed: %+v", s.Right)
			}
		}},
		{"status-button-toggles-bottom-dock", 1200, WorkspaceState{}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			tt.Click("Terminal")
			if !s.Bottom.Open || s.Bottom.ActivePanel != "terminal" || !tt.HasText("Terminal body") || s.Right.Open {
				t.Fatalf("bottom=%+v right=%+v", s.Bottom, s.Right)
			}
			tt.Click("Close")
			if s.Bottom.Open || tt.HasText("Terminal body") {
				t.Fatalf("close button did not close: %+v", s.Bottom)
			}
		}},
		{"default-sizes-and-resize-write-back", 1200, WorkspaceState{SidebarOpen: true, Right: DockState{Open: true, ActivePanel: "outline"}, Bottom: DockState{Open: true, ActivePanel: "terminal"}}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			if s.SidebarWidth != 240 || s.Right.Size != 360 || s.Bottom.Size != 240 {
				t.Fatalf("defaults: %+v", s)
			}
			outline, _ := tt.Find("Outline body")
			if outline.X < 1200-362 || outline.X > 1200-340 {
				t.Fatalf("right dock is not 360 wide: x=%v", outline.X)
			}
			terminal, _ := tt.Find("Terminal body")
			if terminal.Y < 760-240-12 || terminal.Y > 760-240+40 {
				t.Fatalf("bottom dock is not 240 high: y=%v", terminal.Y)
			}
			// 拖动右 Dock 分隔线后，宽度写回状态。
			tt.Press(1200-360, 300)
			tt.Move(1200-300, 300)
			tt.Release(1200-300, 300)
			if s.Right.Size > 320 || s.Right.Size < 280 {
				t.Fatalf("right size not written back: %v", s.Right.Size)
			}
		}},
		{"narrow-sidebar-renders-as-sheet", 700, WorkspaceState{}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			s.SidebarOpen = false
			tt.Frame()
			if tt.HasText("Sidebar body") {
				t.Fatal("sidebar visible while closed")
			}
			tt.Click("Toggle sidebar")
			if s.SidebarOpen || !tt.HasText("Sidebar body") {
				t.Fatal("sidebar sheet missing or written back to state")
			}
			body, _ := tt.Find("Sidebar body")
			center, _ := tt.Find("Center body")
			if body.X > 240 || center.X > 100 {
				// 抽屉覆盖在中间区上，不挤压中间区。
				t.Fatalf("sidebar=%v center=%v", body, center)
			}
			tt.Key(0, ui.KeyEscape)
			if tt.HasText("Sidebar body") {
				t.Fatal("Escape did not close sheet")
			}
		}},
		{"narrow-drawer-starts-below-title-bar", 700, WorkspaceState{}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			// 抽屉和遮罩从顶栏下缘开始：顶栏文字始终可见，点遮罩上方的顶栏区域不会关闭抽屉。
			if !s.IsNarrow() {
				t.Fatal("IsNarrow is false at 700 wide")
			}
			s.ShowSidebar()
			tt.Frame()
			title, _ := tt.Find("Title")
			body, _ := tt.Find("Sidebar body")
			if body.Y < OneLight().Metrics.HeaderHeight || title.Y+title.H > OneLight().Metrics.HeaderHeight {
				t.Fatalf("title=%v body=%v", title, body)
			}
			tt.ClickAt(600, 10)
			if !tt.HasText("Sidebar body") {
				t.Fatal("click on the title bar closed the drawer")
			}
			tt.ClickAt(600, 300)
			if tt.HasText("Sidebar body") {
				t.Fatal("click on the scrim did not close the drawer")
			}
		}},
		{"wide-window-is-not-narrow", 1200, WorkspaceState{}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			if s.IsNarrow() {
				t.Fatal("IsNarrow is true at 1200 wide")
			}
		}},
		{"narrow-right-dock-renders-as-sheet", 700, WorkspaceState{Right: DockState{Open: true, ActivePanel: "info"}}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			// 保存的布局里右 Dock 打开，窄窗口启动时不自动弹出抽屉。
			if tt.HasText("Info body") {
				t.Fatal("drawer opened on startup")
			}
			tt.Click("Info")
			body, ok := tt.Find("Info body")
			if !ok || body.X < 700-360 {
				t.Fatalf("right sheet missing or misplaced: %v %v", body, ok)
			}
			tt.Key(0, ui.KeyEscape)
			if !s.Right.Open || tt.HasText("Info body") {
				t.Fatal("Escape did not close right sheet or changed saved layout")
			}
		}},
		{"narrow-methods-drive-drawers", 700, WorkspaceState{}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			// 命令入口：窄窗口里 ShowPanel 打开右抽屉但不改写 Right.Open。
			s.ShowPanel(&s.Right, "info")
			tt.Frame()
			if !s.IsPanelShown(&s.Right, "info") || s.IsPanelShown(&s.Right, "outline") || s.Right.Open || !tt.HasText("Info body") {
				t.Fatalf("ShowPanel: right=%+v shown=%v", s.Right, s.IsPanelShown(&s.Right, "info"))
			}
			s.ToggleSidebar()
			tt.Frame()
			if s.IsPanelShown(&s.Right, "info") || tt.HasText("Info body") || !tt.HasText("Sidebar body") {
				t.Fatal("ToggleSidebar did not swap drawers")
			}
			s.ToggleDock(&s.Right)
			tt.Frame()
			if !s.IsPanelShown(&s.Right, "info") || tt.HasText("Sidebar body") || s.SidebarOpen {
				t.Fatal("ToggleDock did not open right drawer")
			}
			s.TogglePanel(&s.Right, "info")
			tt.Frame()
			if s.IsPanelShown(&s.Right, "info") || tt.HasText("Info body") {
				t.Fatal("TogglePanel did not close drawer")
			}
			// 底部 Dock 在窄窗口里仍然按 Dock 显示。
			s.ShowPanel(&s.Bottom, "terminal")
			tt.Frame()
			if !s.Bottom.Open || !s.IsPanelShown(&s.Bottom, "terminal") || !tt.HasText("Terminal body") {
				t.Fatalf("bottom: %+v", s.Bottom)
			}
		}},
		{"wide-methods-drive-docks", 1200, WorkspaceState{}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			s.ShowPanel(&s.Right, "info")
			s.ToggleSidebar()
			tt.Frame()
			if !s.Right.Open || !s.SidebarOpen || !s.IsPanelShown(&s.Right, "info") || !tt.HasText("Info body") {
				t.Fatalf("state=%+v", s)
			}
			s.ToggleDock(&s.Right)
			tt.Frame()
			if s.Right.Open || s.IsPanelShown(&s.Right, "info") {
				t.Fatal("ToggleDock did not close dock")
			}
		}},
		{"wide-sidebar-stays-in-layout", 1200, WorkspaceState{SidebarOpen: true}, func(t *testing.T, tt *ui.Tester, s *WorkspaceState) {
			body, _ := tt.Find("Sidebar body")
			center, _ := tt.Find("Center body")
			if center.X <= body.X+100 {
				t.Fatalf("sidebar=%v center=%v", body, center)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.state
			tt := ui.NewTester(workspaceFixture(&state), test.width, 760)
			test.run(t, tt, &state)
		})
	}
}

// TestWorkspaceTabStrip 验证 Dock 标题栏的标签条在标签放不下时横向滚动：滚轮移动标签条，新选中的标签被滚进可见区。
func TestWorkspaceTabStrip(t *testing.T) {
	// 1. 窄窗口里放二十个标签，选中项和下 Dock 状态由测试持有。
	const width = 600
	tests := []struct {
		name string
		run  func(t *testing.T, tt *ui.Tester, selected *int)
	}{
		{"wheel-moves-last-tab-into-view", func(t *testing.T, tt *ui.Tester, selected *int) {
			if last, _ := tt.Find("Tab 19"); last.X < width {
				t.Fatalf("last tab already visible: %v", last)
			}
			for range 80 {
				tt.Wheel(100, 510, 0, 120)
				tt.Advance(100 * time.Millisecond)
			}
			if last, _ := tt.Find("Tab 19"); last.X >= width {
				t.Fatalf("wheel did not scroll the strip: %v", last)
			}
			if first, _ := tt.Find("Tab 0"); first.X+first.W > 0 {
				t.Fatalf("first tab did not scroll away: %v", first)
			}
			// 反向滚动回到起点。
			for range 80 {
				tt.Wheel(100, 510, 0, -120)
				tt.Advance(100 * time.Millisecond)
			}
			if first, _ := tt.Find("Tab 0"); first.X < 0 {
				t.Fatalf("strip did not scroll back: %v", first)
			}
		}},
		{"selecting-last-tab-reveals-it", func(t *testing.T, tt *ui.Tester, selected *int) {
			*selected = 19
			tt.Frame()
			tt.Frame()
			if last, _ := tt.Find("Tab 19"); last.X >= width || last.X < 0 {
				t.Fatalf("selected tab not revealed: %v", last)
			}
			*selected = 0
			tt.Frame()
			tt.Frame()
			if first, _ := tt.Find("Tab 0"); first.X < 0 {
				t.Fatalf("first tab not revealed again: %v", first)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selected := 0
			state := WorkspaceState{Bottom: DockState{Open: true, ActivePanel: "terminal"}}
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				Workspace(c, &state, WorkspaceOpt{
					Center: func() { ui.Text(c, "Center body") },
					BottomPanels: []Panel{{ID: "terminal", Title: "Terminal", Icon: IconTerminal, Content: func() { ui.Text(c, "Terminal body") }, Header: func() {
						for i := range 20 {
							Tab(c.Key(fmt.Sprint("tab-", i)), TabOpt{Title: fmt.Sprint("Tab ", i), Selected: i == selected, CloseLabel: "Close tab"})
						}
					}}},
				})
			}, width, 760)
			test.run(t, tt, &selected)
		})
	}
}

// TestWorkspaceDockCloseIcon 验证 Dock 关闭按钮的图标跟随所在边：下 Dock 用向下收起、右 Dock 用向右收起，Panel.CloseIcon 可覆盖。
func TestWorkspaceDockCloseIcon(t *testing.T) {
	// 1. 每项只打开一个 Dock，把“关闭”按钮的像素与三种候选图标的独立按钮比对，最接近的就是实际图标。
	candidates := map[string]*ui.SVG{"bottom-close": IconPanelBottomClose, "right-close": IconPanelRightClose, "x": IconX}
	tests := []struct {
		name  string
		state WorkspaceState
		close *ui.SVG
		want  string
	}{
		{"bottom-dock-collapses-down", WorkspaceState{Bottom: DockState{Open: true, ActivePanel: "terminal"}}, nil, "bottom-close"},
		{"right-dock-collapses-right", WorkspaceState{Right: DockState{Open: true, ActivePanel: "outline"}}, nil, "right-close"},
		{"panel-override-wins", WorkspaceState{Right: DockState{Open: true, ActivePanel: "outline"}}, IconX, "x"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.state
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				Workspace(c, &state, WorkspaceOpt{
					Center:       func() {},
					RightPanels:  []Panel{{ID: "outline", Title: "Outline", Icon: IconFiles, CloseIcon: test.close, Content: func() {}}},
					BottomPanels: []Panel{{ID: "terminal", Title: "Terminal", Icon: IconTerminal, CloseIcon: test.close, Content: func() {}}},
				})
			}, 1200, 760)
			tt.Frame()
			button, ok := tt.Find("Close")
			if !ok {
				t.Fatalf("no close button: %v", tt.Texts())
			}
			got := cropImage(tt.Image(), button)
			best, bestDiff := "", -1
			for name, icon := range candidates {
				ref := ui.NewTester(func(c *ui.Context) {
					ApplyTheme(c, OneLight())
					// 参照按钮放在与实际按钮相同的小数偏移处，抗锯齿结果才可逐像素比较。
					ui.Column(c).Fill().Padding(button.Y-float32(int(button.Y)), 0, 0, button.X-float32(int(button.X))).Background(OneLight().Colors.Panel).Children(func() { IconButton(c, icon, "Close") })
				}, 200, 100)
				ref.Frame()
				refButton, _ := ref.Find("Close")
				if diff := imageDiff(got, cropImage(ref.Image(), refButton)); bestDiff < 0 || diff < bestDiff {
					best, bestDiff = name, diff
				}
			}
			if best != test.want {
				t.Fatalf("close icon = %s want %s", best, test.want)
			}
		})
	}
}

// cropImage 返回 img 中 r 范围的像素，坐标取整。
func cropImage(img *image.RGBA, r ui.Rect) *image.RGBA {
	return img.SubImage(image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))).(*image.RGBA)
}

// imageDiff 返回两张同尺寸图片逐像素通道差的总和。
func imageDiff(a, b *image.RGBA) int {
	total := 0
	for y := range a.Bounds().Dy() {
		for x := range a.Bounds().Dx() {
			pa := a.RGBAAt(a.Bounds().Min.X+x, a.Bounds().Min.Y+y)
			pb := b.RGBAAt(b.Bounds().Min.X+x, b.Bounds().Min.Y+y)
			total += abs(int(pa.R)-int(pb.R)) + abs(int(pa.G)-int(pb.G)) + abs(int(pa.B)-int(pb.B))
		}
	}
	return total
}
