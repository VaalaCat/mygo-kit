---
title: Disclosure
description: 无边框的整行折叠标题，展开时内容带高度过渡。
---

# Disclosure

一整行可点击的标题，左边是箭头，点击展开或收起下面的内容。键盘和展开语义沿用 MyGo 的折叠组件。

<Preview page="Disclosure" :height="220">

<<< @/../cmd/gallery/disclosure.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
open := true // 展开状态由应用持有

kit.Disclosure(c, "更多设置", &open, func() {
	ui.Text(c, "展开后显示的内容")
})
```

内容区缩进到标题文字的位置。展开用 200ms、收起用 150ms 的高度过渡；系统或主题要求减少动画时没有过渡。收起中的内容立即禁用，不会再响应点击。

## API

### Disclosure

`func Disclosure(c *ui.Context, label string, open *bool, content func()) ui.Element`

返回标题行元素，可以继续链式设置，例如 `.Disabled(true)`。

- **label** `string`：标题文字，同时是可访问名称。
- **open** `*bool`：是否展开，由调用方持有。
- **content** `func()`：展开后绘制内容的函数。完全收起后不会调用。
