package kit

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// shimmerWidth 是高亮带两侧各覆盖的字素数。
const shimmerWidth = 4

// Shimmer 显示次要色单行文字，并让高亮从左向右扫过；减少动画时保持静态。
func Shimmer(c *ui.Context, value string) ui.Element {
	// 1. 减少动画时用整段文字，保持可读与可查找。
	t := GetTheme(c)
	if t.Motion.Reduced || c.Preferences().ReduceMotion {
		return t.TextRole(ui.Text(c, value).NoWrap().Shrink(0), TextUI).TextColor(t.Colors.Muted)
	}
	// 2. 每个字素按到高亮中心的距离在次要色与正文色之间混合。
	glyphs := textGraphemes(value)
	root := ui.Row(c).Gap(0).AlignItems(ui.Center).Shrink(0).Label(value)
	phase := root.Loop("kit-shimmer", 2*time.Second, ui.Linear)
	center := phase*float32(len(glyphs)+2*shimmerWidth) - shimmerWidth
	root.Children(func() {
		for i, glyph := range glyphs {
			// 单独的空格文本会被折叠，改用 4 像素空白保持词距。
			if glyph == " " {
				ui.Box(c).Width(t.Space(1)).Shrink(0)
				continue
			}
			distance := float32(i) - center
			if distance < 0 {
				distance = -distance
			}
			strength := max(0, 1-distance/shimmerWidth)
			t.TextRole(ui.Text(c, glyph).NoWrap().Shrink(0), TextUI).TextColor(t.Colors.Muted.Mix(t.Colors.Text, strength))
		}
	})
	return root
}
