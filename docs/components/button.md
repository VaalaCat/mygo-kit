---
title: Button
description: A family of text buttons, icon buttons, tool buttons, and links, with no border by default and a background on hover.
---

# Button

A family of buttons. Normally borderless, with a background on hover. Keyboard and focus behavior come from MyGo's buttons.

<Preview page="Button" :height="560">

<<< @/../cmd/gallery/button.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
if kit.PrimaryButton(c, "Save").Clicked() {
	// Save
}
```

Every button returns `ui.Element`. Check clicks with `.Clicked()` and disable with `.Disabled(true)`. When placing multiple buttons of the same kind in a row, distinguish them with `c.Key("name")` first.

## Examples

### Text buttons

`Button` is the default ghost button, `PrimaryButton` is a solid primary action, and `DangerButton` changes only the text to the danger color.

```go
ui.Row(c).Gap(8).Children(func() {
	kit.Button(c, "Cancel")
	kit.PrimaryButton(c, "Save")
	kit.DangerButton(c, "Delete")
	kit.Button(c, "Unavailable").Disabled(true)
})
```

### Icon buttons

`IconButton` is a 28×28 ghost button with a 14-sized icon. Pass the icon and accessible name separately. The name is for screen readers and tests and is not shown in the interface.

```go
kit.IconButton(c.Key("settings"), kit.IconSettings, "Settings")
```

For a hover tooltip, use `IconButtonTip`. The name also serves as the tooltip text:

```go
if kit.IconButtonTip(c.Key("copy"), kit.IconCopy, "Copy").Clicked() {
	// Copy
}
```

Use `SmallIconButton` in less prominent places, such as below messages or for closing tabs. It is 20×20 with a 12-sized icon and normally uses a muted color.

### Tool buttons

`ToolButton` combines an icon and text. When `chevron` is true, a dropdown arrow appears on the right. It is often used with [Menu](/components/menu).

```go
kit.MenuButton(c, kit.ToolButton(c.Key("branch"), kit.IconBranch, "main", true), items)
```

### Text buttons aligned to the left edge

`FlushButton` has no horizontal padding, so its text aligns with the left edge of its column. Only the text color changes on hover. It works well below explanatory text, aligned with the text above.

```go
ui.Column(c).Gap(0).AlignItems(ui.Start).Children(func() {
	ui.Text(c, "Provider list")
	kit.FlushButton(c.Key("refresh"), "Refresh")
})
```

### Links

`Link` uses accent-colored text and adds an underline on hover. In a vertical container, it stretches to the full width by default. Add `AlignSelf(ui.Start)` to make only the text clickable.

```go
kit.Link(c, "MyGo documentation").AlignSelf(ui.Start)
```

## API

### Button

`func Button(c *ui.Context, label string) ui.Element`

A 28-high ghost button. The text starts at the left padding and stays left-aligned when the button is stretched.

- **label** `string`: Button text, also used as the accessible name.

### PrimaryButton

`func PrimaryButton(c *ui.Context, label string) ui.Element`

A solid accent-colored button that gets progressively darker on hover and press. Usually only one is placed in an area.

### DangerButton

`func DangerButton(c *ui.Context, label string) ui.Element`

Uses danger-colored text with no outline. For actions such as deleting or clearing.

### FlushButton

`func FlushButton(c *ui.Context, label string) ui.Element`

No horizontal padding or background. Uses the secondary color normally and the body text color on hover.

### IconButton

`func IconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element`

A 28×28 ghost icon button with a 14-sized icon (`IconSizeDefault`).

- **icon** `*ui.SVG`: Icon; see [Typography and icons](/typography).
- **label** `string`: Accessible name, not shown in the interface.

### IconButtonTip

`func IconButtonTip(c *ui.Context, icon *ui.SVG, label string) ui.Element`

The same as `IconButton`, with a [Tooltip](/components/tooltip) attached. `label` is both the accessible name and the tooltip text.

### SmallIconButton

`func SmallIconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element`

A compact 20×20 ghost icon button with a 12-sized icon. Uses the placeholder color normally and the secondary color on hover.

### ToolButton

`func ToolButton(c *ui.Context, icon *ui.SVG, label string, chevron bool) ui.Element`

Height 28, with an icon and text. The icon and arrow use the secondary color.

- **icon** `*ui.SVG`: Leading icon. Pass `nil` to show text only.
- **label** `string`: Button text.
- **chevron** `bool`: Shows a dropdown arrow at the right edge when true.

### Link

`func Link(c *ui.Context, label string) ui.Element`

An accent-colored link button with text that always stays left-aligned and an underline on hover. The caller decides what a click does, such as opening a URL with `Clicked()`.
