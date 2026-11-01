package main

import (
	"fmt"
	"strings"
	"time"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

// terminalDemo 保存画廊终端页的状态：示例输出、进度行和本地回显。
type terminalDemo struct {
	state    *kit.TerminalState
	progress int
	lastTick time.Time
	resizes  int
	// mouse 表示示例程序已开启鼠标上报：点击、拖动和滚轮会作为转义序列回显。
	mouse bool
}

var terminalPage = newTerminalDemo()

// newTerminalDemo 创建终端并写入带颜色、256 色、真彩色、样式和光标移动的示例。
func newTerminalDemo() *terminalDemo {
	d := &terminalDemo{state: kit.NewTerminalState(24, 80), progress: -1}
	var b strings.Builder
	b.WriteString("\x1b[1mkit.Terminal\x1b[0m  \x1b[2mvt emulator, theme ANSI palette\x1b[0m\r\n\r\n")
	// 1. 十六色：前景一行，背景一行。
	for i := 0; i < 16; i++ {
		fmt.Fprintf(&b, "\x1b[%dm%2d \x1b[0m", map[bool]int{true: 30 + i, false: 90 + i - 8}[i < 8], i)
	}
	b.WriteString("\r\n")
	for i := 0; i < 16; i++ {
		fmt.Fprintf(&b, "\x1b[48;5;%dm   \x1b[0m", i)
	}
	// 2. 256 色立方体与灰阶，真彩色渐变。
	b.WriteString("\r\n\r\n")
	for row := 0; row < 6; row++ {
		for i := 16 + row*36; i < 16+(row+1)*36; i++ {
			fmt.Fprintf(&b, "\x1b[48;5;%dm \x1b[0m", i)
		}
		b.WriteString("\r\n")
	}
	for i := 232; i < 256; i++ {
		fmt.Fprintf(&b, "\x1b[48;5;%dm  \x1b[0m", i)
	}
	b.WriteString("\r\n")
	for i := 0; i < 72; i++ {
		fmt.Fprintf(&b, "\x1b[48;2;%d;%d;%dm \x1b[0m", 255-i*3, i*3, 128+i)
	}
	// 3. 文字样式与光标移动：写完整行后左移覆盖其中几个字符。
	b.WriteString("\r\n\r\n\x1b[1mbold\x1b[0m \x1b[3mitalic\x1b[0m \x1b[4munderline\x1b[0m \x1b[7mreverse\x1b[0m \x1b[2mfaint\x1b[0m \x1b[9mstrike\x1b[0m 中文 ✓ ─│┌┐\r\n")
	b.WriteString("cursor: abcdefgh\x1b[4D\x1b[32mWXYZ\x1b[0m\r\n")
	d.state.Write([]byte(b.String()))
	return d
}

// restart 清屏并重新开始进度条。
func (d *terminalDemo) restart() {
	d.state.Write([]byte("\x1b[2J\x1b[H"))
	d.state.Write([]byte("$ "))
	d.progress = 0
}

// toggleMouse 开关鼠标上报：开启后点击、拖动和滚轮以 SGR 序列回显，按住 Shift 仍可选择文字和回看。
func (d *terminalDemo) toggleMouse() {
	d.mouse = !d.mouse
	if d.mouse {
		d.state.Write([]byte("\r\n\x1b[2mmouse reporting on: click, drag or scroll to see SGR reports; hold Shift to select\x1b[0m\r\n$ "))
		d.state.Write([]byte("\x1b[?1002h\x1b[?1006h"))
		return
	}
	d.state.Write([]byte("\x1b[?1002l\x1b[?1006l\r\n$ "))
}

// tick 每次推进一格进度，用回车覆盖同一行。
func (d *terminalDemo) tick() {
	const width = 40
	if d.progress < 0 {
		return
	}
	filled := d.progress * width / 100
	d.state.Write([]byte(fmt.Sprintf("\rdownloading [\x1b[32m%s\x1b[0m%s] %3d%%", strings.Repeat("#", filled), strings.Repeat(".", width-filled), d.progress)))
	if d.progress >= 100 {
		d.state.Write([]byte("\r\n$ "))
		d.progress = -1
		return
	}
	d.progress += 2
}

// echo 把键盘输入本地回显：回车换行并给出新提示符，退格擦除，控制键显示成 ^X。
func (d *terminalDemo) echo(data []byte) {
	var out []byte
	for _, v := range data {
		switch {
		case v == '\r':
			out = append(out, "\r\n$ "...)
		case v == 0x7f:
			out = append(out, "\b \b"...)
		case v == 0x1b:
			out = append(out, "^["...)
		case v < 0x20:
			out = append(out, '^', v+64)
		default:
			out = append(out, v)
		}
	}
	d.state.Write(out)
}

func init() {
	register(page{name: "Terminal", order: 41, view: func(c *ui.Context) {
		d := terminalPage
		// 1. 进度行每 100ms 推进一次，结束后停止请求帧。
		if d.progress >= 0 && c.Now().Sub(d.lastTick) >= 100*time.Millisecond {
			d.lastTick = c.Now()
			d.tick()
		}
		if d.progress >= 0 {
			c.After(100 * time.Millisecond)
		}
		// 2. 顶栏显示行列数，终端占满页面其余空间。
		ui.Column(c).Grow(1).Children(func() {
			rows, cols := d.state.Size()
			kit.Header(c, fmt.Sprintf(t(c, "%d rows × %d cols · resized %d times"), rows, cols, d.resizes), func() {
				if kit.Button(c, t(c, "Run progress")).Clicked() {
					d.restart()
				}
				label := "Mouse reporting: off"
				if d.mouse {
					label = "Mouse reporting: on"
				}
				if kit.Button(c, t(c, label)).Clicked() {
					d.toggleMouse()
				}
			})
			kit.Terminal(c, d.state, kit.TerminalOpt{
				OnInput:  d.echo,
				OnResize: func(int, int) { d.resizes++ },
			})
		})
	}})
}
