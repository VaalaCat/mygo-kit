---
title: Shimmer
description: A single line of secondary text with a highlight sweeping from left to right for ongoing tasks such as thinking or generating a title.
---

# Shimmer

A single line of text in the secondary color, with a highlight band sweeping from left to right every 2 seconds. Use it for ongoing states such as “Thinking” or “Generating title”.

<Preview page="Shimmer" :height="200">

<<< @/../cmd/gallery/shimmer.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Shimmer(c, "Thinking…")
```

When placing multiple Shimmers on the same page, distinguish them with `c.Key(...)`. Each keeps its own timer.

The `Shimmer` field on [ActivityCard](/components/activity-card) also uses this component to animate the title.

## Reduced motion

When the system or theme requests reduced motion, Shimmer becomes static text in the secondary color. It remains readable and searchable.

## API

### Shimmer

`func Shimmer(c *ui.Context, value string) ui.Element`

A single line that does not wrap. The highlight band spans 4 graphemes on each side. Each grapheme transitions between the secondary and foreground colors according to its distance from the center of the highlight.

- **value** `string`: The text to display, also used as the accessible name.
