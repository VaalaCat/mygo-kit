package kit

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

// completionMaxRows 是补全浮层同时显示的最多行数，超出后滚动。
const completionMaxRows = 8

// CompletionItem 是补全浮层的一行候选。
type CompletionItem struct {
	Icon  *ui.SVG
	Title string
	// Description 显示在标题右侧，次要色，放不下时截断。
	Description string
}

// CompletionOpt 是补全浮层的数据和状态。
type CompletionOpt struct {
	// Open 为真时才显示浮层；调用方按输入内容决定，例如光标所在的词以 / 开头。
	Open  bool
	Items []CompletionItem
	// Loading 在候选后面（没有候选时单独）显示加载行。
	Loading bool
	// Error 非空时显示一行错误。
	Error string
	// Empty 是没有候选、没有加载也没有错误时的提示；为空用库内置的「没有匹配项」。
	Empty string
}

type completionKey struct{}

// completionState 是浮层的交互状态，保存在锚点元素上；按键处理在帧之间运行，通过它把结果交给下一帧。
type completionState struct {
	// highlighted 是高亮候选的序号。
	highlighted int
	// chosen 是按键选中、等待下一帧返回的候选序号，-1 表示没有。
	chosen int
	// dismissed 为真表示用户按 Esc 关闭了浮层，直到 Open 变回 false 才会再次出现。
	dismissed bool
	// revealHighlighted 让下一帧把高亮行滚入视野。
	revealHighlighted bool
	// hoveredItem 是指针停留的候选，用来区分「指针移到了新候选」和「指针一直停在原处」。
	hoveredItem int
	// signature 是候选标题的拼接，变化时高亮回到第一项。
	signature string
	// isComposing 为真表示输入法正在组字，按键归输入法所有，浮层不接管；每帧由输入框的状态刷新。
	isComposing bool
}

// newCompletionState 返回没有选中任何候选的初始状态。
func newCompletionState() completionState { return completionState{chosen: -1, hoveredItem: -1} }

// CompletionMenu 在 anchor 上方显示补全浮层，输入框 input 保持焦点；返回本帧选中的候选序号，没有选中时返回 -1。
// 有候选时接管 Up/Down（循环移动）、Enter 和 Tab（选中），Esc 关闭到 Open 变回 false；没有候选时只接管 Esc；输入法组字期间不接管任何键。
// 最多显示 8 行（候选、加载行、错误行合计），超出滚动，行内容的左缘与 input 里的文字对齐。
// 接管按键用的是 input 上的 HandleInput，调用方不要再给同一个输入框设置 HandleInput；OnTextKey 不受影响。
func CompletionMenu(c *ui.Context, anchor, input ui.Element, opt CompletionOpt) (selected int) {
	// 1. 状态随锚点保存；关闭时清空，让下一次打开从头开始。
	t := GetTheme(c)
	state := ui.Local(anchor, completionKey{}, newCompletionState)
	selected = state.chosen
	state.chosen = -1
	if !opt.Open {
		*state = newCompletionState()
		return -1
	}
	if selected >= 0 {
		state.dismissed = true
	}
	if state.dismissed {
		return selected
	}
	// 2. 候选变化时高亮回到第一项；高亮序号始终落在候选范围内。
	var titles strings.Builder
	for _, item := range opt.Items {
		titles.WriteString(item.Title + "\n")
	}
	if state.signature != titles.String() {
		state.signature, state.highlighted = titles.String(), 0
	}
	state.highlighted = min(state.highlighted, max(len(opt.Items)-1, 0))
	// 3. 键盘：编辑器之前先让浮层看到按键；没有候选时只接管 Esc。
	count := len(opt.Items)
	// 手测：Linux 拼音输入法。1. 在输入框输入 /，浮层出现；2. 切到拼音输入法，输入 ni 不上屏，按 Enter 应确认拼音候选，浮层不选中高亮项；
	// 3. 再按 Enter（组字结束后）才选中高亮项；4. 组字时按 Up/Down/Tab 应归输入法所有，浮层高亮不动。离屏测试器不模拟输入法吞键，只能手测。
	state.isComposing = input.Composing()
	input.HandleInput(func(ev ui.InputEvent) bool { return !state.isComposing && handleCompletionInput(state, count, ev) })
	// 4. 浮层贴在锚点上方，宽度与锚点相同，点击候选不抢走输入框焦点；
	// 行内容的起点与输入框里输入的文字对齐：输入框相对锚点的偏移加上它自己的文字内缩（本库输入框记录在元素上），再减去浮层的 1 边线。
	start := t.Space(3)
	if box, field := anchor.Bounds(), input.Bounds(); box.W > 0 && field.W > 0 {
		inset := *ui.Local(input, inputTextInsetKey{}, func() float32 { return 0 })
		start = max(field.X+inset-box.X-1, 0)
	}
	// 候选之后的加载行、错误行和没有候选时的提示各占一行，滚动区高度要把它们算进去。
	rowCount := count
	if opt.Loading {
		rowCount++
	}
	if opt.Error != "" {
		rowCount++
	}
	rowCount = min(max(rowCount, 1), completionMaxRows)
	ui.Overlay(c, func() {
		panel := ui.Column(c).Key("kit-completion").Role(ui.RoleList).AttachTo(anchor, ui.AnchorTopStart, ui.AnchorBottomStart).Margin(0, 0, t.Space(1), 0).KeepFocus().Gap(0).Padding(t.Space(1), 0).Radius(t.Metrics.RadiusLarge).Clip().Background(t.Colors.Elevated).Border(1, t.Colors.Border)
		if width := anchor.Bounds().W; width > 0 {
			panel.Width(width)
		}
		panel.Children(func() {
			ui.Scroll(c).KeepFocus().MaxHeight(float32(rowCount) * t.Metrics.ControlHeight).Children(func() {
				ui.Column(c).Gap(0).Children(func() { completionRows(c, state, opt, start, &selected) })
			})
		})
	})
	return selected
}

// completionRows 绘制候选行，其后依次是加载、错误行；没有任何内容时显示空提示；start 是行内容左侧的内边距。
func completionRows(c *ui.Context, state *completionState, opt CompletionOpt, start float32, selected *int) {
	// 1. 候选行：指针移到新候选时跟随高亮，点击选中。
	t := GetTheme(c)
	hasIcon := false
	for _, item := range opt.Items {
		hasIcon = hasIcon || item.Icon != nil
	}
	pointerOnItem := false
	for i, item := range opt.Items {
		row := ui.Row(c.Key(i)).KeepFocus().Height(t.Metrics.ControlHeight).Shrink(0).Padding(0, t.Space(3), 0, start).Gap(t.Space(2)).AlignItems(ui.Center).Radius(0).Role(ui.RoleListItem).Label(item.Title)
		if row.Hovered() {
			pointerOnItem = true
			if state.hoveredItem != i {
				state.hoveredItem, state.highlighted = i, i
			}
		}
		highlighted := state.highlighted == i
		// Selected 会换成强调色底和文字，所以先设置它，再用控件状态底色和正文色覆盖。
		row.Selected(highlighted).TextColor(t.Colors.Text)
		row.Background(controlBackground(&t, controlState{Pressed: row.Pressed(), Hovered: highlighted}))
		if highlighted && state.revealHighlighted {
			row.ScrollIntoView()
		}
		if row.Clicked() {
			*selected = i
			state.dismissed = true
		}
		// 标题和说明按基线对齐，图标位在行内垂直居中。
		row.Children(func() {
			switch {
			case item.Icon != nil:
				iconSlot(c, item.Icon, t.Colors.Muted)
			case hasIcon:
				ui.Box(c).Size(IconSizeLarge, IconSizeLarge).Shrink(0)
			}
			ui.Row(c).Grow(1).Shrink(1).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Baseline).Children(func() {
				t.TextRole(ui.Text(c, item.Title).SingleLine().Shrink(0), TextUI)
				if item.Description != "" {
					t.TextRole(ui.Text(c, item.Description).SingleLine().Shrink(1).MinWidth(0), TextCaption).TextColor(t.Colors.Muted)
				}
			})
		})
	}
	if !pointerOnItem {
		state.hoveredItem = -1
	}
	state.revealHighlighted = false
	// 2. 加载、错误和空提示都是一整行，样式与列表里的 InlineHint、FieldMessage 一致。
	if opt.Loading {
		InlineHint(c, "", true).PaddingStart(start)
	}
	if opt.Error != "" {
		ui.Row(c).Height(t.Metrics.ControlHeight).Shrink(0).Padding(0, t.Space(3), 0, start).AlignItems(ui.Center).Children(func() {
			FieldMessage(c, MessageError, opt.Error).SingleLine()
		})
	}
	if len(opt.Items) == 0 && !opt.Loading && opt.Error == "" {
		empty := opt.Empty
		if empty == "" {
			empty = text(c, "noResults")
		}
		InlineHint(c, empty, false).PaddingStart(start)
	}
}

// handleCompletionInput 处理输入框收到的按键，返回是否已被浮层接管；它在帧之间运行，只修改 state。
func handleCompletionInput(state *completionState, count int, ev ui.InputEvent) bool {
	// 1. 只处理没有修饰键的按键；Esc 在任何状态下都关闭浮层。
	if ev.Kind != ui.InputKeyDown || ev.Mods != 0 {
		return false
	}
	if ev.Key == ui.KeyEscape {
		state.dismissed = true
		return true
	}
	// 2. 其余键只在有候选时接管，没有候选时交给编辑器（光标移动、Enter 发送）。
	if count == 0 {
		return false
	}
	switch ev.Key {
	case ui.KeyDown:
		state.highlighted = (state.highlighted + 1) % count
		state.revealHighlighted = true
	case ui.KeyUp:
		state.highlighted = (state.highlighted + count - 1) % count
		state.revealHighlighted = true
	case ui.KeyEnter, ui.KeyTab:
		state.chosen = state.highlighted
	default:
		return false
	}
	return true
}
