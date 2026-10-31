---
title: Image
description: 显示调用方准备好的图片，带加载中和出错状态；ImagePreview 可适应窗口、缩放和滚动。
---

# Image

显示图片。本库不读磁盘，也不下载：图片数据、加载状态、错误和重试都由你的应用提供，用 `ImageResource` 交给组件。

<Preview page="Images" :height="420">

<<< @/../cmd/gallery/image.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Image(c, kit.ImageResource{Source: source, Alt: "架构图"})
```

`Source` 是 MyGo 的 `ui.ImageSource`，例如 `ui.NewBitmap(img)` 从 Go 的 `image.Image` 创建。图片按原始宽高比显示，最宽不超过容器。

## 示例

### 加载中和出错

`Source` 为 `nil` 时，组件显示一块说明，按顺序取第一个有的：加载中、`Error` 的文字、`Alt`，都没有时显示“图片不可用”。有 `Error` 并且设了 `Retry` 时，说明下面会有“重试图片”按钮。

```go
kit.Image(c, kit.ImageResource{Alt: "架构图", Loading: true})

kit.Image(c, kit.ImageResource{
	Alt:   "架构图",
	Error: "下载失败：404",
	Retry: func() { download(url) },
})
```

### 预览

`ImagePreview` 是一个带工具栏的预览区：“适应窗口”“原始大小”“−”“+”，图片可以双向滚动。缩放状态由你持有。

```go
var preview = kit.ImagePreviewState{Fit: true} // 由应用持有

kit.ImagePreview(c, resource, &preview).Grow(1)
```

## API

### Image

`func Image(c *ui.Context, image ImageResource) ui.Element`

圆角 4 的图片，`Alt` 是可访问名称。

### ImageResource

- **Source** `ui.ImageSource`：图片数据。`nil` 表示还没有。
- **Alt** `string`：替代文字，同时是可访问名称。
- **Loading** `bool`：正在加载。`Source` 为 `nil` 时显示“正在加载图片…”。
- **Error** `string`：加载失败的说明，优先于 `Loading`。
- **Retry** `func()`：点“重试图片”按钮时调用。只在 `Error` 非空时显示按钮。

### ImagePreview

`func ImagePreview(c *ui.Context, image ImageResource, state *ImagePreviewState) ui.Element`

工具栏标题是 `image.Alt`。“−”和“+”每次缩放 1.25 倍，范围 0.1 到 8。

- **image** `ImageResource`：要预览的图片。
- **state** `*ImagePreviewState`：缩放状态，由调用方持有。

### ImagePreviewState

- **Fit** `bool`：为真时按可用尺寸适应窗口。
- **Zoom** `float32`：缩放比例，`Fit` 为假时有效。小于等于 0 时按 1 处理。
