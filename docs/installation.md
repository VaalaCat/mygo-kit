---
title: Installation
description: Requirements, go get, a runnable minimal example, and how to run the gallery.
---

# Installation

## Requirements

- Go 1.27.3.
- MyGo v0.7.3, downloaded automatically when you run `go get` for this library.
- Linux: GTK 3 is required for native windows. Offscreen text rendering still needs Pango and Cairo, which are usually installed with GTK.
- macOS: uses the system’s native frameworks. No extra installation is needed.

The library depends only on MyGo. It has no dependencies on application code.

## Installation

```sh
go get github.com/VaalaCat/mygo-kit
```

Import it as follows. The package name is `kit`:

```go
import kit "github.com/VaalaCat/mygo-kit"
```

## Minimal example

Create a directory, add the following `main.go`, run `go mod init` and `go mod tidy`, then run `go run .`:

```go
package main

import (
	"fmt"
	"log"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// The application owns clicks. Components do not store application state.
var clicks int

func view(c *ui.Context) {
	// Apply the theme before drawing the UI each frame.
	kit.ApplyTheme(c, kit.OneLight())
	t := kit.GetTheme(c)
	ui.Column(c).Fill().Gap(0).Background(t.Colors.Editor).Children(func() {
		kit.Header(c, "Hello mygo-kit", nil)
		ui.Column(c).Padding(t.Space(6)).Gap(t.Space(3)).Children(func() {
			ui.Text(c, fmt.Sprintf("Clicked %d times", clicks))
			if kit.PrimaryButton(c, "Click me").AlignSelf(ui.Start).Clicked() {
				clicks++
			}
		})
	})
}

func main() {
	// Register fonts only once.
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

A few notes:

- Your application provides windows, the event loop, and `mygo.App.Run`. The library does not create them.
- `ui.View(view)` calls `view` once per frame, so call `ApplyTheme` each frame too. See [Theming](/theming).
- `kit.RegisterFonts()` registers the built-in IBM Plex Sans and Lilex fonts. See [Typography & Icons](/typography).
- To display the interface in Chinese and translate MyGo’s built-in menus, see [Localization](/i18n).

## Run the gallery

The repository’s `cmd/gallery` is a native application that showcases every component. It also powers the previews on each docs page.

```sh
go run ./cmd/gallery
```

- Show a single page without the gallery shell, as in the docs previews:

  ```sh
  go run ./cmd/gallery -page Button -embed
  ```

  Use the page name from the gallery, such as `Button`, `Settings form`, or `Virtual list`. You can also set the starting page with the environment variable `KIT_PAGE=Button`.

- Start in dark mode: `KIT_DARK=1 go run ./cmd/gallery`.
- Start with a narrow window, 700 pixels wide: `KIT_WIDTH=700 go run ./cmd/gallery`. Below a window width of 900, the [Workspace](/components/workspace) sidebar becomes a drawer.

## Run in a browser

MyGo’s native UI can compile to WebAssembly and run in a browser, with windows drawn on the page’s canvas. The library’s components work without changes.

```sh
GOOS=js GOARCH=wasm go build -o app.wasm .
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
```

Create a page that loads `wasm_exec.js` and `app.wasm`, and serve it over HTTP. Serve `.wasm` files as `application/wasm`. For the complete page setup and limitations, see MyGo’s [Native UI in browsers](https://github.com/egoist/mygo/blob/main/docs/ui/browsers.md). The docs previews use this approach. The page is at `website/public/preview/index.html` in the repository.

Two things differ from desktop builds:

- Browsers distinguish fonts only by family name, so browser builds register only regular faces. The browser synthesizes bold text.
- The browser lays out text, and system menus appear as web menus.

## Checks

```sh
go test ./...
go test -race ./...
go vet ./...
```

Only builds with `-tags e2e` import MyGo’s automation package. Production builds do not expose an automation socket. Offscreen tests cannot replace checks of windows, system input methods, and the GPU.
