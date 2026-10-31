package kit

import (
	"github.com/egoist/mygo/ui"
	"image"
	"testing"
)

// TestImage 验证资源错误重试和预览缩放。
func TestImage(t *testing.T) {
	// 1. 错误图片重试只调用应用回调。
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		Image(c, ImageResource{Alt: "Photo", Error: "unavailable", Retry: func() { retries++ }})
	}, 400, 200)
	tt.Click("Retry image")
	if retries != 1 {
		t.Fatal(retries)
	}
	// 2. 预览状态由调用方持有，恢复原始大小不替换资源。
	state := ImagePreviewState{Fit: true}
	bitmap := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 100, 60)))
	tt = ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneDark())
		ImagePreview(c, ImageResource{Source: bitmap, Alt: "Preview"}, &state).Fill()
	}, 500, 300)
	tt.Click("100%")
	if state.Fit || state.Zoom != 1 {
		t.Fatal(state)
	}
	tt.Click("+")
	if state.Zoom <= 1 {
		t.Fatal(state)
	}
	tt.Click("Fit")
	if !state.Fit {
		t.Fatal(state)
	}
}
