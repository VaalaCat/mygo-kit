---
title: CompletionMenu
description: A completion popup above an input that keeps focus in the input and handles the up and down arrows, Enter, Tab, and Esc.
---

# CompletionMenu

A completion popup above an input, such as a list of commands after typing `/`. The input keeps focus so the user can continue typing; the popup handles only a few keys.

<Preview page="CompletionMenu" :height="480">

<<< @/../cmd/gallery/completion.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
box := ui.Column(c.Key("composer")).Padding(8).Radius(8).Border(1, t.Colors.Border)
var input ui.Element
box.Children(func() {
	input = kit.TextAreaWith(c.Key("input"), &text, kit.InputOpt{Borderless: true}).Label("Message")
})

opt := kit.CompletionOpt{Open: strings.HasPrefix(text, "/"), Items: matches(text)}
if i := kit.CompletionMenu(c, box, input, opt); i >= 0 {
	text = opt.Items[i].Title + " "
}
```

The caller is responsible for:

- Deciding when to open (`Open`) and supplying the current candidates every frame. The component does not filter or sort them.
- Handling selection. A return value `>= 0` is the index of the candidate selected this frame; write it back to the input.

## Keys

When candidates are available:

- The up and down arrows cycle through the highlight.
- Enter or Tab selects.
- Esc closes the popup. It will not appear again until `Open` returns to false.

When there are no candidates, only Esc is handled. No keys are handled during IME composition.

The `input` element's `HandleInput` handles these keys. Do not set another `HandleInput` on the same input; `OnTextKey` is unaffected.

## Examples

### Loading and errors

Set `Loading` while candidates are loading. The loading row appears after the candidates, or on its own when there are none. Set `Error` on failure. When neither is set and there are no candidates, `Empty` is shown.

```go
kit.CompletionOpt{Open: true, Loading: true}
kit.CompletionOpt{Open: true, Error: "Unable to load skills: connection refused"}
kit.CompletionOpt{Open: true, Empty: "No matching skills"}
```

## API

### CompletionMenu

`func CompletionMenu(c *ui.Context, anchor, input ui.Element, opt CompletionOpt) (selected int)`

The popup appears above `anchor`, with at most 8 rows (including candidates, loading, and error rows), and scrolls when it overflows. The left edge of row content aligns with the text in `input`.

- **anchor** `ui.Element`: Element the popup is attached to, usually the container wrapping the input.
- **input** `ui.Element`: Input element that keeps focus.
- **opt** `CompletionOpt`: Popup data and state.
- **Return value** `int`: Index of the candidate selected this frame, or `-1` if none was selected.

### CompletionOpt

- **Open** `bool`: Shows the popup only when true. The caller decides based on the input, such as whether the word at the cursor starts with `/`.
- **Items** `[]CompletionItem`: Candidate list.
- **Loading** `bool`: Shows a loading row after the candidates.
- **Error** `string`: Shows an error row when non-empty.
- **Empty** `string`: Message shown when there are no candidates, loading, or error. Defaults to the built-in “No matches” when empty.

### CompletionItem

- **Icon** `*ui.SVG`: Leading icon.
- **Title** `string`: Title.
- **Description** `string`: Description in the secondary color, shown to the right of the title and truncated when it does not fit.
