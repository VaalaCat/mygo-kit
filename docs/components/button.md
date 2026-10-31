---
title: Button
description: 按钮家族：文字按钮、图标按钮、工具按钮和链接，默认无边框，悬停才出现底色。
---

# Button

按钮家族。平时没有边框，悬停才出现底色；键盘和焦点行为沿用 MyGo 的按钮。

<Preview page="Button" :height="560">

<<< @/../cmd/gallery/button.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
if kit.PrimaryButton(c, "保存").Clicked() {
	// 保存
}
```

每个按钮都返回 `ui.Element`。点击用 `.Clicked()` 判断，禁用用 `.Disabled(true)`，同一行放多个同类按钮时先用 `c.Key("名字")` 区分。

## 示例

### 文字按钮

`Button` 是默认的 ghost 按钮，`PrimaryButton` 是实色主操作，`DangerButton` 只把文字变成危险色。

```go
ui.Row(c).Gap(8).Children(func() {
	kit.Button(c, "取消")
	kit.PrimaryButton(c, "保存")
	kit.DangerButton(c, "删除")
	kit.Button(c, "不可用").Disabled(true)
})
```

### 图标按钮

`IconButton` 是 28 见方的 ghost 按钮，图标 14。图标和可访问名称分别传入，名称给读屏和测试用，界面上不显示。

```go
kit.IconButton(c.Key("settings"), kit.IconSettings, "设置")
```

要悬停提示就用 `IconButtonTip`，名称同时当提示文字：

```go
if kit.IconButtonTip(c.Key("copy"), kit.IconCopy, "复制").Clicked() {
	// 复制
}
```

消息下方、标签关闭这类不抢眼的位置用 `SmallIconButton`，20 见方，图标 12，平时是弱化色。

### 工具按钮

`ToolButton` 是图标加文字的按钮，`chevron` 为真时右侧带下拉箭头，常配合 [Menu](/components/menu) 使用。

```go
kit.MenuButton(c, kit.ToolButton(c.Key("branch"), kit.IconBranch, "main", true), items)
```

### 贴齐左线的文字按钮

`FlushButton` 没有水平内边距，文字紧贴所在列的左对齐线，悬停只改变文字颜色。适合放在一段说明文字下面，和上面的文字对齐。

```go
ui.Column(c).Gap(0).AlignItems(ui.Start).Children(func() {
	ui.Text(c, "Provider 列表")
	kit.FlushButton(c.Key("refresh"), "刷新")
})
```

### 链接

`Link` 是强调色文字，悬停出现下划线。在纵向容器里，它默认会拉满宽度；只想让文字可点就加 `AlignSelf(ui.Start)`。

```go
kit.Link(c, "MyGo 文档").AlignSelf(ui.Start)
```

## API

### Button

`func Button(c *ui.Context, label string) ui.Element`

28 高的 ghost 按钮。文字从左内边距起，被拉宽时也不居中。

- **label** `string`：按钮文字，同时是可访问名称。

### PrimaryButton

`func PrimaryButton(c *ui.Context, label string) ui.Element`

强调色实底按钮，悬停和按下逐级加深。一个区域里通常只放一个。

### DangerButton

`func DangerButton(c *ui.Context, label string) ui.Element`

文字用危险色，不加外框。用于删除、清空这类操作。

### FlushButton

`func FlushButton(c *ui.Context, label string) ui.Element`

无水平内边距，不画底色。平时是次要色，悬停变正文色。

### IconButton

`func IconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element`

28×28 的 ghost 图标按钮，图标 14（`IconSizeDefault`）。

- **icon** `*ui.SVG`：图标，见[字体与图标](/typography)。
- **label** `string`：可访问名称，不显示在界面上。

### IconButtonTip

`func IconButtonTip(c *ui.Context, icon *ui.SVG, label string) ui.Element`

同 `IconButton`，另外挂上 [Tooltip](/components/tooltip)。`label` 同时是可访问名称和提示文字。

### SmallIconButton

`func SmallIconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element`

20×20 的紧凑 ghost 图标按钮，图标 12。平时用占位色，悬停用次要色。

### ToolButton

`func ToolButton(c *ui.Context, icon *ui.SVG, label string, chevron bool) ui.Element`

28 高，图标加文字，图标和箭头用次要色。

- **icon** `*ui.SVG`：行首图标，传 `nil` 则只显示文字。
- **label** `string`：按钮文字。
- **chevron** `bool`：为真时右侧显示下拉箭头，箭头贴右缘。

### Link

`func Link(c *ui.Context, label string) ui.Element`

强调色链接按钮，文字始终靠左，悬停加下划线。点击后做什么由调用方决定，例如用 `Clicked()` 打开网址。
