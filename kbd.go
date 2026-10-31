package kit

import "github.com/egoist/mygo/ui"

// Kbd 创建显示快捷键的小标签，例如 "⌘N"；只负责外观，不注册快捷键。
func Kbd(c *ui.Context, keys string) ui.Element {
	// 1. 使用比所在面板略深的底色和细边线，文字取说明字号。
	t := GetTheme(c)
	return ui.Row(c).Height(t.Metrics.CompactHeight).Shrink(0).Padding(0, t.Space(1)).Justify(ui.Center).AlignItems(ui.Center).Radius(t.Metrics.Radius).Background(t.Colors.Background).Border(1, t.Colors.BorderSubtle).Children(func() {
		t.TextRole(ui.Text(c, keys).SingleLine(), TextCaption).TextColor(t.Colors.Muted)
	})
}
