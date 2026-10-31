package kit

import (
	"embed"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"io/fs"
)

//go:embed locales/*.json
var localeFiles embed.FS

type localeKey struct{}

// text 返回 MyGo 当前语言下的库文案；应用目录可用 mygokit.<key> 覆盖。
func text(c *ui.Context, key string) string {
	// 1. 文案目录归窗口保存，不注册应用语言或更改全局 locale。
	catalog := ui.Local(c.Root(), localeKey{}, func() *i18n.Catalog {
		files, err := fs.Sub(localeFiles, "locales")
		if err != nil {
			panic(err)
		}
		catalog, err := i18n.NewLibraryCatalog(files, "en", "mygokit")
		if err != nil {
			panic(err)
		}
		return catalog
	})
	return (*catalog).Format(c.Locale(), key, nil)
}
