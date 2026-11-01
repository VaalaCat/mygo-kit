---
title: ListRow
description: Full-width list rows with one or two lines, indentation, trailing status, and hover actions.
---

# ListRow

A row in a sidebar, conversation list, or file tree. The entire row is a button. The theme supplies its selected, hovered, and pressed background colors.

<Preview page="ListRow" :height="520">

<<< @/../cmd/gallery/listrow.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
for _, name := range names {
	row := kit.ListRow(c.Key(name), kit.ListRowOpt{
		Icon:     kit.IconMessageSquare,
		Title:    name,
		Selected: name == current,
	})
	if row.Clicked() {
		current = name
	}
}
```

Use `c.Key(...)` to distinguish rows in the same list. The caller passes the selection state through `Selected`; the component does not store it.

## Examples

### Two lines

When `Description` is nonempty, the row has two lines and a height of 44. The description uses 12/16 sizing in the secondary color and is truncated to one line.

```go
kit.ListRow(c.Key("a"), kit.ListRowOpt{Title: "Add client retries", Description: "openvaala · 3 minutes ago"})
```

### Trailing status

`Trailing` draws a badge, spinner, or status dot. It aligns with the title's baseline in a single-line row and is vertically centered in a two-line row.

```go
kit.ListRow(c.Key("b"), kit.ListRowOpt{
	Title:       "Update changelog",
	Description: "docs · Yesterday",
	Trailing:    func() { kit.StatusDot(c, kit.GetTheme(c).Colors.Accent).Label("Unread") },
})
```

### Hover actions

`Actions` draws 28-square icon buttons at the end of the row. They appear only when the row is hovered, selected, or contains keyboard focus. Their space stays reserved when hidden, so the row's content does not shift.

```go
kit.ListRow(c.Key("c"), kit.ListRowOpt{
	Title: "A conversation",
	Actions: func() {
		kit.IconButtonTip(c.Key("more"), kit.IconEllipsis, "More")
	},
})
```

### Indentation

`Depth` indents each level 16 units to the right for directory trees.

```go
kit.ListRow(c.Key("file"), kit.ListRowOpt{Icon: kit.IconFileText, Title: "composer.go", Depth: 2})
```

### Disabled

```go
kit.ListRow(c.Key("off"), kit.ListRowOpt{Icon: kit.IconBan, Title: "Unavailable"}).Disabled(true)
```

## API

### ListRow

`func ListRow(c *ui.Context, opt ListRowOpt) ui.Element`

Has a height of 28, or 44 with a description. The returned element is a button. Use `Clicked()` to handle clicks, or chain `FillWidth`, `Key`, and `Label`.

### ListRowOpt

- **Icon** `*ui.SVG`: A leading icon sized 14 and centered in a slot 16 wide. No icon space is reserved when `nil`.
- **Title** `string`: The title, truncated to one line. Also serves as the accessible name.
- **Description** `string`: When nonempty, makes the row two lines and 44 high.
- **Selected** `bool`: Shows the selected background color.
- **Depth** `int`: Indentation level, with 16 units per level. Defaults to 0.
- **Trailing** `func()`: Draws a trailing status or badge.
- **Actions** `func()`: Draws trailing action buttons, visible only when the row is hovered, selected, or contains keyboard focus. Setting it reduces the right padding to 4, aligning the buttons' right edge with icon buttons in the header.
