---
title: Dialog
description: 模态对话框：普通的标题加正文，以及带取消和确认按钮的 DialogWith。
---

# Dialog

模态对话框。焦点陷阱、Esc 关闭由 MyGo 处理。`Dialog` 只有标题栏和正文，需要取消和确认按钮时用 `DialogWith`。

<Preview page="Dialog" :height="460">

<<< @/../cmd/gallery/dialog.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // 由应用持有，置 true 就打开

if kit.Button(c, "重命名").Clicked() {
	open = true
}
if kit.DialogWith(c.Key("rename"), &open, kit.DialogOpt{Title: "重命名", ConfirmLabel: "保存"}, func() {
	kit.Field(c, kit.FieldOpt{Label: "标题"}, func() ui.Element {
		return kit.TextInput(c.Key("title"), &title)
	})
}) {
	// 用户点了确认或按了 Enter：保存 title
}
```

对话框每帧都要调用。`open` 为 false 时它什么也不画。`DialogWith` 返回 true 的那一帧，`*open` 已经被置回 false。

## 示例

### 普通对话框

`Dialog` 的正文完全由你画，包括关闭按钮。

```go
kit.Dialog(c, &open, "原生对话框", func() {
	ui.Text(c, "焦点留在对话框里，Esc 关闭。")
	if kit.PrimaryButton(c, "完成").Clicked() {
		open = false
	}
})
```

### 确认和危险操作

没有正文时，`body` 传一个空函数。`Danger` 让确认按钮变成危险色。

```go
if kit.DialogWith(c.Key("delete"), &open, kit.DialogOpt{
	Title:        "删除会话",
	Description:  "删除后无法恢复。",
	ConfirmLabel: "删除",
	Danger:       true,
}, func() {}) {
	// 删除
}
```

### 带输入的对话框

确认条件不满足时，用 `ConfirmDisabled` 禁用确认按钮，Enter 也不会确认。

```go
kit.DialogWith(c.Key("rename"), &open, kit.DialogOpt{Title: "重命名", ConfirmDisabled: title == ""}, func() { /* 输入框 */ })
```

初始焦点和 Enter 的规则：

- 打开时焦点落在正文里第一个可编辑的输入，只读输入除外。页签、按钮不会拿到初始焦点。
- 正文里没有可编辑输入时，焦点落在确认按钮上，并且不画焦点环；Enter 仍然确认。
- 焦点在单行输入里按 Enter 就是确认；在多行输入里是换行；在按钮上是按下这个按钮。单行输入必须用本库的 `TextInput` 或 `TextInputWith` 才会把 Enter 交给对话框。
- 取消按钮和 Esc 只关闭，返回 false。

## API

### Dialog

`func Dialog(c *ui.Context, open *bool, title string, content func()) ui.Element`

宽 480（最多占窗口 90%）的模态对话框，标题栏加正文。

- **open** `*bool`：是否打开，由调用方持有。
- **title** `string`：标题，同时是可访问名称。
- **content** `func()`：绘制正文的函数。

### DialogWith

`func DialogWith(c *ui.Context, open *bool, opt DialogOpt, body func()) (confirmed bool)`

标题栏、说明、正文，加右对齐的“取消”和确认按钮。点击确认或按 Enter 时关闭并返回 true。

- **body** `func()`：绘制正文的函数。

### DialogOpt

- **Title** `string`：标题。
- **Description** `string`：标题栏下方、正文之前的次要色说明。为空不显示。
- **ConfirmLabel** `string`：确认按钮文字。为空用内置的“确定”（英文 Confirm）。
- **Danger** `bool`：确认按钮用 `DangerButton`，否则用 `PrimaryButton`。
- **ConfirmDisabled** `bool`：禁用确认按钮，Enter 也不会确认。
