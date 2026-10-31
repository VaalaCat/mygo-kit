// Package kit 为 MyGo 原生界面提供连续面板、直角控件和通用内容阅读组件。
package kit

import (
	"github.com/egoist/mygo/ui"
	"time"
)

// Colors 按界面职责保存颜色，不携带应用状态；AccentSubtle 是带透明度的强调色浅底，用于选中的分段等强调态。
type Colors struct {
	Background, Panel, Editor, Elevated                  ui.Color
	Text, Muted, Placeholder, Disabled                   ui.Color
	Border, BorderSubtle, Focus, Hover, Active, Selected ui.Color
	Accent, AccentText, AccentSubtle                     ui.Color
	Danger, Warning, Success                             ui.Color
	Input, Selection, Added, Deleted                     ui.Color
	ANSI                                                 [16]ui.Color
}

// SyntaxColors 保存代码语法角色的颜色。
type SyntaxColors struct{ Keyword, String, Function, Type, Number, Comment, Property, Operator ui.Color }

// Typography 保存界面、正文和代码字体及字号。
type Typography struct {
	Font, MonoFont                               string
	BodySize, CaptionSize, HeadingSize, CodeSize float32
}

// TextRole 是组件选择的文字角色，字号、行高、字体、字重和数字样式由主题统一给出。
type TextRole int

// 文字角色：每个角色对应一组固定字号与行高，都落在四像素网格上。
const (
	TextProse     TextRole = iota // 14/24：Markdown 正文、详情说明
	TextUI                        // 14/20：行标题、按钮、输入框、工具名
	TextCaption                   // 12/16：时间、提示、次要状态
	TextToken                     // 12/20：token 数，等宽数字
	TextCode                      // 12/20：代码、JSON、工具输出、RPC 行，等宽字体
	TextCodeLabel                 // 12/16：工具详情里的小节标签和内容框标题，等宽字体，调用方配次要色
	TextTitle                     // 16/24：小节和页面标题，字重 400
)

// Metrics 保存四像素网格、三档高度（20/28/36）及圆角尺寸。
type Metrics struct {
	Grid                                       float32
	CompactHeight, ControlHeight, HeaderHeight float32 // 紧凑控件 20、普通控件和列表行 28、标题栏 36
	StatusHeight                               float32
	Radius, RadiusLarge                        float32 // 控件 4，浮层和卡片 8
}

// 图标尺寸三档，照 Zed agent_ui：14 是 28 高按钮、行首图标、标签和状态栏的默认值；
// 12 用于 20 高紧凑控件和行尾指示（下拉箭头）；16 用于空状态和 36 高标题栏的主操作。
const (
	IconSizeSmall   float32 = 12
	IconSizeDefault float32 = 14
	IconSizeLarge   float32 = 16
)

// RadiusFull 是圆形控件的圆角，大于任何控件边长的一半，实际绘制时贴合成圆。
const RadiusFull float32 = 9999

// Motion 保存通用界面过渡时长；Reduced 关闭本库的动画。
type Motion struct {
	Fade, DisclosureOpen, DisclosureClose time.Duration
	Reduced                               bool
}

// Theme 是可复制的主题值，由应用保存，每帧交给 ApplyTheme。
type Theme struct {
	Dark       bool
	Colors     Colors
	Typography Typography
	Metrics    Metrics
	Motion     Motion
	Syntax     SyntaxColors
}

type themeKey struct{}
type frameTheme struct {
	value  Theme
	native ui.Theme
}

// OneLight 返回基于 Zed One Light 的明色工作台主题。
func OneLight() Theme {
	// 1. 读取语义颜色并设置统一尺寸。
	return makeTheme(false)
}

// OneDark 返回基于 Zed One Dark 的暗色工作台主题。
func OneDark() Theme {
	// 1. 读取语义颜色并设置统一尺寸。
	return makeTheme(true)
}

// makeTheme 组装明暗颜色、字体、尺寸和动效。
func makeTheme(dark bool) Theme {
	// 1. 统一文字与紧凑尺寸。
	t := Theme{Dark: dark, Typography: Typography{Font: FontSans, MonoFont: FontMono, BodySize: 14, CaptionSize: 12, HeadingSize: 16, CodeSize: 12}, Metrics: Metrics{Grid: 4, CompactHeight: 20, ControlHeight: 28, HeaderHeight: 36, StatusHeight: 28, Radius: 4, RadiusLarge: 8}, Motion: Motion{Fade: 120 * time.Millisecond, DisclosureOpen: 200 * time.Millisecond, DisclosureClose: 150 * time.Millisecond}}
	// 2. 按语义角色映射 Zed One 配色。
	if dark {
		t.Colors = Colors{Background: ui.Hex("#3b414d"), Panel: ui.Hex("#2f343e"), Editor: ui.Hex("#282c33"), Elevated: ui.Hex("#2f343e"), Text: ui.Hex("#dce0e5"), Muted: ui.Hex("#a9afbc"), Placeholder: ui.Hex("#878a98"), Disabled: ui.Hex("#878a98"), Border: ui.Hex("#464b57"), BorderSubtle: ui.Hex("#363c46"), Focus: ui.Hex("#47679e"), Hover: ui.Hex("#363c46"), Active: ui.Hex("#454a56"), Selected: ui.Hex("#454a56"), Accent: ui.Hex("#74ade8"), AccentText: ui.Hex("#282c33"), AccentSubtle: ui.Hex("#74ade833"), Danger: ui.Hex("#d07277"), Warning: ui.Hex("#dec184"), Success: ui.Hex("#a1c181"), Input: ui.Hex("#282c33"), Selection: ui.Hex("#74ade83d"), Added: ui.Hex("#a1c1811a"), Deleted: ui.Hex("#d072771a")}
		t.Syntax = SyntaxColors{Keyword: ui.Hex("#b477cf"), String: ui.Hex("#a1c181"), Function: ui.Hex("#73ade9"), Type: ui.Hex("#6eb4bf"), Number: ui.Hex("#bf956a"), Comment: ui.Hex("#5d636f"), Property: ui.Hex("#d07277"), Operator: ui.Hex("#6eb4bf")}
		values := []string{"#282c34", "#e06c75", "#98c379", "#e5c07b", "#61afef", "#c678dd", "#56b6c2", "#abb2bf", "#636d83", "#ea858b", "#aad581", "#ffd885", "#85c1ff", "#d398eb", "#6ed5de", "#fafafa"}
		for i, v := range values {
			t.Colors.ANSI[i] = ui.Hex(v)
		}
	} else {
		t.Colors = Colors{Background: ui.Hex("#dcdcdd"), Panel: ui.Hex("#ebebec"), Editor: ui.Hex("#fafafa"), Elevated: ui.Hex("#ebebec"), Text: ui.Hex("#242529"), Muted: ui.Hex("#58585a"), Placeholder: ui.Hex("#7e8086"), Disabled: ui.Hex("#7e8086"), Border: ui.Hex("#c9c9ca"), BorderSubtle: ui.Hex("#dfdfe0"), Focus: ui.Hex("#7d82e8"), Hover: ui.Hex("#dfdfe0"), Active: ui.Hex("#cacaca"), Selected: ui.Hex("#cacaca"), Accent: ui.Hex("#5c78e2"), AccentText: ui.Hex("#fafafa"), AccentSubtle: ui.Hex("#5c78e226"), Danger: ui.Hex("#d1383d"), Warning: ui.Hex("#b48a00"), Success: ui.Hex("#3f953a"), Input: ui.Hex("#fafafa"), Selection: ui.Hex("#5c78e23d"), Added: ui.Hex("#3f953a1a"), Deleted: ui.Hex("#d1383d1a")}
		t.Syntax = SyntaxColors{Keyword: ui.Hex("#a449ab"), String: ui.Hex("#649f57"), Function: ui.Hex("#5b79e3"), Type: ui.Hex("#3882b7"), Number: ui.Hex("#ad6e25"), Comment: ui.Hex("#a2a3a7"), Property: ui.Hex("#d3604f"), Operator: ui.Hex("#3882b7")}
		values := []string{"#000000", "#de3e35", "#3f953a", "#d2b67b", "#2f5af3", "#950095", "#1d8991", "#bbbbbb", "#555555", "#de3e35", "#3f953a", "#d2b67b", "#2f5af3", "#950095", "#1d8991", "#ffffff"}
		for i, v := range values {
			t.Colors.ANSI[i] = ui.Hex(v)
		}
	}
	return t
}

// Space 返回 n 个网格单位的间距；只用 1、2、3、4、6、8 档，即 4、8、12、16、24、32，5（20）只用于行高。
func (t Theme) Space(n int) float32 {
	return float32(n) * t.Metrics.Grid
}

// TextRole 把文字角色的字号、固定行高、字体和字重应用到文本或容器元素上，颜色和对齐仍由调用方决定。
func (t Theme) TextRole(e ui.Element, role TextRole) ui.Element {
	// 1. 每个角色给出完整样式，调用方不再单独写字号和行高。
	ty := t.Typography
	switch role {
	case TextProse:
		return e.Font(ty.Font).FontSize(ty.BodySize).FixedLineHeight(t.Space(6)).FontWeight(400)
	case TextCaption:
		return e.Font(ty.Font).FontSize(ty.CaptionSize).FixedLineHeight(t.Space(4)).FontWeight(400)
	case TextToken:
		return e.Font(ty.Font).FontSize(ty.CaptionSize).FixedLineHeight(t.Space(5)).FontWeight(400).FontFeatures("tnum")
	case TextCode:
		return e.Font(ty.MonoFont).FontSize(ty.CodeSize).FixedLineHeight(t.Space(5)).FontWeight(400)
	case TextCodeLabel:
		return e.Font(ty.MonoFont).FontSize(ty.CaptionSize).FixedLineHeight(t.Space(4)).FontWeight(400)
	case TextTitle:
		return e.Font(ty.Font).FontSize(ty.HeadingSize).FixedLineHeight(t.Space(6)).FontWeight(400)
	default:
		return e.Font(ty.Font).FontSize(ty.BodySize).FixedLineHeight(t.Space(5)).FontWeight(400)
	}
}

// scrollbarAlpha 是滚动条滑块在次要色上的不透明度。
const scrollbarAlpha = 0.5

// ApplyTheme 在本帧根节点保存主题并应用 MyGo 的文字、焦点和滚动条颜色。
func ApplyTheme(c *ui.Context, t Theme) {
	// 1. 保存单份主题值，并生成框架所需的基础样式。
	frame := ui.Local(c.Root(), themeKey{}, func() frameTheme { return frameTheme{} })
	frame.value = t
	n := ui.LightTheme()
	if t.Dark {
		n = ui.DarkTheme()
	}
	n.Dark = t.Dark
	n.Background = t.Colors.Editor
	n.Surface = t.Colors.Panel
	n.SurfaceHover = t.Colors.Hover
	n.SurfacePressed = t.Colors.Active
	n.Border = t.Colors.Border
	n.Text = t.Colors.Text
	n.TextMuted = t.Colors.Muted
	n.Accent = t.Colors.Accent
	n.AccentText = t.Colors.AccentText
	n.Danger = t.Colors.Danger
	n.Warning = t.Colors.Warning
	n.Success = t.Colors.Success
	n.Selection = t.Colors.Selection
	n.Focus = t.Colors.Focus
	n.Radius = t.Metrics.Radius
	n.Spacing = t.Metrics.Grid
	// 2. 滚动条细一档且跟随主题次要色，所有面板共用。
	n.Scrollbar = t.Colors.Muted.Alpha(scrollbarAlpha)
	n.ScrollbarWidth = t.Space(1)
	n.Font = t.Typography.Font
	n.FontSize = t.Typography.BodySize
	frame.native = *n
	c.SetTheme(&frame.native)
}

// GetTheme 返回本窗口的主题值；未设置时采用 One Light。
func GetTheme(c *ui.Context) Theme {
	// 1. 从根节点读取窗口独立的主题。
	return ui.Local(c.Root(), themeKey{}, func() frameTheme { return frameTheme{value: OneLight()} }).value
}

// focusRingWidth 是内嵌焦点环的线宽。
const focusRingWidth float32 = 2

// paintFocusRing 在控件自己的范围内画 2 宽的焦点色环，radius 跟随控件圆角；
// 调用方先 FocusRing(false) 关掉框架的环，再在 DrawOver 里于 FocusVisible 为真时调用本函数。
func paintFocusRing(p *ui.Painter, r ui.Rect, t *Theme, radius float32) {
	// 1. 框架自带的环画在元素外侧 3 像素，会被零间距的邻居盖住或被裁剪容器切掉；
	// Tabs、Tab、Segmented、ListRow、命令行、StatusButton 这类控件改用本函数，Stroke 的线宽画在 r 以内，环不会越出控件；
	// 独立的按钮和输入框沿用框架的环。
	p.Stroke(r, t.Colors.Focus, radius, focusRingWidth)
}

// controlState 是控件当前的交互状态，供 controlBackground 取底色。
type controlState struct{ Disabled, Pressed, Selected, Hovered bool }

// controlBackground 返回控件在给定状态下的底色，优先级为禁用 > 按下 > 选中 > 悬停。
func controlBackground(t *Theme, state controlState) ui.Color {
	// 1. 禁用交给框架的 0.5 透明度，不再换色，也不响应其他状态。
	switch {
	case state.Disabled:
		return ui.Transparent
	case state.Pressed:
		return t.Colors.Active
	case state.Selected:
		return t.Colors.Selected
	case state.Hovered:
		return t.Colors.Hover
	}
	return ui.Transparent
}
