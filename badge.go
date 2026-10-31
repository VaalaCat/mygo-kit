package kit

import "github.com/egoist/mygo/ui"

// BadgeKind 选择徽标表达的语义颜色。
type BadgeKind int

// 徽标语义：中性、强调、成功、警告和危险。
const (
	BadgeNeutral BadgeKind = iota
	BadgeAccent
	BadgeSuccess
	BadgeWarning
	BadgeDanger
)

// Badge 创建一枚紧凑的状态徽标，文字颜色和浅底色由 kind 决定。
func Badge(c *ui.Context, text string, kind BadgeKind) ui.Element {
	// 1. 按语义取前景色，背景为同色的低透明度。
	t := GetTheme(c)
	fg := t.Colors.Muted
	switch kind {
	case BadgeAccent:
		fg = t.Colors.Accent
	case BadgeSuccess:
		fg = t.Colors.Success
	case BadgeWarning:
		fg = t.Colors.Warning
	case BadgeDanger:
		fg = t.Colors.Danger
	}
	bg := fg.Alpha(.16)
	// 2. 单行文字居中，不被父级拉伸。
	return ui.Row(c).Height(t.Metrics.CompactHeight).Shrink(0).Padding(0, t.Space(2)).Justify(ui.Center).AlignItems(ui.Center).Radius(t.Metrics.Radius).Background(bg).Children(func() {
		t.TextRole(ui.Text(c, text).SingleLine(), TextCaption).TextColor(fg)
	})
}

// StatusDot 创建 8 见方的圆点，用于连接状态和未读标记，颜色由调用方决定；调用方用 Label 给出可访问名称。
func StatusDot(c *ui.Context, color ui.Color) ui.Element {
	// 1. 固定尺寸，不被所在行拉伸或压缩。
	t := GetTheme(c)
	return ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(RadiusFull).Background(color)
}
