package kit

import (
	"github.com/egoist/mygo/ui"
)

// MessageKind 是 FieldMessage 的语义：说明、警告或错误。
type MessageKind int

// 提示文字的语义，分别用次要色、警告色和危险色。
const (
	MessageDescription MessageKind = iota
	MessageWarning
	MessageError
)

// settingPageMaxWidth 是 SettingPage 内容区的最大宽度。
const settingPageMaxWidth = 720

// FieldOpt 是表单字段的标签、说明和错误；说明和错误都为空时控件下方不留空。
type FieldOpt struct {
	Label, Description, Error string
}

// Field 创建竖排表单字段：标签在上、控件在中、说明或错误在下，各段间距 4。
// 标签和说明接到 control 返回元素的可访问名称与描述上，错误还会把它标为无效；点击标签聚焦控件。
// 说明和错误同时存在时只显示错误。输入类控件的危险色边线由调用方用 InputOpt.Invalid 设置。
func Field(c *ui.Context, opt FieldOpt, control func() ui.Element) {
	// 1. 字段是一列：标签用界面文字次要色，控件由调用方绘制。
	t := GetTheme(c)
	ui.Column(c).Gap(t.Space(1)).Children(func() {
		var label ui.Element
		if opt.Label != "" {
			label = t.TextRole(ui.Text(c, opt.Label), TextUI).TextColor(t.Colors.Muted)
		}
		element := control()
		// 2. 把标签、说明、错误接到控件上，让读屏和测试都能按名称找到它。
		if opt.Label != "" {
			element.Label(opt.Label)
			if label.Clicked() {
				element.Focus()
			}
		}
		element.Description(opt.Description).Error(opt.Error)
		// 3. 说明和错误在控件下方，错误优先。
		switch {
		case opt.Error != "":
			FieldMessage(c, MessageError, opt.Error)
		case opt.Description != "":
			FieldMessage(c, MessageDescription, opt.Description)
		}
	})
}

// FieldMessage 创建一行 12/16 的说明、警告或错误文字，颜色由 kind 决定，可换行。
func FieldMessage(c *ui.Context, kind MessageKind, text string) ui.Element {
	// 1. 说明用次要色，警告和错误用语义色。
	t := GetTheme(c)
	color := t.Colors.Muted
	switch kind {
	case MessageWarning:
		color = t.Colors.Warning
	case MessageError:
		color = t.Colors.Danger
	}
	return t.TextRole(ui.Text(c, text), TextCaption).TextColor(color)
}

// InlineHintOpt 设置 InlineHint 的左对齐线：列表里用 12（零值），对话列用 16。
type InlineHintOpt struct{ Inset float32 }

// InlineHint 创建列表里的一行次要提示，用于加载中或空列表；loading 为真时行首显示转圈，message 为空时用库内置的「加载中」。
func InlineHint(c *ui.Context, message string, loading bool) ui.Element {
	return InlineHintWith(c, message, loading, InlineHintOpt{})
}

// InlineHintWith 同 InlineHint，按 opt 指定左对齐线。
func InlineHintWith(c *ui.Context, message string, loading bool, opt InlineHintOpt) ui.Element {
	// 1. 一行 28 高，文字左对齐线默认与列表行的文字起点一致（12）。
	t := GetTheme(c)
	inset := opt.Inset
	if inset == 0 {
		inset = t.Space(3)
	}
	if message == "" && loading {
		message = text(c, "loading")
	}
	return ui.Row(c).Height(t.Metrics.ControlHeight).Shrink(0).Padding(0, inset).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		if loading {
			Spinner(c, IconSizeDefault)
		}
		t.TextRole(ui.Text(c, message).SingleLine().Shrink(1), TextCaption).TextColor(t.Colors.Placeholder)
	})
}

// SettingPage 创建设置页：页头是 36 高的标题栏，右侧是 actions，内容区在页头下方滚动，页头始终可见；
// 内容左右上下边距 24、最大宽度 720、靠左，小节之间间距 24。
func SettingPage(c *ui.Context, title string, actions, content func()) ui.Element {
	// 1. 标题栏不随内容滚动，标题文字起点与内容左边距同在 24。
	t := GetTheme(c)
	return ui.Column(c).Grow(1).Gap(0).Background(t.Colors.Editor).Children(func() {
		HeaderWith(c, title, HeaderOpt{Inset: t.Space(6)}, actions)
		// 2. 内容区限制宽度并向左对齐。
		ui.Scroll(c).Grow(1).Children(func() {
			ui.Column(c).MaxWidth(settingPageMaxWidth).Padding(t.Space(6)).Gap(t.Space(6)).Children(content)
		})
	})
}

// SettingSection 创建设置小节：上方是小节标题，下方是圆角框，框内是 rows 绘制的若干 SettingRow，行之间 1 像素分隔线。
func SettingSection(c *ui.Context, title string, rows func()) ui.Element {
	// 1. 标题与圆角框相距 8，框裁切子行并画行间分隔线。
	t := GetTheme(c)
	return ui.Column(c).Gap(t.Space(2)).Children(func() {
		t.TextRole(ui.Text(c, title).SingleLine(), TextTitle)
		ui.Column(c).Gap(0).Radius(t.Metrics.RadiusLarge).Clip().Background(t.Colors.Panel).Border(1, t.Colors.Border).Dividers(1, t.Colors.BorderSubtle).Children(rows)
	})
}

// SettingRow 创建一行设置：左侧是标题和可选说明，右侧是 control 绘制的控件，行高至少 36。
func SettingRow(c *ui.Context, title, description string, control func()) ui.Element {
	// 1. 左右两侧垂直居中，文字侧可收缩换行，控件侧保持原宽。
	t := GetTheme(c)
	return ui.Row(c).MinHeight(t.Metrics.HeaderHeight).Padding(t.Space(2), t.Space(3)).Gap(t.Space(4)).AlignItems(ui.Center).Children(func() {
		ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Children(func() {
			t.TextRole(ui.Text(c, title), TextUI)
			if description != "" {
				t.TextRole(ui.Text(c, description), TextCaption).TextColor(t.Colors.Muted)
			}
		})
		ui.Row(c).Shrink(0).Gap(t.Space(2)).AlignItems(ui.Center).Children(control)
	})
}
