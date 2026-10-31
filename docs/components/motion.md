---
title: Motion
description: 减少动画的统一规则，以及淡入淡出的 Fade。
---

# Motion

本库的动画都很短，并且有同一套开关：`Theme.Motion.Reduced` 为真，或系统要求减少动画（`c.Preferences().ReduceMotion`）时，所有动画直接到达终点。

动画时长放在主题里，默认值如下：

- `Motion.Fade`：120ms，淡入淡出、开关滑块、MorphText。
- `Motion.DisclosureOpen`：200ms，Disclosure 和 ActivityCard 展开。
- `Motion.DisclosureClose`：150ms，收起。

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

### 淡入淡出

`Fade` 按 `visible` 把元素的不透明度渐变到 1 或 0。

```go
panel := ui.Column(c).Children(func() { ui.Text(c, "提示") })
kit.Fade(c, panel, showing)
```

`visible` 为 false 时，元素立刻被禁用，不再响应点击和键盘，透明度只影响绘制。要真正把它从布局里去掉，由调用方不再创建它。

### 关闭本库的动画

```go
theme := kit.OneLight()
theme.Motion.Reduced = true
kit.ApplyTheme(c, theme)
```

主题里的开关不会自动跟随系统，不过各组件会同时检查系统偏好，所以系统要求减少动画时即使 `Reduced` 为 false 也会静止。

## 相关组件

- [Spinner](/components/spinner)：每秒一圈的加载图标。
- [Shimmer](/components/shimmer)：高亮扫过的文字。
- [MorphText](/components/morph-text)：字素级的文字变化。
- [Disclosure](/components/disclosure)：展开和收起的高度过渡。

## API

### Fade

`func Fade(c *ui.Context, element ui.Element, visible bool) ui.Element`

更新元素的可见透明度，返回同一个元素。

- **element** `ui.Element`：要淡入淡出的元素。
- **visible** `bool`：目标状态，true 渐变到完全可见，false 渐变到完全透明并禁用。

### Motion

`type Motion struct`，在 `Theme.Motion` 里。

- **Fade** `time.Duration`：淡入淡出时长，默认 120ms。
- **DisclosureOpen** `time.Duration`：展开时长，默认 200ms。
- **DisclosureClose** `time.Duration`：收起时长，默认 150ms。
- **Reduced** `bool`：为真时关闭本库的动画。默认 `false`。
