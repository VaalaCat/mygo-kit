---
title: Terminal
description: 终端视图：显示 PTY 的输出，把键盘输入回传给应用，支持颜色、历史、选择和鼠标上报。
---

# Terminal

一个完整的终端视图。它只负责显示和输入：你的应用把 PTY（或任何进程）的输出写进 `TerminalState`，通过 `OnInput` 收到用户的按键，再写回 PTY。启动进程、PTY 和重连都由你的应用负责。

<Preview page="Terminal" :height="520">

<<< @/../cmd/gallery/terminal.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
// 创建一次，状态由应用持有；不再使用时调用 state.Close()
state := kit.NewTerminalState(24, 80)

// 把 PTY 的输出写进来，可以在任意 goroutine 调用
go func() {
	buf := make([]byte, 4096)
	for {
		n, err := pty.Read(buf)
		state.Write(buf[:n])
		if err != nil {
			return
		}
	}
}()
```

```go
// 每帧
kit.Terminal(c, state, kit.TerminalOpt{
	OnInput:  func(data []byte) { pty.Write(data) },
	OnResize: func(rows, cols int) { resizePTY(rows, cols) },
})
```

`TerminalState.Write` 没有返回值，不是 `io.Writer`，所以不能直接用 `io.Copy`，像上面这样循环读写。行列数由面板尺寸换算，变化时调用 `OnResize`，调用方要把新尺寸同步给 PTY。

## 能显示什么

- 等宽 12/16（Lilex）。16 色用主题里的 `Colors.ANSI`，另外支持 256 色和真彩色。
- 粗体、斜体、下划线、反显、淡色、删除线，光标、滚动历史。
- 选择复制：双击选词，三击选行。macOS 上 Cmd+C / Cmd+V，其他系统是 Ctrl+Shift+C / Ctrl+Shift+V，普通的 Ctrl+C 仍然发给程序当中断。

## 鼠标上报

程序开启鼠标上报（DEC 模式 9、1000 到 1003，SGR 编码）后，点击、拖动、移动和滚轮都交给程序。按住 Shift 仍然可以本地选择文字、回看历史。Shift+PageUp / PageDown 翻看历史，不发给程序。

## 要注意

- 元素默认填满父容器剩余空间，背景是编辑器色，点击获得焦点，焦点只体现在光标上。
- `OnInput` 的数据包括键盘、输入法上屏的文字和粘贴编码后的字节。程序查询光标位置等自动应答也会从这里发出，那时是在后台 goroutine 调用的，写 PTY 要能接受这种调用。
- `NewTerminalState` 会启动一个转发模拟器自动应答的 goroutine，所以不再使用时一定要调用 `Close`。关闭后的 `Write` 被忽略。
- `Write` 在回看历史时不会把视图拉回底部。

## API

### NewTerminalState

`func NewTerminalState(rows, cols int) *TerminalState`

创建 `rows` 行 `cols` 列的空终端。

### TerminalState

由应用持有，可以在任意 goroutine 里调用 `Write`。

- **Write(data []byte)**：喂入 PTY 输出。
- **Close()**：结束转发自动应答的 goroutine，之后的 `Write` 被忽略。
- **Size() (rows, cols int)**：当前行数和列数。
- **Text() string**：当前屏幕的纯文本，每行去掉行尾空格，末尾空行去掉。
- **Title() string**：程序通过 OSC 设置的窗口标题。

### Terminal

`func Terminal(c *ui.Context, state *TerminalState, opt TerminalOpt) ui.Element`

屏幕文字同时作为元素的无障碍值，读屏和自动化测试都能读到当前内容。

### TerminalOpt

- **OnInput** `func(data []byte)`：接收键盘、输入法上屏文字和粘贴编码后的字节，写给 PTY。可能在后台 goroutine 里调用。
- **OnResize** `func(rows, cols int)`：面板尺寸换算出的行列数变化时调用。
