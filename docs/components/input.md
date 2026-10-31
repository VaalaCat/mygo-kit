---
title: Input
description: 单行输入和多行输入，支持只读值、无边框、前缀图标和错误边线。
---

# Input

单行输入 `TextInput` 和多行输入 `TextArea`。内容、撤销、输入法组字和选区都由 MyGo 管理，组件只负责外观。

<Preview page="Input" :height="520">

<<< @/../cmd/gallery/input.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var name string // 由应用持有

kit.TextInput(c.Key("name"), &name).Label("名称").Placeholder("项目名称")
```

输入控件必须先用 `c.Key(...)` 给身份，否则界面重建后焦点和光标会丢。内容变化用返回元素的 `.Changed()` 判断，单行输入按下 Enter 用 `.Submitted()` 判断。

## 示例

### 多行输入

`TextArea` 最小两行高（56），可以用 `.Height(100)` 或 `.MinHeight(...)` 调整。

```go
kit.TextArea(c.Key("desc"), &desc).Label("描述").Placeholder("支持多行和输入法").Height(100)
```

### 错误状态

`Invalid` 让边线变成危险色，优先于焦点色。错误文字用 [Field](/components/field) 显示。

```go
kit.TextInputWith(c.Key("name"), &name, kit.InputOpt{Invalid: name == ""}).Label("名称")
```

### 只读值

`ReadOnly` 把输入显示成一段可以选中复制的文字：没有边线和底色，但文字起点和 28 高与普通输入一致，放进 `Field` 时和可编辑字段对齐。

```go
kit.Field(c, kit.FieldOpt{Label: "Provider ID"}, func() ui.Element {
	return kit.TextInputWith(c.Key("id"), &providerID, kit.InputOpt{ReadOnly: true})
})
```

需要“看起来可以编辑”时不要用它。

### 前缀图标

`Leading` 在输入框左侧画一个 14 的次要色图标，长文字会被裁在图标右边，不会盖住图标。

```go
kit.TextInputWith(c.Key("search"), &query, kit.InputOpt{Leading: kit.IconSearch, Placeholder: "搜索"}).Label("搜索")
```

### 无边框

`Borderless` 去掉边框、底色、内边距和焦点环，外观由外面的容器决定，用于侧栏搜索、消息输入区。

```go
box := ui.Column(c).Padding(kit.GetTheme(c).Space(3)).Radius(8).Border(1, kit.GetTheme(c).Colors.Border)
box.Children(func() {
	kit.TextAreaWith(c.Key("msg"), &msg, kit.InputOpt{Borderless: true, Placeholder: "消息"}).Label("消息")
})
```

### 在对话框里

在 [DialogWith](/components/dialog) 的正文里，只有本库的 `TextInput` / `TextInputWith` 会把 Enter 报告给对话框，Enter 即确认。对话框打开时，焦点落在正文里第一个可编辑的输入。

## API

### TextInput

`func TextInput(c *ui.Context, value *string) ui.Element`

单行输入，等于零值 `InputOpt` 的 `TextInputWith`。

- **value** `*string`：输入内容，由调用方持有，输入时直接写回。

### TextInputWith

`func TextInputWith(c *ui.Context, value *string, opt InputOpt) ui.Element`

28 高的单行输入，文字在高度内垂直居中。

### TextArea

`func TextArea(c *ui.Context, value *string) ui.Element`

多行输入，等于零值 `InputOpt` 的 `TextAreaWith`。

### TextAreaWith

`func TextAreaWith(c *ui.Context, value *string, opt InputOpt) ui.Element`

最小 56 高，保留 MyGo 的组字、选区和受控值处理。

### InputOpt

`type InputOpt struct`

零值就是普通的可编辑输入。

- **Placeholder** `string`：没有内容时显示的提示。
- **Invalid** `bool`：边线显示危险色，优先于焦点色。
- **ReadOnly** `bool`：显示成可选中复制的纯值，没有边线和底色。
- **Borderless** `bool`：去掉边框、底色、内边距和焦点环。
- **Leading** `*ui.SVG`：左侧前缀图标，14，次要色。默认 `nil`，不画图标。
