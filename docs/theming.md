---
title: Theming
description: The Theme value type, One Light / One Dark, color roles, metrics, text roles, icon sizes, and spacing steps.
---

# Theming

A theme is a plain Go value, `kit.Theme`. Your application owns it and passes it to `ApplyTheme` each frame. The library has no global mutable theme. Each window keeps its own theme, so multiple windows can use different light and dark themes.

## Usage

```go
func view(c *ui.Context) {
	theme := kit.OneLight()
	theme.Metrics.Radius = 0 // Optional: use square corners for controls
	kit.ApplyTheme(c, theme)

	t := kit.GetTheme(c) // Retrieve the current window’s theme anywhere after this
	ui.Text(c, "Hello").TextColor(t.Colors.Muted)
}
```

- `kit.OneLight()` and `kit.OneDark()` each return a complete theme based on Zed’s One color schemes.
- Call `kit.ApplyTheme(c, t)` **every frame**, before drawing any components. It stores the theme on the window’s root node for the frame and sets MyGo’s text, focus, selection, and scrollbar colors.
- `kit.GetTheme(c)` returns the current window’s theme. It returns One Light if `ApplyTheme` has not been called.

`Theme` is a value type. Changing a copy does not affect other copies. To create your own theme, start with `OneLight()` or `OneDark()` and change the fields you need.

```go
theme := kit.OneDark()
theme.Colors.Accent = ui.Hex("#e5836b")
```

## Colors

The fields in `Theme.Colors` describe UI roles rather than hues. Components read only these fields. Use them when drawing your own UI too.

Background layers, from outside to inside:

- **Background**: the outermost layer. The background for title bars, status bars, and Kbd.
- **Panel**: the background for sidebars, Header, tab bars, and drawers.
- **Editor**: the background for body content and editing areas. Terminal, Diff, and content boxes also use it.
- **Elevated**: overlays such as menus, dialogs, Toast, Tooltip, and Popover.

Text:

- **Text**: body text.
- **Muted**: secondary text and icons.
- **Placeholder**: input hints, placeholders, and descriptions.
- **Disabled**: reserved. No components currently read it. The framework makes disabled controls translucent.

Lines and states:

- **Border**: borders for controls and panels.
- **BorderSubtle**: lighter dividers, such as Separator and lines between rows.
- **Focus**: focus rings and focused input borders.
- **Hover**: hover backgrounds and progress bar tracks.
- **Active**: pressed backgrounds.
- **Selected**: selected row backgrounds.

Accent and semantic colors:

- **Accent**: primary buttons, links, and switches that are on.
- **AccentText**: text on accent backgrounds.
- **AccentSubtle**: a translucent accent background for highlighted states such as selected segments.
- **Danger**, **Warning**, **Success**: danger, warning, and success.

Other colors:

- **Input**: backgrounds for inputs and bordered Select controls.
- **Selection**: text selection backgrounds, also used for Terminal selections.
- **Added**, **Deleted**: light backgrounds for added and deleted patch lines.
- **ANSI**: `[16]ui.Color`, the 16-color palette for Terminal and ANSIOutput.

Code syntax colors are stored separately in `Theme.Syntax`, of type `SyntaxColors`, with `Keyword`, `String`, `Function`, `Type`, `Number`, `Comment`, `Property`, and `Operator` fields. Markdown code blocks and Diff use these colors. JSON tree keys use `Property`.

## Metrics

`Theme.Metrics` holds the 4-pixel grid, three control heights, and two corner radii:

- **Grid**: 4, one grid unit.
- **CompactHeight**: 20, for compact controls such as Badge, Kbd, SmallIconButton, and StatusButton.
- **ControlHeight**: 28, for standard controls and list rows.
- **HeaderHeight**: 36, for title bars and tabs.
- **StatusHeight**: 28, for status bars.
- **Radius**: 4, for control corners.
- **RadiusLarge**: 8, for overlay and card corners.

Setting `Radius` to 0 is valid. It gives square corners and does not fall back to the default.

`kit.RadiusFull` is the corner radius for round controls: 9999. It exceeds half the size of any control, producing a circle. Switches and status dots use it.

## Text roles

Use `Theme.TextRole(element, role)` to apply a text role instead of setting font size and line height yourself. It sets the font, size, fixed line height, and weight. You still choose the color and alignment.

```go
t.TextRole(ui.Text(c, "Context"), kit.TextCaption).TextColor(t.Colors.Muted)
```

- **TextProse**: 14/24, for Markdown body text and detailed descriptions.
- **TextUI**: 14/20, for row titles, buttons, inputs, and tool names.
- **TextCaption**: 12/16, for timestamps, hints, and secondary status text.
- **TextToken**: 12/20, for token counts, with monospaced digits.
- **TextCode**: 12/20, for code, JSON, and tool output, in a monospace font.
- **TextCodeLabel**: 12/16, for section labels and content box titles, in a monospace font. The caller applies a secondary color.
- **TextTitle**: 16/24, for section and page titles, with a weight of 400.

Fonts and sizes come from `Theme.Typography`: `Font`, `MonoFont`, `BodySize` (14), `CaptionSize` (12), `HeadingSize` (16), and `CodeSize` (12).

## Icon sizes

Icons have three sizes:

- **IconSizeSmall**: 12, for compact controls with a height of 20 and trailing indicators such as dropdown arrows.
- **IconSizeDefault**: 14, the default for buttons with a height of 28, leading row icons, and status bars.
- **IconSizeLarge**: 16, for empty states and primary actions in title bars with a height of 36.

```go
ui.Icon(c, kit.IconSettings).Size(kit.IconSizeDefault, kit.IconSizeDefault)
```

## Spacing

`Theme.Space(n)` returns `n` grid units, or `4*n`. Use only these steps:

- `Space(1)` = 4
- `Space(2)` = 8
- `Space(3)` = 12
- `Space(4)` = 16
- `Space(6)` = 24
- `Space(8)` = 32

`Space(5)` = 20 is used only for row heights, not spacing.

```go
ui.Column(c).Padding(t.Space(6)).Gap(t.Space(3))
```

## Motion

`Theme.Motion` holds animation durations and the reduced motion setting. See [Motion](/components/motion).

## Dark mode

For light and dark themes, see [Dark Mode](/dark-mode).
