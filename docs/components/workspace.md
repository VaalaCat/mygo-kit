---
title: Workspace
description: 整个窗口的工作台骨架：顶栏、侧栏、中间区、右 Dock、下 Dock 和状态栏，窄窗口自动改成抽屉。
---

# Workspace

一个窗口的完整骨架：顶栏、左侧栏、中间区、右 Dock、下 Dock 和状态栏，分隔线可以拖动。窗口宽度小于 900 时，侧栏和右 Dock 自动变成抽屉。

<Preview page="Workspace" :height="520" :width="1180">

<<< @/../cmd/gallery/workspace.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
// 由应用持有，并负责持久化
var layout = kit.WorkspaceState{
	SidebarOpen: true,
	Right:       kit.DockState{Open: true, ActivePanel: "outline"},
}

kit.Workspace(c, &layout, kit.WorkspaceOpt{
	TitleBar: func() { ui.Text(c, "~/code/project") },
	Sidebar:  func() { kit.Header(c, "会话", nil) },
	Center:   func() { /* 主内容 */ },
	RightPanels: []kit.Panel{
		{ID: "outline", Title: "大纲", Icon: kit.IconFiles, Content: func() { /* ... */ }},
	},
})
```

内容区域用函数传入，每帧调用。值为 `nil` 的区域不显示。`Workspace` 自己画了状态栏，里面有侧栏开关和每个面板的按钮。

## 布局状态

`WorkspaceState` 保存侧栏和两个 Dock 的布局，字段带 `json` 标签，应用可以直接序列化保存。尺寸为 0 时用默认值：侧栏 240、右 Dock 360、下 Dock 240。拖动分隔线后，新尺寸会写回 `state`。

窄窗口里，侧栏和右 Dock 的抽屉开关只存在于内存，不写入持久化数据，窗口变宽后恢复原来的布局。

用状态上的方法开关面板，不要直接改字段，这样窄窗口的抽屉也会处理对：

```go
layout.ToggleSidebar()
layout.TogglePanel(&layout.Bottom, "terminal")
layout.ShowPanel(&layout.Right, "info")
```

## 示例

### 面板的标题栏

`Panel.Header` 可以替换 Dock 标题栏里的标题文字，例如放一排 [Tab](/components/tabs)。标签从标题栏左缘开始，放不下会横向滚动。`Actions` 画在标题栏右端、关闭按钮之前。

```go
kit.Panel{
	ID:    "terminal",
	Title: "终端",
	Icon:  kit.IconTerminal,
	Header: func() {
		for _, s := range sessions { /* kit.Tab(...) */ }
	},
	Actions: func() { kit.IconButtonTip(c.Key("new"), kit.IconPlus, "新建") },
	Content: func() { /* kit.Terminal(...) */ },
}
```

### 窗口按钮的位置

用 `TitleBar` 区域时，窗口按钮的避让由 `Workspace` 处理。创建窗口的细节见 [TitleBar](/components/title-bar)。

## API

### Workspace

`func Workspace(c *ui.Context, state *WorkspaceState, opt WorkspaceOpt) ui.Element`

### WorkspaceOpt

- **TitleBar** `func()`：顶栏左侧内容。
- **Sidebar** `func()`：侧栏内容。
- **Center** `func()`：中间区内容。
- **RightPanels** `[]Panel`：右 Dock 的面板，同一时间只显示一个。
- **BottomPanels** `[]Panel`：下 Dock 的面板。
- **StatusLeft** `func()`：状态栏左侧额外项，排在侧栏开关之后。
- **StatusRight** `func()`：状态栏右侧额外项，排在面板按钮之前。

### Panel

- **ID** `string`：面板标识，对应 `DockState.ActivePanel`。
- **Title** `string`：标题，也是状态栏按钮的提示。
- **Icon** `*ui.SVG`：状态栏按钮的图标。
- **Content** `func()`：面板内容。
- **Header** `func()`：替换 Dock 标题栏里的标题文字。为空时显示 `Title`。
- **Actions** `func()`：标题栏右端、关闭按钮之前的按钮。
- **CloseIcon** `*ui.SVG`：覆盖关闭按钮的图标。为空时跟随 Dock 所在边：下 Dock 用 `IconPanelBottomClose`，右 Dock 用 `IconPanelRightClose`。

### WorkspaceState

- **SidebarOpen** `bool`：侧栏是否打开（宽窗口）。
- **SidebarWidth** `float32`：侧栏宽度。
- **Right**、**Bottom** `DockState`：右 Dock 和下 Dock。

### DockState

- **Open** `bool`：Dock 是否打开。
- **ActivePanel** `string`：当前面板的 `ID`。无效时用第一个面板。
- **Size** `float32`：右 Dock 是宽度，下 Dock 是高度。

### WorkspaceState 的方法

- **IsNarrow() bool**：最近一帧窗口是否处于窄窗口（宽度小于 900）？
- **IsPanelShown(dock *DockState, panelID string) bool**：指定面板此刻是否正在显示？宽窗口里看 Dock，窄窗口里看右侧抽屉。
- **ShowPanel(dock *DockState, panelID string)**：打开 Dock 并显示指定面板。窄窗口里的右 Dock 打开抽屉并收起侧栏抽屉，不改写 `Right.Open`。
- **TogglePanel(dock *DockState, panelID string)**：显示指定面板；它已经在显示时，关闭 Dock（窄窗口里关闭抽屉）。
- **ToggleDock(dock *DockState)**：切换 Dock。
- **ShowSidebar()**：显示侧栏。窄窗口里打开侧栏抽屉并收起右 Dock 抽屉。
- **ToggleSidebar()**：切换侧栏。
