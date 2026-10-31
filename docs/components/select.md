---
title: Select
description: 下拉选择器，支持搜索、分组标题、选项说明和自定义触发器文字。
---

# Select

下拉选择器。默认是无边框的工具栏样式，`Bordered` 让它和表单输入框长得一样。选项保存的值和显示的文字是分开的。

<Preview page="Select" :height="460">

<<< @/../cmd/gallery/select.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
choice := "One Light" // 由应用持有，保存的是选项的 Value

options := []kit.SelectOption{
	{Value: "One Light", Label: "One Light"},
	{Value: "One Dark", Label: "One Dark"},
}
kit.Select(c.Key("theme"), &choice, kit.SelectOpt{Options: options}).Label("外观")
```

用户选中后，`*value` 直接被改成选项的 `Value`，用返回元素的 `.Changed()` 判断本帧是否变化。

## 示例

### 带搜索

`Searchable` 在弹层里放一个搜索框，输入的词匹配选项的 `Label`、`Value` 和 `Description`。上下键移动高亮，Enter 选中高亮项（没有高亮时选第一个匹配项）并关闭，Esc 只关闭弹层。

```go
kit.Select(c.Key("model"), &model, kit.SelectOpt{Searchable: true, Options: options}).Label("模型")
```

### 表单样式

`Bordered` 让触发器带边框和输入底色，弹层打开时边线变成焦点色。放进 [Field](/components/field) 时用它。

```go
kit.Select(c.Key("protocol"), &protocol, kit.SelectOpt{Bordered: true, Options: options}).Width(240)
```

### 说明和分组标题

`Description` 显示在选项行右侧，触发器上不显示。`Header` 让一项变成不可选的分组标题，这时只用 `Label`，`Value` 被忽略；搜索时，组里没有匹配项，标题也会隐藏。

```go
models := []kit.SelectOption{
	{Label: "Anthropic", Header: true},
	{Value: "claude-sonnet-5-5", Label: "Claude Sonnet 5.5", Description: "Anthropic"},
	{Label: "OpenAI", Header: true},
	{Value: "gpt-5", Label: "GPT-5", Description: "OpenAI"},
	{Value: "gpt-5-mini", Label: "GPT-5 mini", Description: "OpenAI", Disabled: true},
}
```

### 自定义触发器文字

选项只写“自动 / 低 / 中 / 高”，触发器上显示“推理：自动”：

```go
kit.Select(c.Key("reasoning"), &level, kit.SelectOpt{
	Options:      levels,
	TriggerLabel: func(selected kit.SelectOption) string { return "推理：" + selected.Label },
})
```

### 禁用

整组禁用用 **`SelectOpt.Disabled`**，它同时禁用触发器和已展开的弹层。只给返回的元素加 `.Disabled(true)` 不够，因为弹层脱离触发器布局，禁用不会传过去。

```go
kit.Select(c.Key("model"), &model, kit.SelectOpt{Options: options, Disabled: true})
```

## API

### Select

`func Select(c *ui.Context, value *string, opt SelectOpt) ui.Element`

- **value** `*string`：当前选中的 `Value`，由调用方持有。如果它不匹配任何选项且不为空，触发器直接显示这个值；为空时显示 `Placeholder`。
- **opt** `SelectOpt`：见下。

弹层最小宽 160，最高 320，超出滚动。

### SelectOpt

- **Options** `[]SelectOption`：选项列表，按顺序显示。
- **Placeholder** `string`：没有选中任何选项时触发器上的提示，用占位色显示。
- **Label** `string`：触发器的可访问名称。也可以链式 `.Label(...)`。
- **Disabled** `bool`：禁用触发器和弹层。
- **Searchable** `bool`：弹层里显示搜索框。
- **SearchPlaceholder** `string`：搜索框提示。默认用内置的“搜索…”。
- **Bordered** `bool`：触发器带边框，和表单输入一致。默认无边框。
- **TriggerLabel** `func(selected SelectOption) string`：决定触发器文字，参数是当前选中项；选项行仍显示各自的 `Label`。没有选中项时不会调用。

### SelectOption

- **Value** `string`：选中后写回 `*value` 的值。
- **Label** `string`：显示文字。
- **Description** `string`：选项行右侧的 12/16 次要色说明，单行，放不下就截断。
- **Disabled** `bool`：选项变淡，不可选，方向键会跳过。
- **Header** `bool`：这一项是分组标题，不可选，`Value` 被忽略。
