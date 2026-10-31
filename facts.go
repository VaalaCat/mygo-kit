package kit

import (
	"github.com/egoist/mygo/ui"
)

// Fact 是一行事实：标签与值。
type Fact struct{ Label, Value string }

// factLabelWidth 是 Facts 标签列的宽度。
const factLabelWidth float32 = 112

// Facts 创建两列事实组（标签列宽 112，同一组事实的值共用一条左对齐线）：左列标签用 12/16 次要色，右列值用等宽 12/20 并可选择，两列按首行基线对齐。
func Facts(c *ui.Context, facts ...Fact) ui.Element {
	// 1. 行与行之间不留间距，列宽固定让所有值共用左对齐线。
	t := GetTheme(c)
	return ui.Column(c).Gap(0).Children(func() {
		for index, fact := range facts {
			ui.Row(c.Key(index)).Gap(t.Space(3)).AlignItems(ui.Baseline).Children(func() {
				t.TextRole(ui.Text(c, fact.Label).SingleLine().Width(factLabelWidth).Shrink(0), TextCaption).TextColor(t.Colors.Muted)
				t.TextRole(ui.Text(c, fact.Value).Selectable().Grow(1).Shrink(1).MinWidth(0), TextCode)
			})
		}
	})
}
