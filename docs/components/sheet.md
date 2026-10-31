---
title: Sheet
description: 贴靠窗口左侧或右侧的模态抽屉，点遮罩或按 Esc 关闭。
---

# Sheet

从窗口左边或右边滑出来的模态抽屉，带半透明遮罩。点遮罩或按 Esc 会把 `*open` 置为 false。

<Preview page="Sheet" :height="360">

<<< @/../cmd/gallery/sheet.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // 由应用持有

if kit.Button(c, "打开详情").Clicked() {
	open = true
}
kit.Sheet(c.Key("details"), &open, kit.SheetRight, 360, func() {
	kit.Header(c, "详情", nil)
	ui.Text(c, "抽屉里的内容").Padding(12)
})
```

抽屉每帧都要调用，`open` 为 false 时不画。

## 和 TitleBar 的关系

窗口里用了 [TitleBar](/components/title-bar) 时，抽屉和遮罩从顶栏下缘开始，顶栏和系统窗口按钮不会被盖住。没有 TitleBar 时铺满整个窗口。`TitleBar` 的高度按本帧或上一帧记录，调用顺序只影响第一帧。

## API

### Sheet

`func Sheet(c *ui.Context, open *bool, side SheetSide, width float32, content func()) ui.Element`

抽屉用面板底色，只在靠近内容的一侧画一条分隔线，最宽占窗口 90%。

- **open** `*bool`：是否打开，由调用方持有。
- **side** `SheetSide`：贴靠的边。
- **width** `float32`：抽屉宽度。
- **content** `func()`：绘制抽屉内容的函数。

### SheetSide

- **SheetLeft**：贴靠窗口左侧。
- **SheetRight**：贴靠窗口右侧。
