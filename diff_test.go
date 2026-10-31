package kit

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"testing"
)

// TestDiff 验证统一与并排差异使用原生比较结果。
func TestDiff(t *testing.T) {
	// 1. 文本差异保持真实的增删计数和行号布局。
	for _, layout := range []editor.DiffLayout{editor.DiffUnified, editor.DiffSplit} {
		state := editor.DiffState{Old: editor.NewBuffer("one\ntwo\n"), New: editor.NewBuffer("one\nthree\n"), Layout: layout, ExpandUnchanged: true}
		tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Diff(c, &state).Fill() }, 600, 300)
		added, deleted := state.Stats()
		if added != 1 || deleted != 1 {
			t.Fatalf("stats=%d,%d", added, deleted)
		}
		if tt.Image().Bounds().Dx() != 600 {
			t.Fatal("no frame")
		}
	}
}
