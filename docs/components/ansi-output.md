---
title: ANSIOutput
description: Displays command output with colors, strips other terminal control codes, and supports virtual scrolling, selection, and copying.
---

# ANSIOutput

A content area for command-line output. It interprets only SGR colors and styles (bold, underline, 16/256 colors, and true color). All other terminal control codes (cursor movement, OSC, and so on) are stripped and never executed. It uses a virtual list to keep long output smooth, with horizontal scrolling.

For a complete, interactive terminal, use [Terminal](/components/terminal).

<Preview page="Output" :height="320">

<<< @/../cmd/gallery/output.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows = ui.ListState{FollowEnd: true} // Owned by the application; FollowEnd follows the bottom

kit.ANSIOutput(c, output, &rows, kit.ANSIOutputOpt{}).Grow(1)
```

`output` is the entire output so far, passed in every frame. `ListState.FollowEnd` determines whether to follow the latest output.

## Examples

### Displaying the original text

The default mode interprets colors and strips control codes. When `Plain` is true, the text is displayed without interpretation. Copying also preserves the original text, byte for byte.

```go
kit.ANSIOutput(c, log, &rows, kit.ANSIOutputOpt{Plain: true})
```

## Selection and copying

Drag to select across visible rows that have been built. Historical rows that have not been built (those scrolled out of view) cannot be selected. Use the copy button in the toolbar to copy the full text. In the default mode, the copied full text has control codes removed.

## API

### ANSIOutput

`func ANSIOutput(c *ui.Context, source string, state *ui.ListState, opt ANSIOutputOpt) ui.Element`

- **source** `string`: Full output text. Parsed results are cached while the content and theme remain unchanged.
- **state** `*ui.ListState`: List state, owned by the caller.
- **opt** `ANSIOutputOpt`: See below.

### ANSIOutputOpt

- **Title** `string`: Title bar text. Defaults to the localized “Output” label.
- **CopyLabel** `string`: Copy button tooltip and accessible name. Defaults to the localized “Copy output” label.
- **Plain** `bool`: Displays and copies the original text unchanged. Defaults to `false`.
