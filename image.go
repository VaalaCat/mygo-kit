package kit

import "github.com/egoist/mygo/ui"

// ImageResource 是调用方准备好的图片及加载状态；本库不读磁盘或下载。
type ImageResource struct {
	Source  ui.ImageSource
	Alt     string
	Loading bool
	Error   string
	Retry   func()
}

// Image 显示直角图片或明确的加载／错误内容。
func Image(c *ui.Context, image ImageResource) ui.Element {
	// 1. 不可用资源保留可访问说明与重试入口。
	t := GetTheme(c)
	if image.Source == nil {
		box := ui.Column(c).Radius(t.Metrics.Radius).Padding(t.Space(2)).Gap(t.Space(1)).Background(t.Colors.Panel).Role(ui.RoleImage).Label(image.Alt)
		box.Children(func() {
			label := image.Alt
			if image.Loading {
				label = text(c, "loadingImage")
			}
			if image.Error != "" {
				label = image.Error
			}
			if label == "" {
				label = text(c, "imageUnavailable")
			}
			t.TextRole(ui.Text(c, label), TextCaption).TextColor(t.Colors.Muted)
			if image.Error != "" && image.Retry != nil && Button(c, text(c, "retryImage")).Clicked() {
				image.Retry()
			}
		})
		return box
	}
	// 2. 按原始宽高比绘制，不添加隐式圆角。
	w, h := ui.ImageSize(image.Source)
	e := ui.Image(c, image.Source).Radius(t.Metrics.Radius).Label(image.Alt).MaxWidthPercent(100)
	if w > 0 && h > 0 {
		e.AspectRatio(w / h)
	}
	return e
}

// ImagePreviewState 保存预览缩放，Fit 为 true 时按可用尺寸适应。
type ImagePreviewState struct {
	Fit  bool
	Zoom float32
}

// imageViewportMinHeight 是图片预览滚动区的最小高度。
const imageViewportMinHeight float32 = 160

// ImagePreview 显示可切换适应／原尺寸及双向滚动的图片。
func ImagePreview(c *ui.Context, image ImageResource, state *ImagePreviewState) ui.Element {
	// 1. 工具栏只改变显示尺寸，不复制图片数据。
	t := GetTheme(c)
	if state.Zoom <= 0 {
		state.Zoom = 1
	}
	root := ui.Column(c).Background(t.Colors.Editor)
	root.Children(func() {
		Header(c, image.Alt, func() {
			if Button(c, text(c, "fit")).Clicked() {
				state.Fit = true
			}
			if Button(c, text(c, "originalSize")).Clicked() {
				state.Fit = false
				state.Zoom = 1
			}
			if Button(c, "−").Clicked() {
				state.Fit = false
				state.Zoom = max(.1, state.Zoom/1.25)
			}
			if Button(c, "+").Clicked() {
				state.Fit = false
				state.Zoom = min(8, state.Zoom*1.25)
			}
		})
		viewport := ui.ScrollBoth(c).Grow(1).MinHeight(imageViewportMinHeight)
		viewport.Children(func() {
			if image.Source == nil {
				Image(c, image)
				return
			}
			w, h := ui.ImageSize(image.Source)
			scale := state.Zoom
			if state.Fit {
				bounds := viewport.Bounds()
				if bounds.W > 0 && bounds.H > 0 && w > 0 && h > 0 {
					scale = min(bounds.W/w, bounds.H/h)
				}
			}
			ui.Image(c, image.Source).Size(w*scale, h*scale).Radius(t.Metrics.Radius).Shrink(0).Label(image.Alt)
		})
	})
	return root
}
