package kit

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// Diff 显示原生 diff、增删行号、选择和复制，调用方持有原框架 DiffState。
func Diff(c *ui.Context, state *editor.DiffState) ui.Element {
	// 1. 保留框架的差异算法与编辑状态，只设置工作台外观。
	t := GetTheme(c)
	colors := codeTheme(c)
	colors.Background = t.Colors.Editor
	colors.Text = t.Colors.Text
	colors.LineNumber = t.Colors.Placeholder
	colors.Added = t.Colors.Success
	colors.Deleted = t.Colors.Danger
	colors.Modified = t.Colors.Warning
	state.Theme = colors
	state.Font = ui.Font{Family: t.Typography.MonoFont, Size: t.Typography.CodeSize}
	state.LineHeight = 1.5
	return editor.DiffView(c, state).Radius(0).Background(t.Colors.Editor)
}

// codeTheme 将语义代码颜色映射到 MyGo 编辑器。
func codeTheme(c *ui.Context) *editor.Theme {
	// 1. 只替换已有语法颜色，不维护另一份文本或高亮状态。
	t := GetTheme(c)
	colors := editor.DefaultTheme(c.Theme())
	colors.Background = t.Colors.Editor
	colors.Text = t.Colors.Text
	for _, entry := range []struct {
		kind  editor.TokenKind
		color ui.Color
	}{{editor.TokenKeyword, t.Syntax.Keyword}, {editor.TokenString, t.Syntax.String}, {editor.TokenFunction, t.Syntax.Function}, {editor.TokenType, t.Syntax.Type}, {editor.TokenNumber, t.Syntax.Number}, {editor.TokenComment, t.Syntax.Comment}, {editor.TokenProperty, t.Syntax.Property}, {editor.TokenOperator, t.Syntax.Operator}} {
		colors.Tokens[entry.kind].Color = entry.color
	}
	return colors
}
