---
title: Facts
description: 两列事实组：左列标签，右列等宽的值，所有值共用一条左对齐线。
---

# Facts

“标签 加 值”的两列列表，适合详情面板里的模型、路径、用量这类信息。标签列固定宽 112，所有值左对齐在同一条线上。值可以选中复制。

画廊里没有单独的 Facts 页。

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Facts(c,
	kit.Fact{Label: "模型", Value: "claude-sonnet-5-5"},
	kit.Fact{Label: "目录", Value: "~/code/project"},
	kit.Fact{Label: "上下文", Value: "42%"},
)
```

左列是 12/16 次要色的标签，右列是 12/20 的等宽值，两列按首行基线对齐。值太长会换行。

## API

### Facts

`func Facts(c *ui.Context, facts ...Fact) ui.Element`

行与行之间不留间距。

- **facts** `...Fact`：按顺序显示的每一行。

### Fact

`type Fact struct{ Label, Value string }`

- **Label** `string`：标签，单行。
- **Value** `string`：值，可选中，等宽字体。
