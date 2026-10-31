package kit

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

// TestPatchLines 验证未知行号不补造，复制保留原始片段。
func TestPatchLines(t *testing.T) {
	// 1. 省略行号与明确行号分别渲染，正文仍可读取。
	original := "*** Begin Patch\n-old\n+new\n*** End Patch\n"
	number := 42
	lines := []PatchLine{{Kind: PatchDeleted, Text: "old"}, {Kind: PatchAdded, Text: "new", NewLine: &number}}
	state := ui.ListState{}
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		PatchLines(c, &state, lines, PatchOpt{Source: original, CopyLabel: "Copy patch"}).Fill()
	}, 500, 300)
	tt.Click("Copy patch")
	if tt.Clipboard() != original {
		t.Fatalf("copy changed: %q", tt.Clipboard())
	}
	if !tt.HasText("42") || !tt.HasText("new") {
		t.Fatal(tt.Texts())
	}
}

// TestPatchSharedSelection 验证可见多行可以拖选，行号和装饰符不混入正文。
func TestPatchSharedSelection(t *testing.T) {
	// 1. 共享选区只复制实际正文，不伪造原始补丁符号。
	state := ui.ListState{}
	lines := []PatchLine{{Kind: PatchDeleted, Text: "alpha"}, {Kind: PatchAdded, Text: "beta"}, {Kind: PatchContext, Text: "gamma"}}
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); PatchLines(c, &state, lines, PatchOpt{}).Fill() }, 400, 240)
	a, _ := tt.Find("alpha")
	b, _ := tt.Find("gamma")
	tt.Press(a.X, a.Y+a.H/2)
	tt.Move(b.X+b.W, b.Y+b.H/2)
	tt.Release(b.X+b.W, b.Y+b.H/2)
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "alpha\nbeta\ngamma" {
		t.Fatalf("selection=%q", tt.Clipboard())
	}
}
