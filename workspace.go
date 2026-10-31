package kit

import "github.com/egoist/mygo/ui"

// Panel 是 Dock 里的一个面板，应用每帧传入；同一 Dock 同时只显示一个。
type Panel struct {
	ID      string
	Title   string
	Icon    *ui.SVG
	Content func()
	// Header 替换 Dock 标题栏里的标题文字，例如放一排 Tab；Tab 从标题栏左缘开始、与标题栏同高。为空时显示 Title。
	Header func()
	// Actions 是标题栏右端、关闭按钮之前的按钮，例如新建。
	Actions func()
	// CloseIcon 覆盖关闭按钮图标；为空时跟随 Dock 所在边（下 Dock 用 IconPanelBottomClose，右 Dock 用 IconPanelRightClose）。
	CloseIcon *ui.SVG
}

// DockState 保存一个 Dock 是否打开、当前面板和尺寸。
type DockState struct {
	Open        bool    `json:"open"`
	ActivePanel string  `json:"active_panel"`
	Size        float32 `json:"size"`
}

// WorkspaceState 保存侧栏和两个 Dock 的布局；应用持有并负责持久化。
// 窄窗口里侧栏和右 Dock 的抽屉开关只存在于内存，不写入持久化数据。
type WorkspaceState struct {
	SidebarOpen  bool      `json:"sidebar_open"`
	SidebarWidth float32   `json:"sidebar_width"`
	Right        DockState `json:"right"`
	Bottom       DockState `json:"bottom"`

	sidebarDrawerOpen bool // 窄窗口里侧栏抽屉是否打开
	rightDrawerOpen   bool // 窄窗口里右 Dock 抽屉是否打开
	narrow            bool // 最近一帧窗口是否处于窄窗口
}

// WorkspaceOpt 指定各区域内容；nil 的区域不显示。
type WorkspaceOpt struct {
	TitleBar     func() // 顶栏左侧内容（窗口按钮避让由 Workspace 处理）
	Sidebar      func()
	Center       func()
	RightPanels  []Panel
	BottomPanels []Panel
	StatusLeft   func() // 状态栏左侧额外项，排在侧栏开关之后
	StatusRight  func() // 状态栏右侧额外项，排在面板按钮之前
}

// 工作台默认尺寸与窄窗口阈值；尺寸为零时采用默认值。
const (
	workspaceDefaultSidebarWidth = 240
	workspaceDefaultRightSize    = 360
	workspaceDefaultBottomSize   = 240
	workspaceNarrowWidth         = 900
)

// Workspace 绘制顶栏、侧栏、中间区、右 Dock、下 Dock 和状态栏；拖动分隔线后把新尺寸写回 state。
func Workspace(c *ui.Context, state *WorkspaceState, opt WorkspaceOpt) ui.Element {
	// 1. 尺寸为零时采用默认值；窄窗口里侧栏和右 Dock 改用抽屉。
	t := GetTheme(c)
	if state.SidebarWidth == 0 {
		state.SidebarWidth = workspaceDefaultSidebarWidth
	}
	if state.Right.Size == 0 {
		state.Right.Size = workspaceDefaultRightSize
	}
	if state.Bottom.Size == 0 {
		state.Bottom.Size = workspaceDefaultBottomSize
	}
	width, height := c.Size()
	narrow := width < workspaceNarrowWidth
	state.narrow = narrow
	if !narrow {
		state.sidebarDrawerOpen, state.rightDrawerOpen = false, false
	}
	rightPanel := getActivePanel(&state.Right, opt.RightPanels)
	bottomPanel := getActivePanel(&state.Bottom, opt.BottomPanels)
	// 2. 分隔线占 1 像素；用窗口尺寸换算各 Split 第一块的尺寸。
	bodyWidth := width
	showSidebar := opt.Sidebar != nil && state.SidebarOpen && !narrow
	if showSidebar {
		bodyWidth -= state.SidebarWidth + 1
	}
	bodyHeight := height - t.Metrics.StatusHeight
	if opt.TitleBar != nil {
		bodyHeight -= max(c.TitleBar().Height, t.Metrics.HeaderHeight)
	}
	root := ui.Column(c).Fill().Gap(0).Background(t.Colors.Editor)
	root.Children(func() {
		if opt.TitleBar != nil {
			TitleBar(c, opt.TitleBar)
		}
		// 3. 主体左侧是侧栏，右侧自上而下是「中间区 + 右 Dock」和下 Dock。
		main := func() {
			workspaceMain(c, state, opt.Center, rightPanel, !narrow, bodyWidth, bodyHeight, bottomPanel)
		}
		ui.Column(c).Grow(1).Shrink(1).Gap(0).Children(func() {
			if !showSidebar {
				main()
				return
			}
			Split(c.Key("workspace-sidebar-split"), &state.SidebarWidth, func() {
				ui.Column(c).Fill().Background(t.Colors.Panel).Children(opt.Sidebar)
			}, main).Grow(1)
		})
		workspaceStatusBar(c, state, opt)
	})
	// 4. 窄窗口里侧栏和右 Dock 是临时抽屉，开关不写回 state，变宽后恢复原布局；抽屉从顶栏下缘开始，没有顶栏时铺满窗口。
	if narrow {
		drawerTop := float32(0)
		if opt.TitleBar != nil {
			drawerTop = max(c.TitleBar().Height, t.Metrics.HeaderHeight)
		}
		if opt.Sidebar != nil {
			sheetBelow(c.Key("workspace-sidebar-sheet"), &state.sidebarDrawerOpen, SheetLeft, state.SidebarWidth, drawerTop, func() {
				ui.Column(c).Fill().Children(opt.Sidebar)
			})
		}
		if panel := getDrawerPanel(state.Right.ActivePanel, opt.RightPanels); panel != nil {
			sheetBelow(c.Key("workspace-right-sheet"), &state.rightDrawerOpen, SheetRight, state.Right.Size, drawerTop, func() { dockView(c, &state.rightDrawerOpen, *panel, IconPanelRightClose) })
		}
	}
	return root
}

// IsNarrow 最近一帧窗口是否处于窄窗口（宽度小于 900，侧栏和右 Dock 改用抽屉）？
func (s *WorkspaceState) IsNarrow() bool { return s.narrow }

// ToggleSidebar 切换侧栏：宽窗口改写 SidebarOpen，窄窗口切换侧栏抽屉并收起右 Dock 抽屉。
func (s *WorkspaceState) ToggleSidebar() {
	if !s.narrow {
		s.SidebarOpen = !s.SidebarOpen
		return
	}
	if s.sidebarDrawerOpen {
		s.sidebarDrawerOpen = false
		return
	}
	s.ShowSidebar()
}

// ShowSidebar 显示侧栏：宽窗口里置 SidebarOpen，窄窗口里打开侧栏抽屉并收起右 Dock 抽屉。
func (s *WorkspaceState) ShowSidebar() {
	if !s.narrow {
		s.SidebarOpen = true
		return
	}
	s.sidebarDrawerOpen, s.rightDrawerOpen = true, false
}

// ToggleDock 切换 Dock：窄窗口里的右 Dock 切换抽屉并收起侧栏抽屉，其余情况改写 Dock.Open。
func (s *WorkspaceState) ToggleDock(dock *DockState) {
	if !s.isDrawer(dock) {
		dock.Open = !dock.Open
		return
	}
	s.rightDrawerOpen = !s.rightDrawerOpen
	if s.rightDrawerOpen {
		s.sidebarDrawerOpen = false
	}
}

// ShowPanel 打开 Dock 并显示指定面板；窄窗口里的右 Dock 打开抽屉并收起侧栏抽屉，不改写 Right.Open。
func (s *WorkspaceState) ShowPanel(dock *DockState, panelID string) {
	if !s.isDrawer(dock) {
		dock.Open, dock.ActivePanel = true, panelID
		return
	}
	dock.ActivePanel = panelID
	s.rightDrawerOpen, s.sidebarDrawerOpen = true, false
}

// TogglePanel 显示指定面板；它已经是显示中的面板时关闭 Dock（窄窗口里关闭抽屉）。
func (s *WorkspaceState) TogglePanel(dock *DockState, panelID string) {
	if s.IsPanelShown(dock, panelID) {
		if s.isDrawer(dock) {
			s.rightDrawerOpen = false
			return
		}
		dock.Open = false
		return
	}
	s.ShowPanel(dock, panelID)
}

// IsPanelShown 指定面板此刻是否正在显示（宽窗口里作为 Dock，窄窗口里作为右侧抽屉）？
func (s *WorkspaceState) IsPanelShown(dock *DockState, panelID string) bool {
	if s.isDrawer(dock) {
		return s.rightDrawerOpen && dock.ActivePanel == panelID
	}
	return dock.Open && dock.ActivePanel == panelID
}

// isDrawer 此刻这个 Dock 是否以抽屉形式显示（窄窗口里的右 Dock）？
func (s *WorkspaceState) isDrawer(dock *DockState) bool {
	return s.narrow && dock == &s.Right
}

// workspaceMain 绘制中间区、右 Dock 和下 Dock，并把拖动结果写回对应尺寸。
func workspaceMain(c *ui.Context, state *WorkspaceState, center func(), rightPanel *Panel, rightAsSplit bool, bodyWidth, bodyHeight float32, bottomPanel *Panel) {
	// 1. 上半部分：中间区，宽窗口里右侧并排右 Dock。
	top := func() {
		centerView := func() { ui.Column(c).Fill().Background(GetTheme(c).Colors.Editor).Children(center) }
		if rightPanel == nil || !rightAsSplit {
			centerView()
			return
		}
		first := bodyWidth - 1 - state.Right.Size
		split := Split(c.Key("workspace-right-split"), &first, centerView, func() { dockView(c, &state.Right.Open, *rightPanel, IconPanelRightClose) }).Grow(1)
		if split.Changed() {
			state.Right.Size = bodyWidth - 1 - first
		}
	}
	if bottomPanel == nil {
		top()
		return
	}
	// 2. 下 Dock 横跨中间区和右 Dock，用竖向分隔线调整高度。
	first := bodyHeight - 1 - state.Bottom.Size
	split := ui.SplitVertical(c.Key("workspace-bottom-split"), &first, top, func() { dockView(c, &state.Bottom.Open, *bottomPanel, IconPanelBottomClose) }).Fill().Gap(0).Radius(0).Background(GetTheme(c).Colors.Editor)
	if split.Changed() {
		state.Bottom.Size = bodyHeight - 1 - first
	}
}

// dockView 绘制 Dock 顶部的面板标题栏和面板内容；关闭按钮图标取 panel.CloseIcon，为空时用 dockCloseIcon（跟随 Dock 所在边）。
func dockView(c *ui.Context, open *bool, panel Panel, dockCloseIcon *ui.SVG) {
	// 1. 标题栏与侧栏、中间区的标题栏同高同字号，底线连成一条；面板自带标题栏内容时替换标题文字。
	t := GetTheme(c)
	ui.Column(c).Fill().Gap(0).Background(t.Colors.Panel).Children(func() {
		closeIcon := dockCloseIcon
		if panel.CloseIcon != nil {
			closeIcon = panel.CloseIcon
		}
		closeButton := func() {
			if IconButton(c.Key("dock-close-"+panel.ID), closeIcon, text(c, "close")).Clicked() {
				*open = false
			}
		}
		if panel.Header == nil {
			Header(c, panel.Title, func() {
				if panel.Actions != nil {
					panel.Actions()
				}
				closeButton()
			})
		} else {
			// 自定义标题栏不画整行底线：Tab 自己画底线，选中的 Tab 断开底线；其余空白和右端按钮区补上底线。
			ui.Row(c).Height(t.Metrics.HeaderHeight).Shrink(0).Gap(0).AlignItems(ui.Center).Background(t.Colors.Panel).Children(func() {
				dockTabStrip(c, panel.Header)
				ui.Row(c).Grow(1).Height(t.Metrics.HeaderHeight).Padding(0, t.Space(1)).Gap(t.Space(1)).Justify(ui.End).AlignItems(ui.Center).BorderWidth(0, 0, 1, 0).BorderColor(t.Colors.Border).Children(func() {
					if panel.Actions != nil {
						panel.Actions()
					}
					closeButton()
				})
			})
		}
		// 2. 内容占满剩余空间并裁切。
		ui.Column(c).Grow(1).Shrink(1).Gap(0).Clip().Children(panel.Content)
	})
}

// dockTabStripKey 标识标签条保存的滚动位置。
type dockTabStripKey struct{}

// dockTabStrip 绘制 Dock 标题栏里的标签条：标签放不下时横向滚动，竖直滚轮换成横向滚动；滚动条只在悬停时浮出。
func dockTabStrip(c *ui.Context, tabs func()) {
	// 1. 滚动位置由元素保存，标签条高度固定为标题栏高度。
	t := GetTheme(c)
	strip := ui.ScrollHorizontal(c).Height(t.Metrics.HeaderHeight).Shrink(1).MinWidth(0).Gap(0)
	scroll := ui.Local(strip, dockTabStripKey{}, func() ui.ScrollState { return ui.ScrollState{} })
	strip.TrackScroll(scroll)
	// 2. 鼠标滚轮只有竖直增量，换成横向移动；触控板的横向增量交给容器，到头后继续向外传递。
	strip.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll || ev.DY == 0 || ev.DX != 0 {
			return false
		}
		next := min(max(scroll.X+ev.DY, 0), scroll.MaxX)
		if next == scroll.X {
			return false
		}
		scroll.X = next
		return true
	})
	strip.Children(tabs)
}

// workspaceStatusBar 绘制底部状态栏：左侧侧栏开关，右侧右 Dock 与下 Dock 的面板按钮。
func workspaceStatusBar(c *ui.Context, state *WorkspaceState, opt WorkspaceOpt) {
	t := GetTheme(c)
	StatusBar(c, func() {
		// 1. 左侧：侧栏开关，其后是应用的额外项。
		if opt.Sidebar != nil {
			title := text(c, "toggleSidebar")
			open := state.SidebarOpen
			if state.narrow {
				open = state.sidebarDrawerOpen
			}
			button := StatusButton(c.Key("workspace-sidebar-toggle"), IconPanelLeft, title, open)
			Tooltip(c, button, title)
			if button.Clicked() {
				state.ToggleSidebar()
			}
		}
		if opt.StatusLeft != nil {
			opt.StatusLeft()
		}
		ui.Spacer(c)
		// 2. 右侧：应用的额外项、右 Dock 按钮、竖分隔、下 Dock 按钮。
		if opt.StatusRight != nil {
			opt.StatusRight()
		}
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			dockButtons(c, state, &state.Right, opt.RightPanels, "right")
			if len(opt.RightPanels) > 0 && len(opt.BottomPanels) > 0 {
				ui.Box(c).Width(1).Height(t.Space(4)).Margin(0, t.Space(1)).Shrink(0).Background(t.Colors.Border)
			}
			dockButtons(c, state, &state.Bottom, opt.BottomPanels, "bottom")
		})
	})
}

// dockButtons 为 Dock 的每个面板绘制状态栏按钮；正在显示的面板使用选中底色。
func dockButtons(c *ui.Context, state *WorkspaceState, dock *DockState, panels []Panel, name string) {
	for _, panel := range panels {
		button := StatusButton(c.Key("workspace-"+name+"-"+panel.ID), panel.Icon, panel.Title, state.IsPanelShown(dock, panel.ID))
		Tooltip(c, button, panel.Title)
		if button.Clicked() {
			state.TogglePanel(dock, panel.ID)
		}
	}
}

// getDrawerPanel 返回抽屉里显示的右 Dock 面板；当前面板 ID 无效时采用第一个，没有面板时返回 nil。
func getDrawerPanel(panelID string, panels []Panel) *Panel {
	for i := range panels {
		if panels[i].ID == panelID {
			return &panels[i]
		}
	}
	if len(panels) == 0 {
		return nil
	}
	return &panels[0]
}

// getActivePanel 返回 Dock 当前显示的面板；Dock 关闭或没有面板时返回 nil，当前面板 ID 无效时采用第一个面板。
func getActivePanel(dock *DockState, panels []Panel) *Panel {
	// 1. 只有打开且有面板的 Dock 才显示内容。
	if !dock.Open || len(panels) == 0 {
		return nil
	}
	for i := range panels {
		if panels[i].ID == dock.ActivePanel {
			return &panels[i]
		}
	}
	return &panels[0]
}
