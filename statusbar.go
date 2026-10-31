package kit

import "github.com/egoist/mygo/ui"

// StatusBar 创建窗口底部的状态栏，与标题栏同色并共用顶部分隔线。
func StatusBar(c *ui.Context, content func()) ui.Element {
	// 1. 状态文字使用说明字号和次要色。
	t := GetTheme(c)
	return t.TextRole(ui.Row(c), TextCaption).Height(t.Metrics.StatusHeight).Shrink(0).Padding(0, t.Space(1)).Gap(t.Space(1)).AlignItems(ui.Center).Background(t.Colors.Background).BorderWidth(1, 0, 0, 0).BorderColor(t.Colors.Border).TextColor(t.Colors.Muted).Children(content)
}

// StatusButton 创建状态栏里 20 高的图标按钮，上下左右各留 4；active 为真时使用选中底色，label 作为可访问名称，无图标时显示 label 文字。
func StatusButton(c *ui.Context, icon *ui.SVG, label string, active bool) ui.Element {
	// 1. 复用原生按钮行为，只画选中、按下和悬停底色。
	t := GetTheme(c)
	b := ui.ButtonBase(c).Height(t.Metrics.CompactHeight).Padding(0, t.Space(1)).Gap(t.Space(1)).Justify(ui.Center).AlignItems(ui.Center).Radius(t.Metrics.Radius).Label(label)
	if icon != nil {
		b.Width(t.Metrics.CompactHeight).Padding(0)
	} else {
		b.Padding(0, t.Space(2))
	}
	fg := t.Colors.Muted
	if active {
		fg = t.Colors.Text
	}
	b.TextColor(fg)
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		bg := controlBackground(&t, controlState{Disabled: b.IsDisabled(), Pressed: b.Pressed(), Selected: active, Hovered: b.Hovered()})
		p.Fill(r, bg, t.Metrics.Radius)
	})
	b.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if b.FocusVisible() {
			paintFocusRing(p, r, &t, t.Metrics.Radius)
		}
	})
	// 2. 状态栏按钮图标取默认 14，没有图标时退回文字。
	b.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(IconSizeDefault, IconSizeDefault).Shrink(0)
			return
		}
		t.TextRole(ui.Text(c, label).SingleLine(), TextCaption)
	})
	return b
}
