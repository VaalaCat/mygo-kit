---
title: Kbd
description: 显示快捷键的小标签，只管外观，不注册快捷键。
---

# Kbd

一枚显示快捷键的小标签，例如 `⌘K`。它只是外观，真正的快捷键要你自己用 MyGo 注册。

<Preview page="Kbd" :height="180">

<<< @/../cmd/gallery/kbd.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
	ui.Text(c, "打开命令面板")
	kit.Kbd(c, "⌘K")
})
```

[CommandPalette](/zh/components/command-palette) 的 `Command.Keys` 也用 Kbd 显示快捷键。[Menu](/zh/components/menu) 的 `MenuItem.Shortcut` 只是右侧的次要色文字，没有边框和底色。

## API

### Kbd

`func Kbd(c *ui.Context, keys string) ui.Element`

20 高的小标签，圆角 4，比所在面板略深的底色加细边线，文字 12/16 次要色。

- **keys** `string`：要显示的文字，例如 `"⌘N"`、`"Ctrl+N"`。组件不解析它。
