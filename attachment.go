package kit

import (
	"github.com/egoist/mygo/ui"
)

// AttachmentChip 显示一个待发送附件：20×20 缩略图、截断的文件名和 20 见方的移除按钮，整体高 28。
// thumb 为 nil 时显示文件图标；onRemove 为 nil 时不显示移除按钮。
func AttachmentChip(c *ui.Context, name string, thumb ui.ImageSource, onRemove func()) ui.Element {
	t := GetTheme(c)
	// 1. 内边距 4 加 20 的内容正好 28 高（边线算在高度内）。
	chip := ui.Row(c).Height(t.Metrics.ControlHeight).Shrink(0).MaxWidth(attachmentMaxWidth).Padding(t.Space(1)).Gap(t.Space(2)).AlignItems(ui.Center).
		Radius(t.Metrics.RadiusLarge).Background(t.Colors.Elevated).Border(1, t.Colors.Border).Label(name)
	chip.Children(func() {
		// 2. 缩略图或文件图标，圆角 4 与外框 8 同心。
		if thumb != nil {
			ui.Image(c, thumb).Size(t.Metrics.CompactHeight, t.Metrics.CompactHeight).Shrink(0).Fit(ui.Cover).Radius(t.Metrics.Radius).Clip().Label(name)
		} else {
			ui.Row(c).Size(t.Metrics.CompactHeight, t.Metrics.CompactHeight).Shrink(0).Center().Radius(t.Metrics.Radius).Background(t.Colors.Panel).Children(func() {
				ui.Icon(c, IconFileText).Size(IconSizeSmall, IconSizeSmall).TextColor(t.Colors.Muted)
			})
		}
		t.TextRole(ui.Text(c, name).SingleLine().Shrink(1).MinWidth(0), TextUI)
		// 3. 移除按钮是 20 见方的紧凑图标按钮。
		if onRemove != nil && SmallIconButton(c, IconX, text(c, "removeAttachment")).Clicked() {
			onRemove()
		}
	})
	return chip
}

// attachmentMaxWidth 是附件 chip 的最大宽度。
const attachmentMaxWidth = 240
