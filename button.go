package kit

import (
	"github.com/egoist/mygo/ui"
)

// Button 创建无常驻边框的工具按钮，保留 MyGo 键盘与焦点行为。
func Button(c *ui.Context, label string) ui.Element {
	// 1. 使用透明背景和完整交互颜色。
	return button(c, label, buttonGhost)
}

// FlushButton 创建无水平内边距的文字按钮，文字贴齐所在列的左对齐线；悬停只改变文字颜色，不画底色。
func FlushButton(c *ui.Context, label string) ui.Element {
	// 1. 使用 flush 变体，其余行为与 Button 相同。
	return button(c, label, buttonFlush)
}

// PrimaryButton 创建突出主要操作的实色按钮。
func PrimaryButton(c *ui.Context, label string) ui.Element {
	// 1. 使用强调色显示主要操作。
	return button(c, label, buttonPrimary)
}

// DangerButton 创建文字强调危险操作的按钮。
func DangerButton(c *ui.Context, label string) ui.Element {
	// 1. 使用危险色，悬停前不添加外框。
	return button(c, label, buttonDanger)
}

// buttonVariant 选择按钮的外观变体。
type buttonVariant int

// 按钮外观：ghost 悬停才有底色，flush 无内边距无底色，primary 强调色底，danger 危险色文字。
const (
	buttonGhost buttonVariant = iota
	buttonFlush
	buttonPrimary
	buttonDanger
)

// button 绘制按钮文字、尺寸和指针状态。
func button(c *ui.Context, label string, variant buttonVariant) ui.Element {
	// 1. 复用原生基础按钮并设置密度；文字从内边距起点开始，被拉宽时也不居中；flush 变体没有水平内边距。
	t := GetTheme(c)
	b := ui.ButtonBase(c).Height(t.Metrics.ControlHeight).Padding(0, t.Space(2)).Gap(t.Space(1)).Radius(t.Metrics.Radius).Justify(ui.Start).Label(label)
	if variant == buttonFlush {
		b.Padding(0)
	}
	fg := t.Colors.Text
	switch variant {
	case buttonPrimary:
		fg = t.Colors.AccentText
	case buttonDanger:
		fg = t.Colors.Danger
	case buttonFlush:
		fg = t.Colors.Muted
		if b.Hovered() && !b.IsDisabled() {
			fg = t.Colors.Text
		}
	}
	b.TextColor(fg)
	// 2. 底色由状态函数给出；禁用由框架统一半透明，flush 不画底色。
	if variant != buttonFlush {
		b.Draw(func(p *ui.Painter, r ui.Rect) {
			state := controlState{Disabled: b.IsDisabled(), Pressed: b.Pressed(), Hovered: b.Hovered()}
			bg := controlBackground(&t, state)
			if variant == buttonPrimary {
				bg = t.Colors.Accent
				switch {
				case state.Disabled:
				case state.Pressed:
					bg = t.Colors.Accent.Mix(t.Colors.Text, .22)
				case state.Hovered:
					bg = t.Colors.Accent.Mix(t.Colors.Text, .12)
				}
			}
			p.Fill(r, bg, t.Metrics.Radius)
		})
	}
	if label != "" {
		b.Children(func() { t.TextRole(ui.Text(c, label).SingleLine(), TextUI) })
	}
	return b
}

// IconButton 创建方形 ghost 工具按钮，图标与可访问名称分别传入。
func IconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element {
	// 1. 沿用按钮状态样式并填入图标。
	t := GetTheme(c)
	b := button(c, "", buttonGhost).Size(t.Metrics.ControlHeight, t.Metrics.ControlHeight).Padding(0).Justify(ui.Center).AlignItems(ui.Center).Label(label)
	b.Children(func() { ui.Icon(c, icon).Size(IconSizeDefault, IconSizeDefault) })
	return b
}

// SmallIconButton 创建 20 见方的紧凑 ghost 图标按钮，图标 12，平时用弱化色、悬停用次要色，用于消息下方的复制、调试和标签关闭。
func SmallIconButton(c *ui.Context, icon *ui.SVG, label string) ui.Element {
	// 1. 尺寸取紧凑档，颜色按悬停切换。
	t := GetTheme(c)
	b := button(c, "", buttonGhost).Size(t.Metrics.CompactHeight, t.Metrics.CompactHeight).Padding(0).Justify(ui.Center).AlignItems(ui.Center).Label(label)
	b.TextColor(t.Colors.Placeholder)
	if b.Hovered() {
		b.TextColor(t.Colors.Muted)
	}
	b.Children(func() { ui.Icon(c, icon).Size(IconSizeSmall, IconSizeSmall) })
	return b
}

// ToolButton 创建图标加文字的 ghost 工具按钮：28 高、图标 14，可选尾部下拉箭头。
func ToolButton(c *ui.Context, icon *ui.SVG, label string, chevron bool) ui.Element {
	// 1. 图标与箭头使用次要色，文字保持正文色。
	t := GetTheme(c)
	b := button(c, "", buttonGhost).Padding(0, t.Space(2)).AlignItems(ui.Center).Label(label)
	// 2. 图标和文字是一组靠左；有下拉箭头时箭头贴右缘，被拉宽时文字不居中。
	b.Justify(ui.SpaceBetween)
	b.Children(func() {
		ui.Row(c).Shrink(1).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			if icon != nil {
				ui.Icon(c, icon).Size(IconSizeDefault, IconSizeDefault).Shrink(0).TextColor(t.Colors.Muted)
			}
			t.TextRole(ui.Text(c, label).SingleLine().Shrink(1), TextUI)
		})
		if chevron {
			ui.Icon(c, IconChevronDown).Size(IconSizeSmall, IconSizeSmall).Shrink(0).TextColor(t.Colors.Muted)
		}
	})
	return b
}

// IconButtonTip 创建带悬停提示的 IconButton：label 同时是可访问名称和提示文字，省去每次手写 Tooltip 和 Label；返回按钮，调用方用 Clicked 响应。
func IconButtonTip(c *ui.Context, icon *ui.SVG, label string) ui.Element {
	// 1. 提示挂在按钮上，保留按钮自己的点击和焦点行为。
	b := IconButton(c, icon, label)
	Tooltip(c, b, label)
	return b
}
