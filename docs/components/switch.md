---
title: Switch
description: A 36×20 toggle with an accent-colored track when on and an animated thumb.
---

# Switch

A toggle switch. When on, the track uses the accent color and the thumb slides across. When the system or theme requests reduced motion, the thumb moves directly to its final position.

<Preview page="Switch" :height="360">

<<< @/../cmd/gallery/switch.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var reduceMotion bool // Held by the application

kit.Switch(c.Key("motion"), &reduceMotion, "Reduce motion")
```

## Examples

### Track only

Pass an empty string for `label` to draw only the track, then provide an accessible name with `.Label(...)`. This is common inside a [SettingRow](/components/settings) or [Field](/components/field), where the adjacent row title provides the visible text.

```go
kit.Switch(c.Key("auto-save"), &autoSave, "").Label("Auto save")
```

### Disabled

The framework makes disabled controls translucent.

```go
kit.Switch(c.Key("off"), &on, "Unavailable").Disabled(true)
```

## API

### Switch

`func Switch(c *ui.Context, on *bool, label string) ui.Element`

- **on** `*bool`: The switch state, held by the caller. Clicking or pressing Space toggles it.
- **label** `string`: Text displayed to the right of the track. The entire row (28 high) is clickable. When empty, only the track (20 high) is shown, with no text.
