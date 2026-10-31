package kit

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
)

// ProgressRing 改用警告色、危险色的比例。
const (
	progressWarning = 0.85
	progressDanger  = 0.95
)

// progressRingMinArc 是环形进度弧的最小可见比例（8°）：很低的非零占比也要画出一段弧，读起来才是仪表而不是转圈。
const progressRingMinArc = 8.0 / 360

// ProgressRing 创建环形进度，value 为 0 到 1 的比例，用于上下文用量；达到 0.85 时改用警告色，达到 0.95 时改用危险色。
// value 大于 0 时弧长至少 8°，value 为 0 时只有底环；可访问名称始终是真实百分比。
func ProgressRing(c *ui.Context, value float32, size float32) ui.Element {
	// 1. 比例限制在 0 到 1，达到阈值后切换颜色。
	t := GetTheme(c)
	value = min(max(value, 0), 1)
	fg := t.Colors.Accent
	if value >= progressDanger {
		fg = t.Colors.Danger
	} else if value >= progressWarning {
		fg = t.Colors.Warning
	}
	// 2. 底环用边线色，进度弧从顶部顺时针绘制。
	stroke := float32(2)
	return ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleProgress).Label(fmt.Sprintf("%d%%", int(math.Round(float64(value)*100)))).Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy, radius := r.X+r.W/2, r.Y+r.H/2, (min(r.W, r.H)-stroke)/2
		p.StrokePath((&ui.Path{}).Circle(cx, cy, radius), stroke, t.Colors.Border)
		if value <= 0 {
			return
		}
		arc := &ui.Path{}
		sweep := max(float64(value), progressRingMinArc)
		steps := max(int(math.Ceil(sweep*64)), 2)
		for i := 0; i <= steps; i++ {
			angle := -math.Pi/2 + 2*math.Pi*sweep*float64(i)/float64(steps)
			x, y := cx+radius*float32(math.Cos(angle)), cy+radius*float32(math.Sin(angle))
			if i == 0 {
				arc.MoveTo(x, y)
			} else {
				arc.LineTo(x, y)
			}
		}
		p.StrokePath(arc.Cap(ui.RoundCap), stroke, fg)
	})
}

// ProgressBar 创建 4 高的线性进度条，value 为 0 到 1 的比例：悬停色轨道，强调色填充，两端圆角，宽度由所在容器决定。
func ProgressBar(c *ui.Context, value float32) ui.Element {
	// 1. 比例限制在 0 到 1，填充按轨道宽度换算。
	t := GetTheme(c)
	value = min(max(value, 0), 1)
	return ui.Box(c).Height(t.Space(1)).Shrink(0).Radius(RadiusFull).Background(t.Colors.Hover).Role(ui.RoleProgress).Label(fmt.Sprintf("%d%%", int(math.Round(float64(value)*100)))).Draw(func(p *ui.Painter, r ui.Rect) {
		if value > 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * value, H: r.H}, t.Colors.Accent, r.H/2)
		}
	})
}
