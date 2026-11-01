---
title: Select
description: A dropdown selector with search, group headings, option descriptions, and custom trigger text.
---

# Select

A dropdown selector. The default style is a borderless toolbar control; `Bordered` makes it look like a form input. An option's stored value is separate from its display text.

<Preview page="Select" :height="460">

<<< @/../cmd/gallery/select.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
choice := "One Light" // Held by the application; stores the option's Value

options := []kit.SelectOption{
	{Value: "One Light", Label: "One Light"},
	{Value: "One Dark", Label: "One Dark"},
}
kit.Select(c.Key("theme"), &choice, kit.SelectOpt{Options: options}).Label("Appearance")
```

When the user selects an option, `*value` is set directly to its `Value`. Use `.Changed()` on the returned element to check whether it changed this frame.

## Examples

### Search

`Searchable` adds a search field to the popup. The search text matches each option's `Label`, `Value`, and `Description`. The up and down arrow keys move the highlight. Enter selects the highlighted option (or the first match if none is highlighted) and closes the popup. Esc only closes the popup.

```go
kit.Select(c.Key("model"), &model, kit.SelectOpt{Searchable: true, Options: options}).Label("Model")
```

### Form style

`Bordered` gives the trigger a border and an input background. The border uses the focus color while the popup is open. Use this style inside a [Field](/components/field).

```go
kit.Select(c.Key("protocol"), &protocol, kit.SelectOpt{Bordered: true, Options: options}).Width(240)
```

### Descriptions and group headings

`Description` appears on the right side of an option row, but not on the trigger. `Header` turns an option into a group heading that cannot be selected. Only `Label` is used in this case; `Value` is ignored. During search, a group heading is hidden when none of its options match.

```go
models := []kit.SelectOption{
	{Label: "Anthropic", Header: true},
	{Value: "claude-sonnet-5-5", Label: "Claude Sonnet 5.5", Description: "Anthropic"},
	{Label: "OpenAI", Header: true},
	{Value: "gpt-5", Label: "GPT-5", Description: "OpenAI"},
	{Value: "gpt-5-mini", Label: "GPT-5 mini", Description: "OpenAI", Disabled: true},
}
```

### Custom trigger text

The options can read “Auto / Low / Medium / High” while the trigger shows “Reasoning: Auto”:

```go
kit.Select(c.Key("reasoning"), &level, kit.SelectOpt{
	Options:      levels,
	TriggerLabel: func(selected kit.SelectOption) string { return "Reasoning: " + selected.Label },
})
```

### Disabled

Use **`SelectOpt.Disabled`** to disable the whole control. It disables both the trigger and an open popup. Adding only `.Disabled(true)` to the returned element is not enough: the popup is outside the trigger's layout, so it does not inherit the disabled state.

```go
kit.Select(c.Key("model"), &model, kit.SelectOpt{Options: options, Disabled: true})
```

## API

### Select

`func Select(c *ui.Context, value *string, opt SelectOpt) ui.Element`

- **value** `*string`: The currently selected `Value`, held by the caller. If it is nonempty and does not match any option, the trigger shows the value directly. If it is empty, the trigger shows `Placeholder`.
- **opt** `SelectOpt`: See below.

The popup has a minimum width of 160 and a maximum height of 320. Overflow scrolls.

### SelectOpt

- **Options** `[]SelectOption`: The options, displayed in order.
- **Placeholder** `string`: Text shown on the trigger when no option is selected, using the placeholder color.
- **Label** `string`: The trigger's accessible name. You can also chain `.Label(...)`.
- **Disabled** `bool`: Disables the trigger and popup.
- **Searchable** `bool`: Shows a search field in the popup.
- **SearchPlaceholder** `string`: The search field's placeholder. Defaults to the built-in “Search…”.
- **Bordered** `bool`: Gives the trigger a border to match form inputs. Borderless by default.
- **TriggerLabel** `func(selected SelectOption) string`: Sets the trigger text, receiving the currently selected option. Option rows still show their own `Label`. Not called when no option is selected.

### SelectOption

- **Value** `string`: The value written to `*value` when selected.
- **Label** `string`: The display text.
- **Description** `string`: A description on the right side of the option row, in 12/16 text using the secondary color. It stays on one line and is truncated if it does not fit.
- **Disabled** `bool`: Dims the option and prevents selection. Arrow key navigation skips it.
- **Header** `bool`: Marks the option as a group heading. It cannot be selected, and `Value` is ignored.
