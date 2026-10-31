---
title: CommandPalette
description: 可搜索的命令面板，按分组显示命令，选中后返回命令 ID。
---

# CommandPalette

居中偏上的命令面板：顶部是搜索框，下面按分组列出命令。搜索框自动聚焦，上下键移动高亮，Enter 选中，Esc 关闭。

<Preview page="CommandPalette" :height="460">

<<< @/../cmd/gallery/command.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // 由应用持有，快捷键里置 true 就打开

commands := []kit.Command{
	{ID: "conversation.new", Title: "新建会话", Group: "会话", Keys: "⌘N", Icon: kit.IconPlus},
	{ID: "view.sidebar", Title: "切换侧栏", Group: "视图", Keys: "⌘B", Icon: kit.IconPanelLeft},
}
if id := kit.CommandPalette(c, &open, commands); id != "" {
	run(id)
}
```

选中命令后，面板在下一帧关闭，并返回命令的 `ID`；其余帧返回空字符串。分派到哪个动作由调用方按 `ID` 决定。每个命令的 `ID` 要唯一，它也是这一行的身份。

## 搜索规则

- 搜索词不区分大小写，只匹配命令的 `Title`，不匹配 `Group` 和 `ID`。
- 分组按命令首次出现的顺序排列，`Group` 为空的命令不显示分组标题。
- 禁用的命令会显示，但不能高亮，也选不中。
- 没有结果时显示“没有匹配项”。

## API

### CommandPalette

`func CommandPalette(c *ui.Context, open *bool, commands []Command) string`

宽 560（最多占窗口 90%），结果列表最高 360，超出滚动。

- **open** `*bool`：是否打开，由调用方持有。选中命令、按 Esc 或点遮罩时被置为 false。
- **commands** `[]Command`：全部命令，每帧传入。
- **返回值** `string`：本帧选中的命令 `ID`，没有则为空字符串。

### Command

- **ID** `string`：命令的唯一标识，选中后返回它。
- **Title** `string`：显示和搜索用的标题。
- **Group** `string`：分组名。
- **Keys** `string`：行尾显示的快捷键提示，例如 `"⌘N"`，用 [Kbd](/components/kbd) 画。只用于显示，不注册快捷键。
- **Icon** `*ui.SVG`：行首图标。只要有一条命令有图标，没有图标的命令就会留出图标位。
- **Disabled** `bool`：命令变淡，不可选。
