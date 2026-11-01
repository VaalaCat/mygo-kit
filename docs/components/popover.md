---
title: Popover
description: A floating panel beside an anchor element that reuses MyGo's positioning, dismissal, and focus handling.
---

# Popover

A small floating panel beside an element. MyGo handles positioning, dismissal on outside clicks or Esc, and focus. The component adds styling: a corner radius of 8, a floating-panel background, and a thin border.

There is no separate Popover page in the gallery. Its example appears with Tooltip at the bottom of the Dialog page.

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // Held by the application

anchor := kit.Button(c.Key("popover"), "More information")
if anchor.Clicked() {
	open = !open
}
kit.Popover(c, anchor, &open, func() {
	ui.Text(c, "Content inside the popover")
})
```

Call both the anchor and the popover every frame. You decide how clicking the anchor toggles it; inverting the value as above is enough. Clicking outside the popover or pressing Esc makes MyGo set `*open` to false.

## API

### Popover

`func Popover(c *ui.Context, anchor ui.Element, open *bool, content func()) ui.Element`

The popover appears below the anchor with a gap of 4, padding of 8, and child spacing of 8.

- **anchor** `ui.Element`: The element the popover is attached to.
- **open** `*bool`: Whether it is open, held by the caller.
- **content** `func()`: Draws the popover's content.

Returns the popover element. Chain `.Dismissed()` to check whether it was dismissed in the current frame.
