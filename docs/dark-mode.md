---
title: 暗色模式
description: 用 OneDark 做暗色，由应用自己切换，或跟随系统。
---

# 暗色模式

本库内置两份主题：`kit.OneLight()` 和 `kit.OneDark()`。明暗由应用决定，把哪一份交给 `ApplyTheme` 就是哪一份。

<Preview page="Overview" :height="360">

<<< @/../cmd/gallery/overview.go

</Preview>

文档站的明暗开关和上面的预览是联动的，切换右上角的外观试试。

## 应用自己切换

应用持有一个 `dark` 状态，每帧选主题：

```go
var dark bool // 由应用持有，可以存进配置

func view(c *ui.Context) {
	theme := kit.OneLight()
	if dark {
		theme = kit.OneDark()
	}
	kit.ApplyTheme(c, theme)

	if kit.Button(c, "切换明暗").Clicked() {
		dark = !dark
	}
}
```

## 跟随系统

MyGo 会把系统的明暗报告给每一帧。要跟随它，在 `ApplyTheme` 之前读 `c.Theme().Dark`：

```go
func view(c *ui.Context) {
	theme := kit.OneLight()
	if c.Theme().Dark { // 在 ApplyTheme 之前，这里还是系统给的主题
		theme = kit.OneDark()
	}
	kit.ApplyTheme(c, theme)
}
```

要注意顺序。一旦调用了 `ApplyTheme`，`c.Theme()` 就变成你交给它的主题，之后读到的是你自己的选择，不是系统的。

画廊的 `cmd/gallery/main.go` 在嵌入模式下就是这样做的：嵌入文档页时，页面通过 MyGo 的 `data-mygo-theme` 属性告诉它当前是明是暗。

## 三选一：系统、明、暗

常见的设置项是“系统 / 明 / 暗”，用 [Segmented](/components/segmented) 选，再按选择决定：

```go
var appearance int // 0 系统，1 明，2 暗

dark := c.Theme().Dark // 先读系统
switch appearance {
case 1:
	dark = false
case 2:
	dark = true
}
theme := kit.OneLight()
if dark {
	theme = kit.OneDark()
}
kit.ApplyTheme(c, theme)
```

## 自定义配色

暗色主题和明色一样是普通值，可以改字段后再用，见[主题](/theming)。`Theme.Dark` 字段会告诉 MyGo 用深色的基础样式，所以自己拼主题时要保持它和配色一致。
