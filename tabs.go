package kit

import (
	"github.com/egoist/mygo/ui"
)

// paintTab 绘制页签的底色和分隔线，Tabs 与 Tab 共用：选中时用正文底色、断开底线，与下方内容连成一片；
// 未选中时用面板底色，悬停和按下换成控件状态底色，并画底线。右侧始终是竖线。
func paintTab(p *ui.Painter, r ui.Rect, t *Theme, selected bool, state controlState) {
	// 1. 底色：选中用正文底色，其余按状态取色，没有状态时是面板底色。
	bg := controlBackground(t, state)
	if selected {
		bg = t.Colors.Editor
	} else if bg == ui.Transparent {
		bg = t.Colors.Panel
	}
	p.Fill(r, bg, 0)
	// 2. 分隔线：右侧竖线，未选中再加底线。
	p.Fill(ui.Rect{X: r.X + r.W - 1, Y: r.Y, W: 1, H: r.H}, t.Colors.Border, 0)
	if !selected {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, t.Colors.Border, 0)
	}
}

// tabsSelectionKey 标识页签条上一帧的选中序号。
type tabsSelectionKey struct{}

// Tabs 绘制连续页签条，方向键导航由 MyGo TabsBase 提供。
// 焦点环由页签自己在边界内绘制（见 paintFocusRing），只在键盘聚焦时显示，并始终跟随选中的页签。
func Tabs(c *ui.Context, selected *int, labels []string) ui.Element {
	// 1. 页签之间用竖线分隔，选中项取正文底色并断开底线，与下方内容连成一片。
	t := GetTheme(c)
	parts := ui.TabsBase(c, selected, len(labels))
	parts.List.Height(t.Metrics.HeaderHeight).Gap(0).Background(t.Colors.Panel)
	tabs := make([]ui.Element, len(labels))
	parts.List.Children(func() {
		for i, label := range labels {
			tab := parts.Tab(i).Height(t.Metrics.HeaderHeight).Padding(0, t.Space(3)).Radius(0).AlignItems(ui.Center).Label(label).FocusRing(false)
			tabs[i] = tab
			active := i == *selected
			if active {
				tab.TextColor(t.Colors.Text)
			} else {
				tab.TextColor(t.Colors.Muted)
			}
			tab.Draw(func(p *ui.Painter, r ui.Rect) {
				paintTab(p, r, &t, active, controlState{Pressed: tab.Pressed(), Hovered: tab.Hovered()})
			})
			tab.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if tab.FocusVisible() {
					paintFocusRing(p, r, &t, 0)
				}
			})
			tab.Children(func() { t.TextRole(ui.Text(c, label).SingleLine(), TextUI) })
		}
		ui.Box(c).Grow(1).Height(t.Metrics.HeaderHeight).BorderWidth(0, 0, 1, 0).BorderColor(t.Colors.Border)
	})
	// 2. 选中项被点击、方向键以外的方式改变时，停在旧页签上的焦点跟到选中的页签，环不会留在原处。
	previous := ui.Local(parts.List, tabsSelectionKey{}, func() int { return *selected })
	if *previous != *selected && *selected < len(tabs) {
		for i, tab := range tabs {
			if i != *selected && tab.Focused() {
				tabs[*selected].Focus()
				break
			}
		}
	}
	*previous = *selected
	return parts.List
}

// tabSelectedKey 标识标签上一帧是否处于选中状态。
type tabSelectedKey struct{}

// TabOpt 指定一个可关闭标签的图标、标题和状态。
type TabOpt struct {
	Icon     *ui.SVG
	Title    string
	Selected bool
	// Faded 让标题用弱化色，例如进程已结束的终端。
	Faded bool
	// CloseLabel 是关闭按钮的可访问名称。
	CloseLabel string
}

// tabMaxWidth 是标签标题的最大宽度，超出用省略号截断。
const tabMaxWidth float32 = 192

// Tab 绘制 Zed 风格的可关闭标签，返回本帧是否点了标签和关闭按钮。
// 标签与所在标题栏同高、无圆角：选中时用正文底色、左右竖线、不画底线，与下方内容连成一片；
// 未选中时用面板底色、右侧竖线和底线。关闭按钮在悬停或选中时出现。
func Tab(c *ui.Context, opt TabOpt) (clicked, closed bool) {
	// 1. 外框按选中状态画底色和竖线、底线；新选中的标签（包括刚打开的）滚进所在的标签条。
	t := GetTheme(c)
	row := ui.Row(c).Height(t.Metrics.HeaderHeight).Shrink(0).Padding(0, t.Space(1), 0, t.Space(3)).Gap(t.Space(1)).AlignItems(ui.Center)
	var label ui.Element
	hovered := row.Hovered()
	wasSelected := ui.Local(row, tabSelectedKey{}, func() bool { return false })
	if opt.Selected && !*wasSelected {
		row.ScrollIntoView()
	}
	*wasSelected = opt.Selected
	row.Draw(func(p *ui.Painter, r ui.Rect) {
		paintTab(p, r, &t, opt.Selected, controlState{Hovered: row.Hovered()})
	})
	row.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if label.FocusVisible() {
			paintFocusRing(p, r, &t, 0)
		}
	})
	// 2. 图标和标题是一个按钮，关闭按钮是它旁边的独立按钮。
	row.Children(func() {
		color := t.Colors.Muted
		if opt.Selected || hovered {
			color = t.Colors.Text
		}
		if opt.Faded {
			color = t.Colors.Placeholder
		}
		label = ui.ButtonBase(c).Height(t.Metrics.HeaderHeight).Gap(t.Space(2)).AlignItems(ui.Center).Radius(0).Label(opt.Title).TextColor(color).FocusRing(false)
		label.Children(func() {
			if opt.Icon != nil {
				ui.Icon(c, opt.Icon).Size(IconSizeDefault, IconSizeDefault).Shrink(0)
			}
			t.TextRole(ui.Text(c, opt.Title), TextUI).SingleLine().Ellipsis("…").MaxWidth(tabMaxWidth)
		})
		clicked = label.Clicked()
		closeButton := SmallIconButton(c, IconX, opt.CloseLabel)
		if !opt.Selected && !hovered && !closeButton.FocusVisible() {
			closeButton.Opacity(0)
		}
		closed = closeButton.Clicked()
	})
	return clicked, closed
}
