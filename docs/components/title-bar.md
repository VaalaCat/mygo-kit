---
title: TitleBar
description: 窗口顶栏，可以拖动窗口，并自动避开系统窗口按钮。
---

# TitleBar

窗口顶栏。整条可以拖动窗口，左右边距会按系统窗口按钮（macOS 红绿灯、Windows 按钮）占用的位置留出。

<Preview page="TitleBar" :height="260">

<<< @/../cmd/gallery/titlebar.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.TitleBar(c, func() {
	kit.IconButton(c.Key("folder"), kit.IconFolder, "工作区")
	ui.Text(c, "~/code/project")
})
```

要让系统标题栏让位、由 `TitleBar` 接管，创建窗口时设置 `mygo.WindowOptions{TitleBarStyle: mygo.TitleBarHidden}`。这时窗口按钮浮在内容上方，`TitleBar` 按 `c.TitleBar()` 给出的占位尺寸留出左右边距。窗口保持系统标题栏时，`c.TitleBar()` 为 0，`TitleBar` 按默认的 12 留边。

## 和别的组件配合

`TitleBar` 把自己的高度记在窗口上。同一个窗口里的 [Sheet](/components/sheet) 和 [ToastHost](/components/toast) 会据此从顶栏下缘开始，不盖住顶栏和系统按钮。高度按本帧或上一帧记录，所以调用先后只影响第一帧。[Workspace](/components/workspace) 已经替你做好了。

## API

### TitleBar

`func TitleBar(c *ui.Context, content func()) ui.Element`

高度取系统标题栏高度和 36 的较大值。背景是 `Colors.Background`，底部一条分隔线，子元素间距 8。

- **content** `func()`：绘制顶栏内容的函数。
