package kit

import (
	"github.com/egoist/mygo/ui"
)

// Disclosure 绘制无边框整行标题及可折叠内容。
func Disclosure(c *ui.Context, label string, open *bool, content func()) ui.Element {
	// 1. 标题保留框架的展开语义及键盘处理。
	t := GetTheme(c)
	parts := ui.CollapsibleBase(c, open)
	tr := parts.Trigger.Height(t.Metrics.ControlHeight).Padding(0, t.Space(1)).Gap(t.Space(1)).AlignItems(ui.Center).Radius(t.Metrics.Radius).Label(label).TextColor(t.Colors.Muted)
	if tr.Hovered() {
		tr.TextColor(t.Colors.Text)
	}
	// 底色在绘制时读状态，调用方随后链式设置的 Disabled 也会生效。
	tr.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, controlBackground(&t, controlState{Disabled: tr.IsDisabled(), Pressed: tr.Pressed(), Hovered: tr.Hovered()}), t.Metrics.Radius)
	})
	tr.Children(func() {
		arrow := IconChevronRight
		if *open {
			arrow = IconChevronDown
		}
		ui.Icon(c, arrow).Size(IconSizeDefault, IconSizeDefault).Shrink(0)
		t.TextRole(ui.Text(c, label).SingleLine().Shrink(1), TextUI)
	})
	// 2. 展开内容交给共用的高度过渡。
	disclosureBody(c, tr, open, [4]float32{t.Space(1), 0, t.Space(2), t.Space(6)}, content)
	return tr
}

// disclosureBody 在触发行之后绘制随 open 变高的内容；padding 依次为上、右、下、左。
func disclosureBody(c *ui.Context, trigger ui.Element, open *bool, padding [4]float32, content func()) {
	// 1. 高度只负责视觉过渡；关闭中的内容立即禁用。
	t := GetTheme(c)
	target := float32(0)
	if *open {
		target = 1
	}
	progress := target
	if !t.Motion.Reduced && !c.Preferences().ReduceMotion {
		duration := t.Motion.DisclosureClose
		if *open {
			duration = t.Motion.DisclosureOpen
		}
		progress = trigger.Animate("kit-disclosure", target, duration)
	}
	if progress > 0 {
		clip := ui.Column(c).ClipY().Disabled(!*open)
		var inner ui.Element
		clip.Children(func() {
			inner = ui.Column(c).Shrink(0).Gap(t.Space(2)).Padding(padding[0], padding[1], padding[2], padding[3])
			inner.Children(content)
		})
		if progress < 1 {
			clip.Height(inner.Bounds().H * progress)
		}
	}
}
