package kit

import (
	"github.com/egoist/mygo/ui"
	"strings"
	"testing"
)

// TestMarkdown 验证流式正文、代码复制、任务列表与选区。
func TestMarkdown(t *testing.T) {
	// 1. 相同元素跨流式与完成状态保留可读内容。
	for _, streaming := range []bool{false, true} {
		t.Run(map[bool]string{false: "complete", true: "streaming"}[streaming], func(t *testing.T) {
			source := "# Title\n\nSome **bold** text.\n\n- [x] **Done**\n- [ ] Pending\n\n```go\nfunc main() {}\n```\n"
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Markdown(c, source, MarkdownOpt{Streaming: streaming}) }, 600, 700)
			for _, want := range []string{"Title", "Some bold text.", "Done", "Pending", "func main() {}"} {
				if !tt.HasText(want) {
					t.Fatalf("missing %q in %q", want, tt.Texts())
				}
			}
			if err := tt.Click("Copy code"); err != nil {
				if err = tt.Click("Copy"); err != nil {
					t.Fatal(err)
				}
			}
			if tt.Clipboard() != "func main() {}\n" {
				t.Fatal(tt.Clipboard())
			}
			r, ok := tt.Find("Some bold text.")
			if !ok {
				t.Fatal("paragraph missing")
			}
			tt.ClickAt(r.X+2, r.Y+r.H/2)
			tt.Command("selectAll")
			tt.Command("copy")
			if !strings.Contains(tt.Clipboard(), "Done") || !strings.Contains(tt.Clipboard(), "func main() {}") {
				t.Fatal(tt.Clipboard())
			}
		})
	}
}

// TestMarkdownContext 验证图片父级链接和任务项文档级引用不丢失。
func TestMarkdownContext(t *testing.T) {
	// 1. 原 Document 的链接目标可从真实控件激活。
	tests := []struct{ name, source, label, url string }{{"linked-image", "[![alt](demo://image)](https://example.com)", "alt", "https://example.com"}, {"task-reference", "- [x] [Go][docs]\n\n[docs]: https://go.dev", "Go", "https://go.dev"}, {"task-footnote", "- [x] note[^1]\n\n[^1]: Footnote **body**", "", ""}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Markdown(c, test.source, MarkdownOpt{}) }, 600, 400)
			if test.label != "" {
				if test.name == "task-reference" {
					r, ok := tt.Find(test.label)
					if !ok {
						t.Fatal("missing link")
					}
					tt.ClickAt(r.X+4, r.Y+r.H/2)
				} else if err := tt.Click(test.label); err != nil {
					t.Fatal(err)
				}
				urls := tt.OpenedURLs()
				if len(urls) != 1 || urls[0] != test.url {
					r, _ := tt.Find(test.label)
					t.Fatalf("URLs=%v texts=%v rect=%v", urls, tt.Texts(), r)
				}
			} else if !tt.HasText("Footnote body") {
				t.Fatal(tt.Texts())
			}
			for _, text := range tt.Texts() {
				if strings.Contains(text, "](https:") || strings.Contains(text, "[docs]") {
					t.Fatalf("raw inline leaked: %q", text)
				}
			}
		})
	}
	// 2. 图片前后的强调仍以相同字重排版，不把 Markdown 标记显示为正文。
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		Markdown(c, "**before ![alt](demo://image) after**", MarkdownOpt{})
	}, 600, 400)
	for _, text := range tt.Texts() {
		if strings.Contains(text, "**") {
			t.Fatalf("raw strong marker: %q", text)
		}
	}
	if !tt.HasText("before") || !tt.HasText("after") {
		t.Fatal(tt.Texts())
	}
}

// TestNestedCodeFences 验证任务项代码块中的 fence 作为原文保留。
func TestNestedCodeFences(t *testing.T) {
	// 1. 重新承载代码时 fence 长度必须超过内容中的任何反引号串。
	for _, code := range []string{"before\n```\nafter\n", "a\n``````\nb\n", "literal `tick`\n"} {
		t.Run(code, func(t *testing.T) {
			fence := strings.Repeat("`", 8)
			source := "- [x] Example\n\n  " + fence + "text\n  " + strings.ReplaceAll(strings.TrimSuffix(code, "\n"), "\n", "\n  ") + "\n  " + fence + "\n"
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Markdown(c, source, MarkdownOpt{}) }, 600, 500)
			if err := tt.Click("Copy code"); err != nil {
				t.Fatal(err)
			}
			if tt.Clipboard() != code {
				t.Fatalf("copied %q want %q", tt.Clipboard(), code)
			}
		})
	}
}

// TestMarkdownSharedSelection 验证正文、任务列表和代码共享一个可见选区。
func TestMarkdownSharedSelection(t *testing.T) {
	// 1. 工具栏不是正文，跨块拖选只复制实际内容。
	source := "alpha\n\n- [x] bravo\n\n```text\ncharlie\n```\n"
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); Markdown(c, source, MarkdownOpt{}).Fill() }, 500, 400)
	a, _ := tt.Find("alpha")
	b, _ := tt.Find("charlie")
	tt.Press(a.X+1, a.Y+a.H/2)
	tt.Move(b.X+100, b.Y+b.H/2)
	tt.Release(b.X+100, b.Y+b.H/2)
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "alpha\nbravo\ncharlie" {
		t.Fatalf("selection=%q", tt.Clipboard())
	}
}

// TestMarkdownProseRhythm 验证正文使用 Prose 行高 24 和段间距 8。
func TestMarkdownProseRhythm(t *testing.T) {
	// 1. 两个相邻段落的高度和间隔都落在网格上。
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		ui.Column(c).Width(300).Children(func() { Markdown(c, "First\n\nSecond", MarkdownOpt{}) })
	}, 400, 300)
	first, _ := tt.Find("First")
	second, _ := tt.Find("Second")
	if first.H != 24 || second.Y-first.Y-first.H != 8 {
		t.Fatalf("first=%v second=%v", first, second)
	}
}

// TestMarkdownOpenCodeBlockCopyDisabled 验证围栏还没闭合的代码块不能复制，闭合后可以。
func TestMarkdownOpenCodeBlockCopyDisabled(t *testing.T) {
	// 1. 围栏未闭合（MyGo 报告该块仍在书写）时复制被禁用，和换成自绘标题行之前一致；围栏闭合后可复制并显示「已复制」。
	tests := []struct {
		name      string
		source    string
		streaming bool
		want      string
	}{
		{"open-fence-streaming", "```go\nfunc main() {}", true, ""},
		{"closed-fence-streaming", "```go\nfunc main() {}\n```\n", true, "func main() {}\n"},
		{"open-fence-not-streaming", "```go\nfunc main() {}", false, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				ApplyTheme(c, OneLight())
				Markdown(c, test.source, MarkdownOpt{Streaming: test.streaming, CodeCopyLabel: "Copy code"})
			}, 600, 300)
			tt.Click("Copy code")
			if tt.Clipboard() != test.want {
				t.Fatalf("clipboard=%q want %q", tt.Clipboard(), test.want)
			}
			if copied := tt.HasText("Copied"); copied != (test.want != "") {
				t.Fatalf("copied feedback=%v texts=%v", copied, tt.Texts())
			}
		})
	}
}
