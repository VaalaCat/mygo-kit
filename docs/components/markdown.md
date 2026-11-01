---
title: Markdown
description: A selectable Markdown reader with streaming updates, copyable code blocks, tables, task lists, and images.
---

# Markdown

A native view for reading Markdown: headings, lists, tables, blockquotes, footnotes, inline code, and syntax-highlighted code blocks. Supports selection across paragraphs and streaming updates as the model produces output.

<Preview page="Markdown" :height="520">

<<< @/../cmd/gallery/markdown.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Markdown(c, source, kit.MarkdownOpt{})
```

`source` is the complete Markdown source, passed in each frame. Place it inside `ui.Scroll` to make it scrollable.

## Examples

### Streaming output

Set `Streaming` while the model is still producing output, and pass the full text accumulated so far each frame. The copy button is disabled for an unfinished code block.

```go
kit.Markdown(c, reply.Text, kit.MarkdownOpt{Streaming: reply.IsRunning})
```

### Links

By default, clicking a link asks MyGo to open it through the system. Only safe web URLs are opened; arbitrary protocols are not executed. Set `OnLink` to handle links yourself. Once set, it receives all links.

```go
kit.MarkdownOpt{OnLink: func(url string) { openInBrowser(url) }}
```

### Images

Markdown images appear as separate blocks. The library does not read files or download images. Supply image data by URL: `Image` receives the image URL and returns a prepared `ImageResource`.

```go
kit.MarkdownOpt{
	Image: func(url string) kit.ImageResource {
		if img, ok := cache[url]; ok {
			return kit.ImageResource{Source: img}
		}
		return kit.ImageResource{Loading: true} // The application downloads the image and adds it to cache
	},
	OnImage: func(url string) { openPreview(url) },
}
```

Without `Image`, only the alternative text is shown. If `OnImage` is set, clicking an image calls it, for example to open an [ImagePreview](/components/image).

### Task lists

`- [x]` and `- [ ]` appear as read-only checkboxes. They represent the document's state and cannot be toggled by clicking.

### Selection and copying

By default, Markdown is its own selectable container. Selected text is copied as plain text, while “Copy Markdown” in the context menu copies the entire source. To share a selection with surrounding content, set `ShareSelection` and let the outer `Selectable` container handle selection and the context menu.

```go
ui.Column(c).Selectable().Children(func() {
	kit.Markdown(c, a, kit.MarkdownOpt{ShareSelection: true})
	kit.Markdown(c, b, kit.MarkdownOpt{ShareSelection: true})
})
```

### Formulas and diagrams

Math is displayed as TeX source. The library does not include a math typesetter or render Mermaid.

## API

### Markdown

`func Markdown(c *ui.Context, source string, opt MarkdownOpt) ui.Element`

Body text uses 14/24 sizing, code uses 12/20, and block spacing is 8. Code blocks include a language header and copy button.

- **source** `string`: Markdown source.
- **opt** `MarkdownOpt`: See below.

### MarkdownOpt

- **Streaming** `bool`: Displays streaming output. Defaults to `false`.
- **ShareSelection** `bool`: Leaves selection and the context menu to an outer selectable container instead of creating its own. Defaults to `false`.
- **OnLink** `func(string)`: Called when a link is clicked, with its URL as the argument. When set, the system no longer opens links.
- **Image** `func(string) ImageResource`: Returns an image resource for an image URL.
- **OnImage** `func(string)`: Called when an image is clicked, with its URL as the argument.
- **CopyLabel** `string`: The context menu label for copying the entire source. Defaults to “Copy Markdown” (localized in Chinese).
- **CodeCopyLabel** `string`: The code block copy button's tooltip and accessible name. Defaults to “Copy code” (localized in Chinese).
