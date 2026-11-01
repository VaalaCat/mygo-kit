---
title: Terminal
description: A terminal view that displays PTY output and sends keyboard input back to the application, with colors, scrollback, selection, and mouse reporting.
---

# Terminal

A complete terminal view. It handles display and input: your application writes output from a PTY (or any process) into `TerminalState`, receives user keystrokes through `OnInput`, and writes them back to the PTY. Your application manages process startup, the PTY, and reconnection.

<Preview page="Terminal" :height="520">

<<< @/../cmd/gallery/terminal.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
// Create once and keep the state in the application; call state.Close() when no longer needed
state := kit.NewTerminalState(24, 80)

// Write PTY output here; callable from any goroutine
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
// Every frame
kit.Terminal(c, state, kit.TerminalOpt{
	OnInput:  func(data []byte) { pty.Write(data) },
	OnResize: func(rows, cols int) { resizePTY(rows, cols) },
})
```

`TerminalState.Write` has no return value and does not implement `io.Writer`, so you cannot use `io.Copy` directly. Read and write in a loop as shown above. Row and column counts are calculated from the panel dimensions. When they change, `OnResize` is called, and the caller must apply the new size to the PTY.

## Display features

- Monospaced 12/16 text (Lilex). The 16-color palette uses the theme's `Colors.ANSI`. 256 colors and true color are also supported.
- Bold, italic, underline, reverse video, dim text, strikethrough, a cursor, and scrollback.
- Selection and clipboard: double-click to select a word, triple-click to select a line. Use Cmd+C / Cmd+V on macOS and Ctrl+Shift+C / Ctrl+Shift+V on other systems. Plain Ctrl+C is still sent to the program as an interrupt.

## Mouse reporting

When the program enables mouse reporting (DEC modes 9 and 1000 through 1003, with SGR encoding), clicks, drags, pointer movements, and scrolling are sent to the program. Hold Shift to select text locally or browse scrollback. Shift+PageUp / PageDown navigates scrollback without being sent to the program.

## Notes

- By default, the element fills the remaining space in its parent, uses the editor background color, and receives focus on click. Only the cursor indicates focus.
- `OnInput` receives keyboard input, committed IME text, and encoded paste bytes. Automatic responses to program queries, such as cursor position queries, also go through this callback. Those calls happen on a background goroutine, so writing to the PTY must support calls from that goroutine.
- `NewTerminalState` starts a goroutine that forwards the emulator's automatic responses. Always call `Close` when the state is no longer needed. Calls to `Write` after closure are ignored.
- `Write` does not scroll the view back to the bottom while browsing scrollback.

## API

### NewTerminalState

`func NewTerminalState(rows, cols int) *TerminalState`

Creates an empty terminal with `rows` rows and `cols` columns.

### TerminalState

Held by the application. `Write` can be called from any goroutine.

- **Write(data []byte)**: Feeds in PTY output.
- **Close()**: Stops the goroutine forwarding automatic responses. Subsequent `Write` calls are ignored.
- **Size() (rows, cols int)**: The current row and column counts.
- **Text() string**: Plain text from the current screen, with trailing spaces removed from each line and trailing blank lines removed.
- **Title() string**: The window title set by the program through OSC.

### Terminal

`func Terminal(c *ui.Context, state *TerminalState, opt TerminalOpt) ui.Element`

The screen text is also the element's accessible value, so screen readers and automated tests can read its current content.

### TerminalOpt

- **OnInput** `func(data []byte)`: Receives keyboard input, committed IME text, and encoded paste bytes to write to the PTY. May be called from a background goroutine.
- **OnResize** `func(rows, cols int)`: Called when the row and column counts calculated from the panel dimensions change.
