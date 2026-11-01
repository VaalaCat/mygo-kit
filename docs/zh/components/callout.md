---
title: Callout
description: 带色调的提示条：图标、标题、说明，右侧可放操作。
---

# Callout

放在内容里的提示条，淡色背景加细边框。有错误、警告、信息三种。

<Preview page="Callout" :height="300">

<<< @/../cmd/gallery/callout.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Callout(c, kit.CalloutError, "模型请求失败", "服务返回 429：请求太频繁。几秒后重试。", func() {
	if kit.Button(c, "重试").Clicked() {
		retry()
	}
})
```

## 示例

### 只有标题

`detail` 传空字符串时只显示标题，整条内容垂直居中。

```go
kit.Callout(c, kit.CalloutInfo, "会话已分叉", "", nil)
```

### 带操作

`actions` 传 `nil` 时不留操作区，见上例。传了函数时，按钮画在右侧，可以放几个：

```go
kit.Callout(c, kit.CalloutInfo, "会话已分叉", "", func() {
	kit.Button(c, "打开")
	kit.IconButton(c, kit.IconX, "关闭")
})
```

## API

### Callout

`func Callout(c *ui.Context, kind CalloutKind, title, detail string, actions func()) ui.Element`

圆角 8，内边距 12。标题是正文色，详情是 12/16 次要色。

- **kind** `CalloutKind`：种类，决定图标和色调。
- **title** `string`：标题，同时是可访问名称。
- **detail** `string`：次要色详情，可以为空。
- **actions** `func()`：绘制右侧操作的函数，可以为 `nil`。

### CalloutKind

- **CalloutError**：错误，危险色。
- **CalloutWarning**：警告，警告色。
- **CalloutInfo**：信息，强调色。
