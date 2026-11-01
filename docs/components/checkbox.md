---
title: Checkbox
description: A square checkbox with a clickable label and state held in the caller's bool.
---

# Checkbox

A 16×16 checkbox with a clickable label. Checking, the Space key, and accessibility behavior come from MyGo's checkbox.

<Preview page="Checkbox" :height="200">

<<< @/../cmd/gallery/checkbox.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var autoSave bool // Owned by the application

kit.Checkbox(c.Key("auto-save"), &autoSave, "Auto-save")
```

Clicking directly toggles `*checked`. Use `.Changed()` on the returned element to check whether it changed this frame.

## Examples

### Disabled

```go
kit.Checkbox(c.Key("disabled"), &autoSave, "Unavailable").Disabled(true)
```

## API

### Checkbox

`func Checkbox(c *ui.Context, checked *bool, label string) ui.Element`

The whole row is 28 high. The box has a corner radius of 4 and an accent-colored background with a checkmark when selected.

- **checked** `*bool`: Checked state, owned by the caller.
- **label** `string`: Label text, also used as the accessible name. Clicking the label toggles the state just like clicking the box.
