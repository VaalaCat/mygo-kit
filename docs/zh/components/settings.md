---
title: Settings
description: 设置页：SettingPage 页面，SettingSection 小节，SettingRow 一行设置。
---

# Settings

三个组件拼出一个设置页：`SettingPage` 是整页（固定页头加滚动的内容），`SettingSection` 是带标题的圆角框，`SettingRow` 是框里的一行“标题 加 控件”。

<Preview page="Settings form" :height="520">

<<< @/../cmd/gallery/form.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.SettingPage(c, "设置", func() {
	kit.Button(c.Key("revert"), "还原")
	kit.PrimaryButton(c.Key("apply"), "应用")
}, func() {
	kit.SettingSection(c.Key("general"), "通用", func() {
		kit.SettingRow(c.Key("motion"), "减少动画", "关闭过渡和转圈。", func() {
			kit.Switch(c.Key("motion-switch"), &reduceMotion, "").Label("减少动画")
		})
		kit.SettingRow(c.Key("theme"), "配色", "", func() {
			kit.Select(c.Key("theme-select"), &theme, kit.SelectOpt{Bordered: true, Options: options}).Label("配色")
		})
	})
})
```

同一页里的小节和行都要用 `c.Key(...)` 区分。右侧的控件没有可见标题，请用 `.Label(...)` 给出可访问名称。

## 示例

### 小节里放表单

`SettingSection` 的内容不一定是 `SettingRow`，也可以放一列 [Field](/zh/components/field)，自己加内边距：

```go
kit.SettingSection(c.Key("provider"), "模型服务", func() {
	ui.Column(c).Padding(12).Gap(16).Children(func() {
		kit.Field(c, kit.FieldOpt{Label: "Base URL", Description: "包含 /v1 后缀。"}, func() ui.Element {
			return kit.TextInput(c.Key("url"), &baseURL)
		})
	})
})
```

## API

### SettingPage

`func SettingPage(c *ui.Context, title string, actions, content func()) ui.Element`

页头是 36 高的标题栏，不随内容滚动。内容区在页头下方滚动，左右上下边距 24，最宽 720，靠左，小节之间间距 24。

- **title** `string`：页面标题。
- **actions** `func()`：页头右侧的操作，可以为 `nil`。
- **content** `func()`：绘制各个小节的函数。

### SettingSection

`func SettingSection(c *ui.Context, title string, rows func()) ui.Element`

上方是小节标题（16/24），下方是圆角 8 的框，行之间有 1 像素分隔线。

- **title** `string`：小节标题。
- **rows** `func()`：绘制若干 `SettingRow` 的函数。

### SettingRow

`func SettingRow(c *ui.Context, title, description string, control func()) ui.Element`

行高至少 36。左侧是标题和可选说明，文字太长会换行；右侧是控件，保持原宽。

- **title** `string`：标题。
- **description** `string`：12/16 的次要色说明，可以为空。
- **control** `func()`：绘制右侧控件的函数。
