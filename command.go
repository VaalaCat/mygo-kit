package kit

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

// Command 是命令面板里的一条命令；Keys 只用于显示，不注册快捷键。
type Command struct {
	ID, Title, Group string
	Keys             string
	Icon             *ui.SVG
	Disabled         bool
}

// 命令面板的宽度和结果列表的最大高度。
const (
	commandPaletteWidth      = 560
	commandPaletteListHeight = 360
)

type commandPaletteKey struct{}

// commandPalette 保存一次面板打开期间的搜索词、高亮行和按键选定的命令。
type commandPalette struct {
	query             string
	filteredQuery     string
	highlighted       int
	hoveredID         string
	chosenID          string
	revealHighlighted bool
}

// commandGroup 是搜索结果里同一分组的命令，保持命令首次出现的顺序。
type commandGroup struct {
	name     string
	commands []Command
}

// commandPaletteTop 是命令面板距窗口顶部的距离。
const commandPaletteTop float32 = 96

// CommandPalette 显示可搜索的命令列表，选中后返回命令 ID 并关闭；本帧没有选中则返回空字符串。
func CommandPalette(c *ui.Context, open *bool, commands []Command) string {
	// 1. 遮罩把面板固定在窗口居中偏上的位置。
	t := GetTheme(c)
	picked := ""
	ui.DialogBase(c, open, func(back, panel ui.Element) {
		back.Background(ui.RGBA(0, 0, 0, 0.35)).Justify(ui.Start).Padding(commandPaletteTop, 0, 0, 0)
		panel.Width(commandPaletteWidth).MaxWidthPercent(90).Radius(t.Metrics.RadiusLarge).Clip().Background(t.Colors.Elevated).Border(1, t.Colors.Border).Gap(0).Label(text(c, "commandSearch"))
		// 2. 上一次 Enter 已选定命令时，本帧返回它并关闭面板。
		state := ui.Local(panel, commandPaletteKey{}, func() commandPalette { return commandPalette{} })
		if state.chosenID != "" {
			picked = state.chosenID
			*state = commandPalette{}
			*open = false
			return
		}
		if state.query != state.filteredQuery {
			state.filteredQuery, state.highlighted = state.query, 0
		}
		groups, selectable := filterCommands(commands, state.query)
		hasIcons := false
		for _, command := range commands {
			hasIcons = hasIcons || command.Icon != nil
		}
		state.highlighted = min(state.highlighted, max(len(selectable)-1, 0))
		// 3. 搜索框自动聚焦，方向键和 Enter 在按键间隙更新状态。
		ui.Row(c).Height(t.Metrics.HeaderHeight).Shrink(0).Padding(0, t.Space(3)).Gap(t.Space(2)).AlignItems(ui.Center).BorderWidth(0, 0, 1, 0).BorderColor(t.Colors.BorderSubtle).Children(func() {
			ui.Icon(c, IconSearch).Size(IconSizeDefault, IconSizeDefault).Shrink(0).TextColor(t.Colors.Muted)
			input := t.TextRole(ui.TextInputBase(c.Key("command-search"), &state.query), TextUI).Grow(1).Height(t.Metrics.ControlHeight).Padding(0).TextColor(t.Colors.Text).Label(text(c, "commandSearch")).Placeholder(text(c, "commandSearch")).AutoFocus()
			input.OnTextKey(func(key ui.TextKey) bool {
				if key.Mods != 0 {
					return false
				}
				switch key.Key {
				case ui.KeyDown:
					state.highlighted, state.revealHighlighted = min(state.highlighted+1, max(len(selectable)-1, 0)), true
				case ui.KeyUp:
					state.highlighted, state.revealHighlighted = max(state.highlighted-1, 0), true
				case ui.KeyEnter:
					if state.highlighted < len(selectable) {
						state.chosenID = selectable[state.highlighted].ID
					}
				default:
					return false
				}
				return true
			})
		})
		// 4. 命令按组显示，禁用项不可选，键盘移动的高亮行滚入视野。
		ui.Scroll(c).MaxHeight(commandPaletteListHeight).Shrink(1).Children(func() {
			ui.Column(c).Padding(t.Space(1)).Gap(0).Children(func() {
				if len(groups) == 0 {
					t.TextRole(ui.Text(c, text(c, "noResults")), TextUI).Padding(t.Space(2)).TextColor(t.Colors.Muted)
					return
				}
				for _, group := range groups {
					if group.name != "" {
						t.TextRole(ui.Text(c, group.name), TextCaption).TextColor(t.Colors.Muted).Padding(t.Space(2), t.Space(2), t.Space(1), t.Space(2))
					}
					for _, command := range group.commands {
						commandRow(c, state, command, selectable, hasIcons)
					}
				}
			})
		})
		state.revealHighlighted = false
	})
	return picked
}

// commandRow 绘制一条命令；点击立即选定，鼠标移入会把高亮移到该行；reserveIcon 让没有图标的命令与有图标的对齐。
func commandRow(c *ui.Context, state *commandPalette, command Command, selectable []Command, reserveIcon bool) {
	// 1. 高亮行按可选命令的顺序判断，禁用项永远不高亮。
	t := GetTheme(c)
	index := -1
	for i, candidate := range selectable {
		if candidate.ID == command.ID {
			index = i
			break
		}
	}
	row := ui.ButtonBase(c.Key("command-"+command.ID)).Height(t.Metrics.ControlHeight).Padding(0, t.Space(2)).Gap(t.Space(2)).Justify(ui.Start).AlignItems(ui.Center).Radius(t.Metrics.Radius).Label(command.Title).Disabled(command.Disabled)
	if index >= 0 && row.Hovered() && state.hoveredID != command.ID {
		state.hoveredID, state.highlighted = command.ID, index
	}
	highlighted := index >= 0 && index == state.highlighted
	row.Background(controlBackground(&t, controlState{Disabled: command.Disabled, Pressed: row.Pressed(), Hovered: highlighted}))
	row.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if row.FocusVisible() {
			paintFocusRing(p, r, &t, t.Metrics.Radius)
		}
	})
	if highlighted && state.revealHighlighted {
		row.ScrollIntoView()
	}
	if row.Clicked() {
		state.chosenID = command.ID
	}
	// 2. 图标和标题在左，快捷键在右；禁用由框架统一变淡。
	row.TextColor(t.Colors.Text)
	row.Children(func() {
		if command.Icon != nil {
			ui.Icon(c, command.Icon).Size(IconSizeDefault, IconSizeDefault).Shrink(0).TextColor(t.Colors.Muted)
		} else if reserveIcon {
			ui.Box(c).Size(IconSizeDefault, IconSizeDefault).Shrink(0)
		}
		t.TextRole(ui.Text(c, command.Title).SingleLine().Shrink(1), TextUI)
		ui.Spacer(c)
		if command.Keys != "" {
			Kbd(c, command.Keys)
		}
	})
}

// filterCommands 返回匹配搜索词的分组结果，以及按显示顺序排列的可选命令。
func filterCommands(commands []Command, query string) ([]commandGroup, []Command) {
	// 1. 搜索词不区分大小写，只匹配命令标题。
	query = strings.ToLower(strings.TrimSpace(query))
	var groups []commandGroup
	var selectable []Command
	for _, command := range commands {
		if query != "" && !strings.Contains(strings.ToLower(command.Title), query) {
			continue
		}
		// 2. 分组按首次出现的顺序排列。
		position := -1
		for i := range groups {
			if groups[i].name == command.Group {
				position = i
				break
			}
		}
		if position < 0 {
			groups = append(groups, commandGroup{name: command.Group})
			position = len(groups) - 1
		}
		groups[position].commands = append(groups[position].commands, command)
	}
	// 3. 可选命令的顺序必须与屏幕上的行一致。
	for _, group := range groups {
		for _, command := range group.commands {
			if !command.Disabled {
				selectable = append(selectable, command)
			}
		}
	}
	return groups, selectable
}
