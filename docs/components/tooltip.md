---
title: Tooltip
description: 悬停或键盘聚焦时显示的简短提示，每个锚点独立计时。
---

# Tooltip

给一个元素加悬停提示。指针停在锚点上 600ms 后显示；键盘聚焦到锚点（出现焦点环）时立即显示。Esc 行为沿用 MyGo。

画廊里没有单独的 Tooltip 页，示例在 Dialog 页的最下面，也能在 StatusBar 页里悬停左侧按钮看到。

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
button := kit.Button(c.Key("save"), "保存")
kit.Tooltip(c, button, "保存当前文档")
if button.Clicked() {
	// 保存
}
```

提示文字同时被设为锚点的可访问描述。

## 提示

- 图标按钮想要提示，直接用 [IconButtonTip](/components/button)，不用自己配 `Tooltip` 和 `Label`。
- 每个锚点单独计时，指针从一个按钮移到另一个按钮时，不会立即显示后者的提示。
- 提示是单行短句，不要放操作或可点击的内容，这种内容用 [Popover](/components/popover)。

## API

### Tooltip

`func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element`

提示浮层圆角 8，浮层底色，细边框，文字 12/16。

- **anchor** `ui.Element`：要加提示的元素。
- **text** `string`：提示文字。

返回提示元素。锚点自己仍然是你拿着的那个变量，点击、焦点照常处理。
