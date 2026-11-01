---
title: MorphText
description: Single-line status text where only changed graphemes briefly move upward and fade in.
---

# MorphText

Displays a line of changing status text, such as “Reading files” changing to “Running tests”. Only changed graphemes move slightly upward and fade in; matching graphemes in the prefix, suffix, and middle stay as they are.

<Preview page="MorphText" :height="200">

<<< @/../cmd/gallery/morph.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.MorphText(c, status) // status is the application's current status text
```

## Behavior

- Works with Unicode graphemes, so complete emoji and characters with combining accents are not split apart.
- Common prefixes and suffixes stay still. Appending characters does not animate, and unchanged prefixes and suffixes do not replay their animations.
- Compares only a limited number of graphemes in the changed middle section. If the change is too large, the new text appears without animation.
- Shows the new text directly when the system or theme requests reduced motion.
- Text stays on one line without wrapping. The accessible name for the entire element is the current text.

## API

### MorphText

`func MorphText(c *ui.Context, value string) ui.Element`

- **value** `string`: The current text to display, passed in each frame. The component compares it with the previous text to find changes. Animation duration comes from `Theme.Motion.Fade` (120ms by default).
