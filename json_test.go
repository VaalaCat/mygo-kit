package kit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/egoist/mygo/ui"
	"testing"
)

// TestJSONTree 验证所有合法顶级类型及复制。
func TestJSONTree(t *testing.T) {
	// 1. null、标量、数组和对象统一进入可访问树。
	for _, source := range []string{"null", "true", "1234567890123456789", `"中文"`, `[1,null]`, `{"a":1,"nested":{"b":false}}`} {
		t.Run(source, func(t *testing.T) {
			state, err := NewJSONState([]byte(source))
			if err != nil {
				t.Fatal(err)
			}
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneDark()); JSONTree(c, state, JSONOpt{}).Fill() }, 500, 350)
			if err := tt.Click("Copy JSON"); err != nil {
				t.Fatal(err)
			}
			if tt.Clipboard() == "" {
				t.Fatal("empty clipboard")
			}
		})
	}
	// 2. 错误输入不替换已有树。
	state, _ := NewJSONState([]byte("null"))
	for _, invalid := range []string{"{", "null false", ""} {
		if err := state.Set([]byte(invalid)); err == nil {
			t.Fatalf("accepted %q", invalid)
		}
	}
}

// TestJSONSetKeepsCollapsedRoot 验证刷新内容不会重新展开用户收起的根节点。
func TestJSONSetKeepsCollapsedRoot(t *testing.T) {
	// 1. 根节点首次展开，点击收起后再次 Set，成员行仍然隐藏。
	state, err := NewJSONState([]byte(`{"member":1}`))
	if err != nil {
		t.Fatal(err)
	}
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); JSONTree(c, state, JSONOpt{}).Fill() }, 500, 300)
	tt.Frame()
	if _, ok := tt.Find("member"); !ok {
		t.Fatalf("root not expanded at first: %v", tt.Texts())
	}
	if err := tt.Click("$"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find("member"); ok {
		t.Fatal("root did not collapse")
	}
	if err := state.Set([]byte(`{"member":2}`)); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if _, ok := tt.Find("member"); ok {
		t.Fatalf("Set expanded the collapsed root again: %v", tt.Texts())
	}
}

// TestJSONKeyboard 验证选中节点的展开、收起、跳转父子、切换和精确复制。
func TestJSONKeyboard(t *testing.T) {
	// 1. 每个用例从点击一行开始（点击会选中并切换分支），按键后复制选中节点并检查可见行。
	type press struct {
		mods ui.Modifiers
		key  ui.Key
	}
	tests := []struct {
		name        string
		click       string
		keys        []press
		want        string
		shown       []string
		hidden      []string
		expandFirst bool
	}{
		{"click expands branch", "a", nil, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", []string{"b", "c"}, nil, false},
		{"right goes to first child", "a", []press{{0, ui.KeyRight}}, "1234567890123456789", []string{"b"}, nil, false},
		{"down then left goes to parent", "a", []press{{0, ui.KeyDown}, {0, ui.KeyLeft}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", []string{"b"}, nil, false},
		{"left collapses open branch", "a", []press{{0, ui.KeyLeft}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", []string{"a", "d"}, []string{"b", "c"}, false},
		{"right expands collapsed branch", "a", []press{{0, ui.KeyLeft}, {0, ui.KeyRight}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", []string{"b", "c"}, nil, false},
		{"enter toggles", "a", []press{{0, ui.KeyEnter}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", nil, []string{"b"}, false},
		{"space collapses", "a", []press{{0, ui.KeySpace}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", nil, []string{"b"}, false},
		{"space twice expands again", "a", []press{{0, ui.KeySpace}, {0, ui.KeySpace}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", []string{"b"}, nil, false},
		{"leaf left goes to parent", "b", []press{{0, ui.KeyLeft}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", nil, nil, true},
		{"leaf right and enter do nothing", "d", []press{{0, ui.KeyRight}, {0, ui.KeyEnter}}, "2", []string{"d"}, nil, false},
		{"up moves selection", "d", []press{{0, ui.KeyUp}}, "{\n  \"b\": 1234567890123456789,\n  \"c\": [\n    true\n  ]\n}", nil, nil, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state, err := NewJSONState([]byte(`{"a":{"b":1234567890123456789,"c":[true]},"d":2}`))
			if err != nil {
				t.Fatal(err)
			}
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); JSONTree(c, state, JSONOpt{}).Fill() }, 500, 300)
			if test.expandFirst {
				// 先展开 a，再点击其中的叶子。
				tt.Click("a")
			}
			if err := tt.Click(test.click); err != nil {
				t.Fatal(err)
			}
			for _, k := range test.keys {
				tt.Key(k.mods, k.key)
			}
			tt.Key(ui.Cmd, ui.KeyC)
			if got := tt.Clipboard(); got != test.want {
				t.Fatalf("copy=%q want %q", got, test.want)
			}
			for _, name := range test.shown {
				if _, ok := tt.Find(name); !ok {
					t.Fatalf("missing %q in %v", name, tt.Texts())
				}
			}
			for _, name := range test.hidden {
				if _, ok := tt.Find(name); ok {
					t.Fatalf("unexpected %q in %v", name, tt.Texts())
				}
			}
		})
	}
}

// TestJSONVirtualization 验证大数组只构建可见行，展开、滚动和键盘操作都不随数组长度变慢。
func TestJSONVirtualization(t *testing.T) {
	// 1. 5000 项数组展开后，文字数量只与窗口高度有关。
	const items = 5000
	var source bytes.Buffer
	source.WriteString(`{"items":[`)
	for i := range items {
		if i > 0 {
			source.WriteByte(',')
		}
		fmt.Fprintf(&source, `{"id":%d}`, i)
	}
	source.WriteString(`]}`)
	state, err := NewJSONState(source.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); JSONTree(c, state, JSONOpt{}).Fill() }, 500, 300)
	if err := tt.Click("items"); err != nil {
		t.Fatal(err)
	}
	if len(state.rows) != items+2 {
		t.Fatalf("flattened rows=%d want %d", len(state.rows), items+2)
	}
	// 2. 300 高的窗口最多 15 行，加标题栏和预建行，文字远少于行数。
	const visibleTextLimit = 80
	if got := len(tt.Texts()); got > visibleTextLimit {
		t.Fatalf("built %d texts for %d rows, want virtualized (<= %d)", got, items+2, visibleTextLimit)
	}
	// 3. 展开中间的子节点并移动选中项，仍只构建可见行，且选中项跟随节点。
	const target = 20
	for range target + 1 {
		tt.Key(0, ui.KeyDown)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(ui.Cmd, ui.KeyC)
	if want := fmt.Sprintf("{\n  \"id\": %d\n}", target); tt.Clipboard() != want || len(state.rows) != items+3 {
		t.Fatalf("copy=%q rows=%d want %q with one more row", tt.Clipboard(), len(state.rows), want)
	}
	if got := len(tt.Texts()); got > visibleTextLimit {
		t.Fatalf("built %d texts after expanding, want <= %d", got, visibleTextLimit)
	}
}

// TestJSONText 验证显示和复制的 JSON 文本不转义 HTML 字符且不带结尾换行。
func TestJSONText(t *testing.T) {
	// 1. 逐个值序列化，& < > 保持原样，缩进为两个空格。
	for _, tt := range []struct {
		name  string
		value any
		want  string
	}{
		{"ampersand", "a&b", `"a&b"`},
		{"angle brackets", "<b>x</b>", `"<b>x</b>"`},
		{"escaped quote", `say "hi" & go`, `"say \"hi\" & go"`},
		{"object", map[string]any{"q": "a<b&c>d"}, "{\n  \"q\": \"a<b&c>d\"\n}"},
		{"number", json.Number("1234567890123456789"), "1234567890123456789"},
		{"null", nil, "null"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonText(tt.value); got != tt.want {
				t.Fatalf("jsonText=%q want %q", got, tt.want)
			}
		})
	}
}

// TestJSONRowPitch 验证 JSON 树行距是 20：与代码文字行高相同，没有额外内边距。
func TestJSONRowPitch(t *testing.T) {
	// 1. 相邻两行的文字纵向相距 Metrics.CompactHeight。
	state, _ := NewJSONState([]byte(`{"a":1,"b":{"c":2},"d":null}`))
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneLight()); JSONTree(c, state, JSONOpt{}).Fill() }, 500, 300)
	want := OneLight().Metrics.CompactHeight
	var ys []float32
	for _, name := range []string{"a", "b", "d"} {
		r, ok := tt.Find(name)
		if !ok {
			t.Fatalf("missing %q", name)
		}
		ys = append(ys, r.Y)
	}
	for i := 1; i < len(ys); i++ {
		if got := ys[i] - ys[i-1]; got != want {
			t.Fatalf("row pitch %d = %v, want %v", i, got, want)
		}
	}
}
