package kit

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestProgressBar 验证线性进度条的高度、填充比例和超界裁剪。
func TestProgressBar(t *testing.T) {
	// 1. 每个比例画在 184 宽的槽里，量高度并取填充边缘的像素。
	theme := OneLight()
	tests := []struct {
		name  string
		value float32
		fill  float32
	}{{"empty", 0, 0}, {"quarter", .25, .25}, {"full", 1, 1}, {"above-one-clamps", 3, 1}, {"below-zero-clamps", -1, 0}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, theme)
				ui.Column(c).Width(216).Padding(theme.Space(4)).Children(func() { bounds = ProgressBar(c, test.value).Bounds() })
			}, 300, 80)
			tt.Frame()
			if bounds.H != 4 || bounds.W != 184 {
				t.Fatalf("bar %v", bounds)
			}
			y := bounds.Y + bounds.H/2
			edge := bounds.X + bounds.W*test.fill
			// 填充内侧是强调色，外侧是悬停色轨道；两端圆角，所以取离端点 6 像素处。
			if test.fill > 0 && !pixelNear(tt, min(edge-6, bounds.X+bounds.W-6), y, theme.Colors.Accent, 2) {
				t.Fatalf("fill color inside the edge %v", tt.Image().RGBAAt(int(edge-6), int(y)))
			}
			if test.fill < 1 && !pixelNear(tt, max(edge+6, bounds.X+6), y, theme.Colors.Hover, 2) {
				t.Fatalf("track color outside the edge %v", tt.Image().RGBAAt(int(edge+6), int(y)))
			}
		})
	}
}

// TestProgressRing 验证环形进度是量表而不是加载圈：整圈轨道始终存在，进度弧从 12 点顺时针填充，到 85% 变警告色、95% 变危险色。
func TestProgressRing(t *testing.T) {
	// 1. 每个比例画在 40 的环里，取 12 点、3 点、6 点、9 点四个位置的像素，比对轨道色和进度色。
	theme := OneLight()
	c := theme.Colors
	tests := []struct {
		name  string
		value float32
		// 依次是 12 点、3 点、6 点、9 点位置期望的颜色。
		top, right, bottom, left ui.Color
	}{
		{"empty-is-full-track", 0, c.Border, c.Border, c.Border, c.Border},
		{"tiny-keeps-minimum-arc", .005, c.Accent, c.Border, c.Border, c.Border},
		{"fourteen-percent-starts-at-twelve", .14, c.Accent, c.Border, c.Border, c.Border},
		{"half-fills-the-right-side-clockwise", .5, c.Accent, c.Accent, c.Accent, c.Border},
		{"eighty-five-is-warning", .85, c.Warning, c.Warning, c.Warning, c.Warning},
		{"ninety-seven-is-danger", .97, c.Danger, c.Danger, c.Danger, c.Danger},
		{"above-one-clamps-to-danger", 2, c.Danger, c.Danger, c.Danger, c.Danger},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bounds ui.Rect
			tt := ui.NewTester(func(ctx *ui.Context) {
				ApplyTheme(ctx, theme)
				ui.Column(ctx).Padding(theme.Space(4)).Children(func() { bounds = ProgressRing(ctx, test.value, 40).Bounds() })
			}, 100, 100)
			tt.Frame()
			cx, cy, edge := bounds.X+bounds.W/2, bounds.Y+bounds.H/2, bounds.W/2-1
			points := []struct {
				name string
				x, y float32
				want ui.Color
			}{{"top", cx, cy - edge, test.top}, {"right", cx + edge, cy, test.right}, {"bottom", cx, cy + edge, test.bottom}, {"left", cx - edge, cy, test.left}}
			for _, point := range points {
				if !pixelNear(tt, point.x, point.y, point.want, 6) {
					t.Fatalf("%s pixel %v want %v", point.name, tt.Image().RGBAAt(int(point.x), int(point.y)), point.want)
				}
			}
		})
	}
}
