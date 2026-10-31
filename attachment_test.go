package kit

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

// TestAttachmentChip 验证附件 chip 落在 28 高的控件档，移除按钮是 20 见方的紧凑按钮。
func TestAttachmentChip(t *testing.T) {
	// 1. 有无移除回调各渲染一次，读取外框和移除按钮尺寸。
	tests := []struct {
		name      string
		removable bool
	}{
		{"removable", true},
		{"read-only", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			removed := 0
			var chip ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				var onRemove func()
				if test.removable {
					onRemove = func() { removed++ }
				}
				chip = AttachmentChip(c, "report.png", nil, onRemove).Bounds()
			}, 300, 100)
			tt.Frame()
			if chip.H != 28 {
				t.Fatalf("chip=%v", chip)
			}
			button, ok := tt.Find("Remove attachment")
			if ok != test.removable || ok && (button.W != 20 || button.H != 20) {
				t.Fatalf("remove=%v ok=%v", button, ok)
			}
			if test.removable {
				if err := tt.Click("Remove attachment"); err != nil || removed != 1 {
					t.Fatalf("removed=%d err=%v", removed, err)
				}
			}
		})
	}
}
