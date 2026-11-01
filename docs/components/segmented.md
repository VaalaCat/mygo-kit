---
title: Segmented
description: A row of mutually exclusive options with optional icons and tooltips.
---

# Segmented

Styled like Zed's ToggleButtonGroup: a height of 28, a thin outer border, and a vertical separator between segments. The selected segment uses a light accent background and accent-colored text. Arrow keys switch options, and the whole control is one Tab stop.

<Preview page="Segmented" :height="460">

<<< @/../cmd/gallery/segmented.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
appearance := 0 // Index of the selected segment, held by the application

kit.Segmented(c.Key("appearance"), &appearance, []string{"System", "Light", "Dark"}).Label("Appearance").AlignSelf(ui.Start)
```

The width fits the content by default. Add `AlignSelf(ui.Start)` when placing it inside a column to align it to the left. The control has no visible title; give it an accessible name with `.Label(...)`.

## Examples

### Icons and tooltips

Use `SegmentedWith` and `SegmentedItem`. An empty `Tooltip` displays no tooltip.

```go
sections := []kit.SegmentedItem{
	{Icon: kit.IconFolder, Label: "Workspace", Tooltip: "Executors and directories"},
	{Icon: kit.IconServerCog, Label: "MCP", Tooltip: "Connected MCP servers"},
}
kit.SegmentedWith(c.Key("sections"), &current, sections, kit.SegmentedOpt{}).Label("Section")
```

### Full width

`FullWidth` fills the container and gives every segment equal width.

```go
kit.SegmentedWith(c.Key("full"), &current, sections, kit.SegmentedOpt{FullWidth: true})
```

## API

### Segmented

`func Segmented(c *ui.Context, selected *int, labels []string) ui.Element`

A text-only segmented control whose width fits its content. Equivalent to `SegmentedWith` without icons or tooltips.

- **selected** `*int`: The index of the selected segment, held by the caller.
- **labels** `[]string`: The text for each segment.

### SegmentedWith

`func SegmentedWith(c *ui.Context, selected *int, items []SegmentedItem, opt SegmentedOpt) ui.Element`

### SegmentedItem

- **Icon** `*ui.SVG`: A segment icon sized 12. May be `nil`.
- **Label** `string`: The segment's text and accessible name.
- **Tooltip** `string`: Hover tooltip. May be empty.

### SegmentedOpt

- **FullWidth** `bool`: Fills the container width with equal-width segments. Defaults to `false`, so each segment fits its own content.
