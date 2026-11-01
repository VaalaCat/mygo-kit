---
title: Field
description: A vertical Field with a label above the control and a description or error below, plus standalone FieldMessage text and InlineHint for lists.
---

# Field

Form fields. `Field` stacks a label, control, and description, and connects the label and description to the control's accessible name and description. `FieldMessage` is a standalone message line. `InlineHint` shows messages such as “Loading” or “Empty list” in a list.

<Preview page="Field" :height="520">

<<< @/../cmd/gallery/form.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Field(c, kit.FieldOpt{
	Label:       "Provider ID",
	Description: "Used in config.yaml and as the prefix for model names.",
}, func() ui.Element {
	return kit.TextInputWith(c.Key("id"), &providerID, kit.InputOpt{Placeholder: "openai"})
})
```

`control` must return the control element so that `Field` can set its label and description. Clicking the label focuses the control.

The caller is responsible for:

- Validating input and placing the error text in `Error`.
- Setting danger-colored borders on input controls through `InputOpt.Invalid`. `Field` only displays the error text.

## Examples

### Errors

When `Error` is non-empty, the control is marked invalid and a danger-colored error appears below it. If both a description and an error are present, only the error is shown.

```go
kit.Field(c, kit.FieldOpt{Label: "Name", Error: "Use only letters, numbers, and hyphens."}, func() ui.Element {
	return kit.TextInputWith(c.Key("name"), &name, kit.InputOpt{Invalid: true})
})
```

### Controls other than inputs

Any control can go in `Field`, such as [Select](/components/select) and [Switch](/components/switch):

```go
kit.Field(c, kit.FieldOpt{Label: "Inherit settings"}, func() ui.Element {
	return kit.Switch(c.Key("inherit"), &inherit, "")
})
```

### Standalone message line

```go
kit.FieldMessage(c, kit.MessageWarning, "This value overrides the default configuration.")
```

### Messages in lists

A spinner precedes the message while loading. An empty list shows only text.

```go
kit.InlineHint(c, "Loading…", true)
kit.InlineHint(c, "No conversations yet", false)
```

## API

### Field

`func Field(c *ui.Context, opt FieldOpt, control func() ui.Element)`

Arranged vertically with a gap of 4 between sections. Returns nothing.

- **opt** `FieldOpt`: Label, description, and error.
- **control** `func() ui.Element`: Draws the control and returns its element.

### FieldOpt

- **Label** `string`: Label in the secondary color. When empty, it is hidden and the control's name is not set.
- **Description** `string`: Description.
- **Error** `string`: Error. When both description and error are empty, no space is left below the control.

### FieldMessage

`func FieldMessage(c *ui.Context, kind MessageKind, text string) ui.Element`

A line of 12/16 text that can wrap.

### MessageKind

- **MessageDescription**: Description, using the secondary color.
- **MessageWarning**: Warning color.
- **MessageError**: Danger color.

### InlineHint

`func InlineHint(c *ui.Context, message string, loading bool) ui.Element`

A 28-high row with a secondary message. When `loading` is true, a spinner appears at the start of the row. When `message` is empty and `loading` is true, the built-in “Loading…” is shown.

### InlineHintWith

`func InlineHintWith(c *ui.Context, message string, loading bool, opt InlineHintOpt) ui.Element`

The same as `InlineHint`, with the left alignment set by `opt`.

### InlineHintOpt

- **Inset** `float32`: Distance from the left edge to the text. Use 12 (the zero-value default) for lists and 16 for conversation columns.
