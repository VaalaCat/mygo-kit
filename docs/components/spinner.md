---
title: Spinner
description: A rotating loading icon that follows the text color.
---

# Spinner

A loading icon that makes one rotation per second. It stays still when the system or theme requests reduced motion.

<Preview page="Spinner" :height="180">

<<< @/../cmd/gallery/spinner.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	kit.Spinner(c, kit.IconSizeDefault)
	ui.Text(c, "Thinking…")
})
```

The default is the secondary color. Chain `.TextColor(...)` to change it. For icon sizes, use [`IconSizeSmall`, `IconSizeDefault`, `IconSizeLarge`](/theming#icon-sizes).

## API

### Spinner

`func Spinner(c *ui.Context, size float32) ui.Element`

- **size** `float32`: The icon's side length. Common values are 12, 14, and 16.
