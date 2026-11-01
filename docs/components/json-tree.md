---
title: JSON
description: A virtualized JSON tree with keyboard expansion, selection, and copying.
---

# JSON

Displays a JSON value as an expandable tree. Objects and arrays show only their member counts; leaves show their values. The virtualized list builds only visible rows, so large JSON values remain smooth to browse.

<Preview page="JSON" :height="360">

<<< @/../cmd/gallery/json.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
state, err := kit.NewJSONState(data) // Parse once; the application holds the state
if err != nil {
	// data is not valid JSON
}

kit.JSONTree(c, state, kit.JSONOpt{}).Grow(1)
```

`NewJSONState` parses any top-level value, including `null`, arrays, and scalars. Numbers retain their original representation. Call `state.Set(newData)` when the content changes. On success, it replaces the tree, preserves expanded paths that still exist, and leaves the root collapsed if the user collapsed it.

## Interaction

- Click a row to select it. If it has children, clicking also toggles expansion.
- Use the up and down arrow keys to move the selection. The right arrow expands a node, or moves to its first child if already expanded. The left arrow collapses a node, or moves to its parent if already collapsed. Enter or Space toggles expansion.
- With a node selected, press Cmd+C (Ctrl+C on Linux and Windows) to copy its JSON. The toolbar's copy button copies the entire value.

## API

### NewJSONState

`func NewJSONState(data []byte) (*JSONState, error)`

Parses JSON and creates the state. The root node is expanded by default.

### JSONState.Set

`func (s *JSONState) Set(data []byte) error`

Parses new content. On failure, it returns an error and leaves the old tree unchanged. Content containing multiple JSON values also returns an error.

### JSONTree

`func JSONTree(c *ui.Context, s *JSONState, opt JSONOpt) ui.Element`

A content box with rounded corners and a thin border. A title bar and copy button sit above the tree. Each row is 20 high, with 16 units of indentation per level.

- **s** `*JSONState`: State held by the caller.
- **opt** `JSONOpt`: See below.

### JSONOpt

- **Title** `string`: Title bar text. Defaults to `JSON`.
- **CopyLabel** `string`: The copy button's tooltip and accessible name. Defaults to “Copy JSON” (localized in Chinese).
