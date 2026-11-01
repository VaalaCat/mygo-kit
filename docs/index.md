---
title: Introduction
description: mygo-kit is a native Go component library for MyGo, with Zed-style continuous panels, compact controls, and reading components.
---

# Introduction

mygo-kit is a set of native UI components written in Go, built on the [MyGo](https://github.com/egoist/mygo) native UI framework. It gives MyGo a consistent look: continuous panels, compact toolbars, and a clear text hierarchy. MyGo still handles keyboard input, selection, input methods, and focus.

<Preview page="Overview" :height="560">

<<< @/../cmd/gallery/overview.go

</Preview>

These screenshots show native windows, rather than a browser. They follow the site’s light and dark appearance.

<img class="kit-shot kit-shot-light" src="/screenshots/overview-light.png" alt="Gallery overview, light mode">
<img class="kit-shot kit-shot-dark" src="/screenshots/overview-dark.png" alt="Gallery overview, dark mode">
<img class="kit-shot kit-shot-light" src="/screenshots/settings-light.png" alt="Settings page, light mode">
<img class="kit-shot kit-shot-dark" src="/screenshots/settings-dark.png" alt="Settings page, dark mode">

<style>
.kit-shot { display: block; width: 100%; margin: 12px 0; border: 1px solid var(--vp-c-divider); border-radius: 8px; }
.kit-shot-dark { display: none; }
.dark .kit-shot-light { display: none; }
.dark .kit-shot-dark { display: block; }
</style>

## Design approach

The appearance follows Zed:

- **Continuous panels**: sidebars, title bars, and status bars share dividers instead of wrapping each area in its own border.
- **Ghost toolbar buttons**: no border by default. A background appears on hover.
- **Small corner radii**: controls use a radius of 4, overlays and cards use 8, and structural panels keep square corners.
- **4-pixel grid**: spacing uses multiples of 4. Controls have three heights: 20, 28, and 36.
- **One Light / One Dark**: two built-in Zed One color schemes, with IBM Plex Sans and Lilex fonts.

## What it does not do

- It does not create an application runtime. Your application provides windows, the event loop, and `mygo.App.Run`.
- It does not access the network. Your application provides image data, downloads, file reads, and retries. Components only display the results.
- It does not start background goroutines. Animated components use MyGo frame callbacks. The only exception is `TerminalState`; see [Terminal](/components/terminal).
- It does not use a global mutable theme. Your application owns the theme and passes it to `ApplyTheme` each frame.

## Usage

Each component is a function that takes `*ui.Context` and returns a MyGo `ui.Element`. Chain methods to set its size, label, and focus, or compose it inside `ui.Column` and `ui.Row`.

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.ApplyTheme(c, kit.OneLight())
if kit.PrimaryButton(c, "Save").Clicked() {
	// Save
}
```

Your application owns the state for stateful components such as inputs, switches, and lists, and passes a pointer when calling them. Give input components an identity with `c.Key(...)` first so they keep focus and cursor position across redraws.

## Next steps

- [Installation](/installation): requirements, `go get`, and a minimal example.
- [Theming](/theming), [Dark Mode](/dark-mode), [Typography & Icons](/typography), and [Localization](/i18n).
- Start with [Button](/components/button). Each component page includes an interactive preview and full source code.
