package kit

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestStatusDot 验证状态圆点是 8 见方的圆，用调用方给出的颜色。
func TestStatusDot(t *testing.T) {
	// 1. 量尺寸，圆心取色，包围盒角上是背景（圆角）。
	theme := OneLight()
	var bounds ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, theme)
		ui.Row(c).Padding(theme.Space(4)).Children(func() { bounds = StatusDot(c, theme.Colors.Success).Label("Connected").Bounds() })
	}, 100, 60)
	tt.Frame()
	if bounds.W != 8 || bounds.H != 8 {
		t.Fatalf("dot %v", bounds)
	}
	if !pixelNear(tt, bounds.X+4, bounds.Y+4, theme.Colors.Success, 2) {
		t.Fatalf("center %v", tt.Image().RGBAAt(int(bounds.X+4), int(bounds.Y+4)))
	}
	if pixelNear(tt, bounds.X, bounds.Y, theme.Colors.Success, 20) {
		t.Fatal("corner is filled: dot is not round")
	}
}
