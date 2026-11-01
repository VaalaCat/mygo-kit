---
title: Dialog
description: Modal dialogs with a title and body, plus DialogWith for cancel and confirm buttons.
---

# Dialog

Modal dialogs. MyGo handles the focus trap and closing with Esc. `Dialog` has only a title bar and body; use `DialogWith` when you need cancel and confirm buttons.

<Preview page="Dialog" :height="460">

<<< @/../cmd/gallery/dialog.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // Owned by the application; set to true to open

if kit.Button(c, "Rename").Clicked() {
	open = true
}
if kit.DialogWith(c.Key("rename"), &open, kit.DialogOpt{Title: "Rename", ConfirmLabel: "Save"}, func() {
	kit.Field(c, kit.FieldOpt{Label: "Title"}, func() ui.Element {
		return kit.TextInput(c.Key("title"), &title)
	})
}) {
	// The user clicked confirm or pressed Enter: save title
}
```

Call the dialog every frame. It draws nothing when `open` is false. In the frame when `DialogWith` returns true, `*open` has already been reset to false.

## Examples

### Basic dialog

You draw the entire body of `Dialog`, including its close button.

```go
kit.Dialog(c, &open, "Native dialog", func() {
	ui.Text(c, "Focus stays inside the dialog. Esc closes it.")
	if kit.PrimaryButton(c, "Done").Clicked() {
		open = false
	}
})
```

### Confirmation and dangerous actions

Pass an empty function for `body` when there is no body content. `Danger` gives the confirm button the danger color.

```go
if kit.DialogWith(c.Key("delete"), &open, kit.DialogOpt{
	Title:        "Delete conversation",
	Description:  "This cannot be undone.",
	ConfirmLabel: "Delete",
	Danger:       true,
}, func() {}) {
	// Delete
}
```

### Dialog with input

Use `ConfirmDisabled` to disable the confirm button when the confirmation requirements are not met. Enter will not confirm either.

```go
kit.DialogWith(c.Key("rename"), &open, kit.DialogOpt{Title: "Rename", ConfirmDisabled: title == ""}, func() { /* Input */ })
```

Initial focus and Enter behavior:

- On opening, focus moves to the first editable input in the body, excluding read-only inputs. Tabs and buttons do not receive initial focus.
- When the body has no editable input, focus moves to the confirm button without drawing a focus ring. Enter still confirms.
- Enter confirms in a single-line input, inserts a newline in a multiline input, and presses a focused button. Single-line inputs must use this library's `TextInput` or `TextInputWith` to pass Enter to the dialog.
- The cancel button and Esc only close the dialog and return false.

## API

### Dialog

`func Dialog(c *ui.Context, open *bool, title string, content func()) ui.Element`

A modal dialog with a title bar and body, 480 wide (at most 90% of the window).

- **open** `*bool`: Whether the dialog is open, owned by the caller.
- **title** `string`: Title, also used as the accessible name.
- **content** `func()`: Function that draws the body.

### DialogWith

`func DialogWith(c *ui.Context, open *bool, opt DialogOpt, body func()) (confirmed bool)`

A title bar, description, and body, with right-aligned “Cancel” and confirm buttons. Clicking confirm or pressing Enter closes the dialog and returns true.

- **body** `func()`: Function that draws the body.

### DialogOpt

- **Title** `string`: Title.
- **Description** `string`: Description in the secondary color, below the title bar and before the body. Hidden when empty.
- **ConfirmLabel** `string`: Confirm button text. Defaults to the built-in localized “Confirm” label when empty.
- **Danger** `bool`: Uses `DangerButton` for confirmation; otherwise uses `PrimaryButton`.
- **ConfirmDisabled** `bool`: Disables the confirm button and prevents Enter from confirming.
