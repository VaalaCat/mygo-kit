---
title: Dark Mode
description: Use OneDark for dark mode, switch it in your application, or follow the system appearance.
---

# Dark Mode

The library includes two themes: `kit.OneLight()` and `kit.OneDark()`. Your application chooses the appearance by passing a theme to `ApplyTheme`.

<Preview page="Overview" :height="360">

<<< @/../cmd/gallery/overview.go

</Preview>

The docs site’s appearance switch also controls the preview above. Try changing the appearance in the top right.

## Switch in your application

Keep a `dark` state in your application and choose the theme each frame:

```go
var dark bool // Owned by the application; can be saved in configuration

func view(c *ui.Context) {
	theme := kit.OneLight()
	if dark {
		theme = kit.OneDark()
	}
	kit.ApplyTheme(c, theme)

	if kit.Button(c, "Toggle appearance").Clicked() {
		dark = !dark
	}
}
```

## Follow the system

MyGo reports the system appearance each frame. To follow it, before calling `ApplyTheme`, read `c.Theme().Dark`:

```go
func view(c *ui.Context) {
	theme := kit.OneLight()
	if c.Theme().Dark { // Before ApplyTheme, this is still the system theme
		theme = kit.OneDark()
	}
	kit.ApplyTheme(c, theme)
}
```

Order matters. After you call `ApplyTheme`, `c.Theme()` becomes the theme you passed in. Reading it then returns your choice rather than the system appearance.

The gallery’s `cmd/gallery/main.go` uses this approach in embed mode. When embedded in a docs page, the page reports its current light or dark appearance through MyGo’s `data-mygo-theme` attribute.

## Choose system, light, or dark

A common setting is “System / Light / Dark”. Use [Segmented](/components/segmented) to select an option, then apply it:

```go
var appearance int // 0 system, 1 light, 2 dark

dark := c.Theme().Dark // Read the system appearance first
switch appearance {
case 1:
	dark = false
case 2:
	dark = true
}
theme := kit.OneLight()
if dark {
	theme = kit.OneDark()
}
kit.ApplyTheme(c, theme)
```

## Custom colors

Like a light theme, a dark theme is a plain value. You can change its fields before using it. See [Theming](/theming). The `Theme.Dark` field tells MyGo to use dark base styles, so keep it consistent with the colors when building your own theme.
