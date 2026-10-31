package kit

import (
	"github.com/egoist/mygo/ui"
	"time"
)

// HeaderOpt 设置标题栏的左内边距和左侧内容。
type HeaderOpt struct {
	// Inset 是标题文字距左缘的距离，对齐所在列的左对齐线：侧栏和 Dock 用 12，对话区用 16；零值取 12。
	Inset float32
	// Leading 替换标题文字，例如放一组分段控件；为空时显示 title。内容从 Inset 处开始，标题栏高度和底线不变。
	Leading func()
}

// Header 创建共享底部分隔线的紧凑面板标题栏，标题距左缘 12。
func Header(c *ui.Context, title string, actions func()) ui.Element {
	// 1. 左内边距取侧栏和 Dock 的 12。
	return HeaderWith(c, title, HeaderOpt{}, actions)
}

// HeaderWith 创建共享底部分隔线的紧凑面板标题栏，左内边距和左侧内容由 opt 决定。
func HeaderWith(c *ui.Context, title string, opt HeaderOpt, actions func()) ui.Element {
	// 1. 标题与工具按钮共用连续背景，右侧给图标按钮留出 4 像素。
	t := GetTheme(c)
	inset := opt.Inset
	if inset == 0 {
		inset = t.Space(3)
	}
	row := ui.Row(c).Height(t.Metrics.HeaderHeight).Shrink(0).Padding(0, t.Space(1), 0, inset).Gap(t.Space(1)).AlignItems(ui.Center).Background(t.Colors.Panel).BorderWidth(0, 0, 1, 0).BorderColor(t.Colors.Border)
	row.Children(func() {
		// 2. 左侧是自定义内容或标题文字，其后撑开空白，工具按钮靠右。
		if opt.Leading != nil {
			opt.Leading()
		} else {
			t.TextRole(ui.Text(c, title).SingleLine().Shrink(1), TextUI)
		}
		ui.Spacer(c)
		if actions != nil {
			actions()
		}
	})
	return row
}

// EmptyStateOpt 设置空状态的可选主操作。
type EmptyStateOpt struct {
	// ActionLabel 非空时在文字下方居中显示一个主按钮，点击时调用 OnAction。
	ActionLabel string
	OnAction    func()
}

// EmptyState 在区域中央显示一行提示，可附带一行次要说明。
func EmptyState(c *ui.Context, title, detail string) ui.Element {
	return EmptyStateWith(c, title, detail, EmptyStateOpt{})
}

// EmptyStateWith 同 EmptyState，可在文字下方带一个主操作按钮。
func EmptyStateWith(c *ui.Context, title, detail string, opt EmptyStateOpt) ui.Element {
	// 1. 标题用次要色，说明用占位色，两者居中。
	t := GetTheme(c)
	return ui.Column(c).Grow(1).Center().Gap(t.Space(1)).Padding(t.Space(6)).Children(func() {
		t.TextRole(ui.Text(c, title), TextUI).TextColor(t.Colors.Muted).TextAlign(ui.Center)
		if detail != "" {
			t.TextRole(ui.Text(c, detail), TextCaption).TextColor(t.Colors.Placeholder).TextAlign(ui.Center)
		}
		// 2. 主操作在文字下方居中。
		if opt.ActionLabel != "" && PrimaryButton(c, opt.ActionLabel).AlignSelf(ui.Center).Margin(t.Space(2), 0, 0, 0).Clicked() && opt.OnAction != nil {
			opt.OnAction()
		}
	})
}

// copiedFeedbackDuration 是复制成功后按钮显示「已复制」的时长。
const copiedFeedbackDuration = 2 * time.Second

type frameCopiedKey struct{}

// frameHeader 绘制内容框顶部的小标题与复制按钮，用于 JSON、补丁、Markdown 代码块等嵌在正文中的阅读组件。
// 复制按钮只有图标，提示和可访问名称都是 copyLabel；点击后 2 秒内换成对勾图标和「已复制」；copyDisabled 为真时按钮禁用，例如还在流式输出、没有收尾的代码块。
func frameHeader(c *ui.Context, title, copyLabel string, copyDisabled bool, copy func()) {
	// 1. 标题与 Zed 的小节标签一样用等宽 12/16 次要色，复制按钮是 20 见方的紧凑图标按钮。
	t := GetTheme(c)
	row := ui.Row(c).Unselectable().Height(t.Metrics.ControlHeight).Shrink(0).Padding(0, t.Space(1), 0, t.Space(2)).Gap(t.Space(1)).AlignItems(ui.Center).BorderWidth(0, 0, 1, 0).BorderColor(t.Colors.BorderSubtle)
	copiedAt := ui.Local(row, frameCopiedKey{}, func() time.Time { return time.Time{} })
	remaining := copiedFeedbackDuration - c.Now().Sub(*copiedAt)
	copied := !copiedAt.IsZero() && remaining > 0
	row.Children(func() {
		t.TextRole(ui.Text(c, title).SingleLine().Shrink(1), TextCodeLabel).TextColor(t.Colors.Muted)
		ui.Spacer(c)
		// 2. 复制后 2 秒内图标、提示和可访问名称换成对勾和「已复制」，每帧都按剩余时间重新请求到期唤醒，期间悬停提示的帧不会吞掉它。
		icon, label := IconCopy, copyLabel
		if copied {
			icon, label = IconCheck, text(c, "copied")
		}
		if copied {
			c.After(remaining)
		}
		button := SmallIconButton(c, icon, label).Disabled(copyDisabled)
		Tooltip(c, button, label)
		if button.Clicked() {
			copy()
			*copiedAt = c.Now()
			c.After(copiedFeedbackDuration)
		}
	})
}

// contentFrame 设置嵌入阅读组件的圆角细边框外观。
func contentFrame(c *ui.Context, e ui.Element) ui.Element {
	t := GetTheme(c)
	return e.Radius(t.Metrics.Radius).Clip().Background(t.Colors.Editor).Border(1, t.Colors.Border)
}

// Separator 创建单条水平分隔线。
func Separator(c *ui.Context) ui.Element {
	// 1. 分隔线不包围相邻内容。
	return ui.Box(c).Height(1).Shrink(0).Background(GetTheme(c).Colors.BorderSubtle)
}

// Popover 创建直角浮层并复用框架的定位、关闭和焦点处理。
func Popover(c *ui.Context, anchor ui.Element, open *bool, content func()) ui.Element {
	// 1. 仅为现有浮层添加外观。
	t := GetTheme(c)
	return ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.Margin(t.Space(1), 0, 0, 0).Radius(t.Metrics.RadiusLarge).Padding(t.Space(2)).Gap(t.Space(2)).Background(t.Colors.Elevated).Border(1, t.Colors.Border)
		panel.Children(content)
	})
}

type tooltipKey struct{}
type tooltipState struct {
	since   time.Time
	hovered bool
}

// Tooltip 每个锚点独立停留 600ms 后显示，保留键盘焦点与 Escape 行为。
func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element {
	// 1. 逐锚点计时，不继承其他提示的热启动。
	t := GetTheme(c)
	state := ui.Local(anchor, tooltipKey{}, func() tooltipState { return tooltipState{} })
	hovered := anchor.Hovered()
	if hovered && !state.hovered {
		state.since = c.Now()
	}
	state.hovered = hovered
	ready := anchor.FocusVisible()
	if hovered {
		remaining := 600*time.Millisecond - c.Now().Sub(state.since)
		if remaining <= 0 {
			ready = true
		} else {
			c.After(remaining)
		}
	}
	anchor.Description(text)
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.Hidden(!ready).Radius(t.Metrics.RadiusLarge).Padding(t.Space(1), t.Space(2)).Background(t.Colors.Elevated).Border(1, t.Colors.Border)
		tip.Children(func() { t.TextRole(ui.Text(c, text), TextCaption) })
	})
}

// Split 创建两个可拖宽的连续面板，不添加重复外框。
func Split(c *ui.Context, size *float32, first, second func()) ui.Element {
	// 1. 使用框架分隔线与可访问拖动行为，并约束面板裁切。
	return ui.Split(c, size, first, second).Fill().Gap(0).Radius(0).Background(GetTheme(c).Colors.Editor)
}

// List 创建紧凑虚拟列表，行位置和稳定 Key 由调用方的 ListState 管理。
func List(c *ui.Context, state *ui.ListState, count int, row func(int)) ui.Element {
	// 1. 长列表沿用框架的可见行构建与滚动锚点。
	return ui.List(c, state, count, row).Gap(0).Padding(0).Radius(0).Background(GetTheme(c).Colors.Editor)
}

// Tree 显示紧凑虚拟树，沿用框架的稳定键、展开与方向键导航。
func Tree[K comparable](c *ui.Context, state *ui.OutlineState[K], roots []K, children func(K) []K, row func(K)) ui.Element {
	// 1. 只增加通用工作台外观，不复制树和滚动状态。
	t := GetTheme(c)
	return ui.Outline(c, state, roots, children, row).Gap(0).Radius(0).Background(t.Colors.Editor).TextColor(t.Colors.Text)
}
