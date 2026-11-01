---
title: MessageScroller
description: A virtualized list that follows the latest messages, shows a “Back to bottom” button after scrolling up, and requests earlier history at the top.
---

# MessageScroller

A list for chat messages. It follows the latest messages by default. When the user scrolls up, a “Back to bottom” button appears at the bottom right. Reaching the top notifies you to load earlier history.

<Preview page="MessageScroller" :height="460">

<<< @/../cmd/gallery/scroller.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var rows = ui.ListState{Key: func(i int) any { return messages[i].ID }} // Held by the application

kit.MessageScroller(c, &rows, len(messages), kit.MessageScrollerOpt{
	OnReachTop: loadOlder,
}, func(i int) {
	ui.Text(c, messages[i].Text)
})
```

The caller is responsible for:

- Holding `ui.ListState`.
- Setting `Key` to identify each message with a stable identity, so prepending history preserves the scroll position.
- Loading earlier messages in `OnReachTop` and prepending them to the list.

## Behavior

- The component sets `state.FollowEnd` to true. When at the end, it follows new content to the bottom. When motion is allowed, new content slides in from below.
- `OnReachTop` is called once when the user reaches the top from another scroll position. It is called again only after leaving the top and returning.
- When away from the end, a circular down-arrow button appears at the bottom right. Clicking it scrolls to the bottom. Its accessible name is “Back to bottom”.
- For streaming updates, change the text of the last message directly. The list keeps following while it is at the end.

## API

### MessageScroller

`func MessageScroller(c *ui.Context, state *ui.ListState, count int, opt MessageScrollerOpt, row func(int)) ui.Element`

Row spacing and list padding are both 16.

- **state** `*ui.ListState`: List state, held by the caller.
- **count** `int`: Total number of messages.
- **opt** `MessageScrollerOpt`: See below.
- **row** `func(int)`: Draws message `i`.

### MessageScrollerOpt

- **OnReachTop** `func()`: Called once when scrolling reaches the top. May be `nil`.
