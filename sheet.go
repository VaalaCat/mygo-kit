package kit

import "github.com/egoist/mygo/ui"

// SheetSide 选择侧边抽屉贴靠的窗口边缘。
type SheetSide int

// 抽屉贴靠左侧或右侧。
const (
	SheetLeft SheetSide = iota
	SheetRight
)

// Sheet 创建贴靠窗口左侧或右侧的模态抽屉；点遮罩或按 Escape 把 *open 置为 false。
// 窗口里有 TitleBar 时抽屉和遮罩从顶栏下缘开始，顶栏和系统窗口按钮不被盖住；没有 TitleBar 时铺满整个窗口。
func Sheet(c *ui.Context, open *bool, side SheetSide, width float32, content func()) ui.Element {
	return sheetBelow(c, open, side, width, getTitleBarHeight(c), content)
}

// sheetBelow 同 Sheet，遮罩和抽屉从窗口顶部下移 top 开始；top 为 0 时铺满整个窗口。
func sheetBelow(c *ui.Context, open *bool, side SheetSide, width, top float32, content func()) ui.Element {
	// 1. 遮罩只盖住顶栏以下的内容区，抽屉沿贴靠边拉满遮罩高度。
	t := GetTheme(c)
	return ui.DialogBase(c, open, func(back, panel ui.Element) {
		edge := ui.Start
		var borderRight, borderLeft float32
		if side == SheetRight {
			edge = ui.End
			borderLeft = 1
		} else {
			borderRight = 1
		}
		back.Top(top).Background(ui.RGBA(0, 0, 0, 0.35)).Row().Justify(edge).AlignItems(ui.Stretch)
		// 2. 抽屉使用面板底色，只在靠近内容的一侧画分隔线。
		panel.Column().Gap(0).Width(width).MaxWidthPercent(90).Clip().Background(t.Colors.Panel).BorderWidth(0, borderRight, 0, borderLeft).BorderColor(t.Colors.Border)
		panel.Children(content)
	})
}
