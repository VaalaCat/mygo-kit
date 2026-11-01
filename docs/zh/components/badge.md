---
title: Badge
description: 紧凑的状态徽标 Badge，和表示连接状态、未读的小圆点 StatusDot。
---

# Badge

`Badge` 是一枚带语义颜色的小标签；`StatusDot` 是 8 见方的圆点，用来表示连接状态或未读。

<Preview page="Badge" :height="240">

<<< @/../cmd/gallery/badge.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Badge(c, "本地", kit.BadgeNeutral)
kit.Badge(c, "已连接", kit.BadgeSuccess)
```

## 示例

### 状态点

圆点的颜色由你决定，一般取主题里的语义色。圆点没有文字，要用 `.Label(...)` 给出可访问名称。

```go
t := kit.GetTheme(c)
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	kit.StatusDot(c, t.Colors.Success).Label("已连接")
	ui.Text(c, "已连接")
})
```

## API

### Badge

`func Badge(c *ui.Context, text string, kind BadgeKind) ui.Element`

20 高，圆角 4，文字 12/16，背景是文字色的 16% 透明度。不会被所在行拉伸或压缩。

- **text** `string`：徽标文字，单行。
- **kind** `BadgeKind`：语义颜色。

### BadgeKind

- **BadgeNeutral**：中性，次要色。
- **BadgeAccent**：强调色。
- **BadgeSuccess**：成功。
- **BadgeWarning**：警告。
- **BadgeDanger**：危险。

### StatusDot

`func StatusDot(c *ui.Context, color ui.Color) ui.Element`

8×8 的圆点，不会被所在行拉伸或压缩。

- **color** `ui.Color`：圆点颜色。
