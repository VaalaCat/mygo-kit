package kit

import (
	"github.com/egoist/mygo/ui"
	"strconv"
)

// PatchKind 描述片段行的展示类型，不代表完整文件差异。
type PatchKind string

const (
	PatchContext PatchKind = "context"
	PatchAdded   PatchKind = "added"
	PatchDeleted PatchKind = "deleted"
	PatchHeader  PatchKind = "header"
)

// PatchLine 是调用方已解析的一行片段；缺失行号保持空白，不推测真实位置。
type PatchLine struct {
	Kind             PatchKind
	Text             string
	OldLine, NewLine *int
}

// PatchOpt 保存工具栏名称及完整原文，复制时不从行片段重建补丁。
type PatchOpt struct{ Title, CopyLabel, Source string }

// 补丁行的尺寸：行号列宽 40，两列行号加符号列的最小行宽 120。
const (
	patchLineNumberWidth float32 = 40
	patchGutterWidth     float32 = 120
)

// PatchLines 显示通用补丁片段的颜色、可选行号和全文复制。
func PatchLines(c *ui.Context, state *ui.ListState, lines []PatchLine, opt PatchOpt) ui.Element {
	// 1. 复制只使用调用方提供的原文，行号仅作展示。
	t := GetTheme(c)
	if opt.Title == "" {
		opt.Title = "Patch"
	}
	if opt.CopyLabel == "" {
		opt.CopyLabel = text(c, "copyOutput")
	}
	root := contentFrame(c, ui.Column(c))
	root.Children(func() {
		frameHeader(c, opt.Title, opt.CopyLabel, false, func() { c.WriteClipboard(opt.Source) })
		width := patchGutterWidth
		for _, line := range lines {
			w, _ := c.MeasureText(0, ui.Span{Text: line.Text, Font: t.Typography.MonoFont, Size: t.Typography.CodeSize})
			width = max(width, w+patchGutterWidth)
		}
		list := ui.List(c, state, len(lines)).Grow(1).MinHeight(0).RowWidth(width).Selectable()
		list.Rows(func(row ui.ListRow) {
			c := row.Context
			i := row.Index
			line := lines[i]
			bg := ui.Transparent
			sign := " "
			fg := t.Colors.Text
			switch line.Kind {
			case PatchAdded:
				bg = t.Colors.Added
				sign = "+"
			case PatchDeleted:
				bg = t.Colors.Deleted
				sign = "−"
			case PatchHeader:
				bg = t.Colors.Panel
				fg = t.Colors.Accent
			}
			ui.Row(c).MinHeight(t.Metrics.CompactHeight).Gap(t.Space(2)).Background(bg).Children(func() {
				for _, number := range []*int{line.OldLine, line.NewLine} {
					label := ""
					if number != nil {
						label = strconv.Itoa(*number)
					}
					t.TextRole(ui.Text(c, label).Unselectable().Width(patchLineNumberWidth).TextAlign(ui.End), TextCode).TextColor(t.Colors.Placeholder)
				}
				t.TextRole(ui.Text(c, sign).Unselectable().Width(t.Space(3)), TextCode).TextColor(fg)
				t.TextRole(ui.Text(c, line.Text).NoWrap(), TextCode).TextColor(fg)
			})
		})

	})
	return root
}
