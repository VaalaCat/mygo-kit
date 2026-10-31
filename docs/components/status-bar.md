---
title: StatusBar
description: 窗口底部的状态栏，以及状态栏里 20 高的图标按钮 StatusButton。
---

# StatusBar

窗口底部的状态栏，背景和 [TitleBar](/components/title-bar) 同色，顶部有一条分隔线。里面的文字是 12/16 的次要色。

<Preview page="StatusBar" :height="240">

<<< @/../cmd/gallery/statusbar.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.StatusBar(c, func() {
	button := kit.StatusButton(c.Key("sidebar"), kit.IconPanelLeft, "切换侧栏", sidebarOpen)
	kit.Tooltip(c, button, "切换侧栏")
	if button.Clicked() {
		sidebarOpen = !sidebarOpen
	}
	ui.Text(c, "已连接")
	ui.Spacer(c)
	kit.StatusButton(c.Key("terminal"), kit.IconPanelBottom, "终端", false)
})
```

`StatusBar` 是一行横排，内容靠左；要把一部分推到右边，中间放 `ui.Spacer(c)`。

用 [Workspace](/components/workspace) 时，状态栏已经带好了侧栏开关和面板按钮，只要通过 `StatusLeft`、`StatusRight` 加自己的项。

## API

### StatusBar

`func StatusBar(c *ui.Context, content func()) ui.Element`

28 高，内边距左右 4，子元素间距 4。

- **content** `func()`：绘制状态栏内容的函数。

### StatusButton

`func StatusButton(c *ui.Context, icon *ui.SVG, label string, active bool) ui.Element`

20 高的按钮。有图标时是 20 见方，图标 14；没有图标时显示 `label` 文字。

- **icon** `*ui.SVG`：图标。传 `nil` 则显示 `label` 文字。
- **label** `string`：可访问名称；没有图标时同时是显示的文字。
- **active** `bool`：为真时用选中底色，文字和图标变正文色。
