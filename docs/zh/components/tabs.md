---
title: Tabs
description: 连续页签条，以及 Dock 和终端里用的可关闭标签 Tab。
---

# Tabs

页签条。选中的页签用正文底色并断开底线，和下面的内容连成一片。

<Preview page="Tabs" :height="240">

<<< @/../cmd/gallery/tabs.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
tab := 0 // 选中页签的序号，由应用持有

kit.Tabs(c, &tab, []string{"常规", "编辑器", "高级"})
switch tab {
case 0:
	// 画“常规”的内容
}
```

`Tabs` 只画页签条，下面的内容由你按 `tab` 的值自己画。方向键切换页签由 MyGo 提供，焦点环只在键盘聚焦时出现，并且始终跟着选中的页签。

## 示例

### 可关闭的标签

`Tab` 画一个 Zed 风格的单个标签，常用于终端、编辑器这类可以多开的面板。它返回本帧是否点了标签、点了关闭按钮。选中、关闭后怎么处理都由调用方决定。

```go
for i, session := range sessions {
	clicked, closed := kit.Tab(c.Key(session.ID), kit.TabOpt{
		Icon:       kit.IconTerminal,
		Title:      session.Title,
		Selected:   i == current,
		CloseLabel: "关闭",
	})
	if clicked {
		current = i
	}
	if closed {
		// 移除 session
	}
}
```

放在 [Workspace](/zh/components/workspace) 的 `Panel.Header` 里时，标签和 Dock 标题栏同高，放不下会横向滚动。

## API

### Tabs

`func Tabs(c *ui.Context, selected *int, labels []string) ui.Element`

36 高的页签条，页签之间用竖线分隔，右侧空白补一条底线。

- **selected** `*int`：选中页签的序号，由调用方持有。
- **labels** `[]string`：每个页签的文字，同时是可访问名称。

### Tab

`func Tab(c *ui.Context, opt TabOpt) (clicked, closed bool)`

与标题栏同高（36）、无圆角的可关闭标签。选中时用正文底色加左右竖线，未选中时用面板底色加右竖线和底线。关闭按钮在悬停或选中时出现。新选中的标签会滚进所在的标签条。

### TabOpt

- **Icon** `*ui.SVG`：标题前的图标，14。可以为 `nil`。
- **Title** `string`：标题，最宽 192，超出用省略号截断。
- **Selected** `bool`：是否选中。
- **Faded** `bool`：标题用弱化色，例如进程已结束的终端。
- **CloseLabel** `string`：关闭按钮的可访问名称。
