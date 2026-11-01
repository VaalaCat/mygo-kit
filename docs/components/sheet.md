---
title: Sheet
description: A modal drawer attached to the left or right edge of the window, dismissed by clicking the backdrop or pressing Esc.
---

# Sheet

A modal drawer that slides in from the left or right edge of the window, with a translucent backdrop. Clicking the backdrop or pressing Esc sets `*open` to false.

<Preview page="Sheet" :height="360">

<<< @/../cmd/gallery/sheet.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // Held by the application

if kit.Button(c, "Open details").Clicked() {
	open = true
}
kit.Sheet(c.Key("details"), &open, kit.SheetRight, 360, func() {
	kit.Header(c, "Details", nil)
	ui.Text(c, "Drawer content").Padding(12)
})
```

Call the drawer every frame. It draws nothing when `open` is false.

## Working with TitleBar

When the window uses a [TitleBar](/components/title-bar), the drawer and backdrop start below it, leaving the title bar and system window buttons uncovered. Without a TitleBar, they fill the entire window. The `TitleBar` height is recorded in the current or previous frame, so call order only affects the first frame.

## API

### Sheet

`func Sheet(c *ui.Context, open *bool, side SheetSide, width float32, content func()) ui.Element`

The drawer uses the panel background, with a separator only on the edge facing the content. Its maximum width is 90% of the window.

- **open** `*bool`: Whether the drawer is open, held by the caller.
- **side** `SheetSide`: The edge the drawer attaches to.
- **width** `float32`: The drawer width.
- **content** `func()`: A function that draws the drawer content.

### SheetSide

- **SheetLeft**: Attaches to the left edge of the window.
- **SheetRight**: Attaches to the right edge of the window.
