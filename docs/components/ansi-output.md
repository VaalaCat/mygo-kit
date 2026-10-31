---
title: ANSIOutput
description: 显示命令输出：解释颜色，剥离其他终端控制码，虚拟滚动，可选择复制。
---

# ANSIOutput

显示命令行输出的内容框。只解释 SGR 颜色和样式（粗体、下划线、16/256 色、真彩色），其他终端控制码（光标移动、OSC 等）一律剥离，从不执行。是虚拟列表，长输出也流畅，横向可滚动。

要显示一个完整的、可以交互的终端，用 [Terminal](/components/terminal)。

<Preview page="Output" :height="320">

<<< @/../cmd/gallery/output.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows = ui.ListState{FollowEnd: true} // 由应用持有，FollowEnd 让它跟随底部

kit.ANSIOutput(c, output, &rows, kit.ANSIOutputOpt{}).Grow(1)
```

`output` 是到目前为止的全部输出，每帧传入。是否跟随最新输出由 `ListState.FollowEnd` 决定。

## 示例

### 原样显示

默认模式解释颜色、剥离控制码。`Plain` 为真时，不做任何解释，按原文显示，复制的也是原文，一个字节不差。

```go
kit.ANSIOutput(c, log, &rows, kit.ANSIOutputOpt{Plain: true})
```

## 选择和复制

可以拖选，选区跨过已经构建的可见行。没有构建的历史行（滚出视野的）选不到，要用工具栏的复制按钮复制全文。默认模式下，复制的全文已经去掉了控制码。

## API

### ANSIOutput

`func ANSIOutput(c *ui.Context, source string, state *ui.ListState, opt ANSIOutputOpt) ui.Element`

- **source** `string`：输出全文。内容和主题没变时，解析结果会被缓存。
- **state** `*ui.ListState`：列表状态，由调用方持有。
- **opt** `ANSIOutputOpt`：见下。

### ANSIOutputOpt

- **Title** `string`：标题栏文字。默认“输出”（英文 Output）。
- **CopyLabel** `string`：复制按钮的提示和可访问名称。默认“复制输出”（英文 Copy output）。
- **Plain** `bool`：原样显示并原样复制。默认 `false`。
