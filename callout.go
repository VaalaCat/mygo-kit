package kit

import (
	"github.com/egoist/mygo/ui"
)

// CalloutKind 是提示条的严重程度。
type CalloutKind string

// 提示条的种类。
const (
	CalloutError   CalloutKind = "error"
	CalloutWarning CalloutKind = "warning"
	CalloutInfo    CalloutKind = "info"
)

// Callout 显示带色调的提示条：图标、标题、次要色详情，右侧是 actions 画出的操作。
// detail 为空时只显示标题，actions 为 nil 时不留操作区。
func Callout(c *ui.Context, kind CalloutKind, title, detail string, actions func()) ui.Element {
	// 1. 按种类选择图标和色调。
	t := GetTheme(c)
	tone, icon := t.Colors.Accent, IconInfo
	switch kind {
	case CalloutError:
		tone, icon = t.Colors.Danger, IconAlert
	case CalloutWarning:
		tone, icon = t.Colors.Warning, IconTriangleAlert
	}
	// 2. 淡色背景加混入色调的细边框。
	root := ui.Row(c).Padding(t.Space(3)).Gap(t.Space(2)).AlignItems(ui.Start).Radius(t.Metrics.RadiusLarge).
		Background(tone.Alpha(.1)).Border(1, t.Colors.Border.Mix(tone, .5)).Label(title)
	if detail == "" {
		root.AlignItems(ui.Center)
	}
	root.Children(func() {
		leadingIcon(c, icon, tone)
		ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			t.TextRole(ui.Text(c, title), TextUI).TextColor(t.Colors.Text)
			if detail != "" {
				t.TextRole(ui.Text(c, detail), TextCaption).TextColor(t.Colors.Muted)
			}
		})
		if actions != nil {
			ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Children(actions)
		}
	})
	return root
}

// leadingIcon 绘制通知和提示条行首的状态图标：14 的图标在 20 高的盒子里居中，与 TextUI 的首行垂直对齐。
func leadingIcon(c *ui.Context, icon *ui.SVG, color ui.Color) {
	// 1. 盒子高度取紧凑档，等于一行 TextUI 的行高。
	t := GetTheme(c)
	ui.Row(c).Height(t.Metrics.CompactHeight).Shrink(0).AlignItems(ui.Center).Children(func() {
		ui.Icon(c, icon).Size(IconSizeDefault, IconSizeDefault).TextColor(color)
	})
}
