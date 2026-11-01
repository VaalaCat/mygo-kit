---
title: 安装
description: 环境要求、go get、一个能直接运行的最小示例，以及怎么运行画廊。
---

# 安装

## 环境要求

- Go 1.27.3。
- MyGo v0.7.3，`go get` 本库时会一起拉取。
- Linux：需要 GTK 3 来开原生窗口。离屏的文字绘制仍要用 Pango 和 Cairo，通常随 GTK 一起装好。
- macOS：使用系统自带的原生框架，不用另装。

本库只依赖 MyGo，不依赖任何业务代码。

## 安装

```sh
go get github.com/VaalaCat/mygo-kit
```

代码里这样导入，包名是 `kit`：

```go
import kit "github.com/VaalaCat/mygo-kit"
```

## 最小示例

新建一个目录，放下面的 `main.go`，运行 `go mod init` 和 `go mod tidy`，再 `go run .`：

```go
package main

import (
	"fmt"
	"log"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// clicks 由应用自己持有，组件不保存业务状态。
var clicks int

func view(c *ui.Context) {
	// 每帧先应用主题，再画界面。
	kit.ApplyTheme(c, kit.OneLight())
	t := kit.GetTheme(c)
	ui.Column(c).Fill().Gap(0).Background(t.Colors.Editor).Children(func() {
		kit.Header(c, "Hello mygo-kit", nil)
		ui.Column(c).Padding(t.Space(6)).Gap(t.Space(3)).Children(func() {
			ui.Text(c, fmt.Sprintf("点击了 %d 次", clicks))
			if kit.PrimaryButton(c, "点我").AlignSelf(ui.Start).Clicked() {
				clicks++
			}
		})
	})
}

func main() {
	// 字体只需注册一次。
	if err := kit.RegisterFonts(); err != nil {
		log.Fatal(err)
	}
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{Title: "Hello", Width: 480, Height: 320, Content: ui.View(view)})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
```

几点说明：

- 窗口、事件循环和 `mygo.App.Run` 都是应用自己写的，本库不创建它们。
- `ui.View(view)` 的 `view` 每帧调用一次，所以 `ApplyTheme` 也每帧调用，见[主题](/zh/theming)。
- `kit.RegisterFonts()` 注册内置的 IBM Plex Sans 和 Lilex，见[字体与图标](/zh/typography)。
- 想让文字显示中文界面、翻译 MyGo 自带的菜单，见[多语言](/zh/i18n)。

## 运行画廊

仓库里的 `cmd/gallery` 是一个展示全部组件的原生应用，文档站每页的预览就是它。

```sh
go run ./cmd/gallery
```

- 只看某一页，并且去掉画廊外壳（文档站预览的样子）：

  ```sh
  go run ./cmd/gallery -page Button -embed
  ```

  页名就是页面在画廊里的名字，例如 `Button`、`Settings form`、`Virtual list`。也可以用环境变量 `KIT_PAGE=Button` 指定起始页。

- 用暗色启动：`KIT_DARK=1 go run ./cmd/gallery`。
- 用窄窗口启动（宽 700）：`KIT_WIDTH=700 go run ./cmd/gallery`。窗口宽度小于 900 时，[Workspace](/zh/components/workspace) 的侧栏会改成抽屉。

## 在浏览器里运行

MyGo 的原生界面可以编译成 WebAssembly 跑在浏览器里，窗口画在页面的画布上。本库的组件不用改。

```sh
GOOS=js GOARCH=wasm go build -o app.wasm .
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
```

再写一个页面加载 `wasm_exec.js` 和 `app.wasm`，用 HTTP 服务（`.wasm` 要以 `application/wasm` 提供）打开。完整的页面写法和限制见 MyGo 的 [Native UI in browsers](https://github.com/egoist/mygo/blob/main/docs/ui/browsers.md)。文档站的预览就是这样跑的，页面在仓库的 `website/public/preview/index.html`。

有两点和桌面不同：

- 浏览器只按族名区分字体，所以浏览器构建只注册常规体，粗体由浏览器合成。
- 文字由浏览器排版，系统菜单显示成网页菜单。

## 检查

```sh
go test ./...
go test -race ./...
go vet ./...
```

加 `-tags e2e` 的构建才会导入 MyGo 的 automation；正式构建不开放自动化 socket。窗口、系统输入法和 GPU 的验收不能靠离屏测试代替。
