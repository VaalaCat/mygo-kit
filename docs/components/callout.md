---
title: Callout
description: A tinted notice with an icon, title, and description, plus optional actions on the right.
---

# Callout

An inline notice with a light background and thin border. Available in error, warning, and information variants.

<Preview page="Callout" :height="300">

<<< @/../cmd/gallery/callout.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Callout(c, kit.CalloutError, "Model request failed", "The service returned 429: too many requests. Try again in a few seconds.", func() {
	if kit.Button(c, "Retry").Clicked() {
		retry()
	}
})
```

## Examples

### Title only

Passing an empty string for `detail` shows only the title, with all content vertically centered.

```go
kit.Callout(c, kit.CalloutInfo, "Conversation forked", "", nil)
```

### With actions

With `actions` set to `nil`, no action area is reserved, as in the example above. When a function is supplied, it draws buttons on the right. You can add several:

```go
kit.Callout(c, kit.CalloutInfo, "Conversation forked", "", func() {
	kit.Button(c, "Open")
	kit.IconButton(c, kit.IconX, "Close")
})
```

## API

### Callout

`func Callout(c *ui.Context, kind CalloutKind, title, detail string, actions func()) ui.Element`

Corner radius 8, padding 12. The title uses the body text color; details use 12/16 text in the secondary color.

- **kind** `CalloutKind`: Variant, determining the icon and tone.
- **title** `string`: Title, also used as the accessible name.
- **detail** `string`: Details in the secondary color. May be empty.
- **actions** `func()`: Function that draws the actions on the right. May be `nil`.

### CalloutKind

- **CalloutError**: Error, using the danger color.
- **CalloutWarning**: Warning, using the warning color.
- **CalloutInfo**: Information, using the accent color.
