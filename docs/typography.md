---
title: Typography & Icons
description: The built-in IBM Plex Sans and Lilex fonts, RegisterFonts, all icon constants, and Kbd.
---

# Typography & Icons

## Fonts

The library includes two font families, matching Zed’s default UI and code fonts:

- **IBM Plex Sans** (constant `kit.FontSans`): UI and body text. Regular, semibold, italic, and semibold italic.
- **Lilex** (constant `kit.FontMono`): code, terminals, JSON, and patches. Regular, bold, and italic.

Font files are embedded in the library. Call `RegisterFonts` once at application startup, before creating windows:

```go
if err := kit.RegisterFonts(); err != nil {
	log.Fatal(err)
}
```

Repeated calls still register fonts only once and have no side effects.

The theme uses these two constants as font names: `Theme.Typography.Font` is `FontSans`, and `Theme.Typography.MonoFont` is `FontMono`. To use your own font, register it with MyGo’s `ui.RegisterFont`, then update the theme:

```go
theme := kit.OneLight()
theme.Typography.Font = "My Sans"
kit.ApplyTheme(c, theme)
```

To use a monospace font for your own elements:

```go
t := kit.GetTheme(c)
t.TextRole(ui.Text(c, "go test ./..."), kit.TextCode)
```

For font sizes and line heights, see [Theming](/theming#text-roles).

### In a browser

Browser builds (`GOOS=js GOARCH=wasm`) register only the regular face of each font family. Browsers distinguish fonts only by family name, so bold and italic faces from the same family would overwrite the regular face. The browser synthesizes bold text instead. Desktop builds are unaffected.

## Icons

Icons are Lucide outline icons with a consistent 1.5 stroke width. Their color follows the text color. Each icon is a `*ui.SVG` variable. Pass it directly to a component or draw it with `ui.Icon`:

```go
kit.IconButton(c.Key("settings"), kit.IconSettings, "Settings")

ui.Icon(c, kit.IconSearch).Size(kit.IconSizeDefault, kit.IconSizeDefault).TextColor(t.Colors.Muted)
```

For the size constants `IconSizeSmall` (12), `IconSizeDefault` (14), and `IconSizeLarge` (16), see [Theming](/theming#icon-sizes).

To use your own SVG, parse it with `ui.ParseSVG` or `ui.MustParseSVG`. Use the resulting `*ui.SVG` just like a built-in icon. Set stroke colors to `currentColor` so they follow the text color.

### All icons

Panels and layout:

- `IconPanelLeft`, `IconPanelRight`, `IconPanelBottom`
- `IconPanelBottomClose`, `IconPanelRightClose`
- `IconFold`, `IconFiles`, `IconListTree`, `IconListEnd`

Arrows and directions:

- `IconChevronDown`, `IconChevronLeft`, `IconChevronRight`
- `IconArrowUp`, `IconArrowDown`

Actions:

- `IconPlus`, `IconX`, `IconCheck`, `IconSearch`, `IconSettings`
- `IconCopy`, `IconPencil`, `IconTrash`, `IconDownload`, `IconReplace`
- `IconRefresh`, `IconRotateCcw`
- `IconPlay`, `IconStop`, `IconCirclePause`, `IconCircleStop`
- `IconEllipsis`, `IconImagePlus`

Status and hints:

- `IconInfo`, `IconAlert`, `IconTriangleAlert`
- `IconCircleCheck`, `IconCircleX`, `IconBan`
- `IconLoader`, `IconGauge`, `IconActivity`

Objects:

- `IconFolder`, `IconFileText`, `IconArchive`, `IconBookOpen`
- `IconTerminal`, `IconCommand`, `IconBraces`
- `IconBot`, `IconBrain`, `IconWrench`, `IconBug`
- `IconBranch`, `IconGitFork`, `IconWebhook`, `IconServerCog`
- `IconMessageSquare`

## Kbd

`Kbd` is a small label for keyboard shortcuts. See [Kbd](/components/kbd) for usage. It uses the `TextCaption` font size. When placed beside text, use `AlignItems(ui.Center)` to center both vertically.
