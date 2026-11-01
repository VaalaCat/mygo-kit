---
title: Layout
description: Layout helpers for panel headers, separators, resizable split panels, and empty states.
---

# Layout

A few small layout helpers. They have no separate gallery page and appear in the shells of other pages. For example, [Workspace](/components/workspace) uses `Header` and `Split`, and the center area of the gallery's Workspace page is an `EmptyState`.

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

### Header

A panel header with a height of 36. The title is on the left, and `actions` draws toolbar buttons on the right. Its bottom separator is shared with the content below.

```go
kit.Header(c, "Conversations", func() {
	kit.IconButtonTip(c.Key("new"), kit.IconPlus, "New")
})
```

Use `HeaderWith` to replace the title with other content or align it to a different left edge:

```go
kit.HeaderWith(c, "", kit.HeaderOpt{Inset: 16, Leading: func() {
	kit.Segmented(c, &view, []string{"Conversations", "Archived"})
}}, nil)
```

### Separator

A horizontal separator one pixel high.

```go
kit.Separator(c)
```

### Split

Two adjacent resizable panels, with a MyGo divider between them. `size` is the width of the first panel and is updated directly when the divider is dragged.

```go
sidebarWidth := float32(240) // Held by the application

kit.Split(c, &sidebarWidth, func() {
	// Left side
}, func() {
	// Right side
}).Grow(1)
```

### EmptyState

Displays a message in the center of an area, with an optional secondary description line and primary button:

```go
kit.EmptyState(c, "No conversations yet", "Click the plus button at the top left to create one.")

kit.EmptyStateWith(c, "No conversations yet", "", kit.EmptyStateOpt{
	ActionLabel: "New conversation",
	OnAction:    newConversation,
})
```

## API

### Header

`func Header(c *ui.Context, title string, actions func()) ui.Element`

Has a height of 36, with the title inset 12 from the left edge and truncated to one line.

- **title** `string`: Title text.
- **actions** `func()`: Draws toolbar buttons on the right. May be `nil`.

### HeaderWith

`func HeaderWith(c *ui.Context, title string, opt HeaderOpt, actions func()) ui.Element`

### HeaderOpt

- **Inset** `float32`: Distance from the left edge to the title text, used to align it with the column's left edge. Sidebars and docks use 12; conversation areas use 16. A zero value uses 12.
- **Leading** `func()`: Replaces the title text, for example with a segmented control. Content starts at `Inset`; the header's height and bottom separator stay the same. When unset, displays `title`.

### Separator

`func Separator(c *ui.Context) ui.Element`

Has a height of 1, uses `Colors.BorderSubtle`, and does not shrink.

### Split

`func Split(c *ui.Context, size *float32, first, second func()) ui.Element`

Both panels fill their available space, with no outer border and a corner radius of 0.

- **size** `*float32`: The width of the first panel, held by the caller. Dragging the divider updates it. Use `.Changed()` on the returned element to detect whether it was dragged in the current frame.
- **first**, **second** `func()`: Draw the contents of the two panels, respectively.

### EmptyState

`func EmptyState(c *ui.Context, title, detail string) ui.Element`

Centers its content and fills the parent's remaining space. The title uses the secondary color; the description uses the placeholder color.

- **title** `string`: Title.
- **detail** `string`: Description. May be empty.

### EmptyStateWith

`func EmptyStateWith(c *ui.Context, title, detail string, opt EmptyStateOpt) ui.Element`

### EmptyStateOpt

- **ActionLabel** `string`: When nonempty, displays a primary button centered below the text.
- **OnAction** `func()`: Called when the primary button is clicked.
