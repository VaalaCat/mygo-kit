---
title: Patch
description: 显示已经解析好的补丁片段：增删行着色、可选行号、复制调用方给的原文。
---

# Patch

显示一段补丁片段，例如工具调用返回的修改。它不比较文件，也不推测行号，显示的就是你给的每一行。要比较两份完整文本，用 [Diff](/components/diff)。

<Preview page="Patch" :height="260">

<<< @/../cmd/gallery/patch.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows ui.ListState // 由应用持有

n := 42
kit.PatchLines(c, &rows, []kit.PatchLine{
	{Kind: kit.PatchHeader, Text: "@@ component.go @@"},
	{Kind: kit.PatchDeleted, Text: "return oldValue"},
	{Kind: kit.PatchAdded, Text: "return newValue", NewLine: &n},
}, kit.PatchOpt{Source: patchText}).Grow(1)
```

调用方负责：

- 把补丁解析成一行一行的 `PatchLine`。
- 行号。`OldLine`、`NewLine` 是 `*int`，传 `nil` 就是空白，组件不会去猜。
- 复制用的原文。复制按钮复制的是 `PatchOpt.Source`，不会从行片段拼回补丁，所以不会伪造出完整文件。

## API

### PatchLines

`func PatchLines(c *ui.Context, state *ui.ListState, lines []PatchLine, opt PatchOpt) ui.Element`

带圆角细边框的内容框，顶部是标题栏和复制按钮，下面是虚拟列表，可以选中文字，横向可滚动。

- **state** `*ui.ListState`：列表状态，由调用方持有。
- **lines** `[]PatchLine`：全部行。
- **opt** `PatchOpt`：见下。

### PatchLine

- **Kind** `PatchKind`：行的种类。
- **Text** `string`：行内容，不含 `+`、`-` 前缀，符号由组件按 `Kind` 画。
- **OldLine**、**NewLine** `*int`：旧文件和新文件的行号。`nil` 显示空白。

### PatchKind

- **PatchContext**：上下文行，无底色。
- **PatchAdded**：新增行，`+`，成功色浅底。
- **PatchDeleted**：删除行，`−`，危险色浅底。
- **PatchHeader**：块标题，强调色文字，面板底色。

### PatchOpt

- **Title** `string`：标题栏文字。默认 `Patch`。
- **CopyLabel** `string`：复制按钮的提示和可访问名称。默认“复制输出”（英文 Copy output）。
- **Source** `string`：点击复制时写入剪贴板的完整原文。
