package kit

import "github.com/egoist/mygo/ui"

// titleBarKey 标识窗口根节点上保存的 TitleBar 高度。
type titleBarKey struct{}

// TitleBar 创建窗口顶栏，可拖动窗口，并避开系统窗口按钮占用的位置。
func TitleBar(c *ui.Context, content func()) ui.Element {
	// 1. 按系统窗口按钮的占位留出左右边距，高度不低于标题栏高度，并记录下来供 Sheet 和 ToastHost 避让。
	t := GetTheme(c)
	bar := c.TitleBar()
	height := max(bar.Height, t.Metrics.HeaderHeight)
	*ui.Local(c.Root(), titleBarKey{}, func() float32 { return 0 }) = height
	// 2. 背景与状态栏同色，底部画一条分隔线，整条可拖动窗口。
	return ui.Row(c).Height(height).Shrink(0).Padding(0, max(bar.Right, t.Space(3)), 0, max(bar.Left, t.Space(3))).Gap(t.Space(2)).AlignItems(ui.Center).Background(t.Colors.Background).BorderWidth(0, 0, 1, 0).BorderColor(t.Colors.Border).DragWindow().Children(content)
}

// getTitleBarHeight 返回本窗口里 TitleBar 的高度；窗口没有使用 TitleBar 时为 0。
func getTitleBarHeight(c *ui.Context) float32 {
	// 1. TitleBar 在本帧或上一帧绘制时记录高度，Sheet 和 ToastHost 在其后调用。
	return *ui.Local(c.Root(), titleBarKey{}, func() float32 { return 0 })
}
