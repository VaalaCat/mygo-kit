package kit

import (
	"embed"
	"runtime"
	"sync"

	"github.com/egoist/mygo/ui"
)

//go:embed fonts/*.ttf
var fontFiles embed.FS

// 界面与代码字体族名，与 Zed 默认的 .ZedSans、.ZedMono 一致。
const (
	FontSans = "IBM Plex Sans"
	FontMono = "Lilex"
)

var registerFontsOnce = sync.OnceValue(registerFonts)

// RegisterFonts 注册库内置的 IBM Plex Sans 与 Lilex；重复调用只注册一次。
func RegisterFonts() error {
	return registerFontsOnce()
}

// registerFonts 按字族逐个注册内置字体文件。
func registerFonts() error {
	// 1. 同一字族的各字重注册到同一个族名下。
	files := map[string][]string{
		FontSans: {"IBMPlexSans-Regular", "IBMPlexSans-SemiBold", "IBMPlexSans-Italic", "IBMPlexSans-SemiBoldItalic"},
		FontMono: {"Lilex-Regular", "Lilex-Bold", "Lilex-Italic"},
	}
	for family, names := range files {
		// 2. 浏览器只按族名区分字体，同族的粗体、斜体会盖掉常规体；浏览器构建只注册常规体，粗体由浏览器合成。
		if runtime.GOOS == "js" {
			names = names[:1]
		}
		for _, name := range names {
			data, err := fontFiles.ReadFile("fonts/" + name + ".ttf")
			if err != nil {
				return err
			}
			if err = ui.RegisterFont(data, family); err != nil {
				return err
			}
		}
	}
	return nil
}
