---
title: StatusBar
description: A status bar at the bottom of the window, with StatusButton icon buttons that are 20 high.
---

# StatusBar

A status bar at the bottom of the window. It shares the [TitleBar](/components/title-bar) background color and has a separator along the top. Text uses the secondary color at 12/16.

<Preview page="StatusBar" :height="240">

<<< @/../cmd/gallery/statusbar.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.StatusBar(c, func() {
	button := kit.StatusButton(c.Key("sidebar"), kit.IconPanelLeft, "Toggle sidebar", sidebarOpen)
	kit.Tooltip(c, button, "Toggle sidebar")
	if button.Clicked() {
		sidebarOpen = !sidebarOpen
	}
	ui.Text(c, "Connected")
	ui.Spacer(c)
	kit.StatusButton(c.Key("terminal"), kit.IconPanelBottom, "Terminal", false)
})
```

`StatusBar` is a horizontal row with content aligned to the left. Place `ui.Spacer(c)` between groups to push some content to the right.

When using [Workspace](/components/workspace), the status bar already includes the sidebar toggle and panel buttons. Add your own items through `StatusLeft` and `StatusRight`.

## API

### StatusBar

`func StatusBar(c *ui.Context, content func()) ui.Element`

A height of 28, horizontal padding of 4, and a gap of 4 between children.

- **content** `func()`: A function that draws the status bar content.

### StatusButton

`func StatusButton(c *ui.Context, icon *ui.SVG, label string, active bool) ui.Element`

A button with a height of 20. With an icon, it is a 20-by-20 square containing a 14-sized icon. Without an icon, it shows the `label` text.

- **icon** `*ui.SVG`: The icon. Pass `nil` to show the `label` text.
- **label** `string`: The accessible name. Also the display text when no icon is provided.
- **active** `bool`: When true, uses the selected background and the foreground color for text and icons.
