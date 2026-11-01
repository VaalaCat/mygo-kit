---
title: Spinner
description: 旋转的加载图标，颜色跟随文字。
---

# Spinner

每秒转一圈的加载图标。系统或主题要求减少动画时保持静止。

<Preview page="Spinner" :height="180">

<<< @/../cmd/gallery/spinner.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	kit.Spinner(c, kit.IconSizeDefault)
	ui.Text(c, "思考中…")
})
```

默认是次要色，想换色就链式调用 `.TextColor(...)`。图标尺寸用 [`IconSizeSmall`、`IconSizeDefault`、`IconSizeLarge`](/zh/theming#图标尺寸)。

## API

### Spinner

`func Spinner(c *ui.Context, size float32) ui.Element`

- **size** `float32`：图标的边长。常用 12、14、16。
