---
title: Progress
description: A linear ProgressBar with a height of 4 and a ProgressRing that changes color at usage thresholds.
---

# Progress

`ProgressBar` is a thin linear progress bar; `ProgressRing` is a circular progress indicator, commonly used for context usage. Both take a `value` from 0 to 1 and clamp values outside that range to 0 or 1.

<Preview page="Progress" :height="360">

<<< @/../cmd/gallery/progress.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.ProgressBar(c, 0.42)
```

The container determines the bar's width, so place it inside a container with a width:

```go
ui.Column(c).Width(320).Children(func() {
	kit.ProgressBar(c, done)
})
```

## Examples

### Usage ring

```go
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	kit.ProgressRing(c, used, 16)
	ui.Text(c, "Context window")
})
```

The ring uses the warning color when `value` reaches 0.85 and the danger color at 0.95. When `value` is greater than 0, the arc spans at least 8° so even a very small proportion is visible. When `value` is 0, only the track ring is drawn.

## API

### ProgressBar

`func ProgressBar(c *ui.Context, value float32) ui.Element`

Has a height of 4, a hover-colored track, an accent-colored fill, and rounded ends. The accessible name is the percentage, such as `42%`.

- **value** `float32`: A proportion from 0 to 1.

### ProgressRing

`func ProgressRing(c *ui.Context, value float32, size float32) ui.Element`

The ring has a stroke width of 2 and is drawn clockwise from the top. Its accessible name always uses the actual percentage, unaffected by the minimum arc length.

- **value** `float32`: A proportion from 0 to 1.
- **size** `float32`: The ring's side length. The gallery uses 14, 16, and 20.
