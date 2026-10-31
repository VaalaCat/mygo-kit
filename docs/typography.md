---
title: 字体与图标
description: 内置的 IBM Plex Sans 和 Lilex、RegisterFonts、全部图标常量，以及 Kbd。
---

# 字体与图标

## 字体

本库内置两个字体族，和 Zed 默认的界面字体、代码字体对应：

- **IBM Plex Sans**（常量 `kit.FontSans`）：界面和正文。常规、半粗、斜体、半粗斜体。
- **Lilex**（常量 `kit.FontMono`）：代码、终端、JSON、补丁。常规、粗体、斜体。

字体文件嵌在库里。应用启动时调用一次 `RegisterFonts`，在创建窗口之前：

```go
if err := kit.RegisterFonts(); err != nil {
	log.Fatal(err)
}
```

重复调用只注册一次，也没有副作用。

主题里的字体名就是这两个常量，`Theme.Typography.Font` 是 `FontSans`，`Theme.Typography.MonoFont` 是 `FontMono`。想用自己的字体，先用 MyGo 的 `ui.RegisterFont` 注册，再改主题：

```go
theme := kit.OneLight()
theme.Typography.Font = "My Sans"
kit.ApplyTheme(c, theme)
```

自己的元素要用等宽字体时：

```go
t := kit.GetTheme(c)
t.TextRole(ui.Text(c, "go test ./..."), kit.TextCode)
```

字号和行高的层级见[主题](/theming#文字层级)。

### 浏览器里

浏览器构建（`GOOS=js GOARCH=wasm`）只注册各字体族的常规体。浏览器只按族名区分字体，同族的粗体、斜体会盖掉常规体，所以粗体由浏览器合成。桌面构建不受影响。

## 图标

图标是 Lucide 的线性图标，统一 1.5 描边，颜色跟随文字。每个图标是一个 `*ui.SVG` 变量，直接传给组件，或者用 `ui.Icon` 画：

```go
kit.IconButton(c.Key("settings"), kit.IconSettings, "设置")

ui.Icon(c, kit.IconSearch).Size(kit.IconSizeDefault, kit.IconSizeDefault).TextColor(t.Colors.Muted)
```

尺寸常量 `IconSizeSmall`（12）、`IconSizeDefault`（14）、`IconSizeLarge`（16）见[主题](/theming#图标尺寸)。

用自己的 SVG 时，用 `ui.ParseSVG` 或 `ui.MustParseSVG` 解析，得到的 `*ui.SVG` 和内置图标一样用。线条颜色用 `currentColor` 才会跟随文字颜色。

### 全部图标

面板和布局：

- `IconPanelLeft`、`IconPanelRight`、`IconPanelBottom`
- `IconPanelBottomClose`、`IconPanelRightClose`
- `IconFold`、`IconFiles`、`IconListTree`、`IconListEnd`

箭头和方向：

- `IconChevronDown`、`IconChevronLeft`、`IconChevronRight`
- `IconArrowUp`、`IconArrowDown`

操作：

- `IconPlus`、`IconX`、`IconCheck`、`IconSearch`、`IconSettings`
- `IconCopy`、`IconPencil`、`IconTrash`、`IconDownload`、`IconReplace`
- `IconRefresh`、`IconRotateCcw`
- `IconPlay`、`IconStop`、`IconCirclePause`、`IconCircleStop`
- `IconEllipsis`、`IconImagePlus`

状态和提示：

- `IconInfo`、`IconAlert`、`IconTriangleAlert`
- `IconCircleCheck`、`IconCircleX`、`IconBan`
- `IconLoader`、`IconGauge`、`IconActivity`

对象：

- `IconFolder`、`IconFileText`、`IconArchive`、`IconBookOpen`
- `IconTerminal`、`IconCommand`、`IconBraces`
- `IconBot`、`IconBrain`、`IconWrench`、`IconBug`
- `IconBranch`、`IconGitFork`、`IconWebhook`、`IconServerCog`
- `IconMessageSquare`

## Kbd

`Kbd` 是显示快捷键的小标签，用法见 [Kbd](/components/kbd)。它用 `TextCaption` 的字号，放在文字旁边时要和文字一起用 `AlignItems(ui.Center)` 垂直居中。
