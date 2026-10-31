package kit

import (
	"github.com/egoist/mygo/ui"
	"slices"
	"strings"
	"testing"
	"time"
)

// TestMotion 验证减少动画、纯追加和字素变化。
func TestMotion(t *testing.T) {
	// 1. 不拆开中文组合文字和 emoji。
	for _, test := range []struct {
		input string
		want  []string
	}{{"a👨‍👩‍👧x", []string{"a", "👨‍👩‍👧", "x"}}, {"e\u0301x", []string{"e\u0301", "x"}}, {"中文", []string{"中", "文"}}} {
		if got := textGraphemes(test.input); !slices.Equal(got, test.want) {
			t.Fatalf("graphemes=%q want=%q", got, test.want)
		}
	}
	// 2. 纯追加立即显示，替换渐显在规定时间结束。
	text := "hello"
	theme := OneLight()
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, theme); MorphText(c, text) }, 400, 100)
	text = "hello world"
	tt.Frame()
	if !tt.HasText("hello world") {
		t.Fatal(tt.Texts())
	}
	text = "changed"
	tt.Frame()
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Fatal("animation never stopped")
	}
	theme.Motion.Reduced = true
	text = "静态文字"
	tt.Frame()
	if !tt.HasText(text) || tt.FramePending() {
		t.Fatal("reduced motion still animates")
	}
}

// TestMorphStableGlyphs 验证共同后缀在变化首帧仍保持可见，不随变化字符一起淡入。
func TestMorphStableGlyphs(t *testing.T) {
	// 1. 读取 a.ts 改成 b.ts 时，扩展名保持原有像素。
	value := "读取a.ts"
	prefixWidth := float32(0)
	changedStart, changedWidth := float32(0), float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		theme := OneLight()
		theme.Typography.Font = "monospace"
		ApplyTheme(c, theme)
		prefixWidth, _ = c.MeasureText(0, ui.Span{Text: "读取a.", Font: "monospace", Size: 14})
		changedStart, _ = c.MeasureText(0, ui.Span{Text: "读取", Font: "monospace", Size: 14})
		changedWidth, _ = c.MeasureText(0, ui.Span{Text: "b", Font: "monospace", Size: 14})
		MorphText(c, value)
	}, 400, 100)
	before := tt.Image()
	r, ok := tt.Find(value)
	if !ok {
		t.Fatal(tt.Texts())
	}
	value = "读取b.ts"
	tt.Frame()
	after := tt.Image()
	ink := 0
	for y := 0; y < int(r.H); y++ {
		for x := int(prefixWidth); x < int(r.W); x++ {
			a := before.RGBAAt(int(r.X)+x, int(r.Y)+y)
			b := after.RGBAAt(int(r.X)+x, int(r.Y)+y)
			if a != b {
				t.Fatal("stable suffix animated")
			}
			if a.R < 100 {
				ink++
			}
		}
	}
	if ink == 0 {
		t.Fatal("suffix invisible")
	}
	// 2. 只有新字素从透明到可见，稳定文本无需等动画结束。
	initialInk := 0
	for y := 0; y < int(r.H); y++ {
		for x := int(changedStart); x < int(changedStart+changedWidth); x++ {
			if after.RGBAAt(x, y).R < 100 {
				initialInk++
			}
		}
	}
	tt.Advance(time.Second)
	finished := tt.Image()
	finalInk := 0
	for y := 0; y < int(r.H); y++ {
		for x := int(changedStart); x < int(changedStart+changedWidth); x++ {
			if finished.RGBAAt(x, y).R < 100 {
				finalInk++
			}
		}
	}
	if initialInk != 0 {
		t.Fatalf("changed glyph was not fully transparent at progress zero: %d pixels", initialInk)
	}
	if finalInk <= initialInk {
		t.Fatalf("changed glyph did not fade: %d -> %d", initialInk, finalInk)
	}
	// 3. 超长变化不做二次复杂度匹配或动画，但保留全文。
	value = strings.Repeat("a", 300)
	tt.Frame()
	tt.Advance(time.Second)
	value = strings.Repeat("b", 300)
	tt.Frame()
	if !tt.HasText(value) || tt.FramePending() {
		t.Fatal("long text did not fall back to static display")
	}
}
