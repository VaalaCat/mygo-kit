package kit

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// terminalCell 返回显示行 row、列 col 的格子文字和解析后的样式，供断言使用。
func terminalCell(s *TerminalState, theme Theme, row, col int) (string, terminalStyle) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cell := s.emu.CellAt(col, row)
	if cell == nil {
		return "", terminalStyle{}
	}
	return cell.Content, resolveTerminalStyle(cell.Style, theme)
}

// TestTerminalStateWrite 验证 SGR 颜色、属性、光标移动、回车覆盖、宽字符和滚动历史写入后的格子。
func TestTerminalStateWrite(t *testing.T) {
	// 1. 每项写入一段序列后检查屏幕文字，再检查指定格子的颜色或属性。
	theme := OneDark()
	colors := theme.Colors
	tests := []struct {
		name  string
		write string
		text  string
		row   int
		col   int
		cell  string
		want  terminalStyle
	}{
		{"basic color uses theme ansi", "\x1b[31mred\x1b[0m plain", "red plain", 0, 0, "r", terminalStyle{fg: colors.ANSI[1]}},
		{"reset returns default foreground", "\x1b[31mred\x1b[0m plain", "red plain", 0, 4, "p", terminalStyle{fg: colors.Text}},
		{"bright color maps to upper palette", "\x1b[92mok", "ok", 0, 0, "o", terminalStyle{fg: colors.ANSI[10]}},
		{"256 color uses xterm cube", "\x1b[38;5;196mX", "X", 0, 0, "X", terminalStyle{fg: ansiColor(196, theme)}},
		{"256 color below 16 uses theme", "\x1b[38;5;9mX", "X", 0, 0, "X", terminalStyle{fg: colors.ANSI[9]}},
		{"truecolor foreground and background", "\x1b[38;2;10;20;30;48;2;200;210;220mX", "X", 0, 0, "X", terminalStyle{fg: ui.RGB(10, 20, 30), bg: ui.RGB(200, 210, 220)}},
		{"bold italic underline", "\x1b[1;3;4mX", "X", 0, 0, "X", terminalStyle{fg: colors.Text, bold: true, italic: true, underline: true}},
		{"reverse swaps colors", "\x1b[7mX", "X", 0, 0, "X", terminalStyle{fg: colors.Editor, bg: colors.Text}},
		{"faint blends toward background", "\x1b[2mX", "X", 0, 0, "X", terminalStyle{fg: colors.Text.Mix(colors.Editor, .5)}},
		{"cursor forward then print", "abc\r\x1b[2Cx", "abx", 0, 2, "x", terminalStyle{fg: colors.Text}},
		{"cursor position then print", "\x1b[3;5Hz", "\n\n    z", 2, 4, "z", terminalStyle{fg: colors.Text}},
		{"carriage return overwrites progress", "50%\r100%", "100%", 0, 0, "1", terminalStyle{fg: colors.Text}},
		{"erase line", "hello\r\x1b[2Kx", "x", 0, 0, "x", terminalStyle{fg: colors.Text}},
		{"wide character takes two columns", "中a", "中a", 0, 2, "a", terminalStyle{fg: colors.Text}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := NewTerminalState(5, 20)
			defer state.Close()
			state.Write([]byte(tc.write))
			if got := state.Text(); got != tc.text {
				t.Fatalf("text = %q want %q", got, tc.text)
			}
			got, style := terminalCell(state, theme, tc.row, tc.col)
			if got != tc.cell || style != tc.want {
				t.Fatalf("cell(%d,%d) = %q %+v want %q %+v", tc.row, tc.col, got, style, tc.cell, tc.want)
			}
		})
	}
	// 2. 超出屏幕的行进入滚动历史，Text 只返回屏幕。
	state := NewTerminalState(3, 10)
	defer state.Close()
	for i := 0; i < 10; i++ {
		state.Write([]byte("line" + string(rune('0'+i)) + "\r\n"))
	}
	if got := state.Text(); !strings.Contains(got, "line9") || strings.Contains(got, "line0") {
		t.Fatalf("screen = %q", got)
	}
	if state.emu.ScrollbackLen() < 7 {
		t.Fatalf("history = %d", state.emu.ScrollbackLen())
	}
}

// TestTerminalStateRespondsToQueries 验证程序的光标位置查询由后台 goroutine 转给 OnInput，Write 不会阻塞。
func TestTerminalStateRespondsToQueries(t *testing.T) {
	// 1. 写入设备状态查询，等待应答经 OnInput 返回。
	state := NewTerminalState(5, 20)
	defer state.Close()
	got := make(chan string, 1)
	handler := func(data []byte) { got <- string(data) }
	state.input.Store(&handler)
	state.Write([]byte("ab\x1b[6n"))
	if reply := <-got; reply != "\x1b[1;3R" {
		t.Fatalf("reply = %q", reply)
	}
}

// TestEncodeTerminalKey 验证按键到 xterm 序列的编码。
func TestEncodeTerminalKey(t *testing.T) {
	tests := []struct {
		name  string
		goos  string
		app   bool
		mods  ui.Modifiers
		key   ui.Key
		want  string
		wantK bool
	}{
		{"enter", "linux", false, 0, ui.KeyEnter, "\r", true},
		{"backspace", "linux", false, 0, ui.KeyBackspace, "\x7f", true},
		{"tab", "linux", false, 0, ui.KeyTab, "\t", true},
		{"shift tab", "linux", false, ui.Shift, ui.KeyTab, "\x1b[Z", true},
		{"escape", "linux", false, 0, ui.KeyEscape, "\x1b", true},
		{"up", "linux", false, 0, ui.KeyUp, "\x1b[A", true},
		{"left", "linux", false, 0, ui.KeyLeft, "\x1b[D", true},
		{"up in application mode", "linux", true, 0, ui.KeyUp, "\x1bOA", true},
		{"ctrl right", "linux", false, ui.Ctrl, ui.KeyRight, "\x1b[1;5C", true},
		{"home", "linux", false, 0, ui.KeyHome, "\x1b[H", true},
		{"end in application mode", "linux", true, 0, ui.KeyEnd, "\x1bOF", true},
		{"page up", "linux", false, 0, ui.KeyPageUp, "\x1b[5~", true},
		{"page down", "linux", false, 0, ui.KeyPageDown, "\x1b[6~", true},
		{"insert", "linux", false, 0, ui.KeyInsert, "\x1b[2~", true},
		{"delete", "linux", false, 0, ui.KeyDelete, "\x1b[3~", true},
		{"shift delete", "linux", false, ui.Shift, ui.KeyDelete, "\x1b[3;2~", true},
		{"f1", "linux", false, 0, ui.KeyF1, "\x1bOP", true},
		{"f4", "linux", false, 0, ui.KeyF4, "\x1bOS", true},
		{"f5", "linux", false, 0, ui.KeyF5, "\x1b[15~", true},
		{"f12", "linux", false, 0, ui.KeyF12, "\x1b[24~", true},
		{"ctrl f1", "linux", false, ui.Ctrl, ui.KeyF1, "\x1b[1;5P", true},
		{"ctrl a", "linux", false, ui.Ctrl, ui.KeyA, "\x01", true},
		{"ctrl c", "linux", false, ui.Ctrl, ui.KeyC, "\x03", true},
		{"ctrl z", "linux", false, ui.Ctrl, ui.KeyZ, "\x1a", true},
		{"ctrl space", "linux", false, ui.Ctrl, ui.KeySpace, "\x00", true},
		{"ctrl left bracket", "linux", false, ui.Ctrl, ui.KeyBracketLeft, "\x1b", true},
		{"alt b", "linux", false, ui.Alt, ui.KeyB, "\x1bb", true},
		{"alt shift b", "linux", false, ui.Alt | ui.Shift, ui.KeyB, "\x1bB", true},
		{"alt ctrl b", "linux", false, ui.Alt | ui.Ctrl, ui.KeyB, "\x1b\x02", true},
		{"alt backspace", "linux", false, ui.Alt, ui.KeyBackspace, "\x1b\x7f", true},
		{"mac option types text", "darwin", false, ui.Alt, ui.KeyB, "", false},
		{"windows altgr types text", "windows", false, ui.Ctrl | ui.Alt, ui.KeyQ, "", false},
		{"plain letter waits for text", "linux", false, 0, ui.KeyA, "", false},
		{"shift letter waits for text", "linux", false, ui.Shift, ui.KeyA, "", false},
		{"ctrl digit has no sequence", "linux", false, ui.Ctrl, ui.Key1, "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := encodeTerminalKey(tc.goos, tc.app, tc.mods, tc.key)
			if got != tc.want || ok != tc.wantK {
				t.Fatalf("got %q,%v want %q,%v", got, ok, tc.want, tc.wantK)
			}
		})
	}
}

// TestTerminalClipboardAction 验证复制粘贴快捷键只在各平台约定的修饰键下触发。
func TestTerminalClipboardAction(t *testing.T) {
	tests := []struct {
		name string
		goos string
		mods ui.Modifiers
		key  ui.Key
		want string
	}{
		{"mac cmd c", "darwin", ui.Super, ui.KeyC, "copy"},
		{"mac cmd v", "darwin", ui.Super, ui.KeyV, "paste"},
		{"mac ctrl c is interrupt", "darwin", ui.Ctrl, ui.KeyC, ""},
		{"linux ctrl shift c", "linux", ui.Ctrl | ui.Shift, ui.KeyC, "copy"},
		{"linux ctrl shift v", "linux", ui.Ctrl | ui.Shift, ui.KeyV, "paste"},
		{"linux ctrl c is interrupt", "linux", ui.Ctrl, ui.KeyC, ""},
		{"linux other key", "linux", ui.Ctrl | ui.Shift, ui.KeyX, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := terminalClipboardAction(tc.goos, tc.mods, tc.key); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

// terminalFixture 是挂载终端的离屏窗口，记录输入字节、行列数变化和字符格尺寸。
type terminalFixture struct {
	mu      sync.Mutex
	state   *TerminalState
	inputs  []string
	resizes [][2]int
	metrics terminalMetrics
	tester  *ui.Tester
}

// newTerminalFixture 创建 width×height 的窗口，并点击终端使其获得焦点。
func newTerminalFixture(t *testing.T, width, height int) *terminalFixture {
	t.Helper()
	if err := RegisterFonts(); err != nil {
		t.Fatal(err)
	}
	f := &terminalFixture{state: NewTerminalState(24, 80)}
	t.Cleanup(f.state.Close)
	f.tester = ui.NewTester(func(c *ui.Context) {
		theme := OneLight()
		ApplyTheme(c, theme)
		f.metrics = terminalCellMetrics(c, theme)
		Terminal(c, f.state, TerminalOpt{
			OnInput:  func(data []byte) { f.mu.Lock(); f.inputs = append(f.inputs, string(data)); f.mu.Unlock() },
			OnResize: func(rows, cols int) { f.resizes = append(f.resizes, [2]int{rows, cols}) },
		})
	}, width, height)
	f.tester.ClickAt(float32(width)/2, float32(height)/2)
	return f
}

// typed 返回至今收到的全部输入字节。
func (f *terminalFixture) typed() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return strings.Join(f.inputs, "")
}

// TestTerminalResize 验证元素尺寸换算成行列数后通知应用并调整模拟器，改变窗口大小会再次通知。
func TestTerminalResize(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		resizeTo      [2]int
	}{
		{"initial layout", 640, 400, [2]int{}},
		{"window grows", 640, 400, [2]int{800, 600}},
		{"window shrinks", 640, 400, [2]int{320, 160}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newTerminalFixture(t, tc.width, tc.height)
			width, height := tc.width, tc.height
			if tc.resizeTo != [2]int{} {
				width, height = tc.resizeTo[0], tc.resizeTo[1]
				f.tester.SetSize(width, height)
			}
			wantRows := int((float32(height) - 2*f.metrics.padY) / f.metrics.cellHeight)
			wantCols := int((float32(width) - 2*f.metrics.padX) / f.metrics.cellWidth)
			if len(f.resizes) == 0 || f.resizes[len(f.resizes)-1] != [2]int{wantRows, wantCols} {
				t.Fatalf("resizes = %v want last %d×%d", f.resizes, wantRows, wantCols)
			}
			if rows, cols := f.state.Size(); rows != wantRows || cols != wantCols {
				t.Fatalf("state size = %d×%d want %d×%d", rows, cols, wantRows, wantCols)
			}
			if f.metrics.cellHeight != 16 {
				t.Fatalf("cell height = %v", f.metrics.cellHeight)
			}
		})
	}
	// 一份初始 24×80 的状态在窗口 640×400 下得到 24 行，行数与字符格高度一致。
	f := newTerminalFixture(t, 640, 400)
	if rows, _ := f.state.Size(); rows != 24 {
		t.Fatalf("rows = %d", rows)
	}
}

// TestTerminalResizeKeepsCursorVisible 验证面板变矮后光标所在行仍在可见区底部：多出的顶部行进入滚动历史，回看中的用户位置不被拉回底部。
func TestTerminalResizeKeepsCursorVisible(t *testing.T) {
	// 1. 每项先写入内容，再把窗口缩到 160 高，检查屏幕最后一行、光标行和滚动偏移。
	lines := func(n int) string {
		var b strings.Builder
		for i := range n {
			b.WriteString("line" + strconv.Itoa(i) + "\r\n")
		}
		return b.String()
	}
	tests := []struct {
		name       string
		write      string
		scrollBack int
		wantLast   string
		wantBottom bool
		wantScroll bool
	}{
		{"prompt below the new screen stays in view", lines(20) + "$ ", 0, "$", true, false},
		{"cursor already inside the new screen is untouched", "top\r\nsecond", 0, "second", false, false},
		{"user reading history keeps the position", lines(60) + "$ ", 5, "", true, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newTerminalFixture(t, 640, 400)
			f.state.Write([]byte(tc.write))
			f.tester.Frame()
			if tc.scrollBack > 0 {
				f.tester.Scroll(100, 100, 0, -float32(tc.scrollBack)*f.metrics.cellHeight)
			}
			f.tester.SetSize(640, 160)
			rows, _ := f.state.Size()
			if rows >= 24 {
				t.Fatalf("rows = %d, terminal did not shrink", rows)
			}
			f.state.mu.Lock()
			cursorRow := f.state.emu.CursorPosition().Y
			offset := f.state.scrollOffset
			f.state.mu.Unlock()
			if tc.wantBottom && cursorRow != rows-1 {
				t.Fatalf("cursor row = %d want %d", cursorRow, rows-1)
			}
			if tc.wantLast != "" {
				screen := strings.Split(f.state.Text(), "\n")
				if got := screen[len(screen)-1]; got != tc.wantLast {
					t.Fatalf("last screen line = %q want %q; screen %q", got, tc.wantLast, screen)
				}
			}
			if !tc.wantScroll && offset != 0 {
				t.Fatalf("scroll offset = %d want 0", offset)
			}
			if tc.wantScroll && offset < tc.scrollBack {
				t.Fatalf("scroll offset = %d want at least %d", offset, tc.scrollBack)
			}
		})
	}
}

// TestTerminalInput 验证文字、输入法上屏、控制键、方向键模式和回看后输入回到底部。
func TestTerminalInput(t *testing.T) {
	tests := []struct {
		name string
		run  func(*terminalFixture)
		want string
	}{
		{"typed text", func(f *terminalFixture) { f.tester.Type("ls -la") }, "ls -la"},
		{"unicode text", func(f *terminalFixture) { f.tester.Type("héllo") }, "héllo"},
		{"ime commit", func(f *terminalFixture) { f.tester.Compose("ni", 2); f.tester.Type("你好") }, "你好"},
		{"enter", func(f *terminalFixture) { f.tester.Key(0, ui.KeyEnter) }, "\r"},
		{"ctrl c", func(f *terminalFixture) { f.tester.Key(ui.Ctrl, ui.KeyC) }, "\x03"},
		{"arrow up", func(f *terminalFixture) { f.tester.Key(0, ui.KeyUp) }, "\x1b[A"},
		{"application cursor mode", func(f *terminalFixture) { f.state.Write([]byte("\x1b[?1h")); f.tester.Key(0, ui.KeyUp) }, "\x1bOA"},
		{"typed key does not double send", func(f *terminalFixture) { f.tester.TypeKey(0, ui.KeyA, "a") }, "a"},
		{"paste normalizes newlines", func(f *terminalFixture) { f.tester.SetClipboard("a\nb"); f.tester.Key(ui.Ctrl|ui.Shift, ui.KeyV) }, "a\rb"},
		{"bracketed paste", func(f *terminalFixture) {
			f.state.Write([]byte("\x1b[?2004h"))
			f.tester.SetClipboard("x")
			f.tester.Command("paste")
		}, "\x1b[200~x\x1b[201~"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newTerminalFixture(t, 640, 400)
			tc.run(f)
			if got := f.typed(); got != tc.want {
				t.Fatalf("input = %q want %q", got, tc.want)
			}
		})
	}
}

// TestTerminalSelectionAndScroll 验证拖动选择后复制、回看历史时输入回到底部、取消选择。
func TestTerminalSelectionAndScroll(t *testing.T) {
	f := newTerminalFixture(t, 640, 400)
	m := f.metrics
	f.state.Write([]byte("hello world\r\nsecond line"))
	f.tester.Frame()
	origin := func(row, col int) (float32, float32) {
		return m.padX + (float32(col)+.5)*m.cellWidth, m.padY + (float32(row)+.5)*m.cellHeight
	}
	// 1. 从第一行第 7 列拖到第二行第 5 列，复制得到跨行文字。
	x0, y0 := origin(0, 6)
	x1, y1 := origin(1, 5)
	f.tester.Press(x0, y0)
	f.tester.Move(x1, y1)
	f.tester.Release(x1, y1)
	f.tester.Key(ui.Ctrl|ui.Shift, ui.KeyC)
	if got := f.tester.Clipboard(); got != "world\nsecond" {
		t.Fatalf("clipboard = %q", got)
	}
	// 2. 单击取消选择，复制不覆盖剪贴板。
	f.tester.SetClipboard("kept")
	f.tester.Advance(time.Second)
	f.tester.ClickAt(x0, y0)
	f.tester.Key(ui.Ctrl|ui.Shift, ui.KeyC)
	if got := f.tester.Clipboard(); got != "kept" {
		t.Fatalf("clipboard after click = %q", got)
	}
	// 3. 写满历史后滚轮向上回看，输入后回到底部。
	for i := 0; i < 60; i++ {
		f.state.Write([]byte("history\r\n"))
	}
	f.tester.Frame()
	f.tester.Scroll(100, 100, 0, -5*m.cellHeight)
	f.state.mu.Lock()
	offset := f.state.scrollOffset
	f.state.mu.Unlock()
	if offset != 5 {
		t.Fatalf("scroll offset = %d", offset)
	}
	f.tester.Type("x")
	f.state.mu.Lock()
	offset = f.state.scrollOffset
	f.state.mu.Unlock()
	if offset != 0 {
		t.Fatalf("offset after input = %d", offset)
	}
}

// cellCenter 返回显示行 row、列 col 的格子中心，坐标相对终端元素。
func (f *terminalFixture) cellCenter(row, col int) (float32, float32) {
	return f.metrics.padX + (float32(col)+.5)*f.metrics.cellWidth, f.metrics.padY + (float32(row)+.5)*f.metrics.cellHeight
}

// terminalCursorReport 匹配模拟器对光标位置查询的应答。
var terminalCursorReport = regexp.MustCompile(`\x1b\[\d+;\d+R$`)

// drainInput 返回并清空至今收到的输入：模拟器的应答按顺序经过同一条管道，
// 追加一次光标位置查询，等它的应答到达就说明在它之前的字节都已交给 OnInput。
func (f *terminalFixture) drainInput(t *testing.T) string {
	t.Helper()
	f.state.Write([]byte("\x1b[6n"))
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); time.Sleep(time.Millisecond) {
		f.mu.Lock()
		got := strings.Join(f.inputs, "")
		if terminalCursorReport.MatchString(got) {
			f.inputs = nil
			f.mu.Unlock()
			return terminalCursorReport.ReplaceAllString(got, "")
		}
		f.mu.Unlock()
	}
	t.Fatal("cursor report did not arrive")
	return ""
}

// TestTerminalMouse 验证程序开启鼠标上报后，点击、拖动、移动、滚轮和修饰键按 SGR 编码交给 OnInput，未开启时不产生输入。
func TestTerminalMouse(t *testing.T) {
	// 1. 格子 (行 2, 列 3) 在 SGR 里是 1 起算的 X=4、Y=3；(行 2, 列 5) 是 X=6。
	tests := []struct {
		name  string
		modes string
		run   func(f *terminalFixture)
		want  string
	}{
		{"no mouse mode sends nothing", "", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.ClickAt(x, y)
			f.tester.Scroll(x, y, 0, -f.metrics.cellHeight)
		}, ""},
		{"click and release", "\x1b[?1000h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.ClickAt(x, y)
		}, "\x1b[<0;4;3M\x1b[<0;4;3m"},
		{"secondary button", "\x1b[?1000h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.RightClickAt(x, y)
		}, "\x1b[<2;4;3M\x1b[<2;4;3m"},
		{"x10 reports press only", "\x1b[?9h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.ClickAt(x, y)
		}, "\x1b[<0;4;3M"},
		{"ctrl modifier", "\x1b[?1000h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.ClickAtWith(ui.Ctrl, x, y)
		}, "\x1b[<16;4;3M\x1b[<16;4;3m"},
		{"wheel up and down", "\x1b[?1000h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.Scroll(x, y, 0, -f.metrics.cellHeight)
			f.tester.Scroll(x, y, 0, 2*f.metrics.cellHeight)
		}, "\x1b[<64;4;3M\x1b[<65;4;3M\x1b[<65;4;3M"},
		{"drag reports motion in button event mode", "\x1b[?1002h\x1b[?1006h", func(f *terminalFixture) {
			x0, y0 := f.cellCenter(2, 3)
			x1, y1 := f.cellCenter(2, 5)
			f.tester.Press(x0, y0)
			f.tester.Move(x1, y1)
			f.tester.Move(x1+1, y1)
			f.tester.Release(x1, y1)
		}, "\x1b[<0;4;3M\x1b[<32;6;3M\x1b[<0;6;3m"},
		{"drag is not reported in normal mode", "\x1b[?1000h\x1b[?1006h", func(f *terminalFixture) {
			x0, y0 := f.cellCenter(2, 3)
			x1, y1 := f.cellCenter(2, 5)
			f.tester.Press(x0, y0)
			f.tester.Move(x1, y1)
			f.tester.Release(x1, y1)
		}, "\x1b[<0;4;3M\x1b[<0;6;3m"},
		{"plain motion needs any event mode", "\x1b[?1002h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.Move(x, y)
		}, ""},
		{"plain motion in any event mode", "\x1b[?1003h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.Move(x, y)
			f.tester.Move(x+1, y)
		}, "\x1b[<35;4;3M"},
		{"mouse mode reset stops reporting", "\x1b[?1000h\x1b[?1006h\x1b[?1000l", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.ClickAt(x, y)
		}, ""},
		{"shift click stays local", "\x1b[?1000h\x1b[?1006h", func(f *terminalFixture) {
			x, y := f.cellCenter(2, 3)
			f.tester.ClickAtWith(ui.Shift, x, y)
		}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newTerminalFixture(t, 640, 400)
			f.state.Write([]byte(tc.modes))
			f.tester.Frame()
			f.drainInput(t)
			tc.run(f)
			if got := f.drainInput(t); got != tc.want {
				t.Fatalf("input = %q want %q", got, tc.want)
			}
		})
	}
}

// TestTerminalMouseKeepsLocalActionsWithShift 验证开启鼠标上报时，Shift+拖动仍然本地选择，Shift+滚轮仍然回看历史，且都不发给程序。
func TestTerminalMouseKeepsLocalActionsWithShift(t *testing.T) {
	f := newTerminalFixture(t, 640, 400)
	for i := 0; i < 60; i++ {
		f.state.Write([]byte("history\r\n"))
	}
	f.state.Write([]byte("\x1b[2J\x1b[Hhello world\x1b[?1000h\x1b[?1006h"))
	f.tester.Frame()
	// 测试器的 Press 不带修饰键，直接把带 Shift 的指针和滚轮事件交给终端。
	x0, y0 := f.cellCenter(0, 6)
	x1, y1 := f.cellCenter(0, 10)
	f.drainInput(t)
	for _, ev := range []ui.InputEvent{
		{Kind: ui.InputPointerDown, Mods: ui.Shift, X: x0, Y: y0, Clicks: 1},
		{Kind: ui.InputPointerMove, Mods: ui.Shift, X: x1, Y: y1},
		{Kind: ui.InputPointerUp, Mods: ui.Shift, X: x1, Y: y1},
	} {
		if !f.state.handlePointer(ev, f.metrics) {
			t.Fatalf("pointer %v not taken", ev.Kind)
		}
	}
	f.state.handleScroll(ui.InputEvent{Kind: ui.InputScroll, Mods: ui.Shift, X: x0, Y: y0, DY: -3 * f.metrics.cellHeight}, f.metrics)
	f.tester.Key(ui.Ctrl|ui.Shift, ui.KeyC)
	if got := f.tester.Clipboard(); got != "world" {
		t.Fatalf("clipboard = %q", got)
	}
	f.state.mu.Lock()
	offset := f.state.scrollOffset
	f.state.mu.Unlock()
	if offset != 3 {
		t.Fatalf("scroll offset = %d", offset)
	}
	if got := f.drainInput(t); got != "" {
		t.Fatalf("input = %q", got)
	}
}

// TestTerminalWordAndLineSelection 验证没有鼠标上报时，双击选中单词，三击选中整行。
func TestTerminalWordAndLineSelection(t *testing.T) {
	tests := []struct {
		name   string
		clicks int
		row    int
		col    int
		want   string
	}{
		{"double click selects word", 2, 0, 7, "world"},
		{"double click at word start", 2, 0, 6, "world"},
		{"double click selects path", 2, 1, 5, "/usr/bin/go"},
		{"double click stops at quote", 2, 2, 3, "flag"},
		{"triple click selects line", 3, 1, 5, "/usr/bin/go --flag"},
		{"single click selects nothing", 1, 0, 7, "kept"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newTerminalFixture(t, 640, 400)
			f.state.Write([]byte("hello world\r\n/usr/bin/go --flag\r\n\"flag\" (x)"))
			f.tester.Frame()
			f.tester.SetClipboard("kept")
			x, y := f.cellCenter(tc.row, tc.col)
			for range tc.clicks {
				f.tester.ClickAt(x, y)
			}
			f.tester.Key(ui.Ctrl|ui.Shift, ui.KeyC)
			if got := f.tester.Clipboard(); got != tc.want {
				t.Fatalf("clipboard = %q want %q", got, tc.want)
			}
		})
	}
}

// TestTerminalPaint 验证颜色、背景和光标真实绘制到像素：红色文字所在格有红色像素，聚焦光标是实心块，失焦为空心框。
func TestTerminalPaint(t *testing.T) {
	theme := OneLight()
	tests := []struct {
		name    string
		write   string
		focused bool
		pixel   func(m terminalMetrics) (int, int)
		want    func(*testing.T, [4]uint8)
	}{
		{"background cell is filled", "\x1b[44m \x1b[0m", false, func(m terminalMetrics) (int, int) { return int(m.padX + m.cellWidth/2), int(m.padY + 3) }, func(t *testing.T, got [4]uint8) {
			want := theme.Colors.ANSI[4]
			if got != [4]uint8{want.R, want.G, want.B, 255} {
				t.Fatalf("background pixel = %v want %v", got, want)
			}
		}},
		{"focused block cursor is solid", "", true, func(m terminalMetrics) (int, int) { return int(m.padX + 1), int(m.padY + 8) }, func(t *testing.T, got [4]uint8) {
			want := theme.Colors.Text
			if got != [4]uint8{want.R, want.G, want.B, 255} {
				t.Fatalf("cursor pixel = %v want %v", got, want)
			}
		}},
		{"unfocused cursor is hollow", "", false, func(m terminalMetrics) (int, int) { return int(m.padX + m.cellWidth/2), int(m.padY + 8) }, func(t *testing.T, got [4]uint8) {
			want := theme.Colors.Editor
			if got != [4]uint8{want.R, want.G, want.B, 255} {
				t.Fatalf("hollow cursor center = %v want %v", got, want)
			}
		}},
		{"hidden cursor draws nothing", "\x1b[?25l", true, func(m terminalMetrics) (int, int) { return int(m.padX + 1), int(m.padY + 8) }, func(t *testing.T, got [4]uint8) {
			want := theme.Colors.Editor
			if got != [4]uint8{want.R, want.G, want.B, 255} {
				t.Fatalf("hidden cursor pixel = %v want %v", got, want)
			}
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newTerminalFixture(t, 640, 400)
			if !tc.focused {
				f.tester.SetFocused(false)
			}
			f.state.Write([]byte(tc.write))
			f.tester.Frame()
			x, y := tc.pixel(f.metrics)
			px := f.tester.Image().RGBAAt(x, y)
			tc.want(t, [4]uint8{px.R, px.G, px.B, px.A})
		})
	}
}
