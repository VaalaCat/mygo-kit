---
title: Menu
description: 点击弹出的菜单 MenuButton，和右键菜单 ContextMenu。
---

# Menu

菜单有两种入口：点击一个元素弹出的 `MenuButton`，和右键（或菜单键）弹出的 `ContextMenu`。两者共用同一组 `MenuItem`，键盘行为也一样。

<Preview page="Menu" :height="520">

<<< @/../cmd/gallery/menu.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
items := []kit.MenuItem{
	{Icon: kit.IconPencil, Title: "重命名", Shortcut: "F2", OnSelect: rename},
	{Separator: true},
	{Icon: kit.IconX, Title: "删除", Danger: true, OnSelect: remove},
}
kit.MenuButton(c, kit.IconButton(c.Key("more"), kit.IconEllipsis, "更多"), items)
```

菜单项的动作写在 `OnSelect` 里，选中、菜单关闭的那一帧调用。菜单每帧都要调用，打开状态由组件保存在触发元素上。

## 示例

### 点击菜单

任何元素都能当触发元素，常见的是 `IconButton` 和 `ToolButton`。键盘操作：上下键、Home/End、首字母移动高亮，Enter 选中，Esc 关闭。键盘聚焦时按 Enter 或 Space 打开，会直接高亮第一项。

```go
kit.MenuButton(c, kit.ToolButton(c.Key("ws"), kit.IconFolder, "工作区", true), []kit.MenuItem{
	{Title: "本地", OnSelect: useLocal},
	{Title: "远程", OnSelect: useRemote},
})
```

返回值 `closed` 为真表示菜单在本帧关闭（Esc、点击外部或触发元素、Tab、选中条目），焦点回到触发元素。想把焦点放到别处，就在返回 true 的那一帧自己设置焦点。

### 右键菜单

`ContextMenu` 给目标元素加右键菜单：右键时在指针处弹出；目标有焦点时，按 Shift+F10 或菜单键，在目标下方弹出。目标要能获得焦点，才能用键盘唤出。

```go
area := ui.Column(c.Key("area")).Height(120).Width(360).Focusable().Children(func() {
	ui.Text(c, "在这里右键")
})
kit.ContextMenu(c, area, items)
```

### 文字上的右键菜单

目标是 `Selectable` 的文字容器时，右键弹出的是这个菜单，而不是系统的文字菜单，拖选不受影响。打开时如果有选中的文字，菜单顶部会多一项“复制”，复制的就是选中的那段。

```go
text := ui.Column(c.Key("text")).Selectable().Children(func() {
	ui.Text(c, "选中其中一部分，再右键。")
})
kit.ContextMenu(c, text, items)
```

`Selectable` 容器必须自己就是 `target`。挂在它的祖先上，容器里的文字仍然显示系统菜单。

## API

### MenuButton

`func MenuButton(c *ui.Context, trigger ui.Element, items []MenuItem) (closed bool)`

给 `trigger` 附上点击后在下方弹出的菜单。点击触发元素开关菜单。

- **trigger** `ui.Element`：触发元素。
- **items** `[]MenuItem`：菜单项。
- **返回值** `bool`：菜单在本帧关闭。

### ContextMenu

`func ContextMenu(c *ui.Context, target ui.Element, items []MenuItem) (closed bool)`

给 `target` 附上右键菜单。返回值含义同 `MenuButton`。

### MenuItem

- **Icon** `*ui.SVG`：条目图标，14。任一条目有图标时，没有图标的条目也会留出图标位，标题对齐。
- **Title** `string`：标题。
- **Shortcut** `string`：行尾右对齐显示的快捷键提示，例如 `"⌘N"`。只显示，不注册快捷键。
- **Danger** `bool`：标题和图标用危险色。
- **Disabled** `bool`：条目变淡且不可选，方向键、首字母跳转都会跳过。
- **Separator** `bool`：为真时这一项只画一条分隔线，其余字段忽略。
- **OnSelect** `func()`：条目被选中、菜单关闭的那一帧调用。
