---
title: Layout
description: 排版小件：标题栏 Header、分隔线 Separator、可拖宽的 Split、空状态 EmptyState。
---

# Layout

几个排版用的小件。画廊里没有单独的页面，它们出现在其他页面的外壳里，例如 [Workspace](/components/workspace) 用到了 `Header` 和 `Split`，画廊的 Workspace 页中间区是一个 `EmptyState`。

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

### Header

36 高的面板标题栏，左边是标题，右边是 `actions` 画的工具按钮。底部有一条分隔线，会和下面的内容共用。

```go
kit.Header(c, "会话", func() {
	kit.IconButtonTip(c.Key("new"), kit.IconPlus, "新建")
})
```

需要把标题换成别的内容，或者对齐到另一条左线，用 `HeaderWith`：

```go
kit.HeaderWith(c, "", kit.HeaderOpt{Inset: 16, Leading: func() {
	kit.Segmented(c, &view, []string{"会话", "归档"})
}}, nil)
```

### Separator

一条 1 像素的水平分隔线。

```go
kit.Separator(c)
```

### Split

两个可以拖宽的连续面板，中间是 MyGo 的分隔线。`size` 是第一块的宽度，拖动时直接写回。

```go
sidebarWidth := float32(240) // 由应用持有

kit.Split(c, &sidebarWidth, func() {
	// 左侧
}, func() {
	// 右侧
}).Grow(1)
```

### EmptyState

在区域中央显示一行提示，可以带一行次要说明，也可以带一个主按钮：

```go
kit.EmptyState(c, "还没有会话", "点击左上角的加号新建一个。")

kit.EmptyStateWith(c, "还没有会话", "", kit.EmptyStateOpt{
	ActionLabel: "新建会话",
	OnAction:    newConversation,
})
```

## API

### Header

`func Header(c *ui.Context, title string, actions func()) ui.Element`

36 高，标题距左缘 12，单行截断。

- **title** `string`：标题文字。
- **actions** `func()`：绘制右侧工具按钮的函数，可以为 `nil`。

### HeaderWith

`func HeaderWith(c *ui.Context, title string, opt HeaderOpt, actions func()) ui.Element`

### HeaderOpt

- **Inset** `float32`：标题文字距左缘的距离，用来对齐所在列的左线。侧栏和 Dock 用 12，对话区用 16。零值取 12。
- **Leading** `func()`：替换标题文字，例如放一组分段控件。内容从 `Inset` 处开始，标题栏的高度和底线不变。为空时显示 `title`。

### Separator

`func Separator(c *ui.Context) ui.Element`

1 高，颜色是 `Colors.BorderSubtle`，不会被压缩。

### Split

`func Split(c *ui.Context, size *float32, first, second func()) ui.Element`

两块都铺满，不加外框，圆角 0。

- **size** `*float32`：第一块的宽度，由调用方持有。拖动分隔线时被改写，用返回元素的 `.Changed()` 可以知道本帧是否被拖动。
- **first**、**second** `func()`：分别绘制两块的内容。

### EmptyState

`func EmptyState(c *ui.Context, title, detail string) ui.Element`

居中显示，占满父容器剩余空间。标题是次要色，说明是占位色。

- **title** `string`：标题。
- **detail** `string`：说明，可以为空。

### EmptyStateWith

`func EmptyStateWith(c *ui.Context, title, detail string, opt EmptyStateOpt) ui.Element`

### EmptyStateOpt

- **ActionLabel** `string`：非空时在文字下方居中显示一个主按钮。
- **OnAction** `func()`：点击主按钮时调用。
