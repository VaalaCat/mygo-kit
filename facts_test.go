package kit

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestFacts 验证事实组的标签和值显示，且所有值共用一条左对齐线。
func TestFacts(t *testing.T) {
	// 1. 标签长短不同，值的左缘仍然一致。
	tests := []struct {
		name  string
		facts []Fact
	}{
		{"short and long labels", []Fact{{"cwd", "/work/alpha"}, {"command_id", "cmd-1"}}},
		{"single", []Fact{{"mode", "steer"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Facts(c, test.facts...) }, 400, 120)
			left := float32(-1)
			for _, fact := range test.facts {
				if !tt.HasText(fact.Label) {
					t.Fatalf("missing label %q", fact.Label)
				}
				rect, ok := tt.Find(fact.Value)
				if !ok {
					t.Fatalf("missing value %q", fact.Value)
				}
				if left >= 0 && rect.X != left {
					t.Fatalf("value %q x=%v want %v", fact.Value, rect.X, left)
				}
				left = rect.X
			}
		})
	}
}
