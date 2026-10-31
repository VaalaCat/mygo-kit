package kit

import (
	"github.com/egoist/mygo/ui"
)

// InputOpt 设置输入框的提示、状态和外观变体；零值就是普通的可编辑输入。
type InputOpt struct {
	Placeholder string
	// Invalid 让边线显示危险色，优先于焦点色。
	Invalid bool
	// ReadOnly 把输入显示成一段可选中复制的值：没有边线、没有底色，文字起点和 28 高与普通输入一致（字段对齐不变），
	// 不随焦点和悬停变化，键盘聚焦仍有焦点环；需要「看起来可编辑」时不要用它。放进 Field 作为只读字段的正文。
	ReadOnly bool
	// Borderless 去掉边框、底色、内边距和焦点环，用于侧栏搜索、输入区这类由所在容器决定外观的输入。
	Borderless bool
	// Leading 是输入框左侧的前缀图标，14，次要色；文字区从图标右侧开始，滚动的长文字被裁在图标之外。
	Leading *ui.SVG
}

// TextInput 创建单行输入，返回原生元素以设置提示和键盘处理；等于零值 InputOpt 的 TextInputWith。
func TextInput(c *ui.Context, value *string) ui.Element {
	// 1. 普通输入没有额外选项。
	return TextInputWith(c, value, InputOpt{})
}

// TextInputWith 创建 28 高的单行输入，文字在高度内垂直居中；在 DialogWith 的正文里，按 Enter 会确认对话框。
func TextInputWith(c *ui.Context, value *string, opt InputOpt) ui.Element {
	// 1. 复用原生输入状态并应用统一样式，单行输入不写竖向内边距。
	t := GetTheme(c)
	e := ui.TextInputBase(c, value).Height(t.Metrics.ControlHeight)
	noteDialogInput(c, e, opt, true)
	return inputStyle(&t, e, opt, 0)
}

// TextArea 创建多行输入，内容与撤销状态由 MyGo 管理；等于零值 InputOpt 的 TextAreaWith。
func TextArea(c *ui.Context, value *string) ui.Element {
	// 1. 普通输入没有额外选项。
	return TextAreaWith(c, value, InputOpt{})
}

// textAreaMinHeight 是多行输入的最小高度：两行 20 的文字加上下各 4 的内边距再加边线。
const textAreaMinHeight float32 = 56

// TextAreaWith 创建最小两行高的多行输入，保留原生组字、选区和受控值处理。
func TextAreaWith(c *ui.Context, value *string, opt InputOpt) ui.Element {
	// 1. 竖向内边距 4，使第一行与单行输入的文字位置一致。
	t := GetTheme(c)
	e := ui.TextAreaBase(c, value).MinHeight(textAreaMinHeight)
	noteDialogInput(c, e, opt, false)
	return inputStyle(&t, e, opt, t.Space(1))
}

// inputTextInsetKey 标识输入框左缘到文字起点的距离，保存在输入元素上，供 CompletionMenu 对齐浮层文字。
type inputTextInsetKey struct{}

// inputStyle 设置输入的文字、底色、边线、内边距和前缀图标；vertical 是普通外观的竖向内边距。
func inputStyle(t *Theme, e ui.Element, opt InputOpt, vertical float32) ui.Element {
	// 1. 文字角色固定为界面文字，提示和只读状态交给原生输入。
	t.TextRole(e, TextUI).TextColor(t.Colors.Text)
	if opt.Placeholder != "" {
		e.Placeholder(opt.Placeholder)
	}
	if opt.ReadOnly {
		e.ReadOnly(true)
	}
	// 2. 无边框变体只留文字，外观由所在容器决定。
	pad, edge := float32(0), float32(0) // 水平内边距，边线宽度
	bg, border := ui.Transparent, ui.Transparent
	if opt.Borderless {
		e.Padding(0).Background(bg).Border(0, border).FocusRing(false)
	} else {
		// 3. 边线优先级：只读（无边线无底色）> 危险 > 聚焦 > 悬停；键盘聚焦的外环由框架绘制。
		// 只读值仍占着边线和内边距，文字起点与普通输入一致，字段对齐不变。
		pad, edge = t.Space(2), 1
		if !opt.ReadOnly {
			bg, border = t.Colors.Input, t.Colors.Border
		}
		switch {
		case opt.ReadOnly:
		case opt.Invalid:
			border = t.Colors.Danger
		case e.Focused():
			border = t.Colors.Focus
		case e.Hovered():
			border = t.Colors.Border.Mix(t.Colors.Text, .25)
		}
		e.Radius(t.Metrics.Radius).Padding(vertical, pad).Background(bg).Border(edge, border)
	}
	inset := edge + pad
	// 4. 前缀图标：图标和间距放进左边框，文字区（含裁剪）从图标右侧开始，滚动的文字不会画到图标下面；
	// 边线和图标由 DrawOver 画在边框之上，右侧边框与左侧对称，文字也不会画进右内边距。
	if opt.Leading != nil {
		icon := opt.Leading
		iconX, leftZone := inset, inset+IconSizeDefault+t.Space(2)
		e.Padding(vertical, 0).BorderWidth(0, inset, 0, leftZone).BorderColor(bg)
		e.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if edge > 0 && border != ui.Transparent {
				p.Stroke(r, border, t.Metrics.Radius, edge)
			}
			y := r.Y + (t.Metrics.ControlHeight-IconSizeDefault)/2
			p.Icon(icon, ui.Rect{X: r.X + iconX, Y: y, W: IconSizeDefault, H: IconSizeDefault}, t.Colors.Muted)
		})
		inset = leftZone
	}
	*ui.Local(e, inputTextInsetKey{}, func() float32 { return 0 }) = inset
	return e
}
