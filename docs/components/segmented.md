---
title: Segmented
description: 分段控件，一排互斥的选项，可带图标和提示。
---

# Segmented

Zed 的 ToggleButtonGroup 外观：28 高，一圈细边线，各段之间一条竖线。选中的段用强调色浅底和强调色文字。方向键切换，整体是一个 Tab 停靠点。

<Preview page="Segmented" :height="460">

<<< @/../cmd/gallery/segmented.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
appearance := 0 // 选中段的序号，由应用持有

kit.Segmented(c.Key("appearance"), &appearance, []string{"系统", "明", "暗"}).Label("外观").AlignSelf(ui.Start)
```

默认宽度贴合内容。放进纵向容器时要加 `AlignSelf(ui.Start)` 才靠左。控件没有可见的标题，请用 `.Label(...)` 给出可访问名称。

## 示例

### 带图标和提示

用 `SegmentedWith` 和 `SegmentedItem`。`Tooltip` 为空就不显示提示。

```go
sections := []kit.SegmentedItem{
	{Icon: kit.IconFolder, Label: "工作区", Tooltip: "执行器和目录"},
	{Icon: kit.IconServerCog, Label: "MCP", Tooltip: "已连接的 MCP 服务"},
}
kit.SegmentedWith(c.Key("sections"), &current, sections, kit.SegmentedOpt{}).Label("分区")
```

### 占满宽度

`FullWidth` 让控件占满容器，各段等宽。

```go
kit.SegmentedWith(c.Key("full"), &current, sections, kit.SegmentedOpt{FullWidth: true})
```

## API

### Segmented

`func Segmented(c *ui.Context, selected *int, labels []string) ui.Element`

只有文字的分段控件，宽度贴合内容。等于不带图标和提示的 `SegmentedWith`。

- **selected** `*int`：选中段的序号，由调用方持有。
- **labels** `[]string`：每段的文字。

### SegmentedWith

`func SegmentedWith(c *ui.Context, selected *int, items []SegmentedItem, opt SegmentedOpt) ui.Element`

### SegmentedItem

- **Icon** `*ui.SVG`：段内图标，12。可以为 `nil`。
- **Label** `string`：段的文字，同时是可访问名称。
- **Tooltip** `string`：悬停提示。可以为空。

### SegmentedOpt

- **FullWidth** `bool`：占满容器宽度，各段等宽。默认 `false`，各段贴合自己的内容。
