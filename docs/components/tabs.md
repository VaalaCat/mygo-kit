---
title: Tabs
description: A continuous tab strip and the closable Tab used in Docks and terminals.
---

# Tabs

A tab strip. The selected tab uses the main content background and has no bottom border, joining it visually to the content below.

<Preview page="Tabs" :height="240">

<<< @/../cmd/gallery/tabs.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
tab := 0 // Index of the selected tab, held by the application

kit.Tabs(c, &tab, []string{"General", "Editor", "Advanced"})
switch tab {
case 0:
	// Draw the General content
}
```

`Tabs` only draws the tab strip. Draw the content below it yourself based on `tab`. MyGo provides arrow key navigation. The focus ring appears only with keyboard focus and always follows the selected tab.

## Examples

### Closable tabs

`Tab` draws a single Zed-style tab, typically used for panels that can have multiple instances, such as terminals and editors. It returns whether the tab or its close button was clicked this frame. The caller decides what happens after selection or closure.

```go
for i, session := range sessions {
	clicked, closed := kit.Tab(c.Key(session.ID), kit.TabOpt{
		Icon:       kit.IconTerminal,
		Title:      session.Title,
		Selected:   i == current,
		CloseLabel: "Close",
	})
	if clicked {
		current = i
	}
	if closed {
		// Remove session
	}
}
```

When placed in a [Workspace](/components/workspace) panel's `Panel.Header`, tabs match the Dock title bar height and scroll horizontally when they do not fit.

## API

### Tabs

`func Tabs(c *ui.Context, selected *int, labels []string) ui.Element`

A tab strip with a height of 36, vertical separators between tabs, and a bottom border across the empty space on the right.

- **selected** `*int`: The index of the selected tab, held by the caller.
- **labels** `[]string`: Each tab's text, also used as its accessible name.

### Tab

`func Tab(c *ui.Context, opt TabOpt) (clicked, closed bool)`

A closable tab with square corners and the same height as the title bar (36). A selected tab uses the main content background with vertical borders on both sides. An unselected tab uses the panel background with right and bottom borders. The close button appears on hover or when selected. A newly selected tab scrolls into view within its tab strip.

### TabOpt

- **Icon** `*ui.SVG`: An icon before the title, size 14. Can be `nil`.
- **Title** `string`: The title, with a maximum width of 192. Longer text is truncated with an ellipsis.
- **Selected** `bool`: Whether the tab is selected.
- **Faded** `bool`: Uses a muted color for the title, for example for a terminal whose process has exited.
- **CloseLabel** `string`: The close button's accessible name.
