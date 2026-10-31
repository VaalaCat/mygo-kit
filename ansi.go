package kit

import (
	"github.com/egoist/mygo/ui"
	"strconv"
	"strings"
)

// ansiText 将终端 SGR 颜色转为文本片段，丢弃 OSC、光标和其他控制指令。
func ansiText(source string, theme Theme) []ui.Span {
	// 1. 顺序解码颜色状态，控制指令从不执行。
	var spans []ui.Span
	style := ui.Span{Color: theme.Colors.Text}
	start := 0
	flush := func(end int) {
		if end > start {
			v := style
			v.Text = source[start:end]
			spans = append(spans, v)
		}
	}
	for i := 0; i < len(source); {
		if source[i] != '\x1b' {
			i++
			continue
		}
		flush(i)
		i++
		if i == len(source) {
			start = i
			break
		}
		switch source[i] {
		case '[':
			i++
			begin := i
			for i < len(source) && (source[i] < 0x40 || source[i] > 0x7e) {
				i++
			}
			if i < len(source) {
				if source[i] == 'm' {
					applySGR(&style, source[begin:i], theme)
				}
				i++
			}
		case ']', 'P', '^', '_':
			i++
			for i < len(source) {
				if source[i] == '\a' {
					i++
					break
				}
				if source[i] == '\x1b' && i+1 < len(source) && source[i+1] == '\\' {
					i += 2
					break
				}
				i++
			}
		default:
			i++
		}
		start = i
	}
	flush(len(source))
	return spans
}

// applySGR 更新字重、装饰和前景背景颜色。
func applySGR(style *ui.Span, sequence string, t Theme) {
	// 1. 空参数视为复位，支持基础色、256 色和真彩色。
	parts := strings.Split(sequence, ";")
	codes := make([]int, len(parts))
	for i, p := range parts {
		codes[i], _ = strconv.Atoi(p)
	}
	for i := 0; i < len(codes); i++ {
		code := codes[i]
		switch {
		case code == 0:
			*style = ui.Span{Color: t.Colors.Text}
		case code == 1:
			style.Weight = 700
		case code == 22:
			style.Weight = 400
		case code == 3:
			style.Italic = true
		case code == 23:
			style.Italic = false
		case code == 4:
			style.Underline = true
		case code == 24:
			style.Underline = false
		case code == 9:
			style.Strikethrough = true
		case code == 29:
			style.Strikethrough = false
		case code >= 30 && code <= 37:
			style.Color = t.Colors.ANSI[code-30]
		case code >= 90 && code <= 97:
			style.Color = t.Colors.ANSI[code-90+8]
		case code >= 40 && code <= 47:
			style.Background = t.Colors.ANSI[code-40]
		case code >= 100 && code <= 107:
			style.Background = t.Colors.ANSI[code-100+8]
		case code == 39:
			style.Color = t.Colors.Text
		case code == 49:
			style.Background = ui.Transparent
		case code == 38 || code == 48:
			var color ui.Color
			valid := false
			if i+2 < len(codes) && codes[i+1] == 5 {
				color = ansiColor(codes[i+2], t)
				i += 2
				valid = true
			} else if i+4 < len(codes) && codes[i+1] == 2 {
				color = ui.RGB(uint8(clampByte(codes[i+2])), uint8(clampByte(codes[i+3])), uint8(clampByte(codes[i+4])))
				i += 4
				valid = true
			}
			if valid {
				if code == 38 {
					style.Color = color
				} else {
					style.Background = color
				}
			}
		}
	}
}

// ansiColor 解码 256 色索引。
func ansiColor(index int, t Theme) ui.Color {
	// 1. 前十六色使用当前主题，其余使用终端标准色立方体。
	index = clampByte(index)
	if index < 16 {
		return t.Colors.ANSI[index]
	}
	if index >= 232 {
		v := uint8(8 + (index-232)*10)
		return ui.RGB(v, v, v)
	}
	index -= 16
	levels := [6]uint8{0, 95, 135, 175, 215, 255}
	return ui.RGB(levels[index/36], levels[index/6%6], levels[index%6])
}

// clampByte 限定颜色分量范围。
func clampByte(v int) int { return min(255, max(0, v)) }

type ansiKey struct{}
type ansiCache struct {
	source     string
	colors     Colors
	plain      bool
	typography Typography
	width      float32
	lines      [][]ui.Span
}

// ANSIOutputOpt 设置输出工具栏标题及复制操作名称。
type ANSIOutputOpt struct {
	Title, CopyLabel string
	Plain            bool
}

// ANSIOutput 显示可选择、横向滚动的虚拟输出行；跟随底部由 ListState.FollowEnd 控制。
func ANSIOutput(c *ui.Context, source string, state *ui.ListState, opt ANSIOutputOpt) ui.Element {
	// 1. 标题与复制名称为空时采用本地化默认值；缓存文本解析结果，主题变化时刷新颜色。
	if opt.Title == "" {
		opt.Title = text(c, "output")
	}
	if opt.CopyLabel == "" {
		opt.CopyLabel = text(c, "copyOutput")
	}
	t := GetTheme(c)
	root := contentFrame(c, ui.Column(c))
	cache := ui.Local(root, ansiKey{}, func() ansiCache { return ansiCache{} })
	if cache.lines == nil || cache.source != source || cache.colors != t.Colors || cache.plain != opt.Plain || cache.typography != t.Typography {
		cache.source, cache.colors, cache.plain = source, t.Colors, opt.Plain
		cache.lines = [][]ui.Span{{}}
		cache.typography = t.Typography
		cache.width = 0
		spans := ansiText(source, t)
		if opt.Plain {
			spans = []ui.Span{{Text: source, Color: t.Colors.Text}}
		}
		for _, span := range spans {
			parts := strings.Split(span.Text, "\n")
			for i, part := range parts {
				if i > 0 {
					cache.lines = append(cache.lines, []ui.Span{})
				}
				v := span
				v.Text = part
				v.Font = t.Typography.MonoFont
				v.Size = t.Typography.CodeSize
				cache.lines[len(cache.lines)-1] = append(cache.lines[len(cache.lines)-1], v)
			}
		}
		for _, line := range cache.lines {
			width, _ := c.MeasureText(0, line...)
			cache.width = max(cache.width, width)
		}
	}
	// 2. 只构建可见行，复制入口返回去除控制码的全文。
	root.Children(func() {
		frameHeader(c, opt.Title, opt.CopyLabel, false, func() {
			if opt.Plain {
				c.WriteClipboard(source)
				return
			}
			var text strings.Builder
			for _, span := range ansiText(source, t) {
				text.WriteString(span.Text)
			}
			c.WriteClipboard(text.String())
		})
		list := ui.List(c, state, len(cache.lines)).Grow(1).MinHeight(0).Padding(t.Space(1), t.Space(2)).RowWidth(cache.width).Selectable()
		list.Rows(func(row ui.ListRow) {
			t.TextRole(ui.RichText(row.Context, cache.lines[row.Index]...).NoWrap(), TextCode).MinHeight(t.Space(5))
		})

	})
	return root
}
