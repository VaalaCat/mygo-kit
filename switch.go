package kit

import (
	"github.com/egoist/mygo/ui"
)

// switchKnob 是开关滑块的颜色；开和关两种底色上都用同一个浅色，保证对比。
var switchKnob = ui.Hex("#fafafa")

// switchTrackWidth 是开关轨道的宽度，高度取紧凑档 20。
const switchTrackWidth = 36

// Switch 创建 36×20 的开关：开时轨道用强调色，滑块随状态滑动；label 显示在轨道右侧并随整行点击切换，
// label 为空时只有轨道，调用方用 Label 给出可访问名称。Space 切换，禁用由框架统一半透明。
func Switch(c *ui.Context, on *bool, label string) ui.Element {
	// 1. 复用 SwitchBase 的点击、Space 和无障碍开关语义；带文字的整行 28 高，只有轨道时与轨道同高，焦点环贴着轨道。
	t := GetTheme(c)
	row := ui.SwitchBase(c, on).Gap(t.Space(2)).AlignItems(ui.Center).Radius(t.Metrics.Radius)
	if label == "" {
		// 没有文字时焦点环贴着轨道，圆角与轨道一致。
		row.Height(t.Metrics.CompactHeight).Radius(RadiusFull)
	} else {
		row.Height(t.Metrics.ControlHeight).Label(label)
	}
	// 2. 滑块位置在 0 到 1 之间过渡；减少动画时直接到位。
	target := float32(0)
	if *on {
		target = 1
	}
	position := target
	if !t.Motion.Reduced && !c.Preferences().ReduceMotion {
		position = row.Animate("kit-switch", target, t.Motion.Fade)
	}
	// 3. 轨道底色：开用强调色，关用边线色，悬停和按下逐级加深；滑块留 2 的边距。
	hovered, pressed := row.Hovered(), row.Pressed()
	off := t.Colors.Border.Mix(t.Colors.Text, .15)
	track := off
	if *on {
		track = t.Colors.Accent
	}
	switch {
	case pressed:
		track = track.Mix(t.Colors.Text, .22)
	case hovered:
		track = track.Mix(t.Colors.Text, .12)
	}
	row.Children(func() {
		ui.Box(c).Size(switchTrackWidth, t.Metrics.CompactHeight).Shrink(0).Radius(RadiusFull).Background(track).Draw(func(p *ui.Painter, r ui.Rect) {
			inset := t.Space(1) / 2
			diameter := r.H - 2*inset
			knob := ui.Rect{X: r.X + inset + position*(r.W-r.H), Y: r.Y + inset, W: diameter, H: diameter}
			p.Fill(knob, switchKnob, diameter/2)
		})
		if label != "" {
			t.TextRole(ui.Text(c, label).SingleLine().Shrink(1), TextUI)
		}
	})
	return row
}
