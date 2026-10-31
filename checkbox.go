package kit

import (
	"github.com/egoist/mygo/ui"
)

// checkBoxSize 是复选框方框的边长，Checkbox 与 Markdown 任务列表共用。
const checkBoxSize = 16

// Checkbox 创建直角选择框及可点击标签。
func Checkbox(c *ui.Context, checked *bool, label string) ui.Element {
	// 1. 复用 Space、点击和无障碍勾选行为。
	t := GetTheme(c)
	e := ui.CheckboxBase(c, checked).Height(t.Metrics.ControlHeight).Gap(t.Space(2)).AlignItems(ui.Center).Label(label)
	e.Children(func() {
		box := ui.Box(c).Size(checkBoxSize, checkBoxSize).Radius(t.Metrics.Radius).Border(1, t.Colors.Border)
		if *checked {
			box.Background(t.Colors.Accent).Draw(func(p *ui.Painter, r ui.Rect) { checkMark(p, r, t.Colors.AccentText) })
		}
		t.TextRole(ui.Text(c, label), TextUI)
	})
	return e
}

// checkMark 在 r 内画选中的勾；勾形按 16 见方的方框设计，其他边长等比缩放。
func checkMark(p *ui.Painter, r ui.Rect, color ui.Color) {
	// 1. 三个折点按方框边长换算，描边固定 1.5。
	unit := r.W / checkBoxSize
	path := &ui.Path{}
	path.MoveTo(r.X+4*unit, r.Y+8*unit).LineTo(r.X+7*unit, r.Y+11*unit).LineTo(r.X+12*unit, r.Y+5*unit)
	p.StrokePath(path, 1.5, color)
}
