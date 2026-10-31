package kit

import (
	"strings"
	"unicode"

	"github.com/egoist/mygo/ui"
)

// menuMinWidth 是菜单面板的最小宽度。
const menuMinWidth = 160

// MenuItem 是菜单的一个条目；Separator 为真时只画一条分隔线，其余字段忽略。
type MenuItem struct {
	// Icon 是条目图标，14；任一条目有图标时，没有图标的条目也留出图标位，标题对齐。
	Icon  *ui.SVG
	Title string
	// Shortcut 是行尾右对齐显示的快捷键提示，例如 "⌘N"，只显示不注册。
	Shortcut string
	// Danger 让标题和图标用危险色。
	Danger bool
	// Disabled 让条目变淡且不可选，方向键、首字母跳转都会跳过它。
	Disabled  bool
	Separator bool
	// OnSelect 在条目被选中、菜单关闭的那一帧调用。
	OnSelect func()
}

type menuButtonKey struct{}
type contextMenuKey struct{}

// menuState 是一个菜单的打开状态，保存在触发元素上；按键处理在帧之间运行，通过它把结果交给下一帧。
type menuState struct {
	open bool
	// highlighted 是高亮条目在 items 里的序号，-1 表示没有。
	highlighted int
	// hoveredItem 是指针当前停留的条目，用来区分「指针移到了新条目」和「指针一直停在原处」，后者不抢键盘的高亮。
	hoveredItem int
	// chosen 是按键选中、等待下一帧执行的条目序号，-1 表示没有。
	chosen int
	// shown 记录上一帧菜单是否已显示，用来在打开的那一帧把焦点交给面板。
	shown bool
	// 右键菜单在指针处打开：pointerOpened 为真时，(pointerX, pointerY) 是指针相对目标元素左上角的位置。
	pointerOpened      bool
	pointerX, pointerY float32
	// selectedText 是右键菜单打开时目标 Selectable 容器里选中的文字，非空时菜单顶部带「复制」。
	selectedText string
}

// newMenuState 返回关闭状态的菜单。
func newMenuState() menuState { return menuState{highlighted: -1, hoveredItem: -1, chosen: -1} }

// MenuButton 给 trigger 附上点击后在下方弹出的菜单：点击触发元素开关菜单，上下键、Home/End、首字母移动高亮，
// Enter 选中，Esc 关闭；返回值表示菜单在本帧关闭（Esc、点击外部或触发元素、Tab、选中条目），焦点回到触发元素，
// 想把焦点放到别处的调用方在返回 true 的那一帧自己设置焦点。
func MenuButton(c *ui.Context, trigger ui.Element, items []MenuItem) (closed bool) {
	// 1. 点击切换开关；键盘聚焦时的点击来自 Enter 或 Space，打开后直接高亮第一项。
	state := ui.Local(trigger, menuButtonKey{}, newMenuState)
	if trigger.Clicked() {
		state.open = !state.open
		state.highlighted = -1
		if state.open && trigger.FocusVisible() {
			state.highlighted = nextMenuItem(items, -1, 1)
		}
	}
	// 2. 菜单显示在触发元素下方。
	trigger.Expanded(state.open)
	return menuPanel(c, trigger, state, items)
}

// ContextMenu 给 target 附上右键菜单：右键时在指针处弹出，目标有焦点时 Shift+F10 或菜单键在目标下方弹出，键盘和选中行为与 MenuButton 相同。
// target 是 Selectable 文字容器时，右键同样弹出这个菜单而不是系统的文字菜单，拖选不受影响；
// 打开时有选中文字，菜单顶部多一项「复制」，复制的就是选中的文字。Selectable 容器必须自己就是 target：挂在它的祖先上，
// 容器内的文字仍然显示系统菜单。返回值含义同 MenuButton：菜单在本帧关闭。
func ContextMenu(c *ui.Context, target ui.Element, items []MenuItem) (closed bool) {
	// 1. 向 MyGo 声明 target 有右键菜单，右键就不再交给文字的系统菜单；这个空菜单没有条目、不会弹出，
	// 它的构建函数只在右键或菜单键唤出菜单时运行，借此得知菜单被唤出。
	state := ui.Local(target, contextMenuKey{}, newMenuState)
	summoned := false
	target.ContextMenu(func(*ui.Menu) { summoned = true })
	// 2. 右键记下指针相对目标的位置和选中的文字。
	if summoned {
		x, y, _ := target.PointerPosition()
		*state = newMenuState()
		state.open, state.pointerOpened, state.pointerX, state.pointerY = true, true, x, y
		state.selectedText = selectedText(target)
	}
	// 3. 键盘唤出时菜单贴着目标下缘，第一项高亮；这个快捷键由目标认领，MyGo 不再替它唤出空菜单。
	if target.Shortcut(ui.Shift, ui.KeyF10) || target.Shortcut(0, ui.KeyContextMenu) {
		*state = newMenuState()
		state.open = true
		state.selectedText = selectedText(target)
		state.highlighted = nextMenuItem(contextMenuItems(c, state, items), -1, 1)
	}
	menuItems := contextMenuItems(c, state, items)
	// 4. 指针处打开时，菜单锚在目标内一个零尺寸的点上；其余情况以目标为锚点。
	anchor := target
	if state.open && state.pointerOpened {
		ui.Overlay(c, func() {
			anchor = ui.Box(c).Key("kit-context-menu-point").Size(0, 0).AttachTo(target, ui.AnchorTopLeft, ui.AnchorTopLeft).Margin(state.pointerY, 0, 0, state.pointerX)
		})
	}
	return menuPanel(c, anchor, state, menuItems)
}

// contextMenuItems 返回右键菜单的条目：打开时有选中文字就在顶部加「复制」和分隔线。
func contextMenuItems(c *ui.Context, state *menuState, items []MenuItem) []MenuItem {
	// 1. 没有选中文字时就是调用方的条目。
	if state.selectedText == "" {
		return items
	}
	// 2. 复制的是打开菜单那一刻选中的文字。
	copyText := state.selectedText
	copyItem := MenuItem{Icon: IconCopy, Title: text(c, "copy"), OnSelect: func() { c.WriteClipboard(copyText) }}
	return append([]MenuItem{copyItem, {Separator: true}}, items...)
}

// selectedText 返回 Selectable 容器当前选中的文字，段落之间用换行连接；没有选中时返回空串。
func selectedText(target ui.Element) string {
	// 1. 每段取选区覆盖的字符区间。
	var parts []string
	for _, part := range target.SelectedTexts() {
		parts = append(parts, string([]rune(part.Text)[part.Start:part.End]))
	}
	return strings.Join(parts, "\n")
}

// menuPanel 绘制菜单面板和条目，处理选中、悬停高亮和键盘；面板打开时拿到焦点，关闭时框架把焦点还给原来的元素。
// 返回值表示上一帧显示的面板在本帧关闭了。
func menuPanel(c *ui.Context, anchor ui.Element, state *menuState, items []MenuItem) (closed bool) {
	// 1. 上一帧按键已选定条目：先关闭菜单，再在本帧执行。
	t := GetTheme(c)
	picked := state.chosen
	state.chosen = -1
	if picked >= 0 {
		state.open = false
	}
	if !state.open {
		state.highlighted, state.hoveredItem = -1, -1
	}
	hasIcon := false
	for _, item := range items {
		hasIcon = hasIcon || item.Icon != nil
	}
	// 2. 面板：自己绘制外观，Esc 和点击外部由 PopoverBase 关闭。
	panel := ui.PopoverBase(c, anchor, &state.open, func(panel ui.Element) {
		panel.Column().Role(ui.RoleMenu).Focusable().FocusRing(false).Gap(0).Padding(t.Space(1)).Margin(t.Space(1), 0, 0, 0).MinWidth(menuMinWidth).Radius(t.Metrics.RadiusLarge).Background(t.Colors.Elevated).Border(1, t.Colors.Border)
		pointerOnItem := false
		for i, item := range items {
			// 3. 分隔线上下各留 4，不参与导航。
			if item.Separator {
				ui.Box(c).Height(1).Shrink(0).Margin(t.Space(1), 0).Background(t.Colors.BorderSubtle)
				continue
			}
			row := ui.Row(c.Key(i)).Height(t.Metrics.ControlHeight).Shrink(0).Padding(0, t.Space(2)).Gap(t.Space(2)).AlignItems(ui.Center).Radius(t.Metrics.Radius).Role(ui.RoleMenuItem).Label(item.Title).Disabled(item.Disabled)
			// 4. 指针移到新条目时跟随高亮；一直停着不动则让键盘决定高亮。
			if !item.Disabled && row.Hovered() {
				pointerOnItem = true
				if state.hoveredItem != i {
					state.hoveredItem, state.highlighted = i, i
				}
			}
			highlighted := state.highlighted == i
			row.Background(controlBackground(&t, controlState{Disabled: item.Disabled, Pressed: row.Pressed(), Hovered: highlighted}))
			if highlighted {
				panel.ActiveDescendant(row)
			}
			if row.Clicked() {
				picked = i
				state.open = false
			}
			// 5. 图标位、标题、右对齐的快捷键；危险项整行用危险色。
			color := t.Colors.Text
			if item.Danger {
				color = t.Colors.Danger
			}
			row.TextColor(color)
			row.Children(func() {
				switch {
				case item.Icon != nil:
					iconSlot(c, item.Icon, menuIconColor(&t, item))
				case hasIcon:
					ui.Box(c).Size(IconSizeLarge, IconSizeLarge).Shrink(0)
				}
				t.TextRole(ui.Text(c, item.Title).SingleLine().Grow(1).Shrink(1).MinWidth(0), TextUI)
				if item.Shortcut != "" {
					t.TextRole(ui.Text(c, item.Shortcut).SingleLine().Shrink(0), TextCaption).TextColor(t.Colors.Muted)
				}
			})
		}
		if !pointerOnItem {
			state.hoveredItem = -1
		}
		panel.HandleInput(func(ev ui.InputEvent) bool { return handleMenuInput(state, items, ev) })
	})
	// 6. 面板在打开的那一帧拿到焦点；这要在浮层记下它打开时的焦点之后，关闭时框架才把焦点还给触发元素。
	if panel.Valid() && !state.shown {
		panel.Focus()
	}
	closed = state.shown && !panel.Valid()
	state.shown = panel.Valid()
	// 7. 条目被选中的那一帧调用 OnSelect。
	if picked >= 0 && picked < len(items) && items[picked].OnSelect != nil {
		items[picked].OnSelect()
	}
	return closed
}

// menuIconColor 返回条目图标颜色：危险项用危险色，其余用次要色。
func menuIconColor(t *Theme, item MenuItem) ui.Color {
	if item.Danger {
		return t.Colors.Danger
	}
	return t.Colors.Muted
}

// handleMenuInput 处理菜单面板收到的按键和文字，返回是否已处理；它在帧之间运行，只修改 state，条目的选中留给下一帧执行。
func handleMenuInput(state *menuState, items []MenuItem, ev ui.InputEvent) bool {
	// 1. 方向键、Home/End、Enter、Space 和 Tab 只在没有修饰键时响应。
	switch ev.Kind {
	case ui.InputKeyDown:
		// Shift+Tab 与 Tab 一样关闭菜单，焦点交给框架。
		if ev.Key == ui.KeyTab && ev.Mods == ui.Shift {
			state.open = false
			return true
		}
		if ev.Mods != 0 {
			return false
		}
		switch ev.Key {
		case ui.KeyDown:
			state.highlighted = nextMenuItem(items, state.highlighted, 1)
		case ui.KeyUp:
			state.highlighted = nextMenuItem(items, state.highlighted, -1)
		case ui.KeyHome:
			state.highlighted = nextMenuItem(items, -1, 1)
		case ui.KeyEnd:
			state.highlighted = nextMenuItem(items, -1, -1)
		case ui.KeyEnter, ui.KeySpace:
			if state.highlighted >= 0 {
				state.chosen = state.highlighted
			}
		case ui.KeyTab:
			state.open = false
		default:
			return false
		}
		return true
	// 2. 首字母跳转：从当前高亮之后找标题以该字母开头的可选条目，到末尾回到开头。
	case ui.InputText:
		letter, ok := firstLetter(ev.Text)
		if !ok {
			return false
		}
		for step := 1; step <= len(items); step++ {
			i := (max(state.highlighted, -1) + step + len(items)) % len(items)
			item := items[i]
			if item.Separator || item.Disabled {
				continue
			}
			if first, ok := firstLetter(item.Title); ok && first == letter {
				state.highlighted = i
				return true
			}
		}
		return false
	}
	return false
}

// nextMenuItem 返回从 from 出发、沿 direction（1 向下、-1 向上）的下一个可选条目序号，跳过禁用项和分隔线，到头后绕回；
// from 为 -1 时 direction 为 1 取第一项、为 -1 取最后一项；没有可选条目时返回 -1。
func nextMenuItem(items []MenuItem, from, direction int) int {
	// 1. 最多走一圈。
	n := len(items)
	for step := 1; step <= n; step++ {
		var i int
		if from < 0 && direction < 0 {
			i = n - step
		} else {
			i = ((from+direction*step)%n + n) % n
		}
		if !items[i].Separator && !items[i].Disabled {
			return i
		}
	}
	return -1
}

// firstLetter 返回文字的第一个非空白字符的小写形式。
func firstLetter(s string) (rune, bool) {
	for _, r := range strings.TrimSpace(s) {
		return unicode.ToLower(r), true
	}
	return 0, false
}
