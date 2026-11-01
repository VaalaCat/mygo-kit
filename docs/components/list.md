---
title: List
description: A compact virtualized list that builds only visible rows and supports variable row heights.
---

# List

A virtualized list for long lists. It builds only visible rows, stays smooth with ten thousand rows, and supports different heights for each row. It wraps MyGo's `ui.List` with workbench styling: no corner radius, no padding, and the editor background color.

<Preview page="Virtual list" :height="360">

<<< @/../cmd/gallery/list.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows = ui.ListState{Key: func(i int) any { return items[i].ID }} // Held by the application

kit.List(c, &rows, len(items), func(i int) {
	ui.Text(c, items[i].Title).Padding(0, 16)
}).Grow(1)
```

The caller is responsible for:

- Holding `ui.ListState`, which stores row positions and scrolling.
- Setting `Key` to give each row a stable identity. This keeps the scroll position from jumping when data is inserted or removed at the start.
- Using `row` to draw row `i`.

## Related components

- Use [MessageScroller](/components/message-scroller) for a list that follows the latest messages and has a “Back to bottom” button.
- Use [Tree](/components/tree) for tree structures.
- Use [ListRow](/components/list-row) to draw a clickable item in each row.

## API

### List

`func List(c *ui.Context, state *ui.ListState, count int, row func(int)) ui.Element`

- **state** `*ui.ListState`: List state, held by the caller.
- **count** `int`: Total number of rows.
- **row** `func(int)`: Draws row `i`. Called only for visible rows.
