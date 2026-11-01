---
title: Diff
description: A native full-file diff with added and removed line numbers, selection, and copying, using MyGo's DiffState.
---

# Diff

Compares two texts, either side by side (`editor.DiffSplit`, the default) or in one combined column (`editor.DiffUnified`). The diff algorithm, selection, and copying all come from the MyGo editor's `DiffView`. This library only sets colors and fonts to match the other panels.

To display an already-parsed patch instead of comparing two complete texts, use [PatchLines](/components/patch).

<Preview page="Diff" :height="320">

<<< @/../cmd/gallery/diff.go

</Preview>

## Usage

```go
import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui/editor"
)
```

```go
var diff = editor.DiffState{
	Old:             editor.NewBuffer(oldText),
	New:             editor.NewBuffer(newText),
	Layout:          editor.DiffUnified,
	ExpandUnchanged: true,
}

kit.Diff(c, &diff).Grow(1)
```

The caller owns `DiffState` and configures options such as layout and expansion of unchanged lines on it. See MyGo's [diff viewer documentation](https://github.com/egoist/mygo/blob/main/docs/ui/diff-viewer.md).

Every frame, `Diff` overwrites `state.Theme`, `state.Font`, and `state.LineHeight` (1.5), so do not set these yourself in `DiffState`.

## API

### Diff

`func Diff(c *ui.Context, state *editor.DiffState) ui.Element`

Added lines use the success color, removed lines the danger color, changes the warning color, and line numbers the placeholder color. The font is the theme's monospace font, and the background is the editor background.

- **state** `*editor.DiffState`: MyGo diff state, owned by the caller.
