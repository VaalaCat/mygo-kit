---
title: Kbd
description: A small keyboard shortcut label that controls appearance without registering shortcuts.
---

# Kbd

A small label that displays a keyboard shortcut, such as `⌘K`. It handles appearance only; register the actual shortcut with MyGo yourself.

<Preview page="Kbd" :height="180">

<<< @/../cmd/gallery/kbd.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	ui.Text(c, "Open command palette")
	kit.Kbd(c, "⌘K")
})
```

[CommandPalette](/components/command-palette) also uses Kbd to display `Command.Keys` shortcuts. A [Menu](/components/menu) item's `MenuItem.Shortcut` is secondary-colored text on the right, without a border or background.

## API

### Kbd

`func Kbd(c *ui.Context, keys string) ui.Element`

A small label with a height of 20, a corner radius of 4, a background slightly darker than its panel, and a thin border. Text uses 12/16 sizing and the secondary color.

- **keys** `string`: The text to display, such as `"⌘N"` or `"Ctrl+N"`. The component does not parse it.
