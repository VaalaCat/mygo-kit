---
title: TitleBar
description: A draggable window title bar that automatically leaves room for system window buttons.
---

# TitleBar

A window title bar. The whole bar can be used to drag the window. Left and right margins account for the space occupied by system window buttons (macOS traffic lights and Windows buttons).

<Preview page="TitleBar" :height="260">

<<< @/../cmd/gallery/titlebar.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.TitleBar(c, func() {
	kit.IconButton(c.Key("folder"), kit.IconFolder, "Workspace")
	ui.Text(c, "~/code/project")
})
```

To replace the system title bar with `TitleBar`, set `mygo.WindowOptions{TitleBarStyle: mygo.TitleBarHidden}` when creating the window. The window buttons then float above the content, and `TitleBar` leaves left and right margins based on the reserved space reported by `c.TitleBar()`. When the window keeps the system title bar, `c.TitleBar()` is 0 and `TitleBar` uses the default margin of 12.

## Working with other components

`TitleBar` records its height on the window. [Sheet](/components/sheet) and [ToastHost](/components/toast) in the same window use it to start below the title bar, leaving it and the system buttons uncovered. The height comes from the current or previous frame, so call order only affects the first frame. [Workspace](/components/workspace) handles this for you.

## API

### TitleBar

`func TitleBar(c *ui.Context, content func()) ui.Element`

The height is the greater of the system title bar height and 36. The background is `Colors.Background`, with a bottom separator and a gap of 8 between children.

- **content** `func()`: A function that draws the title bar content.
