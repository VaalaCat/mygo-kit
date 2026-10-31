---
title: Diff
description: 原生的完整文件 diff：增删行号、选择和复制，状态由 MyGo 的 DiffState 持有。
---

# Diff

两份文本的对比，可以左右并排（`editor.DiffSplit`，默认），也可以合并成一列（`editor.DiffUnified`）。差异算法、选择和复制都是 MyGo 编辑器的 `DiffView`，本库只设置颜色和字体，让它和其他面板一致。

只想显示一段已经解析好的补丁，不是比较两份完整文本，用 [PatchLines](/components/patch)。

<Preview page="Diff" :height="320">

<<< @/../cmd/gallery/diff.go

</Preview>

## 用法

```go
import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui/editor"
)
```

```go
var diff = editor.DiffState{
	Old:             editor.NewBuffer(oldText),
	New:             editor.NewBuffer(newText),
	Layout:          editor.DiffUnified,
	ExpandUnchanged: true,
}

kit.Diff(c, &diff).Grow(1)
```

`DiffState` 由调用方持有，布局、展开未改动的行等选项都在它上面设置，见 MyGo 的 [diff viewer 文档](https://github.com/egoist/mygo/blob/main/docs/ui/diff-viewer.md)。

每帧 `Diff` 会改写 `state.Theme`、`state.Font` 和 `state.LineHeight`（1.5），所以这几项请不要在 `DiffState` 里自己设置。

## API

### Diff

`func Diff(c *ui.Context, state *editor.DiffState) ui.Element`

增行用成功色，删行用危险色，修改用警告色，行号用占位色，字体是主题的等宽字体，背景是编辑区底色。

- **state** `*editor.DiffState`：MyGo 的 diff 状态，由调用方持有。
