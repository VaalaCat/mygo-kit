---
title: Markdown
description: 可选择、可流式追加的 Markdown 阅读视图，带可复制的代码块、表格、任务列表和图片。
---

# Markdown

阅读 Markdown 的原生视图：标题、列表、表格、引用、脚注、行内代码、带语法颜色的代码块。可以跨段落拖选，支持流式追加（模型一边输出一边显示）。

<Preview page="Markdown" :height="520">

<<< @/../cmd/gallery/markdown.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Markdown(c, source, kit.MarkdownOpt{})
```

`source` 是完整的 Markdown 原文，每帧传入。放进 `ui.Scroll` 里就能滚动。

## 示例

### 流式输出

模型还在输出时设 `Streaming`，每帧传入当前累积的全文。没收尾的代码块，复制按钮是禁用的。

```go
kit.Markdown(c, reply.Text, kit.MarkdownOpt{Streaming: reply.IsRunning})
```

### 链接

默认点击链接会让 MyGo 用系统打开，只会打开安全的网址，不会执行任意协议。想自己处理就设 `OnLink`，设了之后所有链接都交给它。

```go
kit.MarkdownOpt{OnLink: func(url string) { openInBrowser(url) }}
```

### 图片

Markdown 里的图片单独成块显示。本库不读文件、不下载，图片数据由你按地址给出：`Image` 收到图片地址，返回准备好的 `ImageResource`。

```go
kit.MarkdownOpt{
	Image: func(url string) kit.ImageResource {
		if img, ok := cache[url]; ok {
			return kit.ImageResource{Source: img}
		}
		return kit.ImageResource{Loading: true} // 应用自己去下载，下载完放进 cache
	},
	OnImage: func(url string) { openPreview(url) },
}
```

没设 `Image` 时只显示替代文字。设了 `OnImage`，点击图片会调用它，例如打开 [ImagePreview](/components/image)。

### 任务列表

`- [x]` 和 `- [ ]` 显示成只读的任务框，它是文档的状态，不能点击切换。

### 选区和复制

默认情况下，Markdown 自己就是一个可选择的容器：选中的文字复制为纯文本，右键菜单的“复制 Markdown”复制整份源文。想让它和外面的内容连成一个选区，就设 `ShareSelection`，由外层的 `Selectable` 容器接管选择和右键菜单。

```go
ui.Column(c).Selectable().Children(func() {
	kit.Markdown(c, a, kit.MarkdownOpt{ShareSelection: true})
	kit.Markdown(c, b, kit.MarkdownOpt{ShareSelection: true})
})
```

### 公式和图表

数学式显示成 TeX 源码，本库不自带公式排版器，也不渲染 Mermaid。

## API

### Markdown

`func Markdown(c *ui.Context, source string, opt MarkdownOpt) ui.Element`

正文用 14/24，代码用 12/20，块间距 8。代码块带语言标题栏和复制按钮。

- **source** `string`：Markdown 原文。
- **opt** `MarkdownOpt`：见下。

### MarkdownOpt

- **Streaming** `bool`：流式显示。默认 `false`。
- **ShareSelection** `bool`：不自己建立选区和右键菜单，交给外层可选择的容器。默认 `false`。
- **OnLink** `func(string)`：点击链接时调用，参数是链接地址。设了之后不再由系统打开。
- **Image** `func(string) ImageResource`：按图片地址返回图片资源。
- **OnImage** `func(string)`：点击图片时调用，参数是图片地址。
- **CopyLabel** `string`：右键菜单里复制整份源文那一项的名称。默认“复制 Markdown”（英文 Copy Markdown）。
- **CodeCopyLabel** `string`：代码块复制按钮的提示和可访问名称。默认“复制代码”（英文 Copy code）。
