package main

import (
	"flag"
	"log"
	"os"
	"sort"
	"strconv"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo"
	_ "github.com/egoist/mygo/locales/zh-Hans"
	"github.com/egoist/mygo/ui"
)

// page 是画廊里的一页，每个组件在自己的文件里注册一页。
type page struct {
	name  string
	order int
	// full 为真时页面占据整个窗口，自己负责顶栏和返回按钮。
	full bool
	view func(c *ui.Context)
}

type gallery struct {
	dark    bool
	current string
	nav     float32
	// embed 为真时只画当前页，不画画廊外壳，配色跟随系统或嵌入页面的暗色设置；文档站的预览用这个模式。
	embed bool
}

var (
	pages []page
	app   = &gallery{nav: 188}
)

// register 在包初始化时登记一页；order 决定侧栏顺序。
func register(p page) { pages = append(pages, p) }

// showPage 切换到指定页面。
func showPage(name string) { app.current = name }

// main 创建无业务依赖的原生组件画廊。
//
//	go run ./cmd/gallery                       # 完整画廊
//	go run ./cmd/gallery -page Button -embed   # 只显示 Button 页，文档站预览的样子
func main() {
	// 1. 页面按 order 排序；起始页和嵌入模式来自命令行参数，浏览器里由页面的 go.argv 传入。
	sort.SliceStable(pages, func(i, j int) bool { return pages[i].order < pages[j].order })
	pageName := flag.String("page", os.Getenv("KIT_PAGE"), "起始页名称，例如 Button")
	flag.BoolVar(&app.embed, "embed", false, "只显示起始页，不显示画廊外壳")
	flag.Parse()
	if err := kit.RegisterFonts(); err != nil {
		panic(err)
	}
	app.dark = os.Getenv("KIT_DARK") == "1"
	app.current = pages[0].name
	if *pageName != "" {
		app.current = *pageName
	}
	// 2. 原生循环和系统窗口仍由应用自己持有。
	mygo.App.SetName("mygo-kit-gallery")
	mygo.App.SetLanguages("en", "zh-Hans")
	mygo.App.WhenReady(func() {
		width := 1180
		if value, err := strconv.Atoi(os.Getenv("KIT_WIDTH")); err == nil && value > 0 {
			width = value
		}
		mygo.NewWindow(mygo.WindowOptions{Title: "MyGo Kit — Native workspace", Width: width, Height: 760, Content: ui.View(view)})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

// view 应用主题并绘制当前页面；全窗口页面不带画廊外壳。
func view(c *ui.Context) {
	// 1. 每帧应用窗口主题；嵌入模式跟随框架报告的暗色设置。
	if app.embed {
		app.dark = c.Theme().Dark
	}
	theme := kit.OneLight()
	if app.dark {
		theme = kit.OneDark()
	}
	kit.ApplyTheme(c, theme)
	current := pages[0]
	for _, p := range pages {
		if p.name == app.current {
			current = p
		}
	}
	defer kit.ToastHost(c)
	if current.full {
		current.view(c)
		return
	}
	if app.embed {
		ui.Column(c).Fill().Gap(0).Background(theme.Colors.Editor).Children(func() { current.view(c) })
		return
	}
	// 2. 外壳：顶栏、页面导航、页面内容和状态栏。
	ui.Column(c).Fill().Gap(0).Children(func() {
		kit.Header(c, "mygo-kit  /  Native workspace", func() {
			if kit.Button(c.Key("toggle-theme"), "◐ Theme").Clicked() {
				app.dark = !app.dark
			}
			if kit.Button(c.Key("toggle-language"), "EN / 中文").Clicked() {
				if c.Locale() == "zh-Hans" {
					mygo.App.SetLocale("en")
				} else {
					mygo.App.SetLocale("zh-Hans")
				}
			}
			if kit.Button(c.Key("close-gallery"), "Close").Clicked() {
				mygo.App.Quit()
			}
		})
		kit.Split(c, &app.nav, func() {
			ui.Column(c).Fill().Background(theme.Colors.Panel).Children(func() {
				kit.Header(c, "COMPONENTS", nil)
				ui.Scroll(c).Grow(1).Children(func() {
					for _, p := range pages {
						if kit.ListRow(c.Key(p.name), kit.ListRowOpt{Title: p.name, Selected: p.name == current.name}).Clicked() {
							showPage(p.name)
						}
					}
				})
				theme.TextRole(ui.Text(c, "Go · MyGo 0.7.3"), kit.TextCaption).Padding(12).TextColor(theme.Colors.Muted)
			})
		}, func() {
			ui.Column(c).Fill().Children(func() {
				kit.Header(c, current.name, nil)
				current.view(c)
			})
		}).Grow(1)
		kit.StatusBar(c, func() {
			ui.Text(c, "● Ready").TextColor(theme.Colors.Success)
			ui.Spacer(c)
			ui.Text(c, "Native UI  ·  "+current.name)
		})
	})
}

// pageBody 绘制页面常用的可滚动内容区。
func pageBody(c *ui.Context, content func()) {
	ui.Scroll(c).Grow(1).Children(func() {
		ui.Column(c).Padding(24).Gap(16).Children(content)
	})
}

// sectionTitle 绘制页面里的小节标题。
func sectionTitle(c *ui.Context, title string) {
	kit.GetTheme(c).TextRole(ui.Text(c, title), kit.TextUI)
}
