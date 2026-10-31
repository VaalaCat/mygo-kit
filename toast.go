package kit

import (
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
)

// ToastKind 选择通知的语义图标和颜色。
type ToastKind string

// 通知语义：信息、成功和错误。
const (
	ToastInfo    ToastKind = "info"
	ToastSuccess ToastKind = "success"
	ToastError   ToastKind = "error"
)

// ToastOpt 描述一条通知；Action 非空时显示按钮，点击后执行 OnAction 并关闭通知。
type ToastOpt struct {
	Title       string
	Description string
	Kind        ToastKind
	Action      string
	OnAction    func()
}

// ShowToast 在窗口右上角显示一条通知；需要窗口里已有 ToastHost。空闲窗口也会立即出现，调用方不用自己请求重绘。
func ShowToast(c *ui.Context, opt ToastOpt) {
	// 1. 语义类型交给 ToastHost 取外观，错误保持更久。
	toast := ui.Toast{Title: opt.Title, Description: opt.Description, Type: string(opt.Kind), Action: opt.Action, OnAction: opt.OnAction}
	if opt.Kind == ToastError {
		toast.Timeout = 8 * time.Second
	}
	c.AddToast(toast)
	// 2. 调用时 ToastHost 可能已经在本帧画过，请求下一帧让通知出现；到期关闭由 ToastHost 里的定时唤醒负责。
	c.AnimationFrame()
}

// toastWidth 是通知的宽度。
const toastWidth = 360

// ToastHost 在窗口右上角、顶栏下方绘制通知，最新的在最上面；每个窗口的视图里调用一次。
func ToastHost(c *ui.Context) ui.Element {
	// 1. 视口覆盖窗口并让指针穿过；通知从 TitleBar 下缘（没有 TitleBar 时从窗口顶部）留 Space(2) 起向下排，右缘同样留 Space(2)。
	t := GetTheme(c)
	top := getTitleBarHeight(c) + t.Space(2)
	return ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Padding(top, t.Space(2), 0, t.Space(2)).Justify(ui.Start).AlignItems(ui.End).Gap(t.Space(2))
		for _, toast := range slices.Backward(toasts) {
			// 2. 每条通知是带语义图标、标题、可选说明、动作和关闭按钮的浮层。
			parts := ui.ToastBase(c, toast)
			parts.Root.Row().Width(toastWidth).MaxWidthPercent(100).Gap(t.Space(2)).Padding(t.Space(3)).AlignItems(ui.Start).Radius(t.Metrics.RadiusLarge).Background(t.Colors.Elevated).Border(1, t.Colors.Border).Shadow(0, 4, 16, 0, ui.RGBA(0, 0, 0, 0.2))
			parts.Root.Children(func() {
				icon, color := IconInfo, t.Colors.Accent
				switch ToastKind(toast.Type) {
				case ToastSuccess:
					icon, color = IconCircleCheck, t.Colors.Success
				case ToastError:
					icon, color = IconCircleX, t.Colors.Danger
				}
				leadingIcon(c, icon, color)
				ui.Column(c).Grow(1).Shrink(1).Gap(t.Space(1)).Children(func() {
					t.TextRole(ui.Text(c, toast.Title), TextUI)
					if toast.Description != "" {
						t.TextRole(ui.Text(c, toast.Description), TextCaption).TextColor(t.Colors.Muted)
					}
					if toast.Action != "" {
						action := parts.ActionButton().AlignSelf(ui.Start).Padding(0).Margin(t.Space(1), 0, 0, 0).Radius(0).TextColor(t.Colors.Accent)
						action.Children(func() {
							text := t.TextRole(ui.Text(c, toast.Action).SingleLine(), TextUI)
							if action.Hovered() {
								text.Underline()
							}
						})
					}
				})
				closeButton := parts.CloseButton().Label(text(c, "close")).Size(t.Metrics.CompactHeight, t.Metrics.CompactHeight).Shrink(0).Justify(ui.Center).AlignItems(ui.Center).Radius(t.Metrics.Radius).TextColor(t.Colors.Muted)
				closeButton.Background(controlBackground(&t, controlState{Pressed: closeButton.Pressed(), Hovered: closeButton.Hovered()}))
				closeButton.Children(func() { ui.Icon(c, IconX).Size(IconSizeSmall, IconSizeSmall) })
			})
		}
	})
}
