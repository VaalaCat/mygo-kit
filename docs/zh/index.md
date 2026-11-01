---
title: 介绍
description: mygo-kit 是基于 MyGo 的原生 Go 组件库，Zed 风格的连续面板、紧凑控件和阅读组件。
---

# 介绍

mygo-kit 是一套 Go 写的原生界面组件，跑在 [MyGo](https://github.com/egoist/mygo) 的原生 UI 框架上。它给 MyGo 加上一层统一的外观：连续面板、紧凑工具栏、清晰的文字层级。键盘、选区、输入法、焦点这些行为仍然由 MyGo 负责。

<Preview page="Overview" :height="560">

<<< @/../cmd/gallery/overview.go

</Preview>

下面是原生窗口（不是浏览器）里的截图，会跟着站点的明暗切换。

<img class="kit-shot kit-shot-light" src="/screenshots/overview-light.png" alt="画廊总览，明色">
<img class="kit-shot kit-shot-dark" src="/screenshots/overview-dark.png" alt="画廊总览，暗色">
<img class="kit-shot kit-shot-light" src="/screenshots/settings-light.png" alt="设置页，明色">
<img class="kit-shot kit-shot-dark" src="/screenshots/settings-dark.png" alt="设置页，暗色">

<style>
.kit-shot { display: block; width: 100%; margin: 12px 0; border: 1px solid var(--vp-c-divider); border-radius: 8px; }
.kit-shot-dark { display: none; }
.dark .kit-shot-light { display: none; }
.dark .kit-shot-dark { display: block; }
</style>

## 设计取向

外观照 Zed 做：

- **连续面板**：侧栏、标题栏、状态栏共用分隔线，不给每块区域再套一圈边框。
- **ghost 工具按钮**：平时没有边框，悬停才出现底色。
- **小圆角**：控件圆角 4，浮层和卡片圆角 8；结构面板保持平直。
- **4 网格**：间距按 4 的倍数取，控件高度只有 20、28、36 三档。
- **One Light / One Dark**：内置两套 Zed One 配色，字体是 IBM Plex Sans 和 Lilex。

## 它不做什么

- 不创建应用运行时。窗口、事件循环、`mygo.App.Run` 都由你的应用自己写。
- 不联网。图片数据、下载、文件读取和重试都由你的应用提供，组件只负责显示。
- 不启动后台 goroutine。需要动画的组件靠 MyGo 的帧回调驱动，唯一例外是 `TerminalState`，见 [Terminal](/zh/components/terminal)。
- 不用全局可变主题。主题由应用持有，每帧交给 `ApplyTheme`。

## 怎么用

每个组件都是一个函数，参数是 `*ui.Context`，返回 MyGo 的 `ui.Element`。你可以继续链式设置尺寸、标签、焦点，也可以放进 `ui.Column`、`ui.Row` 里组合。

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.ApplyTheme(c, kit.OneLight())
if kit.PrimaryButton(c, "保存").Clicked() {
	// 保存
}
```

有状态的组件（输入、开关、列表）由你的应用持有状态，调用时传指针。输入类组件要先用 `c.Key(...)` 给身份，这样重绘时才能保住焦点和光标。

## 下一步

- [安装](/zh/installation)：环境要求、`go get` 和最小示例。
- [主题](/zh/theming)、[暗色模式](/zh/dark-mode)、[字体与图标](/zh/typography)、[多语言](/zh/i18n)。
- 从 [Button](/zh/components/button) 开始看组件，每页都带可交互的预览和完整源码。
