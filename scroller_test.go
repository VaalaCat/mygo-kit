package kit

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"testing"
)

// TestMessageScroller 验证跟随尾部、上滚后的回到底部按钮和到顶通知。
func TestMessageScroller(t *testing.T) {
	// 1. 每项用同一份 200 行列表，行高固定便于滚动。
	newScroller := func(state *ui.ListState, top *int) *ui.Tester {
		return ui.NewTester(func(c *ui.Context) {
			theme := OneLight()
			theme.Motion.Reduced = true
			ApplyTheme(c, theme)
			MessageScroller(c, state, 200, MessageScrollerOpt{OnReachTop: func() { *top++ }}, func(i int) {
				ui.Text(c, fmt.Sprintf("message %03d", i)).Height(28)
			})
		}, 400, 300)
	}
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"jump button appears after scrolling up and returns to end", func(t *testing.T) {
			state := &ui.ListState{Key: func(i int) any { return i }}
			top := 0
			tt := newScroller(state, &top)
			if !state.AtEnd() || !tt.HasText("message 199") {
				t.Fatal("not following end", tt.Texts())
			}
			if _, ok := tt.Find("Jump to latest"); ok {
				t.Fatal("jump button visible at end")
			}
			tt.Scroll(200, 150, 0, -400)
			if state.AtEnd() {
				t.Fatal("still at end")
			}
			if _, ok := tt.Find("Jump to latest"); !ok {
				t.Fatal("jump button missing")
			}
			if err := tt.Click("Jump to latest"); err != nil {
				t.Fatal(err)
			}
			if !state.AtEnd() || !tt.HasText("message 199") {
				t.Fatal("did not return to end")
			}
			if _, ok := tt.Find("Jump to latest"); ok {
				t.Fatal("jump button still visible")
			}
		}},
		{"on reach top fires once per arrival", func(t *testing.T) {
			state := &ui.ListState{Key: func(i int) any { return i }}
			top := 0
			tt := newScroller(state, &top)
			if top != 0 {
				t.Fatal("fired at end", top)
			}
			for range 40 {
				tt.Scroll(200, 150, 0, -2000)
			}
			if top != 1 {
				t.Fatal("reach top count", top)
			}
			tt.Scroll(200, 150, 0, 600)
			tt.Scroll(200, 150, 0, -3000)
			if top != 2 {
				t.Fatal("second arrival", top)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
