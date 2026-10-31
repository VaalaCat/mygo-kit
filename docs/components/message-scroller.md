---
title: MessageScroller
description: 跟随最新消息的虚拟列表，上滚后出现“回到底部”按钮，滚到顶时通知加载更早的历史。
---

# MessageScroller

聊天消息用的列表。默认跟随最新消息；用户往上滚后，右下角浮出一个“回到底部”按钮。滚到顶部时通知你加载更早的历史。

<Preview page="MessageScroller" :height="460">

<<< @/../cmd/gallery/scroller.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows = ui.ListState{Key: func(i int) any { return messages[i].ID }} // 由应用持有

kit.MessageScroller(c, &rows, len(messages), kit.MessageScrollerOpt{
	OnReachTop: loadOlder,
}, func(i int) {
	ui.Text(c, messages[i].Text)
})
```

调用方负责：

- 持有 `ui.ListState`。
- 设置 `Key`，用稳定身份标识每条消息。前插历史时，滚动位置才会保持不动。
- 在 `OnReachTop` 里加载更早的消息，加到列表前面。

## 行为

- 组件会把 `state.FollowEnd` 设为 true，在尾部时新增内容跟随到底；允许动画时，新内容从下方滑入。
- `OnReachTop` 在用户从没到顶变为到顶时调用一次；离开顶部后再次回到顶部才会再次调用。
- 不在尾部时，右下角浮出圆形的向下箭头按钮，点击滚到底部。按钮的可访问名称是“回到底部”。
- 流式追加：直接改最后一条消息的文字，列表在尾部就会继续跟随。

## API

### MessageScroller

`func MessageScroller(c *ui.Context, state *ui.ListState, count int, opt MessageScrollerOpt, row func(int)) ui.Element`

行间距 16，列表内边距 16。

- **state** `*ui.ListState`：列表状态，由调用方持有。
- **count** `int`：消息总数。
- **opt** `MessageScrollerOpt`：见下。
- **row** `func(int)`：绘制第 `i` 条消息。

### MessageScrollerOpt

- **OnReachTop** `func()`：滚到顶部时调用一次。可以为 `nil`。
