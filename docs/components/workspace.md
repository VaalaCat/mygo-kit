---
title: Workspace
description: A complete window layout with a title bar, sidebar, center area, right Dock, bottom Dock, and status bar, switching to drawers in narrow windows.
---

# Workspace

A complete window layout: title bar, left sidebar, center area, right Dock, bottom Dock, and status bar, with draggable dividers. When the window width is less than 900, the sidebar and right Dock automatically become drawers.

<Preview page="Workspace" :height="520" :width="1180">

<<< @/../cmd/gallery/workspace.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
// Held and persisted by the application
var layout = kit.WorkspaceState{
	SidebarOpen: true,
	Right:       kit.DockState{Open: true, ActivePanel: "outline"},
}

kit.Workspace(c, &layout, kit.WorkspaceOpt{
	TitleBar: func() { ui.Text(c, "~/code/project") },
	Sidebar:  func() { kit.Header(c, "Sessions", nil) },
	Center:   func() { /* Main content */ },
	RightPanels: []kit.Panel{
		{ID: "outline", Title: "Outline", Icon: kit.IconFiles, Content: func() { /* ... */ }},
	},
})
```

Pass content areas as functions, called every frame. An area set to `nil` is hidden. `Workspace` draws its own status bar with a sidebar toggle and a button for each panel.

## Layout state

`WorkspaceState` stores the layout of the sidebar and both Docks. Its fields have `json` tags, so the application can serialize it directly for persistence. A size of 0 uses the default: 240 for the sidebar, 360 for the right Dock, and 240 for the bottom Dock. Dragging a divider writes the new size back to `state`.

In narrow windows, the open states of the sidebar and right Dock drawers are kept only in memory, outside the persisted data. The original layout is restored when the window becomes wide again.

Use the state's methods to open and close panels instead of changing fields directly, so drawers in narrow windows are handled correctly too:

```go
layout.ToggleSidebar()
layout.TogglePanel(&layout.Bottom, "terminal")
layout.ShowPanel(&layout.Right, "info")
```

## Examples

### Panel title bars

`Panel.Header` can replace the title text in a Dock's title bar, for example with a row of [Tab](/components/tabs) components. Tabs start at the left edge of the title bar and scroll horizontally when they do not fit. `Actions` draws at the right end of the title bar, before the close button.

```go
kit.Panel{
	ID:    "terminal",
	Title: "Terminal",
	Icon:  kit.IconTerminal,
	Header: func() {
		for _, s := range sessions { /* kit.Tab(...) */ }
	},
	Actions: func() { kit.IconButtonTip(c.Key("new"), kit.IconPlus, "New") },
	Content: func() { /* kit.Terminal(...) */ },
}
```

### Window button placement

When using the `TitleBar` area, `Workspace` leaves room for the window buttons. See [TitleBar](/components/title-bar) for window creation details.

## API

### Workspace

`func Workspace(c *ui.Context, state *WorkspaceState, opt WorkspaceOpt) ui.Element`

### WorkspaceOpt

- **TitleBar** `func()`: Content on the left side of the title bar.
- **Sidebar** `func()`: Sidebar content.
- **Center** `func()`: Center area content.
- **RightPanels** `[]Panel`: Panels in the right Dock. Only one is displayed at a time.
- **BottomPanels** `[]Panel`: Panels in the bottom Dock.
- **StatusLeft** `func()`: Extra items on the left side of the status bar, after the sidebar toggle.
- **StatusRight** `func()`: Extra items on the right side of the status bar, before the panel buttons.

### Panel

- **ID** `string`: The panel identifier, corresponding to `DockState.ActivePanel`.
- **Title** `string`: The title, also used as the status bar button's tooltip.
- **Icon** `*ui.SVG`: The status bar button's icon.
- **Content** `func()`: Panel content.
- **Header** `func()`: Replaces the title text in the Dock's title bar. When absent, `Title` is shown.
- **Actions** `func()`: Buttons at the right end of the title bar, before the close button.
- **CloseIcon** `*ui.SVG`: Overrides the close button's icon. When absent, the icon follows the Dock's edge: `IconPanelBottomClose` for the bottom Dock and `IconPanelRightClose` for the right Dock.

### WorkspaceState

- **SidebarOpen** `bool`: Whether the sidebar is open in a wide window.
- **SidebarWidth** `float32`: The sidebar width.
- **Right**, **Bottom** `DockState`: The right and bottom Docks.

### DockState

- **Open** `bool`: Whether the Dock is open.
- **ActivePanel** `string`: The current panel's `ID`. If invalid, the first panel is used.
- **Size** `float32`: Width for the right Dock, height for the bottom Dock.

### WorkspaceState methods

- **IsNarrow() bool**: Was the window narrow (less than 900 wide) in the most recent frame?
- **IsPanelShown(dock *DockState, panelID string) bool**: Is the specified panel currently visible? Checks the Dock in wide windows and the right drawer in narrow windows.
- **ShowPanel(dock *DockState, panelID string)**: Opens the Dock and shows the specified panel. For the right Dock in narrow windows, opens its drawer and closes the sidebar drawer without changing `Right.Open`.
- **TogglePanel(dock *DockState, panelID string)**: Shows the specified panel. If it is already visible, closes the Dock (or the drawer in narrow windows).
- **ToggleDock(dock *DockState)**: Toggles the Dock.
- **ShowSidebar()**: Shows the sidebar. In narrow windows, opens the sidebar drawer and closes the right Dock drawer.
- **ToggleSidebar()**: Toggles the sidebar.
