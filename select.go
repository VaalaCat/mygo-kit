package kit

import (
	"github.com/egoist/mygo/ui"
	"strings"
)

// SelectOption 保存独立的选项值、显示文字及禁用状态。
type SelectOption struct {
	Value, Label string
	// Description 显示在选项行右侧，12/16 次要色单行；触发器不显示它，例如模型选择里的 Provider。
	Description string
	Disabled    bool
	// Header 让这一项成为分组标题：只显示 Label，是 12/16 次要色的小标题，不可选，方向键跳过，Value 被忽略；
	// 搜索时只在组内有匹配项时保留。
	Header bool
}

// SelectOpt 设置选项、空选择提示和可访问名称。
type SelectOpt struct {
	Options            []SelectOption
	Placeholder, Label string
	Disabled           bool
	Searchable         bool
	SearchPlaceholder  string
	// Bordered 让触发器与表单输入框外观一致；工具栏保持无边框。
	Bordered bool
	// TriggerLabel 非空时决定触发器上显示的文字，参数是当前选中的选项；选项行仍显示各自的 Label。
	// 例如选项是「自动 / 低 / 中 / 高」，触发器显示「推理：自动」。没有选中任何选项时仍显示 Placeholder 或原值。
	TriggerLabel func(selected SelectOption) string
}

// 下拉面板的最小宽度和最大高度。
const (
	selectMinWidth  = 160
	selectMaxHeight = 320
)

type selectSearchKey struct{}
type selectFilterKey struct{}

// Select 显示工具栏选择器；整组禁用使用 opt.Disabled，同时禁用触发器和已打开的弹层。
// 带搜索的弹层里：输入过滤词，↑↓ 移动高亮，Enter 选中高亮项（没有高亮时选第一个匹配项）并关闭弹层，Esc 只关闭弹层。
func Select(c *ui.Context, value *string, opt SelectOpt) ui.Element {
	// 1. 使用真实值绑定，显示名称与保存值独立。
	t := GetTheme(c)
	parts := ui.SelectBase(c, value)
	trigger := parts.Trigger.Disabled(opt.Disabled).Height(t.Metrics.ControlHeight).Padding(0, t.Space(2)).Gap(t.Space(1)).Radius(t.Metrics.Radius).AlignItems(ui.Center).Justify(ui.SpaceBetween)
	if opt.Label != "" {
		trigger.Label(opt.Label)
	}
	if opt.Bordered {
		border := t.Colors.Border
		// 弹层打开时边线使用焦点色；键盘焦点环由框架绘制。
		if parts.Open() {
			border = t.Colors.Focus
		}
		trigger.Background(t.Colors.Input).Border(1, border)
	} else {
		trigger.Background(controlBackground(&t, controlState{Disabled: opt.Disabled, Pressed: trigger.Pressed(), Selected: parts.Open(), Hovered: trigger.Hovered()}))
	}
	label := opt.Placeholder
	matched := false
	for _, option := range opt.Options {
		if !option.Header && option.Value == *value {
			label = option.Label
			if opt.TriggerLabel != nil {
				label = opt.TriggerLabel(option)
			}
			matched = true
			break
		}
	}
	if !matched && *value != "" {
		label = *value
	}
	trigger.Children(func() {
		text := t.TextRole(ui.Text(c, label).SingleLine().Shrink(1), TextUI)
		if !matched && *value == "" {
			text.TextColor(t.Colors.Placeholder)
		}
		ui.Icon(c, IconChevronDown).Size(IconSizeSmall, IconSizeSmall).Shrink(0).TextColor(t.Colors.Muted)
	})
	query := ui.Local(trigger, selectSearchKey{}, func() string { return "" })
	filtered := ui.Local(trigger, selectFilterKey{}, func() string { return "" })
	if !parts.Open() {
		*query = ""
		*filtered = ""
	}
	firstMatch, highlighted := "", ""
	hasMatch, submitted := false, false
	var enabled []string
	// 2. 弹层与触发器隔 Space(1)，翻到上方时间距不变；禁用行不注册为可选值，方向键自动跳过。
	// 搜索框是弹层里的输入，不向 DialogWith 报告，所以不用 TextInput。
	parts.Popup(func(panel ui.Element) {
		panel.Disabled(opt.Disabled).Margin(t.Space(1), 0, 0, 0).Radius(t.Metrics.RadiusLarge).Padding(t.Space(1)).Background(t.Colors.Elevated).Border(1, t.Colors.Border).MinWidth(selectMinWidth).MaxHeight(selectMaxHeight)
		if opt.Searchable {
			placeholder := opt.SearchPlaceholder
			if placeholder == "" {
				placeholder = text(c, "search")
			}
			search := ui.TextInputBase(c.Key("select-search"), query).Height(t.Metrics.ControlHeight)
			inputStyle(&t, search, InputOpt{Placeholder: placeholder}, 0)
			search.Label(placeholder).AutoFocus()
			submitted = search.Submitted()
		}
		visible := selectVisibleOptions(opt.Options, query, opt.Searchable)
		for i, option := range visible {
			if option.Header {
				header := ui.Row(c).Height(t.Metrics.CompactHeight).Padding(0, t.Space(2)).AlignItems(ui.Center).Role(ui.RoleListItem).Label(option.Label)
				if i > 0 {
					header.Margin(t.Space(1), 0, 0, 0)
				}
				header.Children(func() { t.TextRole(ui.Text(c, option.Label).SingleLine(), TextCaption).TextColor(t.Colors.Muted) })
				continue
			}
			if option.Disabled {
				ui.Row(c).Height(t.Metrics.ControlHeight).Padding(0, t.Space(2)).Gap(t.Space(3)).Justify(ui.SpaceBetween).AlignItems(ui.Center).Disabled(true).Role(ui.RoleListItem).Label(option.Label).Children(func() { selectOptionText(c, &t, option) })
				continue
			}
			if !hasMatch {
				firstMatch = option.Value
				hasMatch = true
			}
			enabled = append(enabled, option.Value)
			item := parts.Item(option.Value).Radius(t.Metrics.Radius).Height(t.Metrics.ControlHeight).Padding(0, t.Space(2)).Gap(t.Space(3)).Justify(ui.SpaceBetween).AlignItems(ui.Center).Label(option.Label)
			if item.Highlighted() {
				highlighted = option.Value
			}
			item.Background(controlBackground(&t, controlState{Pressed: item.Pressed(), Hovered: item.Highlighted()}))
			if option.Value == *value {
				item.TextColor(t.Colors.Accent)
			}
			item.Children(func() { selectOptionText(c, &t, option) })
		}
		if len(visible) == 0 {
			t.TextRole(ui.Text(c, text(c, "noResults")), TextUI).Padding(t.Space(2)).TextColor(t.Colors.Muted)
		}
		// 搜索框拿着焦点时 ↑↓ 到不了触发器，由弹层接管并移动高亮；其余焦点在触发器上时由框架处理。
		if opt.Searchable {
			switch {
			case panel.OverlayShortcut(0, ui.KeyDown):
				parts.Highlight(selectNeighbor(enabled, highlighted, 1))
			case panel.OverlayShortcut(0, ui.KeyUp):
				parts.Highlight(selectNeighbor(enabled, highlighted, -1))
			}
		}
	})
	// 3. 搜索词变化时高亮回到第一个匹配项。
	if *filtered != *query {
		if hasMatch {
			parts.Highlight(firstMatch)
		}
		*filtered = *query
	}
	// 4. 搜索框按 Enter：选中高亮项并关闭弹层；没有匹配项时什么也不做。
	if submitted && hasMatch {
		if highlighted == "" {
			highlighted = firstMatch
		}
		*value = highlighted
		trigger.SetChanged()
		*ui.Local(trigger, selectOpenKey, func() bool { return false }) = false
	}
	return trigger
}

// selectOpenKey 是 ui.SelectBase 保存弹层开关的元素本地状态键（MyGo v0.7.3），搜索框 Enter 选中后用它关闭弹层。
const selectOpenKey = "open"

// selectNeighbor 返回 values 里 current 前后第 step 个值；到头时停在头尾，current 不在其中时取第一个。
func selectNeighbor(values []string, current string, step int) string {
	// 1. 没有可选项时返回空串，Highlight 对空串不做任何事。
	if len(values) == 0 {
		return ""
	}
	// 2. 找到当前高亮，按 step 移动并限制在范围内。
	at := -1
	for i, v := range values {
		if v == current {
			at = i
			break
		}
	}
	if at < 0 {
		return values[0]
	}
	return values[max(0, min(at+step, len(values)-1))]
}

// selectVisibleOptions 返回搜索后显示的选项：搜索词匹配标签、值和说明；分组标题只在它下面至少有一项匹配时保留。
func selectVisibleOptions(options []SelectOption, query *string, searchable bool) []SelectOption {
	// 1. 没有搜索时全部显示。
	if !searchable || *query == "" {
		return options
	}
	// 2. 逐项过滤，标题先放进去，随后一项都没有匹配就撤回。
	needle := strings.ToLower(*query)
	var visible []SelectOption
	headerAt := -1
	for _, option := range options {
		if option.Header {
			if headerAt >= 0 {
				visible = visible[:headerAt]
			}
			headerAt = len(visible)
			visible = append(visible, option)
			continue
		}
		if !strings.Contains(strings.ToLower(option.Label+" "+option.Value+" "+option.Description), needle) {
			continue
		}
		headerAt = -1
		visible = append(visible, option)
	}
	if headerAt >= 0 {
		visible = visible[:headerAt]
	}
	return visible
}

// selectOptionText 绘制选项行内容：标签在左，说明在右（12/16 次要色，放不下时截断）。
func selectOptionText(c *ui.Context, t *Theme, option SelectOption) {
	// 1. 标签不收缩，说明占剩余宽度。
	t.TextRole(ui.Text(c, option.Label).SingleLine().Shrink(0), TextUI)
	if option.Description != "" {
		t.TextRole(ui.Text(c, option.Description).SingleLine().Shrink(1).MinWidth(0), TextCaption).TextColor(t.Colors.Muted)
	}
}
