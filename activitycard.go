package kit

import (
	"github.com/egoist/mygo/ui"
)

// ActivityStatus 是条目标题行里紧跟标题的状态；空值表示成功或没有状态，不显示图标。
type ActivityStatus string

// 条目标题行的非成功状态。
const (
	ActivityRunning     ActivityStatus = "running"
	ActivityFailed      ActivityStatus = "failed"
	ActivityInterrupted ActivityStatus = "interrupted"
	ActivityBlocked     ActivityStatus = "blocked"
	ActivityStopped     ActivityStatus = "stopped"
)

// ActivityTone 是标题和类型图标的语义色；普通条目用次要色，悬停时变正文色。
type ActivityTone int

// 条目标题行的语义色。
const (
	ActivityNormal ActivityTone = iota
	ActivityWarning
	ActivityDanger
)

// ActivityCardOpt 指定条目标题行的图标、标题、摘要和状态。
type ActivityCardOpt struct {
	Icon    *ui.SVG // 类型图标，可展开的条目在悬停、聚焦和展开时换成箭头
	Title   string
	Summary string // 参数摘要，单行截断，不加悬停提示
	Status  ActivityStatus
	Tone    ActivityTone
	Shimmer bool // 标题闪烁，用于思考、压缩、运行中
	// Extra 是紧跟标题的自定义交互内容，例如 Session 链接；它在展开按钮之外，设置后不再显示 Summary。
	Extra func()
}

// ActivityCard 显示一行 28 高的条目标题行，detail 非空时可展开；展开状态由调用方持有。
//
//	[图标/箭头 16] 标题  状态  摘要 ········ [操作]
//
// actions 是独立于展开按钮的右侧图标按钮，只在悬停或键盘聚焦本行时可见；
// detail 内容与标题文字左对齐，左侧竖线对齐图标中线，块之间间隔 12。
func ActivityCard(c *ui.Context, open *bool, opt ActivityCardOpt, actions func(), detail func()) ui.Element {
	t := GetTheme(c)
	// 1. 外层纵向容纳标题行和详情，二者共用展开动画。
	root := ui.Column(c).Gap(0).Shrink(0)
	root.Children(func() {
		row := ui.Row(c).Height(t.Metrics.ControlHeight).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center)
		var lead ui.Element
		row.Children(func() {
			// 2. 左半是展开按钮；没有详情时是不可交互的普通行，没有箭头。
			expandable := detail != nil
			if expandable {
				lead = ui.CollapsibleBase(c, open).Trigger
			} else {
				lead = ui.Row(c)
			}
			label := opt.Title
			if opt.Summary != "" && opt.Extra == nil {
				label += " " + opt.Summary
			}
			lead.Height(t.Metrics.ControlHeight)
			if opt.Extra == nil {
				lead.Grow(1).Shrink(1).MinWidth(0)
			}
			lead.Padding(0).Gap(t.Space(2)).AlignItems(ui.Center).Radius(t.Metrics.Radius).Label(label)
			active := expandable && (lead.Hovered() || lead.FocusVisible() || *open)
			lead.Children(func() { activityLead(c, t, opt, active, *open) })
			if opt.Extra != nil {
				ui.Row(c).Shrink(1).MinWidth(0).AlignItems(ui.Center).Margin(0, 0, 0, t.Space(1)).Children(opt.Extra)
				ui.Spacer(c)
			}
			// 3. 操作是兄弟节点，占位不变，仅在悬停或聚焦本行时显示。
			if actions != nil {
				wrap := ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center)
				if !row.Hovered() && !row.FocusWithin() {
					wrap.Opacity(0)
				}
				wrap.Children(actions)
			}
		})
		// 4. 详情缩进到标题文字起点，左侧竖线落在 16 宽图标位的中线。
		if detail != nil {
			disclosureBody(c, lead, open, [4]float32{t.Space(1), 0, t.Space(2), t.Space(2)}, func() {
				line := t.TextRole(ui.Column(c).Gap(t.Space(3)).Padding(0, 0, 0, t.Space(4)-1).BorderWidth(0, 0, 0, 1).BorderColor(t.Colors.Border), TextUI)
				line.Children(detail)
			})
		}
	})
	return root
}

// activityLead 绘制展开按钮内的图标位、标题、状态和摘要。
func activityLead(c *ui.Context, t Theme, opt ActivityCardOpt, active, open bool) {
	// 1. 标题色：语义色优先，其余次要色，悬停、聚焦、展开时变正文色。
	fg := t.Colors.Muted
	switch {
	case opt.Tone == ActivityWarning:
		fg = t.Colors.Warning
	case opt.Tone == ActivityDanger:
		fg = t.Colors.Danger
	case active:
		fg = t.Colors.Text
	}
	// 2. 图标位 16，图标 14 居中：可展开且交互中换成箭头；详情竖线和详情文字因此仍落在 16 宽图标位加 8 间距的网格上。
	icon := opt.Icon
	if active {
		icon = IconChevronRight
		if open {
			icon = IconChevronDown
		}
	}
	if icon != nil {
		iconSlot(c, icon, fg)
	}
	if opt.Shimmer {
		Shimmer(c, opt.Title)
	} else {
		t.TextRole(ui.Text(c, opt.Title).SingleLine().Shrink(0), TextUI).TextColor(fg)
	}
	StatusIcon(c, opt.Status)
	// 3. 摘要弱化、单行截断。
	if opt.Summary != "" && opt.Extra == nil {
		t.TextRole(ui.Text(c, opt.Summary).SingleLine().Shrink(1).MinWidth(0), TextUI).TextColor(t.Colors.Placeholder)
	}
}

// StatusIcon 绘制活动状态图标；成功不显示任何内容。
func StatusIcon(c *ui.Context, status ActivityStatus) {
	t := GetTheme(c)
	switch status {
	case ActivityRunning:
		ui.Row(c).Size(IconSizeLarge, IconSizeLarge).Shrink(0).Center().Label(text(c, "statusRunning")).Children(func() { Spinner(c, IconSizeDefault) })
	case ActivityFailed:
		iconSlot(c, IconAlert, t.Colors.Danger).Label(text(c, "statusFailed"))
	case ActivityInterrupted:
		iconSlot(c, IconCirclePause, t.Colors.Warning).Label(text(c, "statusInterrupted"))
	case ActivityBlocked:
		iconSlot(c, IconBan, t.Colors.Warning).Label(text(c, "statusBlocked"))
	case ActivityStopped:
		iconSlot(c, IconCircleStop, t.Colors.Warning).Label(text(c, "statusStopped"))
	}
}
