package kit

import (
	"github.com/egoist/mygo/ui"
)

// dialogWidth 是 Dialog 面板的宽度。
const dialogWidth = 480

// Dialog 创建模态对话框：标题栏加正文，焦点陷阱和 Escape 由框架处理；需要确认和取消按钮时用 DialogWith。
func Dialog(c *ui.Context, open *bool, title string, content func()) ui.Element {
	// 1. 应用统一遮罩和面板样式。
	t := GetTheme(c)
	return ui.DialogBase(c, open, func(back, panel ui.Element) {
		dialogFrame(back, panel, &t, title)
		panel.Children(func() {
			HeaderWith(c, title, HeaderOpt{Inset: t.Space(4)}, nil)
			ui.Column(c).Padding(t.Space(4)).Gap(t.Space(2)).Children(content)
		})
	})
}

type dialogSubmitKey struct{}

// dialogSubmit 记录 DialogWith 构建正文时输入控件报告的事实：单行输入按了 Enter、正文里第一个可编辑输入。
// 单行输入自己吞掉 Enter，所以由 TextInputWith 在构建时报告。
type dialogSubmit struct {
	// buildingBody 为真表示 DialogWith 正在构建正文。
	buildingBody bool
	submitted    bool
	// firstInput 是正文里第一个可编辑（非只读）的输入，hasInput 表示它存在。
	firstInput ui.Element
	hasInput   bool
}

// noteDialogInput 在 DialogWith 构建正文期间，把可编辑输入和单行输入的 Enter 报告给对话框。
func noteDialogInput(c *ui.Context, input ui.Element, opt InputOpt, singleLine bool) {
	// 1. 只有正文构建期间才记录，避免普通界面里多做一次 Submitted。
	state := ui.Local(c.Root(), dialogSubmitKey{}, func() dialogSubmit { return dialogSubmit{} })
	if !state.buildingBody {
		return
	}
	// 2. 第一个可编辑输入拿初始焦点；单行输入按 Enter 提交。
	if !state.hasInput && !opt.ReadOnly {
		state.firstInput, state.hasInput = input, true
	}
	if singleLine && input.Submitted() {
		state.submitted = true
	}
}

// dialogFocusKey 标识对话框面板上保存的初始焦点状态。
type dialogFocusKey struct{}

// dialogFocus 是对话框自己的初始焦点：打开的第一帧只设置一次。
type dialogFocus struct {
	// isPlaced 为真表示初始焦点已经设置过。
	isPlaced bool
	// isRingHidden 为真表示焦点落在确认按钮上，且还没有离开过：这时不画焦点环。
	isRingHidden bool
}

// DialogOpt 描述 DialogWith 的标题、说明和确认按钮。
type DialogOpt struct {
	Title string
	// Description 显示在标题栏下方、正文之前，次要色；为空时不显示。
	Description string
	// ConfirmLabel 是确认按钮的文字；为空用库内置的「确定」。
	ConfirmLabel string
	// Danger 让确认按钮用 DangerButton，否则用 PrimaryButton。
	Danger bool
	// ConfirmDisabled 禁用确认按钮，Enter 也不会确认。
	ConfirmDisabled bool
}

// DialogWith 创建带页脚的模态对话框：标题栏、说明、body 画的正文、右对齐的「取消」和确认按钮。
// 点击确认或按 Enter 时关闭并返回 true（焦点在单行 TextInput 里时 Enter 提交，在多行输入里时是换行，在按钮上时按下该按钮）；
// 单行输入要用本库的 TextInput 或 TextInputWith 才会把 Enter 报告给对话框。
// 点击取消或按 Esc 只关闭，返回 false。打开时焦点落在正文里第一个可编辑的输入（只读输入除外），
// 页签、按钮等其他控件不会拿到初始焦点；正文里没有可编辑输入时落在确认按钮上，且不画焦点环（Enter 仍然确认，焦点移走后环照常显示）。
func DialogWith(c *ui.Context, open *bool, opt DialogOpt, body func()) (confirmed bool) {
	// 1. 面板外观与 Dialog 一致；正文区与页脚之间留 16。
	t := GetTheme(c)
	submit := ui.Local(c.Root(), dialogSubmitKey{}, func() dialogSubmit { return dialogSubmit{} })
	ui.DialogBase(c, open, func(back, panel ui.Element) {
		dialogFrame(back, panel, &t, opt.Title)
		panel.Children(func() {
			HeaderWith(c, opt.Title, HeaderOpt{Inset: t.Space(4)}, nil)
			ui.Column(c).Padding(t.Space(4)).Gap(t.Space(4)).Children(func() {
				if opt.Description != "" {
					t.TextRole(ui.Text(c, opt.Description), TextUI).TextColor(t.Colors.Muted)
				}
				submit.buildingBody, submit.submitted, submit.hasInput = true, false, false
				body()
				submit.buildingBody = false
				// 2. 页脚用反向行：确认按钮先于取消按钮建立，所以没有输入控件时框架的默认焦点落在确认按钮上，而视觉上它仍在最右。
				ui.Row(c).Reverse().Justify(ui.Start).Gap(t.Space(2)).Shrink(0).Children(func() {
					label := opt.ConfirmLabel
					if label == "" {
						label = text(c, "confirm")
					}
					var confirm ui.Element
					if opt.Danger {
						confirm = DangerButton(c.Key("kit-dialog-confirm"), label)
					} else {
						confirm = PrimaryButton(c.Key("kit-dialog-confirm"), label)
					}
					confirm.Disabled(opt.ConfirmDisabled)
					if confirm.Clicked() {
						confirmed = true
					}
					// 3. 初始焦点只在打开的第一帧设置：有可编辑输入给它，否则给确认按钮并先隐藏焦点环。
					focus := ui.Local(panel, dialogFocusKey{}, func() dialogFocus { return dialogFocus{} })
					if !focus.isPlaced {
						focus.isPlaced = true
						if submit.hasInput {
							submit.firstInput.Focus()
						} else if !opt.ConfirmDisabled {
							confirm.Focus()
							focus.isRingHidden = true
						}
					}
					if focus.isRingHidden {
						if confirm.Focused() {
							confirm.FocusRing(false)
						} else {
							focus.isRingHidden = false
						}
					}
					if Button(c.Key("kit-dialog-cancel"), text(c, "cancel")).Clicked() {
						*open = false
					}
				})
			})
		})
		// 4. Enter 没有被输入控件或获得焦点的按钮使用时确认（窗口级快捷键让给按钮和编辑器）；单行输入按 Enter 由正文报告；禁用确认时都忽略。
		if (c.Shortcut(0, ui.KeyEnter) || submit.submitted) && !opt.ConfirmDisabled {
			confirmed = true
		}
		if confirmed {
			*open = false
		}
	})
	return confirmed
}

// dialogFrame 设置对话框遮罩和面板的统一外观。
func dialogFrame(back, panel ui.Element, t *Theme, title string) {
	back.Background(ui.RGBA(0, 0, 0, 0.35))
	panel.Width(dialogWidth).MaxWidthPercent(90).Radius(t.Metrics.RadiusLarge).Clip().Background(t.Colors.Elevated).Border(1, t.Colors.Border).Label(title).Gap(0)
}
