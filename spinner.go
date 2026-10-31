package kit

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Spinner 创建旋转的加载图标，颜色跟随文字；减少动画时保持静止。
func Spinner(c *ui.Context, size float32) ui.Element {
	// 1. 图标使用次要色，调用方可继续覆盖 TextColor。
	t := GetTheme(c)
	icon := ui.Icon(c, IconLoader).Size(size, size).Shrink(0).TextColor(t.Colors.Muted)
	// 2. 仅在允许动画时按每秒一圈旋转。
	if !t.Motion.Reduced && !c.Preferences().ReduceMotion {
		icon.Rotate(icon.Loop("kit-spinner", time.Second, ui.Linear) * 360)
	}
	return icon
}
