---
title: Toast
description: 窗口右上角的通知：一次 ShowToast 调用，加一个 ToastHost。
---

# Toast

窗口右上角的通知，最新的在最上面。有信息、成功、错误三种语义，可以带说明和一个操作按钮。

<Preview page="Toast" :height="300">

<<< @/../cmd/gallery/toast.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

每个窗口的视图里调用一次 `ToastHost`，通知才会显示。一般放在视图函数开头，用 `defer` 让它最后画：

```go
func view(c *ui.Context) {
	kit.ApplyTheme(c, kit.OneLight())
	defer kit.ToastHost(c)
	// 画界面
}
```

要发通知的地方调用 `ShowToast`：

```go
if kit.Button(c, "导出").Clicked() {
	kit.ShowToast(c, kit.ToastOpt{Title: "会话已导出", Kind: kit.ToastInfo})
}
```

窗口空闲时通知也会立刻出现，不用自己请求重绘。

## 示例

### 带操作按钮

`Action` 非空时显示一个文字按钮，点击后先执行 `OnAction`，再关闭通知。

```go
kit.ShowToast(c, kit.ToastOpt{
	Title:    "会话已归档",
	Kind:     kit.ToastSuccess,
	Action:   "撤销",
	OnAction: undoArchive,
})
```

### 错误通知

错误通知显示 8 秒，比其他通知停留更久。

```go
kit.ShowToast(c, kit.ToastOpt{Title: "连接中断", Description: "运行时没有响应。", Kind: kit.ToastError})
```

## API

### ToastHost

`func ToastHost(c *ui.Context) ui.Element`

在窗口右上角、顶栏下方绘制通知。通知宽 360，每条带语义图标、标题、可选说明、操作和关闭按钮。窗口里用了 [TitleBar](/zh/components/title-bar) 时，通知从顶栏下缘留 8 起向下排；没有 TitleBar 时从窗口顶部留 8。每个窗口调用一次。

### ShowToast

`func ShowToast(c *ui.Context, opt ToastOpt)`

窗口里要已经有 `ToastHost`。

### ToastOpt

- **Title** `string`：标题。
- **Description** `string`：次要色说明。可以为空。
- **Kind** `ToastKind`：语义。零值（空字符串）按信息样式显示。
- **Action** `string`：操作按钮文字，非空时显示按钮。
- **OnAction** `func()`：点击操作按钮时调用，之后通知关闭。

### ToastKind

- **ToastInfo**：信息，强调色图标。
- **ToastSuccess**：成功，绿色图标。
- **ToastError**：错误，危险色图标，停留 8 秒。
