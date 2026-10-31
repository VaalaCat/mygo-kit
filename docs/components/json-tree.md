---
title: JSON
description: 可键盘展开、选择和复制的虚拟 JSON 树。
---

# JSON

把一个 JSON 值显示成可展开的树：对象和数组只显示成员数，叶子显示值。是虚拟列表，只构建可见的行，大 JSON 也流畅。

<Preview page="JSON" :height="360">

<<< @/../cmd/gallery/json.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
state, err := kit.NewJSONState(data) // 解析一次，状态由应用持有
if err != nil {
	// data 不是合法的 JSON
}

kit.JSONTree(c, state, kit.JSONOpt{}).Grow(1)
```

`NewJSONState` 解析任意顶级值，包括 `null`、数组和标量；数字保留原文的写法。内容变了就调用 `state.Set(newData)`，成功后替换树，仍然存在的展开路径会保留，用户收起的根节点不会被重新展开。

## 操作

- 点击一行选中它，有子节点的同时切换展开。
- 上下键移动选中。右键展开，已展开时进入第一个子节点；左键收起，已收起时回到父节点。Enter 或 Space 切换展开。
- 选中后按 Cmd+C（Linux 和 Windows 上是 Ctrl+C）复制这个节点的 JSON；工具栏的复制按钮复制整个值。

## API

### NewJSONState

`func NewJSONState(data []byte) (*JSONState, error)`

解析 JSON 并创建状态。根节点默认展开。

### JSONState.Set

`func (s *JSONState) Set(data []byte) error`

解析新内容。出错时返回错误，旧的树保持不变。内容里有多个 JSON 值也会报错。

### JSONTree

`func JSONTree(c *ui.Context, s *JSONState, opt JSONOpt) ui.Element`

带圆角细边框的内容框，顶部是标题栏和复制按钮，下面是树。每行 20 高，每级缩进 16。

- **s** `*JSONState`：由调用方持有的状态。
- **opt** `JSONOpt`：见下。

### JSONOpt

- **Title** `string`：标题栏文字。默认 `JSON`。
- **CopyLabel** `string`：复制按钮的提示和可访问名称。默认“复制 JSON”（英文 Copy JSON）。
