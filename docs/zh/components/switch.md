---
title: Switch
description: 36×20 的开关，开时轨道用强调色，滑块带过渡动画。
---

# Switch

开关。开时轨道是强调色，滑块滑动过去；系统或主题要求减少动画时，滑块直接到位。

<Preview page="Switch" :height="360">

<<< @/../cmd/gallery/switch.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var reduceMotion bool // 由应用持有

kit.Switch(c.Key("motion"), &reduceMotion, "减少动画")
```

## 示例

### 只有轨道

`label` 传空字符串时只画轨道，这时要用 `.Label(...)` 给出可访问名称。常放在 [SettingRow](/zh/components/settings) 或 [Field](/zh/components/field) 里，文字由旁边的行标题负责。

```go
kit.Switch(c.Key("auto-save"), &autoSave, "").Label("自动保存")
```

### 禁用

禁用由框架统一变半透明。

```go
kit.Switch(c.Key("off"), &on, "不可用").Disabled(true)
```

## API

### Switch

`func Switch(c *ui.Context, on *bool, label string) ui.Element`

- **on** `*bool`：开关状态，由调用方持有，点击或按 Space 时被翻转。
- **label** `string`：显示在轨道右侧的文字，整行（28 高）可点击。为空时只有轨道（20 高），不带文字。
