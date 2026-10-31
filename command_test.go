package kit

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestCommandPalette 验证搜索过滤、键盘选择、禁用项、点击和 Escape。
func TestCommandPalette(t *testing.T) {
	commands := []Command{
		{ID: "conversation.new", Title: "New conversation", Group: "Conversation", Keys: "⌘N"},
		{ID: "conversation.archive", Title: "Archive conversation", Group: "Conversation"},
		{ID: "conversation.delete", Title: "Delete conversation", Group: "Conversation", Disabled: true},
		{ID: "view.sidebar", Title: "Toggle sidebar", Group: "View", Keys: "⌘B"},
	}
	// 1. 每个用例都从已打开的面板开始，记录每一帧返回的命令。
	tests := []struct {
		name     string
		run      func(*testing.T, *ui.Tester)
		wantID   []string
		wantOpen bool
	}{
		{"typing-filters", func(t *testing.T, tt *ui.Tester) {
			tt.Type("arch")
			if tt.HasText("New conversation") || !tt.HasText("Archive conversation") {
				t.Fatal(tt.Texts())
			}
		}, nil, true},
		{"typing-ignores-group-name", func(t *testing.T, tt *ui.Tester) {
			tt.Type("view")
			if tt.HasText("Toggle sidebar") || !tt.HasText("No results") {
				t.Fatal(tt.Texts())
			}
		}, nil, true},
		{"no-match-shows-empty-state", func(t *testing.T, tt *ui.Tester) {
			tt.Type("zzz")
			if !tt.HasText("No results") {
				t.Fatal(tt.Texts())
			}
		}, nil, true},
		{"down-enter-returns-id", func(t *testing.T, tt *ui.Tester) {
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
		}, []string{"conversation.archive"}, false},
		{"enter-takes-first-command", func(t *testing.T, tt *ui.Tester) {
			tt.Key(0, ui.KeyEnter)
		}, []string{"conversation.new"}, false},
		{"down-skips-disabled", func(t *testing.T, tt *ui.Tester) {
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
		}, []string{"view.sidebar"}, false},
		{"filter-resets-highlight", func(t *testing.T, tt *ui.Tester) {
			tt.Key(0, ui.KeyDown)
			tt.Type("side")
			tt.Key(0, ui.KeyEnter)
		}, []string{"view.sidebar"}, false},
		{"click-returns-id", func(t *testing.T, tt *ui.Tester) {
			if err := tt.Click("Toggle sidebar"); err != nil {
				t.Fatal(err)
			}
		}, []string{"view.sidebar"}, false},
		{"disabled-click-ignored", func(t *testing.T, tt *ui.Tester) {
			tt.Click("Delete conversation")
		}, nil, true},
		{"escape-closes", func(t *testing.T, tt *ui.Tester) {
			tt.Key(0, ui.KeyEscape)
		}, nil, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			open := true
			var got []string
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				if id := CommandPalette(c, &open, commands); id != "" {
					got = append(got, id)
				}
			}, 800, 600)
			test.run(t, tt)
			tt.Frame()
			if !slices.Equal(got, test.wantID) || open != test.wantOpen {
				t.Fatalf("got=%v open=%v, want=%v open=%v", got, open, test.wantID, test.wantOpen)
			}
		})
	}
}
