package kit

import (
	"github.com/egoist/mygo/ui"
	"go/ast"
	"go/parser"
	"go/token"
	"image/color"
	"path/filepath"
	"strings"
	"testing"
)

// TestTheme 验证窗口隔离、主题切换与直角控件的实际像素。
func TestTheme(t *testing.T) {
	// 1. 两个窗口的主题互不覆盖。
	for _, dark := range []bool{false, true} {
		t.Run(map[bool]string{false: "light", true: "dark"}[dark], func(t *testing.T) {
			theme := OneLight()
			if dark {
				theme = OneDark()
			}
			theme.Metrics.Radius = 0
			tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, theme); PrimaryButton(c, "Primary").Width(100) }, 240, 100)
			r, ok := tt.Find("Primary")
			if !ok {
				t.Fatal("missing")
			}
			got := tt.Image().RGBAAt(int(r.X)+1, int(r.Y)+1)
			want := theme.Colors.Accent
			if got != (color.RGBA{want.R, want.G, want.B, want.A}) {
				t.Fatalf("corner=%v want=%v", got, want)
			}
		})
	}
	// 2. 切换主题不依赖全局可变单例。
	theme := OneLight()
	var shown Theme
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, theme); shown = GetTheme(c); ui.Text(c, "Theme") }, 200, 100)
	ui.NewTester(func(c *ui.Context) { ApplyTheme(c, OneDark()); ui.Text(c, "Other") }, 200, 100)
	tt.Frame()
	if shown.Dark {
		t.Fatal("theme leaked")
	}
	theme = OneDark()
	tt.Frame()
	if !shown.Dark {
		t.Fatal("theme did not update")
	}
}

// TestRadiusAndPrimaryStates 验证显式圆角和 primary 的悬停对比。
func TestRadiusAndPrimaryStates(t *testing.T) {
	// 1. 半径零与非零都作用于实际绘制。
	for _, radius := range []float32{0, 8} {
		theme := OneLight()
		theme.Metrics.Radius = radius
		tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, theme); PrimaryButton(c, "Primary").Width(120) }, 200, 100)
		r, _ := tt.Find("Primary")
		corner := tt.Image().RGBAAt(int(r.X), int(r.Y))
		accent := theme.Colors.Accent
		if radius == 0 && corner.R != accent.R {
			t.Fatal("zero radius ignored")
		}
		if radius > 0 && corner.R == accent.R && corner.G == accent.G {
			t.Fatal("radius token ignored")
		}
	}
	// 2. 悬停仍是强调色而不是低对比灰底。
	theme := OneLight()
	tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, theme); PrimaryButton(c, "Primary").Width(120) }, 200, 100)
	r, _ := tt.Find("Primary")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	got := tt.Image().RGBAAt(int(r.X)+10, int(r.Y)+2)
	want := theme.Colors.Accent.Mix(theme.Colors.Text, .12)
	if got != (color.RGBA{want.R, want.G, want.B, want.A}) {
		t.Fatalf("hover=%v want=%v", got, want)
	}
}

// TestScrollbarTheme 验证滚动条细且颜色来自主题次要色。
func TestScrollbarTheme(t *testing.T) {
	// 1. 明暗主题都把滚动条宽度和颜色交给框架主题。
	for _, dark := range []bool{false, true} {
		theme := OneLight()
		if dark {
			theme = OneDark()
		}
		var got ui.Theme
		tt := ui.NewTester(func(c *ui.Context) { ApplyTheme(c, theme); got = *c.Theme(); ui.Text(c, "x") }, 100, 60)
		tt.Frame()
		if got.ScrollbarWidth != 4 || got.Scrollbar != theme.Colors.Muted.Alpha(scrollbarAlpha) {
			t.Fatalf("dark=%v width=%v color=%v", dark, got.ScrollbarWidth, got.Scrollbar)
		}
	}
}

// literalSpacingMethods 是接收尺寸和间距的元素方法，守卫测试检查它们的参数。
var literalSpacingMethods = map[string]bool{
	"Padding": true, "PaddingStart": true, "PaddingEnd": true, "Margin": true, "Gap": true,
	"Height": true, "Width": true, "MinHeight": true, "MinWidth": true, "MaxHeight": true, "MaxWidth": true, "Size": true,
}

// literalSpacingAllowed 列出允许保留数字字面量的位置，键是「文件:方法:字面量」，值是原因。
// 目前组件文件里没有需要保留的字面量：一次性的宽高都写成了具名常量。
var literalSpacingAllowed = map[string]string{}

// TestKitLiteralSpacing 验证 kit 的组件文件没有写死的间距和尺寸数字。
func TestKitLiteralSpacing(t *testing.T) {
	// 1. 逐个解析非测试源文件，收集间距和尺寸方法里的数字字面量。
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	used := map[string]bool{}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(parsed, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			// 2. 任何位置的 Space(n) 档位只能取 spaceStepsAllowed，大尺寸写成具名常量。
			if ok && selector.Sel.Name == "Space" && len(call.Args) == 1 {
				if literal, isLiteral := call.Args[0].(*ast.BasicLit); isLiteral && !spaceStepsAllowed[literal.Value] {
					t.Errorf("%s: Space(%s) 不在网格档位内，改用具名常量", file, literal.Value)
				}
				return true
			}
			if !ok || !literalSpacingMethods[selector.Sel.Name] {
				return true
			}
			for _, arg := range call.Args {
				for _, literal := range numericLiterals(arg) {
					// 3. 0 和 1 是无间距和一像素线，其余必须来自 Space、Metrics 或具名常量。
					if literal == "0" || literal == "1" {
						continue
					}
					key := file + ":" + selector.Sel.Name + ":" + literal
					if _, ok := literalSpacingAllowed[key]; ok {
						used[key] = true
						continue
					}
					t.Errorf("%s: %s(%s) 写死了数字，改用 Space(n)、Metrics 或具名常量", file, selector.Sel.Name, literal)
				}
			}
			return true
		})
	}
	// 4. 白名单里不再出现的条目说明已经没有用处，直接删除。
	for key := range literalSpacingAllowed {
		if !used[key] {
			t.Errorf("白名单条目 %s 已无对应字面量", key)
		}
	}
}

// numericLiterals 返回表达式里的数字字面量；Space(n) 的参数由 spaceStepsAllowed 另行校验，不算在内。
func numericLiterals(e ast.Expr) []string {
	var literals []string
	ast.Inspect(e, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.IndexExpr:
			// 下标不是尺寸。
			return false
		case *ast.CallExpr:
			if selector, ok := v.Fun.(*ast.SelectorExpr); ok && selector.Sel.Name == "Space" {
				return false
			}
		case *ast.BasicLit:
			if v.Kind == token.INT || v.Kind == token.FLOAT {
				literals = append(literals, v.Value)
			}
		}
		return true
	})
	return literals
}

// spaceStepsAllowed 是 Space(n) 允许的网格档位：1、2、3、4、6、8 即 4 到 32；5 只表示 20 的行高。
var spaceStepsAllowed = map[string]bool{"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "8": true}

// TestControlBackground 验证控件底色的状态优先级：禁用 > 按下 > 选中 > 悬停。
func TestControlBackground(t *testing.T) {
	// 1. 每一行给出一个状态组合和预期的底色来源。
	theme := OneLight()
	c := theme.Colors
	tests := []struct {
		name  string
		state controlState
		want  ui.Color
	}{
		{"idle", controlState{}, ui.Transparent},
		{"hovered", controlState{Hovered: true}, c.Hover},
		{"selected-over-hovered", controlState{Selected: true, Hovered: true}, c.Selected},
		{"pressed-over-selected", controlState{Pressed: true, Selected: true, Hovered: true}, c.Active},
		{"disabled-over-all", controlState{Disabled: true, Pressed: true, Selected: true, Hovered: true}, ui.Transparent},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := controlBackground(&theme, test.state); got != test.want {
				t.Fatalf("got %v want %v", got, test.want)
			}
		})
	}
}

// pixelNear 判断离屏截图在 (x, y) 处的颜色与 want 每个通道相差是否不超过 tolerance。
func pixelNear(tt *ui.Tester, x, y float32, want ui.Color, tolerance int) bool {
	got := tt.Image().RGBAAt(int(x), int(y))
	near := func(a, b uint8) bool { return abs(int(a)-int(b)) <= tolerance }
	return near(got.R, want.R) && near(got.G, want.G) && near(got.B, want.B)
}

// abs 返回整数的绝对值。
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
