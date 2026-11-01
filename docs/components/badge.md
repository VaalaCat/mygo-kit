---
title: Badge
description: A compact status Badge and a small StatusDot for connection status or unread items.
---

# Badge

`Badge` is a small label with a semantic color. `StatusDot` is an 8×8 dot for connection status or unread items.

<Preview page="Badge" :height="240">

<<< @/../cmd/gallery/badge.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Badge(c, "Local", kit.BadgeNeutral)
kit.Badge(c, "Connected", kit.BadgeSuccess)
```

## Examples

### Status dot

Choose the dot's color yourself, usually a semantic color from the theme. The dot has no text, so provide an accessible name with `.Label(...)`.

```go
t := kit.GetTheme(c)
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	kit.StatusDot(c, t.Colors.Success).Label("Connected")
	ui.Text(c, "Connected")
})
```

## API

### Badge

`func Badge(c *ui.Context, text string, kind BadgeKind) ui.Element`

Height 20, corner radius 4, text 12/16, and a background using the text color at 16% opacity. It does not stretch or shrink within its row.

- **text** `string`: Single-line badge text.
- **kind** `BadgeKind`: Semantic color.

### BadgeKind

- **BadgeNeutral**: Neutral, using the secondary color.
- **BadgeAccent**: Accent color.
- **BadgeSuccess**: Success.
- **BadgeWarning**: Warning.
- **BadgeDanger**: Danger.

### StatusDot

`func StatusDot(c *ui.Context, color ui.Color) ui.Element`

An 8×8 dot that does not stretch or shrink within its row.

- **color** `ui.Color`: Dot color.
