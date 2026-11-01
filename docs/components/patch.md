---
title: Patch
description: Displays parsed patch fragments with colored additions and deletions, optional line numbers, and copying of caller-supplied source.
---

# Patch

Displays a patch fragment, such as changes returned by a tool call. It does not compare files or infer line numbers; it displays exactly the lines you supply. Use [Diff](/components/diff) to compare two complete texts.

<Preview page="Patch" :height="260">

<<< @/../cmd/gallery/patch.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows ui.ListState // Held by the application

n := 42
kit.PatchLines(c, &rows, []kit.PatchLine{
	{Kind: kit.PatchHeader, Text: "@@ component.go @@"},
	{Kind: kit.PatchDeleted, Text: "return oldValue"},
	{Kind: kit.PatchAdded, Text: "return newValue", NewLine: &n},
}, kit.PatchOpt{Source: patchText}).Grow(1)
```

The caller is responsible for:

- Parsing the patch into individual `PatchLine` values.
- Providing line numbers. `OldLine` and `NewLine` are `*int`; pass `nil` for a blank. The component does not infer them.
- Providing the source for copying. The copy button copies `PatchOpt.Source` rather than reconstructing a patch from line fragments, so it does not fabricate a complete file.

## API

### PatchLines

`func PatchLines(c *ui.Context, state *ui.ListState, lines []PatchLine, opt PatchOpt) ui.Element`

A content box with rounded corners and a thin border. A title bar and copy button sit above a virtualized list with selectable text and horizontal scrolling.

- **state** `*ui.ListState`: List state, held by the caller.
- **lines** `[]PatchLine`: All lines.
- **opt** `PatchOpt`: See below.

### PatchLine

- **Kind** `PatchKind`: The kind of line.
- **Text** `string`: Line content without a `+` or `-` prefix. The component draws the symbol according to `Kind`.
- **OldLine**, **NewLine** `*int`: Line numbers in the old and new files. `nil` displays a blank.

### PatchKind

- **PatchContext**: A context line with no background.
- **PatchAdded**: An added line with `+` and a light success-colored background.
- **PatchDeleted**: A deleted line with `−` and a light danger-colored background.
- **PatchHeader**: A block header with accent-colored text and a panel background.

### PatchOpt

- **Title** `string`: Title bar text. Defaults to `Patch`.
- **CopyLabel** `string`: The copy button's tooltip and accessible name. Defaults to “Copy output” (localized in Chinese).
- **Source** `string`: The complete source written to the clipboard when Copy is clicked.
