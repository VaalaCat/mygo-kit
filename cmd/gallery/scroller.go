package main

import (
	"fmt"
	"strings"
	"time"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

// scrollerDemo 保存画廊里的 1000 条消息和流式追加状态。
type scrollerDemo struct {
	messages  []string
	streaming bool
	lastTick  time.Time
	topLoads  int
}

var (
	scrollerRows = ui.ListState{Key: func(i int) any { return i }}
	scroller     = newScrollerDemo()
)

// newScrollerDemo 生成长短不一的示例消息。
func newScrollerDemo() *scrollerDemo {
	d := &scrollerDemo{}
	for i := 0; i < 1000; i++ {
		d.messages = append(d.messages, fmt.Sprintf("Message %04d. %s", i, strings.Repeat("The quick brown fox jumps over the lazy dog. ", 1+i%4)))
	}
	return d
}

func init() {
	register(page{name: "MessageScroller", order: 40, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		// 1. 流式开启时每 120ms 向最后一条追加文字，每 12 次另起一条新消息。
		if scroller.streaming && c.Now().Sub(scroller.lastTick) >= 120*time.Millisecond {
			scroller.lastTick = c.Now()
			last := len(scroller.messages) - 1
			if strings.Count(scroller.messages[last], " ") > 24 {
				scroller.messages = append(scroller.messages, "Streaming reply:")
			} else {
				scroller.messages[last] += " token"
			}
		}
		if scroller.streaming {
			c.After(120 * time.Millisecond)
		}
		// 2. 工具栏加消息列表。
		ui.Column(c).Grow(1).Children(func() {
			kit.Header(c, fmt.Sprintf(t(c, "%d messages · top reached %d times"), len(scroller.messages), scroller.topLoads), func() {
				label := "Start streaming"
				if scroller.streaming {
					label = "Stop streaming"
				}
				if kit.Button(c, t(c, label)).Clicked() {
					scroller.streaming = !scroller.streaming
				}
			})
			kit.MessageScroller(c, &scrollerRows, len(scroller.messages), kit.MessageScrollerOpt{OnReachTop: func() { scroller.topLoads++ }}, func(i int) {
				ui.Column(c).Padding(theme.Space(2), theme.Space(3)).Radius(theme.Metrics.RadiusLarge).Background(theme.Colors.Panel).Children(func() {
					ui.Text(c, scroller.messages[i])
				})
			})
		})
	}})
}
