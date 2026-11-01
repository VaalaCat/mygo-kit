---
title: AttachmentChip
description: A chip for a pending attachment, with a thumbnail, filename, and remove button.
---

# AttachmentChip

A small chip above the input area for an attachment waiting to be sent: a 20×20 thumbnail, a truncated filename, and a remove button. It is 28 high and at most 240 wide.

<Preview page="AttachmentChip" :height="240">

<<< @/../cmd/gallery/attachment.go

</Preview>

## Usage

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

The caller prepares both the thumbnails and the attachment list; the component does not read files. Update your own list in `onRemove`.

## Examples

### No thumbnail or removal

With `thumb` set to `nil`, a file icon is shown. With `onRemove` set to `nil`, the remove button is hidden.

```go
kit.AttachmentChip(c, "notes.txt", nil, nil)
```

When displaying several chips, use `ui.Row(c).Wrap()` to let them wrap automatically.

## API

### AttachmentChip

`func AttachmentChip(c *ui.Context, name string, thumb ui.ImageSource, onRemove func()) ui.Element`

- **name** `string`: Filename, truncated to one line. Also used as the accessible name.
- **thumb** `ui.ImageSource`: Thumbnail. Shows a file icon when `nil`.
- **onRemove** `func()`: Called when the remove button is clicked. The button is hidden when `nil`.
