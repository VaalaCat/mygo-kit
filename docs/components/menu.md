---
title: Menu
description: MenuButton for menus opened by clicking, and ContextMenu for context menus.
---

# Menu

Menus have two entry points: `MenuButton` opens when an element is clicked, and `ContextMenu` opens with a right-click or the Menu key. Both use the same set of `MenuItem` values and have the same keyboard behavior.

<Preview page="Menu" :height="520">

<<< @/../cmd/gallery/menu.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
items := []kit.MenuItem{
	{Icon: kit.IconPencil, Title: "Rename", Shortcut: "F2", OnSelect: rename},
	{Separator: true},
	{Icon: kit.IconX, Title: "Delete", Danger: true, OnSelect: remove},
}
kit.MenuButton(c, kit.IconButton(c.Key("more"), kit.IconEllipsis, "More"), items)
```

Put each menu item's action in `OnSelect`. It is called in the frame when the item is selected and the menu closes. Call the menu every frame; the component stores its open state on the trigger element.

## Examples

### Click menu

Any element can be a trigger; `IconButton` and `ToolButton` are common choices. Use the up and down arrows, Home/End, or an initial letter to move the highlight, Enter to select, and Esc to close. Pressing Enter or Space while the trigger has keyboard focus opens the menu with the first item highlighted.

```go
kit.MenuButton(c, kit.ToolButton(c.Key("ws"), kit.IconFolder, "Workspace", true), []kit.MenuItem{
	{Title: "Local", OnSelect: useLocal},
	{Title: "Remote", OnSelect: useRemote},
})
```

A true `closed` return value means the menu closed in this frame, whether through Esc, a click outside or on the trigger, Tab, or selecting an item. Focus returns to the trigger. To focus something else, set focus yourself in the frame that returns true.

### Context menu

`ContextMenu` adds a context menu to a target element. A right-click opens it at the pointer. When the target has focus, Shift+F10 or the Menu key opens it below the target. The target must be focusable for keyboard activation.

```go
area := ui.Column(c.Key("area")).Height(120).Width(360).Focusable().Children(func() {
	ui.Text(c, "Right-click here")
})
kit.ContextMenu(c, area, items)
```

### Context menu on text

When the target is a `Selectable` text container, right-clicking opens this menu instead of the system text menu. Drag selection still works. If text is selected when the menu opens, a “Copy” item appears at the top and copies that selection.

```go
text := ui.Column(c.Key("text")).Selectable().Children(func() {
	ui.Text(c, "Select part of this text, then right-click.")
})
kit.ContextMenu(c, text, items)
```

The `Selectable` container itself must be the `target`. Attaching the menu to an ancestor still leaves the system menu on the container's text.

## API

### MenuButton

`func MenuButton(c *ui.Context, trigger ui.Element, items []MenuItem) (closed bool)`

Attaches a menu to `trigger` that opens below it on click. Clicking the trigger toggles the menu.

- **trigger** `ui.Element`: Trigger element.
- **items** `[]MenuItem`: Menu items.
- **Returns** `bool`: Whether the menu closed in the current frame.

### ContextMenu

`func ContextMenu(c *ui.Context, target ui.Element, items []MenuItem) (closed bool)`

Attaches a context menu to `target`. The return value has the same meaning as for `MenuButton`.

### MenuItem

- **Icon** `*ui.SVG`: An item icon sized 14. If any item has an icon, items without icons also reserve icon space to keep titles aligned.
- **Title** `string`: Title.
- **Shortcut** `string`: A shortcut hint aligned to the right at the end of the row, such as `"⌘N"`. It only displays the hint and does not register a shortcut.
- **Danger** `bool`: Uses the danger color for the title and icon.
- **Disabled** `bool`: Dims the item and makes it unselectable. Arrow key and initial-letter navigation skip it.
- **Separator** `bool`: When true, draws only a separator for this item and ignores the other fields.
- **OnSelect** `func()`: Called in the frame when the item is selected and the menu closes.
