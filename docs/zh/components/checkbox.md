---
title: Checkbox
description: 直角复选框，标签可点击，状态由调用方的 bool 持有。
---

# Checkbox

16 见方的复选框加一个可点击的标签。勾选、Space 键和无障碍行为沿用 MyGo 的复选框。

<Preview page="Checkbox" :height="200">

<<< @/../cmd/gallery/checkbox.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var autoSave bool // 由应用持有

kit.Checkbox(c.Key("auto-save"), &autoSave, "自动保存")
```

点击后 `*checked` 被直接翻转。本帧是否变化用返回元素的 `.Changed()` 判断。

## 示例

### 禁用

```go
kit.Checkbox(c.Key("disabled"), &autoSave, "不可用").Disabled(true)
```

## API

### Checkbox

`func Checkbox(c *ui.Context, checked *bool, label string) ui.Element`

整行 28 高，方框圆角 4，选中时是强调色底加勾。

- **checked** `*bool`：勾选状态，由调用方持有。
- **label** `string`：标签文字，同时是可访问名称。点标签和点方框一样切换。
