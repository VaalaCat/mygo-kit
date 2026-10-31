---
title: Shimmer
description: 次要色单行文字，高亮从左向右扫过，用于思考、生成标题这类进行中的状态。
---

# Shimmer

一行次要色的文字，一条高亮带每 2 秒从左扫到右。用在“思考中”“正在生成标题”这类进行中的状态。

<Preview page="Shimmer" :height="200">

<<< @/../cmd/gallery/shimmer.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Shimmer(c, "思考中…")
```

同一页里放多个 Shimmer 时，用 `c.Key(...)` 区分，每个各自计时。

[ActivityCard](/components/activity-card) 的 `Shimmer` 字段也是用它让标题闪烁。

## 减少动画

系统或主题要求减少动画时，Shimmer 变成一段不动的次要色文字，仍然可读、可查找。

## API

### Shimmer

`func Shimmer(c *ui.Context, value string) ui.Element`

单行、不换行。高亮带两侧各覆盖 4 个字素，字素按到高亮中心的距离，在次要色和正文色之间过渡。

- **value** `string`：要显示的文字，同时是可访问名称。
