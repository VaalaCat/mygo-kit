---
title: Input
description: Single-line and multiline inputs with read-only values, borderless styling, leading icons, and error borders.
---

# Input

Single-line `TextInput` and multiline `TextArea` inputs. MyGo manages content, undo, input method composition, and selection; these components handle appearance.

<Preview page="Input" :height="520">

<<< @/../cmd/gallery/input.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var name string // Held by the application

kit.TextInput(c.Key("name"), &name).Label("Name").Placeholder("Project name")
```

Give each input an identity with `c.Key(...)` first, or it will lose focus and cursor position when the UI is rebuilt. Use `.Changed()` on the returned element to detect content changes, and `.Submitted()` to detect Enter in a single-line input.

## Examples

### Multiline input

`TextArea` has a minimum height of two lines (56). Adjust it with `.Height(100)` or `.MinHeight(...)`.

```go
kit.TextArea(c.Key("desc"), &desc).Label("Description").Placeholder("Supports multiple lines and input methods").Height(100)
```

### Error state

`Invalid` gives the border the danger color, taking precedence over the focus color. Show error text with [Field](/components/field).

```go
kit.TextInputWith(c.Key("name"), &name, kit.InputOpt{Invalid: name == ""}).Label("Name")
```

### Read-only value

`ReadOnly` displays the input as selectable, copyable text, without a border or background. Its text starting position and height of 28 match a regular input, so it aligns with editable fields inside a `Field`.

```go
kit.Field(c, kit.FieldOpt{Label: "Provider ID"}, func() ui.Element {
	return kit.TextInputWith(c.Key("id"), &providerID, kit.InputOpt{ReadOnly: true})
})
```

Do not use it when the value should look editable.

### Leading icon

`Leading` draws a 14-unit icon in the secondary color on the left of the input. Long text is clipped to the right of the icon and cannot cover it.

```go
kit.TextInputWith(c.Key("search"), &query, kit.InputOpt{Leading: kit.IconSearch, Placeholder: "Search"}).Label("Search")
```

### Borderless input

`Borderless` removes the border, background, padding, and focus ring. The surrounding container controls its appearance. Use it for sidebar search and message composers.

```go
box := ui.Column(c).Padding(kit.GetTheme(c).Space(3)).Radius(8).Border(1, kit.GetTheme(c).Colors.Border)
box.Children(func() {
	kit.TextAreaWith(c.Key("msg"), &msg, kit.InputOpt{Borderless: true, Placeholder: "Message"}).Label("Message")
})
```

### Inside a dialog

In the body of [DialogWith](/components/dialog), only this library's `TextInput` / `TextInputWith` report Enter to the dialog, so Enter confirms it. When the dialog opens, focus moves to the first editable input in its body.

## API

### TextInput

`func TextInput(c *ui.Context, value *string) ui.Element`

A single-line input that uses a zero-value `InputOpt` with `TextInputWith`.

- **value** `*string`: The input content, held by the caller and updated directly as the user types.

### TextInputWith

`func TextInputWith(c *ui.Context, value *string, opt InputOpt) ui.Element`

A single-line input with a height of 28 and vertically centered text.

### TextArea

`func TextArea(c *ui.Context, value *string) ui.Element`

A multiline input that uses a zero-value `InputOpt` with `TextAreaWith`.

### TextAreaWith

`func TextAreaWith(c *ui.Context, value *string, opt InputOpt) ui.Element`

Has a minimum height of 56 and preserves MyGo's composition, selection, and controlled value handling.

### InputOpt

`type InputOpt struct`

The zero value gives a regular editable input.

- **Placeholder** `string`: Hint text shown when the input is empty.
- **Invalid** `bool`: Uses the danger color for the border, taking precedence over the focus color.
- **ReadOnly** `bool`: Displays a plain value that can be selected and copied, without a border or background.
- **Borderless** `bool`: Removes the border, background, padding, and focus ring.
- **Leading** `*ui.SVG`: A leading icon on the left, sized 14 and drawn in the secondary color. Defaults to `nil`, which draws no icon.
