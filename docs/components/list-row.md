---
title: ListRow
description: 整行列表项，支持单行、两行、缩进、行尾状态和悬停操作。
---

# ListRow

侧栏、会话列表、文件树里的一行。整行是一个按钮，选中、悬停、按下的底色由主题统一给出。

<Preview page="ListRow" :height="520">

<<< @/../cmd/gallery/listrow.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
for _, name := range names {
	row := kit.ListRow(c.Key(name), kit.ListRowOpt{
		Icon:     kit.IconMessageSquare,
		Title:    name,
		Selected: name == current,
	})
	if row.Clicked() {
		current = name
	}
}
```

同一列表里的行要用 `c.Key(...)` 区分。选中状态由调用方用 `Selected` 传入，组件不保存。

## 示例

### 两行

`Description` 非空时行变成两行、44 高，说明用 12/16 的次要色，单行截断。

```go
kit.ListRow(c.Key("a"), kit.ListRowOpt{Title: "给客户端加重试", Description: "openvaala · 3 分钟前"})
```

### 行尾状态

`Trailing` 画徽标、转圈或状态点。单行时和标题按基线对齐，两行时在整行垂直居中。

```go
kit.ListRow(c.Key("b"), kit.ListRowOpt{
	Title:       "更新变更日志",
	Description: "docs · 昨天",
	Trailing:    func() { kit.StatusDot(c, kit.GetTheme(c).Colors.Accent).Label("未读") },
})
```

### 悬停操作

`Actions` 画行尾的 28 见方图标按钮，只在悬停、选中或键盘聚焦到行内时显示。不显示时仍占位，行内容不会跳。

```go
kit.ListRow(c.Key("c"), kit.ListRowOpt{
	Title: "某个会话",
	Actions: func() {
		kit.IconButtonTip(c.Key("more"), kit.IconEllipsis, "更多")
	},
})
```

### 缩进

`Depth` 每级向右缩进 16，用来画目录树。

```go
kit.ListRow(c.Key("file"), kit.ListRowOpt{Icon: kit.IconFileText, Title: "composer.go", Depth: 2})
```

### 禁用

```go
kit.ListRow(c.Key("off"), kit.ListRowOpt{Icon: kit.IconBan, Title: "不可用"}).Disabled(true)
```

## API

### ListRow

`func ListRow(c *ui.Context, opt ListRowOpt) ui.Element`

28 高（带说明时 44）。返回的元素是按钮，调用方用 `Clicked()` 响应点击，也可以链式调用 `FillWidth`、`Key`、`Label`。

### ListRowOpt

- **Icon** `*ui.SVG`：行首图标，14，画在 16 宽的图标位中央。为 `nil` 时不留图标位。
- **Title** `string`：标题，单行截断，同时是可访问名称。
- **Description** `string`：非空时行变成两行、44 高。
- **Selected** `bool`：显示选中底色。
- **Depth** `int`：缩进层级，每级 16。默认 0。
- **Trailing** `func()`：绘制行尾的状态或徽标。
- **Actions** `func()`：绘制行尾的操作按钮，只在悬停、选中或键盘聚焦到行内时显示。设置后右侧内边距缩到 4，按钮和标题栏里的图标按钮在同一条右对齐线上。
