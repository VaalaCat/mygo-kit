---
title: MorphText
description: 单行状态文字，文字变化时只有变了的字素短暂上推并渐显。
---

# MorphText

显示一行会变化的状态文字，例如“正在读取文件”变成“正在运行测试”。文字变化时，只有变了的字素向上轻推并淡入；没变的前缀、后缀和中间能对上的字素保持原样。

<Preview page="MorphText" :height="200">

<<< @/../cmd/gallery/morph.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.MorphText(c, status) // status 是应用里当前的状态文字
```

## 行为

- 按 Unicode 字素处理，完整的 emoji 和带组合重音的字符不会被拆开。
- 共同的前缀和后缀保持不动。纯追加（后面多出新字符）不播放动画，稳定的前后缀也不会重播。
- 变化的中段最多比较一定数量的字素，变化太大时直接显示新文字，不做动画。
- 系统或主题要求减少动画时，直接显示新文字。
- 文字单行不换行。整段的可访问名称就是当前文字。

## API

### MorphText

`func MorphText(c *ui.Context, value string) ui.Element`

- **value** `string`：当前要显示的文字。每帧传入；组件按上一次的文字算出变化。动画时长取 `Theme.Motion.Fade`（默认 120ms）。
