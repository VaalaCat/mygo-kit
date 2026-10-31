package kit

import (
	"image/color"
	"io"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/vt"
	"github.com/egoist/mygo/ui"
)

// 终端的固定参数：滚动历史行数与下划线、竖线光标的粗细。
const (
	terminalScrollbackLines = 10000
	terminalCursorEdge      = 2
)

// terminalMouseTracking 是程序请求的鼠标上报级别，数值越大上报的事件越多。
type terminalMouseTracking int

const (
	// terminalMouseOff 程序没有开启鼠标上报，指针只用于本地选择和滚动。
	terminalMouseOff terminalMouseTracking = iota
	// terminalMouseX10 只上报按下（含滚轮），对应 DEC 模式 9。
	terminalMouseX10
	// terminalMousePressRelease 上报按下、松开和滚轮，对应模式 1000、1001。
	terminalMousePressRelease
	// terminalMouseDrag 另外上报按住按键时的移动，对应模式 1002。
	terminalMouseDrag
	// terminalMouseMotion 另外上报没有按键时的移动，对应模式 1003。
	terminalMouseMotion
)

// terminalPoint 是选择端点：line 是含滚动历史的绝对行号（0 为最旧一行），col 是列。
type terminalPoint struct{ line, col int }

// TerminalState 保存一个终端的屏幕、滚动历史和选择；由应用持有，可在任意 goroutine 调用 Write。
type TerminalState struct {
	mu   sync.Mutex
	emu  *vt.Emulator
	rows int
	cols int
	// 以下字段随终端输出由模拟器回调更新。
	title           string
	cursorVisible   bool
	cursorStyle     vt.CursorStyle
	applicationKeys bool
	bracketedPaste  bool
	closed          bool
	input           atomic.Pointer[func([]byte)]
	scrollOffset    int
	scrollRemainder float32
	selection       [2]terminalPoint
	hasSelection    bool
	selecting       bool
	preedit         string
	// mouseModes 是程序当前开启的鼠标上报模式；mouseReporting 表示正按住的按键已上报给程序，
	// 其后的移动和松开也上报；mouseButton 是该按键，mouseCell 是最近一次上报的屏幕格子。
	mouseModes     map[ansi.Mode]bool
	mouseReporting bool
	mouseButton    vt.MouseButton
	mouseCell      [2]int
}

// NewTerminalState 创建 rows 行 cols 列的空终端，并启动转发模拟器自动应答的 goroutine；不再使用时调用 Close。
func NewTerminalState(rows, cols int) *TerminalState {
	// 1. 模拟器回调更新光标、键盘模式和标题；自动应答经管道读出后转给 OnInput。
	rows, cols = max(1, rows), max(1, cols)
	s := &TerminalState{emu: vt.NewEmulator(cols, rows), rows: rows, cols: cols, cursorVisible: true, mouseModes: map[ansi.Mode]bool{}}
	s.emu.SetScrollbackSize(terminalScrollbackLines)
	s.emu.SetCallbacks(vt.Callbacks{
		Title:            func(title string) { s.title = title },
		CursorVisibility: func(visible bool) { s.cursorVisible = visible },
		CursorStyle:      func(style vt.CursorStyle, _ bool) { s.cursorStyle = style },
		EnableMode:       func(mode ansi.Mode) { s.setMode(mode, true) },
		DisableMode:      func(mode ansi.Mode) { s.setMode(mode, false) },
	})
	// 2. 模拟器向管道写应答时会阻塞到有人读取，所以必须常驻读取。
	go s.forwardResponses()
	return s
}

// forwardResponses 把模拟器的自动应答（光标位置、设备属性等）交给当前的 OnInput，直到 Close。
func (s *TerminalState) forwardResponses() {
	// 1. 管道关闭后读取返回错误，goroutine 随之结束。
	buf := make([]byte, 1024)
	for {
		n, err := s.emu.Read(buf)
		if n > 0 {
			if send := s.input.Load(); send != nil && *send != nil {
				(*send)(append([]byte(nil), buf[:n]...))
			}
		}
		if err != nil {
			return
		}
	}
}

// setMode 记录影响按键编码、粘贴和鼠标上报的终端模式。
func (s *TerminalState) setMode(mode ansi.Mode, on bool) {
	// 1. 只关心应用光标键、括号粘贴和五种鼠标上报模式；鼠标编码方式由模拟器自己按 SGR 模式选择。
	switch mode {
	case ansi.ModeCursorKeys:
		s.applicationKeys = on
	case ansi.ModeBracketedPaste:
		s.bracketedPaste = on
	case ansi.ModeMouseX10, ansi.ModeMouseNormal, ansi.ModeMouseHighlight, ansi.ModeMouseButtonEvent, ansi.ModeMouseAnyEvent:
		s.mouseModes[mode] = on
	}
}

// mouseTracking 返回程序请求的鼠标上报级别：同时开启多个模式时取最详细的一个。
func (s *TerminalState) mouseTracking() terminalMouseTracking {
	switch {
	case s.mouseModes[ansi.ModeMouseAnyEvent]:
		return terminalMouseMotion
	case s.mouseModes[ansi.ModeMouseButtonEvent]:
		return terminalMouseDrag
	case s.mouseModes[ansi.ModeMouseNormal], s.mouseModes[ansi.ModeMouseHighlight]:
		return terminalMousePressRelease
	case s.mouseModes[ansi.ModeMouseX10]:
		return terminalMouseX10
	}
	return terminalMouseOff
}

// Write 喂入 PTY 输出；在回看历史时保持当前位置不动。
func (s *TerminalState) Write(data []byte) {
	// 1. 记录写入前的历史行数，写入后按新增行数补偿滚动偏移。
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	before := s.emu.ScrollbackLen()
	_, _ = s.emu.Write(data)
	if s.scrollOffset > 0 {
		s.scrollOffset += s.emu.ScrollbackLen() - before
	}
	s.scrollOffset = min(s.scrollOffset, s.historyLines())
}

// Close 结束转发自动应答的 goroutine；之后的 Write 被忽略。
func (s *TerminalState) Close() {
	// 1. 关闭应答管道的写端，读取它的 goroutine 收到 EOF 后退出；不调用 Emulator.Close，它与 Read 共享未加锁的标志。
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	if pipe, ok := s.emu.InputPipe().(io.Closer); ok {
		_ = pipe.Close()
	}
}

// Size 返回终端当前的行数和列数。
func (s *TerminalState) Size() (rows, cols int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rows, s.cols
}

// Title 返回程序通过 OSC 设置的窗口标题。
func (s *TerminalState) Title() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.title
}

// Text 返回当前屏幕的纯文本，每行去掉行尾空格，末尾空行去掉。
func (s *TerminalState) Text() string {
	// 1. 只读屏幕区，不含滚动历史。
	s.mu.Lock()
	defer s.mu.Unlock()
	start := s.historyLines()
	lines := make([]string, s.rows)
	for row := range lines {
		lines[row] = s.lineText(start+row, 0, s.cols-1)
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

// resize 在行列数变化时调整模拟器，返回是否发生了变化；行数变少时光标所在行保持可见。
func (s *TerminalState) resize(rows, cols int) bool {
	// 1. 相同尺寸不触碰模拟器。
	s.mu.Lock()
	defer s.mu.Unlock()
	if rows == s.rows && cols == s.cols {
		return false
	}
	// 2. 行数变少且光标落到新屏幕之外时，先把顶部多出的行推进滚动历史，光标随内容上移，提示符仍在可见区底部。
	if rows < s.rows && !s.emu.IsAltScreen() {
		if drop := s.emu.CursorPosition().Y - (rows - 1); drop > 0 {
			before := s.emu.ScrollbackLen()
			_, _ = s.emu.Write([]byte("\x1b[" + strconv.Itoa(drop) + "S\x1b[" + strconv.Itoa(drop) + "A"))
			if s.scrollOffset > 0 {
				s.scrollOffset += s.emu.ScrollbackLen() - before
			}
		}
	}
	// 3. 调整模拟器尺寸；用户正在回看历史时保持原位置，否则停在底部。
	s.rows, s.cols = rows, cols
	s.emu.Resize(cols, rows)
	s.scrollOffset = min(s.scrollOffset, s.historyLines())
	s.hasSelection = false
	return true
}

// historyLines 返回可回看的历史行数；备用屏幕没有历史。
func (s *TerminalState) historyLines() int {
	if s.emu.IsAltScreen() {
		return 0
	}
	return s.emu.ScrollbackLen()
}

// cellAt 返回绝对行 line、列 col 的格子，越界时为 nil。
func (s *TerminalState) cellAt(line, col int) *uv.Cell {
	history := s.historyLines()
	if line < history {
		return s.emu.ScrollbackCellAt(col, line)
	}
	return s.emu.CellAt(col, line-history)
}

// lineText 返回绝对行 line 中 from 到 to（含）列的文字，去掉行尾空格。
func (s *TerminalState) lineText(line, from, to int) string {
	// 1. 空格子补空格，宽字符后的占位格不输出。
	var b strings.Builder
	for col := from; col <= to; col++ {
		cell := s.cellAt(line, col)
		switch {
		case cell == nil || cell.Content == "":
			if cell == nil || cell.Width != 0 {
				b.WriteByte(' ')
			}
		default:
			b.WriteString(cell.Content)
		}
	}
	return strings.TrimRight(b.String(), " ")
}

// scrollBy 把视图向历史方向滚动 lines 行（负数向最新输出滚动），返回是否移动。
func (s *TerminalState) scrollBy(lines int) bool {
	next := min(max(s.scrollOffset+lines, 0), s.historyLines())
	moved := next != s.scrollOffset
	s.scrollOffset = next
	return moved
}

// selectedText 返回选择区域的文字，行间用换行连接。
func (s *TerminalState) selectedText() string {
	// 1. 端点按阅读顺序排列，首尾行截取，中间行取整行。
	if !s.hasSelection {
		return ""
	}
	start, end := s.selectionBounds()
	lines := make([]string, 0, end.line-start.line+1)
	for line := start.line; line <= end.line; line++ {
		from, to := 0, s.cols-1
		if line == start.line {
			from = start.col
		}
		if line == end.line {
			to = end.col
		}
		lines = append(lines, s.lineText(line, from, to))
	}
	return strings.Join(lines, "\n")
}

// selectionBounds 返回按阅读顺序排列的选择起点和终点。
func (s *TerminalState) selectionBounds() (start, end terminalPoint) {
	start, end = s.selection[0], s.selection[1]
	if end.line < start.line || end.line == start.line && end.col < start.col {
		start, end = end, start
	}
	return start, end
}

// selectedColumns 返回绝对行 line 上被选中的列范围，没有选中时 ok 为假。
func (s *TerminalState) selectedColumns(line int) (from, to int, ok bool) {
	if !s.hasSelection {
		return 0, 0, false
	}
	start, end := s.selectionBounds()
	if line < start.line || line > end.line {
		return 0, 0, false
	}
	from, to = 0, s.cols-1
	if line == start.line {
		from = start.col
	}
	if line == end.line {
		to = end.col
	}
	return from, to, true
}

// TerminalOpt 指定终端向应用回传的输入和尺寸变化。
type TerminalOpt struct {
	// OnInput 接收键盘、输入法上屏文字和粘贴编码后的字节。
	// 程序查询光标位置等自动应答也从这里发出，此时在后台 goroutine 调用。
	OnInput func(data []byte)
	// OnResize 在面板尺寸换算出的行列数变化时调用。
	OnResize func(rows, cols int)
}

// terminalMetrics 是一个字符格的尺寸：宽取等宽字体的真实前进宽度，高落在四像素网格上。
type terminalMetrics struct {
	cellWidth, cellHeight float32
	textOffset            float32
	padX, padY            float32
}

// terminalCellMetrics 测量主题等宽字体的字符格，行高取 12/16。
func terminalCellMetrics(c *ui.Context, t Theme) terminalMetrics {
	// 1. 十个字符取平均，避免小数前进宽度被取整。
	const sample = 10
	width, height := c.MeasureText(0, ui.Span{Text: strings.Repeat("M", sample), Font: t.Typography.MonoFont, Size: t.Typography.CodeSize})
	cellHeight := t.Space(4)
	return terminalMetrics{
		cellWidth:  width / sample,
		cellHeight: cellHeight,
		textOffset: float32(math.Floor(float64((cellHeight-height)/2 + .5))),
		padX:       t.Space(3),
		padY:       t.Space(1),
	}
}

// grid 返回宽 width、高 height 的面板能容纳的行列数，至少一行一列。
func (m terminalMetrics) grid(width, height float32) (rows, cols int) {
	rows = max(1, int((height-2*m.padY)/m.cellHeight))
	cols = max(1, int((width-2*m.padX)/m.cellWidth))
	return rows, cols
}

// cellRect 返回显示行 row、列 col 起宽 cols 格的矩形，边缘取整以免相邻色块出现缝隙。
func (m terminalMetrics) cellRect(origin ui.Rect, row, col, cols int) ui.Rect {
	left := float32(math.Round(float64(origin.X + float32(col)*m.cellWidth)))
	right := float32(math.Round(float64(origin.X + float32(col+cols)*m.cellWidth)))
	return ui.Rect{X: left, Y: origin.Y + float32(row)*m.cellHeight, W: right - left, H: m.cellHeight}
}

// Terminal 绘制终端：等宽 12/16（Lilex），16 色用主题 ANSI，支持 256 色与真彩色、光标、滚动历史、选择复制（双击选词、三击选行）。
// 程序开启鼠标上报（DEC 模式 9、1000-1003，SGR 编码）后，点击、拖动、移动和滚轮交给程序，按住 Shift 仍可本地选择和回看。
// 元素默认填满父容器剩余空间，背景为编辑器色；点击获得焦点，焦点只体现在光标上。
func Terminal(c *ui.Context, state *TerminalState, opt TerminalOpt) ui.Element {
	// 1. 登记当前的输入回调，并测量字符格。
	t := GetTheme(c)
	m := terminalCellMetrics(c, t)
	state.input.Store(&opt.OnInput)
	services := c.Services()
	root := ui.Box(c).Grow(1).MinWidth(0).MinHeight(0).Clip().Background(t.Colors.Editor).Focusable().FocusRing(false).Cursor(ui.CursorText)
	// 屏幕文字作为无障碍值，读屏和自动化测试都能读到终端当前的内容。
	root.Label(text(c, "terminal")).Value(state.Text())
	// 2. 用上一帧的布局尺寸换算行列，变化时同步模拟器并通知应用。
	if bounds := root.LayoutBounds(); bounds.W > 0 && bounds.H > 0 {
		rows, cols := m.grid(bounds.W, bounds.H)
		if state.resize(rows, cols) && opt.OnResize != nil {
			opt.OnResize(rows, cols)
		}
	}
	// 3. 输入法候选框跟随光标；绘制时若尺寸与状态不符，再请求一帧去同步。
	caret := state.caretRect(m)
	focused := root.Focused()
	root.TextCaret(caret).HandleInput(func(ev ui.InputEvent) bool { return state.handleInput(services, ev, m, opt) })
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if rows, cols := m.grid(r.W, r.H); !state.hasSize(rows, cols) {
			c.Invalidate()
		}
		state.paint(p, r, t, m, focused)
	})
	return root
}

// hasSize 判断终端当前行列数是否等于 rows、cols。
func (s *TerminalState) hasSize(rows, cols int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rows == rows && s.cols == cols
}

// caretRect 返回光标格相对元素的矩形，供输入法定位候选框。
func (s *TerminalState) caretRect(m terminalMetrics) ui.Rect {
	s.mu.Lock()
	defer s.mu.Unlock()
	x, y := s.emu.CursorPosition().X, s.emu.CursorPosition().Y
	return ui.Rect{X: m.padX + float32(x)*m.cellWidth, Y: m.padY + float32(y+s.scrollOffset)*m.cellHeight, W: m.cellWidth, H: m.cellHeight}
}

// terminalStyle 是一个格子解析后的显示样式。
type terminalStyle struct {
	fg, bg                          ui.Color
	bold, italic, underline, strike bool
}

// terminalRun 是同一行内样式相同的连续格子；text 为空表示只画背景。
type terminalRun struct {
	col, cols int
	text      string
	style     terminalStyle
}

// terminalColor 把模拟器的颜色换成界面颜色：0-15 取主题，16-255 取标准色表，其余直接取 RGB。
func terminalColor(value color.Color, t Theme) ui.Color {
	// 1. 索引色走主题与 256 色表，真彩色逐分量取高八位。
	switch v := value.(type) {
	case ansi.BasicColor:
		return t.Colors.ANSI[v&15]
	case ansi.IndexedColor:
		return ansiColor(int(v), t)
	}
	r, g, b, _ := value.RGBA()
	return ui.RGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// resolveTerminalStyle 解析格子样式：默认前景取主题文字色，背景取透明；反色互换，弱化向背景靠拢。
func resolveTerminalStyle(style uv.Style, t Theme) terminalStyle {
	// 1. 默认色取主题，显式色按类型转换。
	fg, bg := t.Colors.Text, ui.Transparent
	if style.Fg != nil {
		fg = terminalColor(style.Fg, t)
	}
	if style.Bg != nil {
		bg = terminalColor(style.Bg, t)
	}
	// 2. 反色时透明背景视为编辑器色。
	if style.Attrs&uv.AttrReverse != 0 {
		if bg == ui.Transparent {
			bg = t.Colors.Editor
		}
		fg, bg = bg, fg
	}
	base := bg
	if base == ui.Transparent {
		base = t.Colors.Editor
	}
	if style.Attrs&uv.AttrFaint != 0 {
		fg = fg.Mix(base, .5)
	}
	if style.Attrs&uv.AttrConceal != 0 {
		fg = base
	}
	return terminalStyle{fg: fg, bg: bg, bold: style.Attrs&uv.AttrBold != 0, italic: style.Attrs&uv.AttrItalic != 0, underline: style.Underline != uv.UnderlineNone, strike: style.Attrs&uv.AttrStrikethrough != 0}
}

// isTerminalGridRune 判断字符能和邻格连成一段文字：这些字符在等宽字体里有字形且前进宽度等于一格。
func isTerminalGridRune(r rune) bool {
	return r >= 0x20 && r < 0x7f || r >= 0xa0 && r <= 0xff || r >= 0x2500 && r <= 0x259f
}

// rowRuns 把绝对行 line 切成样式相同的连续段；非网格字符和宽字符单独成段，保证落在格线上。
func (s *TerminalState) rowRuns(line int, t Theme) []terminalRun {
	// 1. 逐格解析样式，相同样式的网格字符并入上一段。
	var runs []terminalRun
	for col := 0; col < s.cols; {
		cell := s.cellAt(line, col)
		if cell == nil {
			col++
			continue
		}
		width := max(1, cell.Width)
		if cell.Width == 0 && cell.Content == "" {
			col++
			continue
		}
		text := cell.Content
		if text == "" {
			text = " "
		}
		style := resolveTerminalStyle(cell.Style, t)
		runes := []rune(text)
		grid := width == 1 && len(runes) == 1 && isTerminalGridRune(runes[0])
		last := len(runs) - 1
		if grid && last >= 0 && runs[last].style == style && runs[last].col+runs[last].cols == col && isGridText(runs[last].text) {
			runs[last].cols++
			runs[last].text += text
		} else {
			runs = append(runs, terminalRun{col: col, cols: width, text: text, style: style})
		}
		col += width
	}
	return runs
}

// isGridText 判断一段文字全部由网格字符组成，可以继续追加相邻格子。
func isGridText(text string) bool {
	for _, r := range text {
		if !isTerminalGridRune(r) {
			return false
		}
	}
	return true
}

// paint 绘制可见的每一行、选择、光标和输入法候选文字。
func (s *TerminalState) paint(p *ui.Painter, r ui.Rect, t Theme, m terminalMetrics, focused bool) {
	// 1. 在锁内按显示行取出绝对行并绘制，最多画到状态里的行数。
	s.mu.Lock()
	defer s.mu.Unlock()
	origin := ui.Rect{X: r.X + m.padX, Y: r.Y + m.padY}
	top := s.historyLines() - min(s.scrollOffset, s.historyLines())
	maxRows := min(s.rows, max(1, int((r.H-2*m.padY)/m.cellHeight)))
	for row := range maxRows {
		line := top + row
		runs := s.rowRuns(line, t)
		// 2. 先画背景，再画选择，最后画文字。
		for _, run := range runs {
			if run.style.bg != ui.Transparent {
				p.Fill(m.cellRect(origin, row, run.col, run.cols), run.style.bg, 0)
			}
		}
		if from, to, ok := s.selectedColumns(line); ok {
			p.Fill(m.cellRect(origin, row, from, to-from+1), t.Colors.Selection, 0)
		}
		for _, run := range runs {
			s.paintRun(p, t, m, m.cellRect(origin, row, run.col, run.cols), run)
		}
	}
	// 3. 光标只在回看历史时随内容移动，移出可见区就不画。
	s.paintCursor(p, t, m, origin, top, maxRows, focused)
}

// paintRun 在 rect 处绘制一段文字；空白段不绘制文字。
func (s *TerminalState) paintRun(p *ui.Painter, t Theme, m terminalMetrics, rect ui.Rect, run terminalRun) {
	// 1. 全空格没有字形，省掉一次排版。
	if strings.TrimLeft(run.text, " ") == "" && !run.style.underline && !run.style.strike {
		return
	}
	weight := 400
	if run.style.bold {
		weight = 700
	}
	span := ui.Span{Text: run.text, Font: t.Typography.MonoFont, Size: t.Typography.CodeSize, Weight: weight, Italic: run.style.italic, Color: run.style.fg, Underline: run.style.underline, Strikethrough: run.style.strike}
	p.RichText(rect.X, rect.Y+m.textOffset, 0, span)
}

// paintCursor 画光标：聚焦时是实心块（或下划线、竖线），失焦时是空心框；程序隐藏光标时不画。
func (s *TerminalState) paintCursor(p *ui.Painter, t Theme, m terminalMetrics, origin ui.Rect, top, rows int, focused bool) {
	// 1. 光标所在的绝对行换算成显示行，不在可见区内就不画。
	pos := s.emu.CursorPosition()
	row := s.historyLines() + pos.Y - top
	if !s.cursorVisible || row < 0 || row >= rows {
		return
	}
	cell := s.emu.CellAt(pos.X, pos.Y)
	cols := 1
	if cell != nil && cell.Width > 1 {
		cols = cell.Width
	}
	rect := m.cellRect(origin, row, pos.X, cols)
	// 2. 输入法组字中的文字盖在光标处并加下划线。
	if s.preedit != "" {
		width := float32(len([]rune(s.preedit)))
		p.Fill(ui.Rect{X: rect.X, Y: rect.Y, W: width * m.cellWidth * 2, H: rect.H}, t.Colors.Editor, 0)
		p.RichText(rect.X, rect.Y+m.textOffset, 0, ui.Span{Text: s.preedit, Font: t.Typography.MonoFont, Size: t.Typography.CodeSize, Color: t.Colors.Text, Underline: true})
		return
	}
	if !focused {
		p.Stroke(rect, t.Colors.Text, 0, 1)
		return
	}
	switch s.cursorStyle {
	case vt.CursorUnderline:
		p.Fill(ui.Rect{X: rect.X, Y: rect.Y + rect.H - terminalCursorEdge, W: rect.W, H: terminalCursorEdge}, t.Colors.Text, 0)
	case vt.CursorBar:
		p.Fill(ui.Rect{X: rect.X, Y: rect.Y, W: terminalCursorEdge, H: rect.H}, t.Colors.Text, 0)
	default:
		p.Fill(rect, t.Colors.Text, 0)
		if cell != nil && strings.TrimSpace(cell.Content) != "" {
			p.RichText(rect.X, rect.Y+m.textOffset, 0, ui.Span{Text: cell.Content, Font: t.Typography.MonoFont, Size: t.Typography.CodeSize, Color: t.Colors.Editor})
		}
	}
}

// send 把编码后的字节交给 OnInput；任何输入都会回到最新输出并取消选择。
func (s *TerminalState) send(opt TerminalOpt, data string) {
	// 1. 先改本地状态，再通知应用，避免回调重入时读到旧偏移。
	s.mu.Lock()
	s.scrollOffset = 0
	s.hasSelection = false
	s.mu.Unlock()
	if opt.OnInput != nil && data != "" {
		opt.OnInput([]byte(data))
	}
}

// handleInput 处理键盘、输入法、编辑命令、指针选择和滚轮；返回是否接管了事件。
func (s *TerminalState) handleInput(services ui.Services, ev ui.InputEvent, m terminalMetrics, opt TerminalOpt) bool {
	switch ev.Kind {
	case ui.InputKeyDown:
		return s.handleKeyDown(services, ev, opt)
	case ui.InputText:
		return s.handleText(ev.Text, opt)
	case ui.InputCompose:
		s.mu.Lock()
		s.preedit = ev.Text
		s.mu.Unlock()
		return true
	case ui.InputCommand:
		return s.handleCommand(services, ev.Text, opt)
	case ui.InputPointerDown, ui.InputPointerMove, ui.InputPointerUp:
		return s.handlePointer(ev, m)
	case ui.InputScroll:
		return s.handleScroll(ev, m)
	}
	return false
}

// handleText 发送键盘输入或输入法上屏的文字；控制字符由按键编码负责，这里丢弃。
func (s *TerminalState) handleText(text string, opt TerminalOpt) bool {
	// 1. 上屏后组字文字消失。
	s.mu.Lock()
	s.preedit = ""
	s.mu.Unlock()
	printable := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, text)
	s.send(opt, printable)
	return true
}

// handleKeyDown 处理一次按键：复制粘贴、翻页滚动、转义序列；会输入文字的按键等待随后的文字事件。
func (s *TerminalState) handleKeyDown(services ui.Services, ev ui.InputEvent, opt TerminalOpt) bool {
	// 1. 剪贴板快捷键。
	if action := terminalClipboardAction(runtime.GOOS, ev.Mods, ev.Key); action != "" {
		return s.handleCommand(services, action, opt)
	}
	// 2. 其余带 Super 的按键留给应用快捷键。
	if ev.Mods&ui.Super != 0 {
		return false
	}
	// 3. Shift+PageUp/PageDown 翻看历史，不发给程序。
	if ev.Mods == ui.Shift && (ev.Key == ui.KeyPageUp || ev.Key == ui.KeyPageDown) {
		s.mu.Lock()
		page := max(1, s.rows-1)
		if ev.Key == ui.KeyPageDown {
			page = -page
		}
		s.scrollBy(page)
		s.mu.Unlock()
		return true
	}
	// 4. 编码成转义序列；文字键没有序列，等待文字事件。
	s.mu.Lock()
	applicationKeys := s.applicationKeys
	s.mu.Unlock()
	if data, ok := encodeTerminalKey(runtime.GOOS, applicationKeys, ev.Mods, ev.Key); ok {
		s.send(opt, data)
		return true
	}
	_, typesText := terminalKeyChars[ev.Key]
	return typesText
}

// handleCommand 执行复制、粘贴和全选；复制与粘贴同时服务编辑菜单和快捷键。
func (s *TerminalState) handleCommand(services ui.Services, name string, opt TerminalOpt) bool {
	switch name {
	case "copy":
		s.mu.Lock()
		text := s.selectedText()
		s.mu.Unlock()
		if text != "" {
			services.WriteClipboard(text)
		}
	case "paste":
		text := services.ReadClipboard()
		if text == "" {
			return true
		}
		text = strings.NewReplacer("\r\n", "\r", "\n", "\r").Replace(text)
		s.mu.Lock()
		bracketed := s.bracketedPaste
		s.mu.Unlock()
		if bracketed {
			text = "\x1b[200~" + strings.ReplaceAll(text, "\x1b[201~", "") + "\x1b[201~"
		}
		s.send(opt, text)
	case "selectAll":
		s.mu.Lock()
		s.selection = [2]terminalPoint{{0, 0}, {s.historyLines() + s.rows - 1, s.cols - 1}}
		s.hasSelection = true
		s.mu.Unlock()
	default:
		return false
	}
	return true
}

// screenCell 把相对元素的坐标换成可见屏幕里的格子，超出可见区时取最近的一格。
func (s *TerminalState) screenCell(x, y float32, m terminalMetrics) (row, col int) {
	row = min(max(int(math.Floor(float64((y-m.padY)/m.cellHeight))), 0), s.rows-1)
	col = min(max(int(math.Floor(float64((x-m.padX)/m.cellWidth))), 0), s.cols-1)
	return row, col
}

// pointAt 把相对元素的坐标换成绝对行列，超出可见区时取最近的一格。
func (s *TerminalState) pointAt(x, y float32, m terminalMetrics) terminalPoint {
	row, col := s.screenCell(x, y, m)
	return terminalPoint{line: s.historyLines() - s.scrollOffset + row, col: col}
}

// isWordCell 绝对行 line、列 col 的格子是否属于一个单词（宽字符的占位格随其字符）？
func (s *TerminalState) isWordCell(line, col int) bool {
	cell := s.cellAt(line, col)
	if cell == nil {
		return false
	}
	if cell.Content == "" {
		return cell.Width == 0
	}
	return !strings.ContainsAny(cell.Content, " \t\"'`<>()[]{}|;,")
}

// selectWord 选中 point 所在的单词；point 落在分隔符上时只选中这一格。
func (s *TerminalState) selectWord(point terminalPoint) {
	from, to := point.col, point.col
	if s.isWordCell(point.line, point.col) {
		for from > 0 && s.isWordCell(point.line, from-1) {
			from--
		}
		for to < s.cols-1 && s.isWordCell(point.line, to+1) {
			to++
		}
	}
	s.selection = [2]terminalPoint{{point.line, from}, {point.line, to}}
	s.hasSelection = true
}

// handlePointer 处理指针：程序开启鼠标上报时交给程序（按住 Shift 除外）；否则左键拖动选择格子，
// 双击选单词，三击选整行，拖到面板外时每个移动事件滚动一行。
func (s *TerminalState) handlePointer(ev ui.InputEvent, m terminalMetrics) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reportPointer(ev, m) {
		return true
	}
	switch ev.Kind {
	case ui.InputPointerDown:
		if ev.Button != 0 {
			return false
		}
		point := s.pointAt(ev.X, ev.Y, m)
		s.selection = [2]terminalPoint{point, point}
		s.hasSelection, s.selecting = false, true
		switch {
		case ev.Clicks == 2:
			s.selectWord(point)
			s.selecting = false
		case ev.Clicks >= 3:
			s.selection = [2]terminalPoint{{point.line, 0}, {point.line, s.cols - 1}}
			s.hasSelection, s.selecting = true, false
		}
		return true
	case ui.InputPointerMove:
		if !s.selecting {
			return false
		}
		if ev.Y < m.padY {
			s.scrollBy(1)
		} else if ev.Y > m.padY+float32(s.rows)*m.cellHeight {
			s.scrollBy(-1)
		}
		s.selection[1] = s.pointAt(ev.X, ev.Y, m)
		s.hasSelection = s.selection[0] != s.selection[1]
		return true
	default:
		selecting := s.selecting
		s.selecting = false
		return selecting
	}
}

// reportPointer 把按下、移动和松开按鼠标上报模式交给程序，返回是否已上报；调用方持有 s.mu。
func (s *TerminalState) reportPointer(ev ui.InputEvent, m terminalMetrics) bool {
	// 1. 没有上报模式、按住 Shift（本地选择）或不认识的按键时，指针留给本地处理。
	tracking := s.mouseTracking()
	row, col := s.screenCell(ev.X, ev.Y, m)
	switch ev.Kind {
	case ui.InputPointerDown:
		button, ok := terminalMouseButton(ev.Button)
		if tracking == terminalMouseOff || ev.Mods&ui.Shift != 0 || !ok {
			return false
		}
		// 2. 按下回到最新输出、取消选择，并记住这个按键，让随后的移动和松开也上报。
		s.scrollOffset, s.hasSelection = 0, false
		s.mouseReporting, s.mouseButton, s.mouseCell = true, button, [2]int{row, col}
		s.emu.SendMouse(vt.MouseClick(uv.Mouse{X: col, Y: row, Button: button, Mod: terminalMouseMods(ev.Mods)}))
		return true
	case ui.InputPointerMove:
		// 3. 按住已上报的按键拖动需要拖动级别；没有按键的移动需要全移动级别；格子没变就不重复上报。
		button := vt.MouseNone
		switch {
		case s.mouseReporting && tracking >= terminalMouseDrag:
			button = s.mouseButton
		case s.mouseReporting:
			return true
		case ev.Button == -1 && tracking == terminalMouseMotion && ev.Mods&ui.Shift == 0:
		default:
			return false
		}
		if s.mouseCell == [2]int{row, col} {
			return true
		}
		s.mouseCell = [2]int{row, col}
		s.emu.SendMouse(vt.MouseMotion(uv.Mouse{X: col, Y: row, Button: button, Mod: terminalMouseMods(ev.Mods)}))
		return true
	case ui.InputPointerUp:
		// 4. 松开只在程序要求按下/松开级别时上报；X10 模式没有松开事件。
		if !s.mouseReporting {
			return false
		}
		s.mouseReporting = false
		if tracking >= terminalMousePressRelease {
			s.emu.SendMouse(vt.MouseRelease(uv.Mouse{X: col, Y: row, Button: s.mouseButton, Mod: terminalMouseMods(ev.Mods)}))
		}
		return true
	}
	return false
}

// terminalMouseButton 把 MyGo 的按键编号换成鼠标上报的按键：0 主键、1 副键、2 中键。
func terminalMouseButton(button int) (vt.MouseButton, bool) {
	switch button {
	case 0:
		return vt.MouseLeft, true
	case 1:
		return vt.MouseRight, true
	case 2:
		return vt.MouseMiddle, true
	}
	return vt.MouseNone, false
}

// terminalMouseMods 把修饰键换成鼠标上报的修饰位。
func terminalMouseMods(mods ui.Modifiers) uv.KeyMod {
	var out uv.KeyMod
	if mods&ui.Shift != 0 {
		out |= uv.ModShift
	}
	if mods&ui.Alt != 0 {
		out |= uv.ModAlt
	}
	if mods&ui.Ctrl != 0 {
		out |= uv.ModCtrl
	}
	return out
}

// handleScroll 用滚轮和触控板回看历史，不足一行的像素累积到下次；程序开启鼠标上报时改为上报滚轮（按住 Shift 除外）。
func (s *TerminalState) handleScroll(ev ui.InputEvent, m terminalMetrics) bool {
	// 1. DY 向下为正，对应向最新输出滚动，所以取反。
	s.mu.Lock()
	s.scrollRemainder -= ev.DY
	lines := int(s.scrollRemainder / m.cellHeight)
	s.scrollRemainder -= float32(lines) * m.cellHeight
	if s.mouseTracking() == terminalMouseOff || ev.Mods&ui.Shift != 0 {
		defer s.mu.Unlock()
		return s.scrollBy(lines) || s.historyLines() > 0
	}
	// 2. 上报时每一行一次滚轮事件，位置取指针所在的格子；每次发送单独加锁，
	// 让 OnInput 里同步调用 Write 的应用不会卡在后一次发送上。
	row, col := s.screenCell(ev.X, ev.Y, m)
	s.mu.Unlock()
	button := vt.MouseWheelUp
	if lines < 0 {
		button, lines = vt.MouseWheelDown, -lines
	}
	for range lines {
		s.mu.Lock()
		s.emu.SendMouse(vt.MouseWheel(uv.Mouse{X: col, Y: row, Button: button, Mod: terminalMouseMods(ev.Mods)}))
		s.mu.Unlock()
	}
	return true
}

// terminalClipboardAction 返回快捷键对应的剪贴板命令：macOS 是 Cmd+C/V，其他系统是 Ctrl+Shift+C/V。
func terminalClipboardAction(goos string, mods ui.Modifiers, key ui.Key) string {
	// 1. 修饰键必须恰好匹配，Ctrl+C 仍然是中断。
	want := ui.Ctrl | ui.Shift
	if goos == "darwin" {
		want = ui.Super
	}
	if mods != want {
		return ""
	}
	switch key {
	case ui.KeyC:
		return "copy"
	case ui.KeyV:
		return "paste"
	}
	return ""
}

// terminalKeyChars 是文字键在美式键盘上不按 Shift 与按 Shift 输出的字符。
var terminalKeyChars = func() map[ui.Key][2]rune {
	chars := map[ui.Key][2]rune{
		ui.KeySpace: {' ', ' '}, ui.KeyMinus: {'-', '_'}, ui.KeyEqual: {'=', '+'}, ui.KeyComma: {',', '<'},
		ui.KeyPeriod: {'.', '>'}, ui.KeySlash: {'/', '?'}, ui.KeySemicolon: {';', ':'}, ui.KeyQuote: {'\'', '"'},
		ui.KeyBracketLeft: {'[', '{'}, ui.KeyBracketRight: {']', '}'}, ui.KeyBackslash: {'\\', '|'}, ui.KeyBackquote: {'`', '~'},
	}
	for i := range 26 {
		chars[ui.KeyA+ui.Key(i)] = [2]rune{'a' + rune(i), 'A' + rune(i)}
	}
	for i, shifted := range ")!@#$%^&*(" {
		chars[ui.Key0+ui.Key(i)] = [2]rune{'0' + rune(i), shifted}
	}
	return chars
}()

// terminalControlChar 返回 Ctrl 加该字符对应的控制码。
func terminalControlChar(r rune) (string, bool) {
	switch {
	case r >= 'a' && r <= 'z':
		return string(rune(r - 'a' + 1)), true
	case r == ' ' || r == '`':
		return "\x00", true
	case r == '[':
		return "\x1b", true
	case r == '\\':
		return "\x1c", true
	case r == ']':
		return "\x1d", true
	case r == '/':
		return "\x1f", true
	}
	return "", false
}

// encodeTerminalKey 把一次按键编码成发给程序的转义序列。
// 返回 ok 为假表示该键没有序列：文字键等输入事件送来文字，其他键忽略。
func encodeTerminalKey(goos string, applicationKeys bool, mods ui.Modifiers, key ui.Key) (data string, ok bool) {
	// 1. Windows 的 AltGr 是 Ctrl+Alt，输出的是文字。
	shift, alt, ctrl := mods&ui.Shift != 0, mods&ui.Alt != 0, mods&ui.Ctrl != 0
	if goos == "windows" && alt && ctrl {
		return "", false
	}
	escape := ""
	if alt {
		escape = "\x1b"
	}
	// xterm 修饰键参数：1 加 Shift(1)、Alt(2)、Ctrl(4)。
	modifier := 1
	if shift {
		modifier += 1
	}
	if alt {
		modifier += 2
	}
	if ctrl {
		modifier += 4
	}
	// 2. 单字节键。
	switch key {
	case ui.KeyEnter:
		return escape + "\r", true
	case ui.KeyBackspace:
		return escape + "\x7f", true
	case ui.KeyEscape:
		return escape + "\x1b", true
	case ui.KeyTab:
		if shift {
			return "\x1b[Z", true
		}
		return escape + "\t", true
	}
	// 3. 方向键与 Home/End：带修饰键用 CSI 1;m，否则按应用光标模式选 SS3 或 CSI。
	if final, found := map[ui.Key]string{ui.KeyUp: "A", ui.KeyDown: "B", ui.KeyRight: "C", ui.KeyLeft: "D", ui.KeyHome: "H", ui.KeyEnd: "F"}[key]; found {
		switch {
		case modifier > 1:
			return "\x1b[1;" + strconv.Itoa(modifier) + final, true
		case applicationKeys:
			return "\x1bO" + final, true
		}
		return "\x1b[" + final, true
	}
	// 4. 编辑键与 F5-F12：CSI n ~，带修饰键追加 ;m。
	tildes := map[ui.Key]int{ui.KeyInsert: 2, ui.KeyDelete: 3, ui.KeyPageUp: 5, ui.KeyPageDown: 6, ui.KeyF5: 15, ui.KeyF6: 17, ui.KeyF7: 18, ui.KeyF8: 19, ui.KeyF9: 20, ui.KeyF10: 21, ui.KeyF11: 23, ui.KeyF12: 24}
	if number, found := tildes[key]; found {
		if modifier > 1 {
			return "\x1b[" + strconv.Itoa(number) + ";" + strconv.Itoa(modifier) + "~", true
		}
		return "\x1b[" + strconv.Itoa(number) + "~", true
	}
	// 5. F1-F4：SS3 P-S，带修饰键用 CSI 1;m。
	if key >= ui.KeyF1 && key <= ui.KeyF4 {
		final := string(rune('P' + key - ui.KeyF1))
		if modifier > 1 {
			return "\x1b[1;" + strconv.Itoa(modifier) + final, true
		}
		return "\x1bO" + final, true
	}
	// 6. 文字键只有带 Ctrl 或 Alt 时才有序列；macOS 的 Option 输出文字，交给输入事件。
	chars, isText := terminalKeyChars[key]
	if !isText {
		return "", false
	}
	if ctrl {
		control, found := terminalControlChar(chars[0])
		return escape + control, found
	}
	if alt && goos != "darwin" {
		if shift {
			return "\x1b" + string(chars[1]), true
		}
		return "\x1b" + string(chars[0]), true
	}
	return "", false
}
