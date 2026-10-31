---
title: 主题
description: Theme 值类型、One Light / One Dark、颜色职责、尺寸、文字层级、图标尺寸和间距档位。
---

# 主题

主题是一个普通的 Go 值 `kit.Theme`，由应用持有，每帧交给 `ApplyTheme`。本库不使用全局可变主题，每个窗口各自保存自己的主题，所以多窗口、明暗各一个也没问题。

## 用法

```go
func view(c *ui.Context) {
	theme := kit.OneLight()
	theme.Metrics.Radius = 0 // 可选：控件用直角
	kit.ApplyTheme(c, theme)

	t := kit.GetTheme(c) // 之后任何地方都可以取回当前窗口的主题
	ui.Text(c, "你好").TextColor(t.Colors.Muted)
}
```

- `kit.OneLight()` 和 `kit.OneDark()` 各返回一份完整的主题，基于 Zed 的 One 配色。
- `kit.ApplyTheme(c, t)` 要**每帧**调用，并且在画任何组件之前。它做两件事：把主题存在本帧的窗口根节点上，并把文字、焦点、选区和滚动条的颜色设给 MyGo。
- `kit.GetTheme(c)` 返回当前窗口的主题。还没调用过 `ApplyTheme` 时，返回 One Light。

`Theme` 是值类型，修改它的副本不会影响别人。要做自己的主题，从 `OneLight()` 或 `OneDark()` 拿一份，改想改的字段。

```go
theme := kit.OneDark()
theme.Colors.Accent = ui.Hex("#e5836b")
```

## 颜色

`Theme.Colors` 按界面里的职责分字段，而不是按色相。组件只读这些字段，你画自己的界面时也应当用它们。

背景层次，从外到内：

- **Background**：最外层，标题栏、状态栏、Kbd 的底色。
- **Panel**：侧栏、Header、页签条、抽屉的底色。
- **Editor**：正文和编辑区的底色，Terminal、Diff、内容框也用它。
- **Elevated**：浮在上面的东西：菜单、对话框、Toast、Tooltip、Popover。

文字：

- **Text**：正文。
- **Muted**：次要文字，图标。
- **Placeholder**：输入提示、占位和说明文字。
- **Disabled**：保留字段。目前没有组件读它：禁用的控件由框架统一变成半透明。

线和状态：

- **Border**：控件和面板的边线。
- **BorderSubtle**：更淡的分隔线，例如 Separator、行与行之间。
- **Focus**：焦点环和获得焦点的输入边线。
- **Hover**：悬停底色，也是进度条的轨道色。
- **Active**：按下底色。
- **Selected**：选中行的底色。

强调和语义色：

- **Accent**：强调色：主按钮、链接、选中的开关。
- **AccentText**：强调色底上的文字。
- **AccentSubtle**：带透明度的强调色浅底，用于选中的分段这类强调态。
- **Danger**、**Warning**、**Success**：危险、警告、成功。

其他：

- **Input**：输入框和带边框 Select 的底色。
- **Selection**：文字选区的底色，也用在 Terminal 的选区。
- **Added**、**Deleted**：补丁里新增和删除行的浅底。
- **ANSI**：`[16]ui.Color`，Terminal 和 ANSIOutput 的 16 色调色板。

代码语法颜色另放在 `Theme.Syntax`（类型 `SyntaxColors`，`Keyword`、`String`、`Function`、`Type`、`Number`、`Comment`、`Property`、`Operator`），用于 Markdown 代码块和 Diff，JSON 树的键名用 `Property`。

## 尺寸

`Theme.Metrics` 保存 4 像素网格、三档控件高度和两档圆角：

- **Grid**：4，一个网格单位。
- **CompactHeight**：20，紧凑控件，例如 Badge、Kbd、SmallIconButton、StatusButton。
- **ControlHeight**：28，普通控件和列表行。
- **HeaderHeight**：36，标题栏和页签。
- **StatusHeight**：28，状态栏。
- **Radius**：4，控件圆角。
- **RadiusLarge**：8，浮层和卡片圆角。

`Radius` 设为 0 是有效的，表示直角，不会回退成默认值。

`kit.RadiusFull` 是圆形控件用的圆角（9999），比任何控件边长的一半都大，画出来就是圆，开关和状态点用它。

## 文字层级

文字不要自己写字号和行高，用 `Theme.TextRole(element, role)` 套上一个文字角色。它设置字体、字号、固定行高、字重，颜色和对齐仍由你决定。

```go
t.TextRole(ui.Text(c, "上下文"), kit.TextCaption).TextColor(t.Colors.Muted)
```

- **TextProse**：14/24，Markdown 正文、详情说明。
- **TextUI**：14/20，行标题、按钮、输入框、工具名。
- **TextCaption**：12/16，时间、提示、次要状态。
- **TextToken**：12/20，token 数，等宽数字。
- **TextCode**：12/20，代码、JSON、工具输出，等宽字体。
- **TextCodeLabel**：12/16，小节标签和内容框标题，等宽字体，调用方配次要色。
- **TextTitle**：16/24，小节和页面标题，字重 400。

字体和字号来自 `Theme.Typography`：`Font`、`MonoFont`、`BodySize`（14）、`CaptionSize`（12）、`HeadingSize`（16）、`CodeSize`（12）。

## 图标尺寸

图标有三档：

- **IconSizeSmall**：12，20 高紧凑控件和行尾指示（下拉箭头）。
- **IconSizeDefault**：14，28 高按钮、行首图标、状态栏的默认值。
- **IconSizeLarge**：16，空状态和 36 高标题栏的主操作。

```go
ui.Icon(c, kit.IconSettings).Size(kit.IconSizeDefault, kit.IconSizeDefault)
```

## 间距

`Theme.Space(n)` 返回 `n` 个网格单位，也就是 `4*n`。只用这几档：

- `Space(1)` = 4
- `Space(2)` = 8
- `Space(3)` = 12
- `Space(4)` = 16
- `Space(6)` = 24
- `Space(8)` = 32

`Space(5)` = 20 只用在行高上，不用作间距。

```go
ui.Column(c).Padding(t.Space(6)).Gap(t.Space(3))
```

## 动效

`Theme.Motion` 保存动画时长和“减少动画”开关，见 [Motion](/components/motion)。

## 暗色

明暗主题见[暗色模式](/dark-mode)。
