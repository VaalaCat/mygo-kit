package kit

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"math"
	"strings"
	"testing"
)

// TestParseANSI 验证颜色保留和终端控制指令剥离。
func TestParseANSI(t *testing.T) {
	// 1. 外部输出只产生文本，不执行 OSC 剪贴板、超链接或光标操作。
	tests := []struct{ name, input, want string }{{"plain", "hello", "hello"}, {"sgr", "\x1b[31mred\x1b[0m plain", "red plain"}, {"osc", "before\x1b]52;c;c2VjcmV0\aafter", "beforeafter"}, {"cursor", "a\x1b[2Jb", "ab"}, {"incomplete", "safe\x1b]52;ignored", "safe"}, {"truecolor", "\x1b[38;2;12;34;56mcolor", "color"}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spans := ansiText(test.input, OneDark())
			var got strings.Builder
			for _, span := range spans {
				got.WriteString(span.Text)
			}
			if got.String() != test.want {
				t.Fatalf("got %q", got.String())
			}
			if test.name == "sgr" && spans[0].Color != OneDark().Colors.ANSI[1] {
				t.Fatal("color lost")
			}
			if test.name == "truecolor" && (spans[0].Color.R != 12 || spans[0].Color.G != 34 || spans[0].Color.B != 56) {
				t.Fatal("truecolor lost")
			}
		})
	}
}

// TestANSIViewport 验证一万行只构建可见部分，上滚后追加不抢阅读位置。
func TestANSIViewport(t *testing.T) {
	// 1. 真实受限 viewport 内保留少量行。
	source := ""
	for i := 0; i < 10000; i++ {
		source += fmt.Sprintf("line-%05d\n", i)
	}
	state := ui.ListState{FollowEnd: true}
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneDark()); ANSIOutput(c, source, &state, ANSIOutputOpt{}).Fill() }, 400, 300)
	count := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "line-") {
			count++
		}
	}
	if count > 40 {
		t.Fatalf("built %d rows", count)
	}
	// 2. 离开底部时保持锚点，回底后恢复跟随。
	tt.Scroll(100, 180, 0, -220)
	anchor := ""
	var before ui.Rect
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "line-") {
			r, _ := tt.Find(text)
			if r.W > 0 && r.H > 0 && r.Y >= 40 && r.Y+r.H < 300 {
				anchor, before = text, r
				break
			}
		}
	}
	if anchor == "" {
		t.Fatalf("anchor missing: %v", tt.Texts())
	}
	source += "line-10000\n"
	tt.Frame()
	after, ok := tt.Find(anchor)
	if !ok || math.Abs(float64(after.Y-before.Y)) > .1 {
		t.Fatalf("append moved anchor %s: %v -> %v", anchor, before, after)
	}
	tt.Scroll(100, 180, 0, 1000000)
	source += "line-10001\n"
	tt.Frame()
	if !tt.HasText("line-10001") {
		t.Fatal("end follow not resumed")
	}
}

// TestPlainOutput 验证原文模式不解释或删除 ESC，复制逐字相同。
func TestPlainOutput(t *testing.T) {
	// 1. 与终端输出模式共用布局，但不经过 ANSI 解码。
	source := "literal \x1b[31mred\x1b[0m\n\x1b]52;c;abc\a"
	state := ui.ListState{}
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		ANSIOutput(c, source, &state, ANSIOutputOpt{Plain: true, CopyLabel: "Copy exact"}).Fill()
	}, 500, 300)
	tt.Click("Copy exact")
	if tt.Clipboard() != source {
		t.Fatalf("copy=%q want=%q", tt.Clipboard(), source)
	}
}

// TestANSIWheelOnBlankArea 验证短输出右侧空白也属于完整滚动区域。
func TestANSIWheelOnBlankArea(t *testing.T) {
	// 1. 先跟随一万行，再从文字右侧空白滚动。
	for _, x := range []float32{100, 350} {
		t.Run(fmt.Sprint(x), func(t *testing.T) {
			var source strings.Builder
			for i := 0; i < 10000; i++ {
				fmt.Fprintf(&source, "row%05d\n", i)
			}
			value := source.String()
			state := ui.ListState{FollowEnd: true}
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); ANSIOutput(c, value, &state, ANSIOutputOpt{}).Fill() }, 400, 300)
			before, _ := state.Visible()
			tt.Scroll(x, 180, 0, -220)
			after, _ := state.Visible()
			if after >= before {
				t.Fatalf("blank area did not scroll: %d -> %d", before, after)
			}
			value += "tail\n"
			tt.Frame()
			appended, _ := state.Visible()
			if appended != after {
				t.Fatalf("append stole reading position: %d -> %d", after, appended)
			}
		})
	}
}

// TestANSISharedSelection 验证可见多行可以拖选复制，而非每行独立选区。
func TestANSISharedSelection(t *testing.T) {
	// 1. 父列表持有共享选区，逐行文本不覆盖它。
	state := ui.ListState{}
	tt := ui.NewTester(func(c *ui.Context) {
		ApplyTheme(c, OneLight())
		ANSIOutput(c, "alpha\nbeta\ngamma", &state, ANSIOutputOpt{}).Fill()
	}, 400, 240)
	a, _ := tt.Find("alpha")
	b, _ := tt.Find("gamma")
	tt.Press(a.X, a.Y+a.H/2)
	tt.Move(b.X+50, b.Y+b.H/2)
	tt.Release(b.X+50, b.Y+b.H/2)
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "alpha\nbeta\ngamma" {
		t.Fatalf("selection=%q", tt.Clipboard())
	}
}
