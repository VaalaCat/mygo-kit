---
title: CommandPalette
description: A searchable command palette that groups commands and returns the selected command ID.
---

# CommandPalette

A command palette centered near the top of the window, with a search box above grouped commands. The search box receives focus automatically. Use the up and down arrows to move the highlight, Enter to select, and Esc to close.

<Preview page="CommandPalette" :height="460">

<<< @/../cmd/gallery/command.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // Owned by the application; set to true in a shortcut handler to open

commands := []kit.Command{
	{ID: "conversation.new", Title: "New conversation", Group: "Conversation", Keys: "⌘N", Icon: kit.IconPlus},
	{ID: "view.sidebar", Title: "Toggle sidebar", Group: "View", Keys: "⌘B", Icon: kit.IconPanelLeft},
}
if id := kit.CommandPalette(c, &open, commands); id != "" {
	run(id)
}
```

Selecting a command returns its `ID` and closes the palette on the next frame. All other frames return an empty string. The caller dispatches the action by `ID`. Each command's `ID` must be unique; it also identifies the row.

## Search rules

- Search is case-insensitive and matches only the command's `Title`, not its `Group` or `ID`.
- Groups appear in the order their commands first appear. Commands with an empty `Group` have no group heading.
- Disabled commands are shown but cannot be highlighted or selected.
- When there are no results, “No matches” is shown.

## API

### CommandPalette

`func CommandPalette(c *ui.Context, open *bool, commands []Command) string`

Width 560 (at most 90% of the window). The results list is at most 360 high and scrolls when it overflows.

- **open** `*bool`: Whether the palette is open, owned by the caller. Set to false when a command is selected, Esc is pressed, or the overlay is clicked.
- **commands** `[]Command`: All commands, passed in every frame.
- **Return value** `string`: The command `ID` selected this frame, or an empty string if none was selected.

### Command

- **ID** `string`: Unique command identifier, returned on selection.
- **Title** `string`: Title used for display and search.
- **Group** `string`: Group name.
- **Keys** `string`: Shortcut hint shown at the end of the row, such as `"⌘N"`, drawn with [Kbd](/components/kbd). Display only; it does not register a shortcut.
- **Icon** `*ui.SVG`: Leading icon. If any command has an icon, commands without one reserve the icon space.
- **Disabled** `bool`: Dims the command and makes it unavailable for selection.
