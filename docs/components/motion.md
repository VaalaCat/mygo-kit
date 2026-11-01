---
title: Motion
description: Shared reduced-motion rules and the Fade transition.
---

# Motion

All animations in this library are brief and share the same controls. When `Theme.Motion.Reduced` is true or the system requests reduced motion (`c.Preferences().ReduceMotion`), every animation moves directly to its end state.

Animation durations are stored in the theme, with these defaults:

- `Motion.Fade`: 120ms for fading, switch thumbs, and MorphText.
- `Motion.DisclosureOpen`: 200ms for expanding Disclosure and ActivityCard.
- `Motion.DisclosureClose`: 150ms for collapsing.

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

### Fading

`Fade` transitions an element's opacity to 1 or 0 according to `visible`.

```go
panel := ui.Column(c).Children(func() { ui.Text(c, "Hint") })
kit.Fade(c, panel, showing)
```

When `visible` is false, the element is immediately disabled and stops responding to clicks and keyboard input. Opacity affects drawing only. To remove it from the layout, the caller must stop creating it.

### Disabling library animations

```go
theme := kit.OneLight()
theme.Motion.Reduced = true
kit.ApplyTheme(c, theme)
```

The theme setting does not automatically follow the system, but components also check system preferences. They therefore remain still when the system requests reduced motion, even if `Reduced` is false.

## Related components

- [Spinner](/components/spinner): A loading icon that rotates once per second.
- [Shimmer](/components/shimmer): Text with a moving highlight.
- [MorphText](/components/morph-text): Grapheme-level text transitions.
- [Disclosure](/components/disclosure): Height transitions when expanding and collapsing.

## API

### Fade

`func Fade(c *ui.Context, element ui.Element, visible bool) ui.Element`

Updates the element's opacity and returns the same element.

- **element** `ui.Element`: The element to fade in or out.
- **visible** `bool`: The target state. True transitions to fully visible; false transitions to fully transparent and disables the element.

### Motion

`type Motion struct`, stored in `Theme.Motion`.

- **Fade** `time.Duration`: Fade duration. Defaults to 120ms.
- **DisclosureOpen** `time.Duration`: Expansion duration. Defaults to 200ms.
- **DisclosureClose** `time.Duration`: Collapse duration. Defaults to 150ms.
- **Reduced** `bool`: Disables library animations when true. Defaults to `false`.
