---
title: Facts
description: A two-column group of facts, with labels on the left and monospace values sharing one left alignment on the right.
---

# Facts

A two-column list of labels and values, suited to information such as models, paths, and usage in a detail panel. The label column has a fixed width of 112, and all values share the same left alignment. Values can be selected and copied.

The gallery has no separate Facts page.

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Facts(c,
	kit.Fact{Label: "Model", Value: "claude-sonnet-5-5"},
	kit.Fact{Label: "Directory", Value: "~/code/project"},
	kit.Fact{Label: "Context", Value: "42%"},
)
```

The left column uses 12/16 labels in the secondary color. The right column uses 12/20 monospace values. Both columns align on the first line's baseline. Long values wrap.

## API

### Facts

`func Facts(c *ui.Context, facts ...Fact) ui.Element`

There is no gap between rows.

- **facts** `...Fact`: Rows to display in order.

### Fact

`type Fact struct{ Label, Value string }`

- **Label** `string`: Single-line label.
- **Value** `string`: Selectable value in a monospace font.
