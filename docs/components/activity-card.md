---
title: ActivityCard
description: 一行条目标题，可展开详情，带状态图标，用于工具调用、思考这类过程记录。
---

# ActivityCard

一行 28 高的条目标题，有详情时可以展开。常用来显示 Agent 的工具调用、思考、错误这类过程。

```
[图标/箭头] 标题  状态  摘要 ········ [操作]
```

<Preview page="ActivityCard" :height="460">

<<< @/../cmd/gallery/activitycard.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // 展开状态由应用持有

kit.ActivityCard(c, &open, kit.ActivityCardOpt{
	Icon:    kit.IconTerminal,
	Title:   "bash",
	Summary: "go test ./...",
	Status:  kit.ActivityRunning,
}, nil, func() {
	ui.Text(c, "$ go test ./...")
})
```

最后两个参数分别是右侧操作和展开后的详情。`detail` 传 `nil` 时这一行不可展开，也没有箭头，只是一行普通文字。每张卡片的 `open` 都要各自持有。

## 示例

### 悬停才出现的操作

`actions` 里的图标按钮独立于展开按钮，只在悬停或键盘聚焦这一行时显示，占位不变。

```go
inspect := func() { kit.IconButton(c, kit.IconBug, "检查") }
kit.ActivityCard(c, &open, opt, inspect, detail)
```

### 状态

`Status` 在标题后面显示一个图标；空值表示成功或没有状态，不显示图标。

```go
kit.ActivityCardOpt{Icon: kit.IconPencil, Title: "apply_patch", Status: kit.ActivityFailed}
```

### 语义色

`Tone` 改变标题和图标的颜色。

```go
kit.ActivityCardOpt{Icon: kit.IconTriangleAlert, Title: "模型调用失败", Tone: kit.ActivityDanger}
```

### 进行中

`Shimmer` 让标题闪烁，用于思考、压缩、运行中。

```go
kit.ActivityCardOpt{Icon: kit.IconBrain, Title: "思考中", Shimmer: true}
```

### 标题后的自定义内容

`Extra` 紧跟标题，在展开按钮之外，可以放可交互的内容，例如会话链接。设置后不再显示 `Summary`。

## API

### ActivityCard

`func ActivityCard(c *ui.Context, open *bool, opt ActivityCardOpt, actions func(), detail func()) ui.Element`

详情缩进到标题文字的位置，左侧竖线对齐图标中线，块之间间隔 12。

- **open** `*bool`：展开状态，由调用方持有。
- **opt** `ActivityCardOpt`：标题行的内容。
- **actions** `func()`：绘制右侧操作按钮的函数，可以为 `nil`。
- **detail** `func()`：绘制展开内容的函数。`nil` 表示不可展开。

### ActivityCardOpt

- **Icon** `*ui.SVG`：类型图标。可展开的条目在悬停、聚焦和展开时换成箭头。
- **Title** `string`：标题。
- **Summary** `string`：参数摘要，单行截断，不加悬停提示，占位色。
- **Status** `ActivityStatus`：状态，默认空，不显示。
- **Tone** `ActivityTone`：语义色，默认 `ActivityNormal`。
- **Shimmer** `bool`：标题闪烁。
- **Extra** `func()`：紧跟标题的自定义内容。

### ActivityStatus

- **ActivityRunning**：运行中，显示转圈。
- **ActivityFailed**：失败，危险色图标。
- **ActivityInterrupted**：已中断，警告色。
- **ActivityBlocked**：被拦截，警告色。
- **ActivityStopped**：已停止，警告色。

### ActivityTone

- **ActivityNormal**：次要色，悬停时变正文色。
- **ActivityWarning**：警告色。
- **ActivityDanger**：危险色。

### StatusIcon

`func StatusIcon(c *ui.Context, status ActivityStatus)`

单独画一个状态图标，不返回元素。成功（空值）什么也不画。图标有对应语言的可访问名称，例如“运行中”。可以在 `ActivityCard` 之外的列表里复用。
