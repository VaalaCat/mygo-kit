---
title: Tooltip
description: A short hint shown on hover or keyboard focus, with an independent timer for each anchor.
---

# Tooltip

Adds a hover hint to an element. It appears after the pointer rests on the anchor for 600ms, or immediately when the anchor receives keyboard focus (with a focus ring). Esc behavior comes from MyGo.

The gallery has no separate Tooltip page. An example appears at the bottom of the Dialog page, and you can also hover over the left button on the StatusBar page to see one.

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
button := kit.Button(c.Key("save"), "Save")
kit.Tooltip(c, button, "Save the current document")
if button.Clicked() {
	// Save
}
```

The hint text is also set as the anchor's accessible description.

## Tips

- For an icon button with a hint, use [IconButtonTip](/components/button) directly instead of combining `Tooltip` and `Label` yourself.
- Each anchor keeps its own timer. Moving the pointer from one button to another does not immediately show the second button's hint.
- Use a short, single-line hint. For actions or clickable content, use a [Popover](/components/popover).

## API

### Tooltip

`func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element`

The hint popup has a corner radius of 8, a popup background, a thin border, and 12/16 text.

- **anchor** `ui.Element`: The element to attach the hint to.
- **text** `string`: The hint text.

Returns the hint element. The anchor remains the element in your original variable; handle clicks and focus on it as usual.
