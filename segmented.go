package kit

import (
	"github.com/egoist/mygo/ui"
)

// SegmentedItem 是分段控件的一段：Icon 可为空，Label 同时是可访问名称，Tooltip 可为空。
type SegmentedItem struct {
	Icon    *ui.SVG
	Label   string
	Tooltip string
}

// SegmentedOpt 指定分段控件的宽度方式。
type SegmentedOpt struct {
	// FullWidth 让整条控件占满所在容器宽度、各段等宽；默认各段贴合自己的内容。
	FullWidth bool
}

// Segmented 创建只有文字的分段控件，宽度贴合内容；见 SegmentedWith。
func Segmented(c *ui.Context, selected *int, labels []string) ui.Element {
	// 1. 文字转成没有图标和提示的分段，交给 SegmentedWith。
	items := make([]SegmentedItem, len(labels))
	for i, label := range labels {
		items[i] = SegmentedItem{Label: label}
	}
	return SegmentedWith(c, selected, items, SegmentedOpt{})
}

// SegmentedWith 创建 Zed ToggleButtonGroup 外观（Outlined）的分段控件：28 高，一圈 60% 透明度的边线，
// 圆角裁剪，各段之间 1 像素竖线分隔；*selected 是选中段序号，选中段用强调色浅底和强调色文字、图标，
// 其余段透明底，悬停和按下取控件状态底色。默认宽度贴合内容，放进纵向容器时调用方加 AlignSelf(ui.Start) 才靠左；
// FullWidth 时占满容器、各段等宽。方向键切换，整体是一个 Tab 停靠点，调用方用 Label 给出可访问名称。
func SegmentedWith(c *ui.Context, selected *int, items []SegmentedItem, opt SegmentedOpt) ui.Element {
	// 1. 复用 SegmentedBase 的单选组、方向键与无障碍语义；轨道画外框、圆角裁剪和段间竖线，竖线占段间的 1 像素间隙。
	t := GetTheme(c)
	border := t.Colors.Border.Alpha(segmentedBorderAlpha)
	parts := ui.SegmentedBase(c, selected, len(items))
	parts.Track.Height(t.Metrics.ControlHeight).Padding(0).Gap(1).AlignItems(ui.Stretch).Radius(t.Metrics.Radius).Clip().Border(1, border).Dividers(1, border)
	if opt.FullWidth {
		parts.Track.WidthPercent(100)
	}
	// 2. 每段是一个单选按钮：图标（12）加文字居中；全宽时各段平分，否则贴合内容。
	parts.Track.Children(func() {
		for i, item := range items {
			active := i == *selected
			segment := parts.Segment(i).Padding(0, t.Space(2)).Gap(t.Space(1)).Radius(0).Justify(ui.Center).AlignItems(ui.Center).Label(item.Label).FocusRing(false)
			if opt.FullWidth {
				segment.Basis(0).Grow(1).Shrink(1).MinWidth(0)
			} else {
				segment.Shrink(0)
			}
			// 3. 底色按状态取；选中段在状态底色上再叠强调色浅底。
			segment.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, controlBackground(&t, controlState{Disabled: segment.IsDisabled(), Pressed: segment.Pressed(), Hovered: segment.Hovered()}), 0)
				if active {
					p.Fill(r, t.Colors.AccentSubtle, 0)
				}
			})
			// 4. 键盘聚焦时在段内画 2 像素焦点环，轨道的圆角裁剪会切掉两端外角。
			segment.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if segment.FocusVisible() {
					paintFocusRing(p, r, &t, 0)
				}
			})
			color := t.Colors.Text
			if active {
				color = t.Colors.Accent
			}
			segment.TextColor(color)
			segment.Children(func() {
				if item.Icon != nil {
					iconColor := t.Colors.Muted
					if active {
						iconColor = t.Colors.Accent
					}
					ui.Icon(c, item.Icon).Size(IconSizeSmall, IconSizeSmall).Shrink(0).TextColor(iconColor)
				}
				t.TextRole(ui.Text(c, item.Label).SingleLine(), TextUI)
			})
			if item.Tooltip != "" {
				Tooltip(c, segment, item.Tooltip)
			}
		}
	})
	return parts.Track
}

// segmentedBorderAlpha 是分段控件外框和段间竖线在边线色上的不透明度，与 Zed 一致。
const segmentedBorderAlpha = 0.6
