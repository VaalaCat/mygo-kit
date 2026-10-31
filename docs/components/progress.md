---
title: Progress
description: 4 高的线性进度条 ProgressBar，和环形进度 ProgressRing（用量到阈值会变色）。
---

# Progress

`ProgressBar` 是细长的线性进度条；`ProgressRing` 是环形进度，常用来显示上下文用量。两者的 `value` 都是 0 到 1 的比例，超出范围会被截到 0 或 1。

<Preview page="Progress" :height="360">

<<< @/../cmd/gallery/progress.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.ProgressBar(c, 0.42)
```

条的宽度由所在容器决定，所以放进有宽度的容器里：

```go
ui.Column(c).Width(320).Children(func() {
	kit.ProgressBar(c, done)
})
```

## 示例

### 用量环

```go
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	kit.ProgressRing(c, used, 16)
	ui.Text(c, "上下文窗口")
})
```

环在 `value` 达到 0.85 时变成警告色，达到 0.95 时变成危险色。`value` 大于 0 时，弧至少画 8°，这样很低的占比也能看出来；`value` 为 0 时只有底环。

## API

### ProgressBar

`func ProgressBar(c *ui.Context, value float32) ui.Element`

4 高，悬停色轨道，强调色填充，两端圆角。可访问名称是百分比，例如 `42%`。

- **value** `float32`：0 到 1 的比例。

### ProgressRing

`func ProgressRing(c *ui.Context, value float32, size float32) ui.Element`

环线宽 2，从顶部顺时针画。可访问名称始终是真实的百分比（不受最小弧长影响）。

- **value** `float32`：0 到 1 的比例。
- **size** `float32`：环的边长。画廊里用 14、16、20。
