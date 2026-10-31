---
title: Popover
description: 挂在锚点元素旁的浮层，复用 MyGo 的定位、关闭和焦点处理。
---

# Popover

挂在某个元素旁边的小浮层。位置、点外面关闭、Esc 关闭、焦点都由 MyGo 处理，组件只加外观：圆角 8、浮层底色、细边框。

画廊里没有单独的 Popover 页，它和 Tooltip 的示例放在 Dialog 页的最下面。

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // 由应用持有

anchor := kit.Button(c.Key("popover"), "更多信息")
if anchor.Clicked() {
	open = !open
}
kit.Popover(c, anchor, &open, func() {
	ui.Text(c, "浮层里的内容")
})
```

锚点和浮层都要每帧调用。点击锚点怎么开关由你决定，像上面这样取反即可；点浮层外面或按 Esc 时，MyGo 会把 `*open` 置为 false。

## API

### Popover

`func Popover(c *ui.Context, anchor ui.Element, open *bool, content func()) ui.Element`

浮层出现在锚点下方，与锚点相距 4，内边距 8，子元素间距 8。

- **anchor** `ui.Element`：浮层依附的元素。
- **open** `*bool`：是否打开，由调用方持有。
- **content** `func()`：绘制浮层内容的函数。

返回浮层元素，可以链式调用 `.Dismissed()` 判断本帧是否被关闭。
