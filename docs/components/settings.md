---
title: Settings
description: Settings pages built from SettingPage, SettingSection, and SettingRow.
---

# Settings

Three components make up a settings page: `SettingPage` is the full page with a fixed header and scrolling content, `SettingSection` is a rounded box with a heading, and `SettingRow` is a row containing a title and a control inside that box.

<Preview page="Settings form" :height="520">

<<< @/../cmd/gallery/form.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.SettingPage(c, "Settings", func() {
	kit.Button(c.Key("revert"), "Revert")
	kit.PrimaryButton(c.Key("apply"), "Apply")
}, func() {
	kit.SettingSection(c.Key("general"), "General", func() {
		kit.SettingRow(c.Key("motion"), "Reduce motion", "Turn off transitions and spinning animations.", func() {
			kit.Switch(c.Key("motion-switch"), &reduceMotion, "").Label("Reduce motion")
		})
		kit.SettingRow(c.Key("theme"), "Color scheme", "", func() {
			kit.Select(c.Key("theme-select"), &theme, kit.SelectOpt{Bordered: true, Options: options}).Label("Color scheme")
		})
	})
})
```

Use `c.Key(...)` to distinguish sections and rows on the same page. Controls on the right have no visible title, so provide an accessible name with `.Label(...)`.

## Examples

### A form inside a section

`SettingSection` does not have to contain `SettingRow` components. It can also contain a column of [Field](/components/field) components with padding you set yourself:

```go
kit.SettingSection(c.Key("provider"), "Model provider", func() {
	ui.Column(c).Padding(12).Gap(16).Children(func() {
		kit.Field(c, kit.FieldOpt{Label: "Base URL", Description: "Include the /v1 suffix."}, func() ui.Element {
			return kit.TextInput(c.Key("url"), &baseURL)
		})
	})
})
```

## API

### SettingPage

`func SettingPage(c *ui.Context, title string, actions, content func()) ui.Element`

The header is a title bar with a height of 36 that does not scroll with the content. Content scrolls below it, with padding of 24 on all sides, a maximum width of 720, left alignment, and a gap of 24 between sections.

- **title** `string`: The page title.
- **actions** `func()`: Actions on the right side of the header. Can be `nil`.
- **content** `func()`: A function that draws the sections.

### SettingSection

`func SettingSection(c *ui.Context, title string, rows func()) ui.Element`

A section heading (16/24) sits above a box with a corner radius of 8. Rows have 1-pixel separators.

- **title** `string`: The section heading.
- **rows** `func()`: A function that draws the `SettingRow` components.

### SettingRow

`func SettingRow(c *ui.Context, title, description string, control func()) ui.Element`

Rows are at least 36 high. The title and optional description sit on the left and wrap when too long. The control sits on the right and keeps its original width.

- **title** `string`: The title.
- **description** `string`: A description in 12/16 text using the secondary color. Can be empty.
- **control** `func()`: A function that draws the control on the right.
