---
title: Disclosure
description: A borderless, full-row collapsible header with a height transition when its content expands.
---

# Disclosure

A full-row clickable header with an arrow on the left. Click to expand or collapse the content below. Keyboard behavior and expansion semantics come from MyGo's collapsible component.

<Preview page="Disclosure" :height="220">

<<< @/../cmd/gallery/disclosure.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
open := true // The application owns the expanded state

kit.Disclosure(c, "More settings", &open, func() {
	ui.Text(c, "Content shown when expanded")
})
```

The content area is indented to align with the title text. Height transitions take 200ms to expand and 150ms to collapse. There is no transition when the system or theme requests reduced motion. Content is disabled as soon as collapsing begins and no longer responds to clicks.

## API

### Disclosure

`func Disclosure(c *ui.Context, label string, open *bool, content func()) ui.Element`

Returns the header row element for further chained settings, such as `.Disabled(true)`.

- **label** `string`: Header text, also used as the accessible name.
- **open** `*bool`: Whether the content is expanded, owned by the caller.
- **content** `func()`: Function that draws the content when expanded. Not called once fully collapsed.
