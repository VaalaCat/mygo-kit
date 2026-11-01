---
title: List
description: 紧凑的虚拟列表，只构建可见的行，行高可以不一样。
---

# List

长列表用的虚拟列表。只构建可见的行，一万行也流畅，每行高度可以不同。它是 MyGo `ui.List` 加上工作台外观（无圆角、无内边距、编辑区底色）。

<Preview page="Virtual list" :height="360">

<<< @/../cmd/gallery/list.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows = ui.ListState{Key: func(i int) any { return items[i].ID }} // 由应用持有

kit.List(c, &rows, len(items), func(i int) {
	ui.Text(c, items[i].Title).Padding(0, 16)
}).Grow(1)
```

调用方负责：

- 持有 `ui.ListState`，行位置和滚动都记在里面。
- 设置 `Key`，给每一行稳定的身份。数据在前面插入或删除时，滚动位置才不会跳。
- 在 `row` 里画第 `i` 行。

## 相关组件

- 带“回到底部”按钮、跟随最新消息的，用 [MessageScroller](/zh/components/message-scroller)。
- 树形结构用 [Tree](/zh/components/tree)。
- 一行一个可点击条目，用 [ListRow](/zh/components/list-row) 来画行。

## API

### List

`func List(c *ui.Context, state *ui.ListState, count int, row func(int)) ui.Element`

- **state** `*ui.ListState`：列表状态，由调用方持有。
- **count** `int`：总行数。
- **row** `func(int)`：绘制第 `i` 行的函数，只会为可见的行调用。
