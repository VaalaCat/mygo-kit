package kit

import (
	"github.com/egoist/mygo/ui"
)

// MessageScrollerOpt 指定消息滚动区的可选通知。
type MessageScrollerOpt struct {
	// OnReachTop 在用户滚到列表顶部时调用一次，调用方据此加载更早的历史；
	// 离开顶部后再次回到顶部才会再次调用。
	OnReachTop func()
}

type scrollerKey struct{}
type scrollerState struct{ atTop bool }

// MessageScroller 创建跟随最新消息的虚拟列表；用户上滚后右下角浮出「回到底部」按钮，opt 指定滚到顶部的通知。
// state 由调用方持有，建议设置 Key 以便前插历史时保持位置。
func MessageScroller(c *ui.Context, state *ui.ListState, count int, opt MessageScrollerOpt, row func(int)) ui.Element {
	// 1. 列表跟随尾部；允许动画时新增内容从下方滑入。
	t := GetTheme(c)
	state.FollowEnd = true
	if state.FollowAnimation == nil && !t.Motion.Reduced {
		follow := ui.Smooth
		state.FollowAnimation = &follow
	}
	root := ui.Column(c).Grow(1).MinHeight(0)
	root.Children(func() {
		ui.List(c, state, count, row).Grow(1).Gap(t.Space(4)).Padding(t.Space(4)).Radius(0)
		first, last := state.Visible()
		laidOut := last >= first
		// 2. 从未到顶变为到顶时通知一次。
		top := ui.Local(root, scrollerKey{}, func() scrollerState { return scrollerState{} })
		if laidOut && first == 0 {
			if !top.atTop && opt.OnReachTop != nil {
				opt.OnReachTop()
			}
			top.atTop = true
		} else if laidOut {
			top.atTop = false
		}
		// 3. 不在尾部时浮出圆形按钮，点击回到底部。
		if laidOut && count > 0 && !state.AtEnd() {
			holder := ui.Box(c).Absolute().Right(t.Space(4)).Bottom(t.Space(4)).Size(t.Metrics.ControlHeight, t.Metrics.ControlHeight).Radius(RadiusFull).Clip().Background(t.Colors.Elevated).Border(1, t.Colors.Border)
			holder.Children(func() {
				if IconButton(c, IconArrowDown, text(c, "jumpToLatest")).Radius(RadiusFull).Clicked() {
					state.ScrollToEnd()
				}
			})
		}
	})
	return root
}
