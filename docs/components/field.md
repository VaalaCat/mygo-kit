---
title: Field
description: 竖排表单字段 Field：标签在上、控件在中、说明或错误在下；另有一行提示 FieldMessage 和列表里的 InlineHint。
---

# Field

表单字段。`Field` 把标签、控件和说明排成一列，并且把标签和说明接到控件的可访问名称和描述上。`FieldMessage` 是单独的一行提示文字，`InlineHint` 是列表里“加载中”“空列表”这类提示。

<Preview page="Field" :height="520">

<<< @/../cmd/gallery/form.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Field(c, kit.FieldOpt{
	Label:       "Provider ID",
	Description: "用在 config.yaml 里，也是模型名的前缀。",
}, func() ui.Element {
	return kit.TextInputWith(c.Key("id"), &providerID, kit.InputOpt{Placeholder: "openai"})
})
```

`control` 要返回控件的元素，`Field` 要在它上面设置标签和描述。点标签会聚焦这个控件。

调用方负责：

- 校验输入，并把错误文字放进 `Error`。
- 输入类控件的危险色边线，由你用 `InputOpt.Invalid` 设置，`Field` 只显示错误文字。

## 示例

### 错误

`Error` 非空时，控件被标为无效，下面显示危险色的错误。说明和错误同时存在时只显示错误。

```go
kit.Field(c, kit.FieldOpt{Label: "名称", Error: "只能用字母、数字和连字符。"}, func() ui.Element {
	return kit.TextInputWith(c.Key("name"), &name, kit.InputOpt{Invalid: true})
})
```

### 不是输入的控件

任何控件都可以放进 `Field`，例如 [Select](/components/select) 和 [Switch](/components/switch)：

```go
kit.Field(c, kit.FieldOpt{Label: "继承设置"}, func() ui.Element {
	return kit.Switch(c.Key("inherit"), &inherit, "")
})
```

### 单独的提示行

```go
kit.FieldMessage(c, kit.MessageWarning, "这个值会覆盖默认配置。")
```

### 列表里的提示

加载中时前面带转圈，列表为空时只有文字。

```go
kit.InlineHint(c, "加载中…", true)
kit.InlineHint(c, "还没有会话", false)
```

## API

### Field

`func Field(c *ui.Context, opt FieldOpt, control func() ui.Element)`

竖排，各段间距 4。没有返回值。

- **opt** `FieldOpt`：标签、说明和错误。
- **control** `func() ui.Element`：画控件并返回它的元素。

### FieldOpt

- **Label** `string`：标签，次要色。为空时不显示，也不设置控件的名称。
- **Description** `string`：说明。
- **Error** `string`：错误。说明和错误都为空时，控件下方不留空。

### FieldMessage

`func FieldMessage(c *ui.Context, kind MessageKind, text string) ui.Element`

一行 12/16 的文字，可以换行。

### MessageKind

- **MessageDescription**：说明，次要色。
- **MessageWarning**：警告色。
- **MessageError**：危险色。

### InlineHint

`func InlineHint(c *ui.Context, message string, loading bool) ui.Element`

一行 28 高的次要提示。`loading` 为真时行首显示转圈。`message` 为空且 `loading` 为真时，显示内置的“加载中…”。

### InlineHintWith

`func InlineHintWith(c *ui.Context, message string, loading bool, opt InlineHintOpt) ui.Element`

同 `InlineHint`，由 `opt` 指定左对齐线。

### InlineHintOpt

- **Inset** `float32`：文字距左缘的距离。列表里用 12（零值），对话列用 16。
