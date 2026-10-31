---
title: Tree
description: 紧凑的虚拟树，沿用 MyGo 的稳定键、展开和方向键导航。
---

# Tree

可展开的树，例如文件目录。展开状态、方向键导航和虚拟滚动都是 MyGo 的 `ui.Outline`，本库只加上外观。

<Preview page="Tree" :height="320">

<<< @/../cmd/gallery/tree.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var state ui.OutlineState[string] // 展开和选中状态，由应用持有

children := map[string][]string{
	"workspace": {"cmd", "internal"},
	"cmd":       {"gallery"},
}

kit.Tree(c, &state, []string{"workspace"}, func(id string) []string { return children[id] }, func(id string) {
	ui.Row(c).Height(28).Padding(0, 8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, id).SingleLine()
	})
}).Grow(1)
```

树的节点用你自己的标识 `K` 表示（可以是字符串、数字等可比较的类型）。组件通过 `children` 向你问每个节点的子节点，通过 `row` 让你画每一行。每一行的缩进、展开箭头由 `ui.Outline` 处理。

## API

### Tree

`func Tree[K comparable](c *ui.Context, state *ui.OutlineState[K], roots []K, children func(K) []K, row func(K)) ui.Element`

- **state** `*ui.OutlineState[K]`：展开状态，由调用方持有。
- **roots** `[]K`：根节点。
- **children** `func(K) []K`：返回某个节点的子节点，没有子节点返回空。
- **row** `func(K)`：绘制某个节点这一行的内容。
