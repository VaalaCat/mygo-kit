package kit

import (
	"github.com/egoist/mygo/ui"
)

// ListRowOpt 描述 ListRow 的内容和状态。
type ListRowOpt struct {
	// Icon 是行首图标，14，画在 16 宽的图标位中央；为空时不留图标位。
	Icon  *ui.SVG
	Title string
	// Description 非空时行变成两行、44 高（标题 20 加说明 16，上下各留 4），说明用 12/16 次要色单行截断。
	Description string
	Selected    bool
	// Depth 是缩进层级，每级向右 16。
	Depth int
	// Trailing 画行尾的状态或徽标；单行时与标题按基线对齐，两行时在整行垂直居中。
	Trailing func()
	// Actions 画行尾的操作按钮（28 见方的 IconButton）；只在悬停、选中或键盘聚焦到行内时显示，不显示时仍占位，行内容不会跳动。
	Actions func()
}

// ListRow 创建整行列表项：28 高（带说明时 44），整行是一个按钮，调用方用 Clicked 响应选择，用 FillWidth、Key、Label 调整。
// 图标在 16 宽的图标位里居中，标题和行尾文字按基线对齐，状态底色由 controlBackground 给出。
func ListRow(c *ui.Context, opt ListRowOpt) ui.Element {
	// 1. 行高、缩进和内边距：文字起点在 12 加每级 16，有操作按钮时右侧只留 4，让按钮与标题栏的图标按钮同在一条右对齐线。
	t := GetTheme(c)
	height, top := t.Metrics.ControlHeight, float32(0)
	if opt.Description != "" {
		// 两行的文字块 36 高，上下各留 Space(1)，行共 44 高。
		height, top = t.Metrics.HeaderHeight+t.Space(2), t.Space(1)
	}
	right := t.Space(3)
	if opt.Actions != nil {
		right = t.Space(1)
	}
	row := ui.ButtonBase(c).Height(height).Padding(top, right, top, t.Space(3)+float32(opt.Depth)*t.Space(4)).Gap(t.Space(2)).Justify(ui.Start).AlignItems(ui.Center).Radius(0).Label(opt.Title).TextColor(t.Colors.Text)
	// 底色在绘制时读状态，调用方随后链式设置的 Disabled 也会生效。
	row.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, controlBackground(&t, controlState{Disabled: row.IsDisabled(), Pressed: row.Pressed(), Selected: opt.Selected, Hovered: row.Hovered()}), 0)
	})
	row.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if row.FocusVisible() {
			paintFocusRing(p, r, &t, 0)
		}
	})
	// 2. 操作按钮是否显示在构建子元素之前确定，避免本帧子元素改变悬停判断。
	showActions := opt.Selected || row.Hovered() || row.FocusWithin()
	row.Children(func() {
		if opt.Icon != nil {
			iconSlot(c, opt.Icon, t.Colors.Muted)
		}
		// 3. 单行：标题与行尾按基线对齐；两行：标题和说明竖排，行尾在整行居中。
		if opt.Description == "" {
			ui.Row(c).Grow(1).Shrink(1).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Baseline).Children(func() {
				t.TextRole(ui.Text(c, opt.Title).SingleLine().Grow(1).Shrink(1).MinWidth(0), TextUI)
				if opt.Trailing != nil {
					opt.Trailing()
				}
			})
		} else {
			ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Gap(0).Children(func() {
				t.TextRole(ui.Text(c, opt.Title).SingleLine(), TextUI)
				t.TextRole(ui.Text(c, opt.Description).SingleLine(), TextCaption).TextColor(t.Colors.Muted)
			})
			if opt.Trailing != nil {
				ui.Row(c).Shrink(0).Gap(t.Space(2)).AlignItems(ui.Center).Children(opt.Trailing)
			}
		}
		// 4. 操作按钮不显示时透明并让指针穿过，仍占着位置。
		if opt.Actions != nil {
			actions := ui.Row(c).Shrink(0).AlignItems(ui.Center).Children(opt.Actions)
			if !showActions {
				actions.Opacity(0).PassThrough()
			}
		}
	})
	return row
}

// iconSlot 绘制行首图标位：14 的图标在 16 见方的盒子里居中，图标位左缘就是所在列的文字起点基准。
func iconSlot(c *ui.Context, icon *ui.SVG, color ui.Color) ui.Element {
	// 1. 盒子固定 16，图标取默认 14，居中。
	return ui.Row(c).Size(IconSizeLarge, IconSizeLarge).Shrink(0).Center().Children(func() {
		ui.Icon(c, icon).Size(IconSizeDefault, IconSizeDefault).Shrink(0).TextColor(color)
	})
}
