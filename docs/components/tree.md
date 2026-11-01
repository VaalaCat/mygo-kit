---
title: Tree
description: A compact virtualized tree using MyGo's stable keys, expansion state, and arrow key navigation.
---

# Tree

An expandable tree, such as a file directory. MyGo's `ui.Outline` handles expansion state, arrow key navigation, and virtual scrolling. This library only adds styling.

<Preview page="Tree" :height="320">

<<< @/../cmd/gallery/tree.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var state ui.OutlineState[string] // Expansion and selection state, held by the application

children := map[string][]string{
	"workspace": {"cmd", "internal"},
	"cmd":       {"gallery"},
}

kit.Tree(c, &state, []string{"workspace"}, func(id string) []string { return children[id] }, func(id string) {
	ui.Row(c).Height(28).Padding(0, 8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, id).SingleLine()
	})
}).Grow(1)
```

Tree nodes use your own identifier `K`, which can be any comparable type, such as a string or number. The component asks for each node's children through `children` and lets you draw each row through `row`. `ui.Outline` handles row indentation and expansion arrows.

## API

### Tree

`func Tree[K comparable](c *ui.Context, state *ui.OutlineState[K], roots []K, children func(K) []K, row func(K)) ui.Element`

- **state** `*ui.OutlineState[K]`: The expansion state, held by the caller.
- **roots** `[]K`: The root nodes.
- **children** `func(K) []K`: Returns a node's children, or an empty result if it has none.
- **row** `func(K)`: Draws the content of a node's row.
