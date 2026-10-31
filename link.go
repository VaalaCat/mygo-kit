package kit

import (
	"github.com/egoist/mygo/ui"
)

// Link 创建强调色链接按钮；文字始终靠左；在纵向容器里要贴合文字宽度（只有文字可点）时，调用方自己加 AlignSelf(ui.Start)。
func Link(c *ui.Context, label string) ui.Element {
	// 1. 悬停加下划线，不设置自身对齐。
	t := GetTheme(c)
	b := ui.ButtonBase(c).Padding(0).Radius(0).Justify(ui.Start).Label(label).TextColor(t.Colors.Accent)
	b.Children(func() {
		text := t.TextRole(ui.Text(c, label).SingleLine().Shrink(1), TextUI)
		if b.Hovered() {
			text.Underline()
		}
	})
	return b
}
