---
title: CompletionMenu
description: 输入框上方的补全浮层，输入框保持焦点，接管上下键、Enter、Tab 和 Esc。
---

# CompletionMenu

输入框上方弹出的补全浮层，例如输入 `/` 后列出命令。输入框保持焦点，用户继续打字，浮层只接管几个按键。

<Preview page="CompletionMenu" :height="480">

<<< @/../cmd/gallery/completion.go

</Preview>

## 用法

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
box := ui.Column(c.Key("composer")).Padding(8).Radius(8).Border(1, t.Colors.Border)
var input ui.Element
box.Children(func() {
	input = kit.TextAreaWith(c.Key("input"), &text, kit.InputOpt{Borderless: true}).Label("消息")
})

opt := kit.CompletionOpt{Open: strings.HasPrefix(text, "/"), Items: matches(text)}
if i := kit.CompletionMenu(c, box, input, opt); i >= 0 {
	text = opt.Items[i].Title + " "
}
```

调用方负责：

- 判断什么时候打开（`Open`），并在每帧给出当前的候选。组件不过滤、不排序。
- 处理选中。返回值 `>= 0` 就是本帧选中的候选序号，把它写回输入框。

## 按键

有候选时：

- 上下键循环移动高亮。
- Enter 或 Tab 选中。
- Esc 关闭，直到 `Open` 变回 false 才会再次出现。

没有候选时只接管 Esc。输入法组字期间不接管任何键。

接管按键用的是 `input` 上的 `HandleInput`，不要再给同一个输入框设置 `HandleInput`；`OnTextKey` 不受影响。

## 示例

### 加载中和出错

候选还在加载时设 `Loading`，它显示在候选后面（没有候选时单独显示）。出错时设 `Error`。两者都没有、候选也为空时，显示 `Empty`。

```go
kit.CompletionOpt{Open: true, Loading: true}
kit.CompletionOpt{Open: true, Error: "无法加载技能：连接被拒绝"}
kit.CompletionOpt{Open: true, Empty: "没有匹配的技能"}
```

## API

### CompletionMenu

`func CompletionMenu(c *ui.Context, anchor, input ui.Element, opt CompletionOpt) (selected int)`

浮层显示在 `anchor` 上方，最多 8 行（候选、加载行、错误行合计），超出滚动。行内容的左缘和 `input` 里的文字对齐。

- **anchor** `ui.Element`：浮层依附的元素，通常是包着输入框的那个容器。
- **input** `ui.Element`：输入框元素，保持焦点。
- **opt** `CompletionOpt`：浮层的数据和状态。
- **返回值** `int`：本帧选中的候选序号，没有选中时是 `-1`。

### CompletionOpt

- **Open** `bool`：为真才显示浮层。由调用方按输入内容决定，例如光标所在的词以 `/` 开头。
- **Items** `[]CompletionItem`：候选列表。
- **Loading** `bool`：在候选后面显示加载行。
- **Error** `string`：非空时显示一行错误。
- **Empty** `string`：没有候选、没有加载也没有错误时的提示。为空用内置的“没有匹配项”。

### CompletionItem

- **Icon** `*ui.SVG`：行首图标。
- **Title** `string`：标题。
- **Description** `string`：显示在标题右侧的次要色说明，放不下就截断。
