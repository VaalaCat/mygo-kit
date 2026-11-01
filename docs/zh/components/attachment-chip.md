---
title: AttachmentChip
description: 待发送附件的一枚标签：缩略图、文件名和移除按钮。
---

# AttachmentChip

输入区上方显示待发送附件的小标签：20×20 缩略图、截断的文件名和一个移除按钮，整体 28 高，最宽 240。

<Preview page="AttachmentChip" :height="240">

<<< @/../cmd/gallery/attachment.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
for i, name := range names {
	kit.AttachmentChip(c.Key(name), name, thumbs[i], func() {
		names = append(names[:i], names[i+1:]...)
	})
}
```

缩略图和附件列表都由调用方准备，组件不读文件。`onRemove` 里改的是你自己的列表。

## 示例

### 没有缩略图或不能移除

`thumb` 传 `nil` 时显示文件图标；`onRemove` 传 `nil` 时不显示移除按钮。

```go
kit.AttachmentChip(c, "notes.txt", nil, nil)
```

放多个时，用 `ui.Row(c).Wrap()` 让它们自动换行。

## API

### AttachmentChip

`func AttachmentChip(c *ui.Context, name string, thumb ui.ImageSource, onRemove func()) ui.Element`

- **name** `string`：文件名，单行截断，同时是可访问名称。
- **thumb** `ui.ImageSource`：缩略图。`nil` 时显示文件图标。
- **onRemove** `func()`：点击移除按钮时调用。`nil` 时不显示按钮。
