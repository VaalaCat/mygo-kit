import fs from "node:fs"
import path from "node:path"
import { defineConfig, type DefaultTheme } from "vitepress"

const docsDir = path.resolve(__dirname, "../../docs")

// componentPages 读取一种语言的 components 目录下每页的 title，按名称排序，作为侧栏「组件」一节。
// prefix 是这种语言的路径前缀：英文是 ""，中文是 "/zh"。
function componentPages(prefix: string): DefaultTheme.SidebarItem[] {
  const dir = path.join(docsDir, prefix, "components")
  return fs
    .readdirSync(dir)
    .filter((file) => file.endsWith(".md"))
    .map((file) => {
      const source = fs.readFileSync(path.join(dir, file), "utf8")
      const title = /^title:\s*(.+)$/m.exec(source)?.[1].trim() ?? file.replace(/\.md$/, "")
      return { text: title, link: `${prefix}/components/${file.replace(/\.md$/, "")}` }
    })
    .sort((a, b) => a.text.localeCompare(b.text))
}

export default defineConfig({
  title: "mygo-kit",
  // GitHub Pages 发布在 /mygo-kit/ 下，由 DOCS_BASE 指定；本地开发用根路径。
  base: process.env.DOCS_BASE ?? "/",
  srcDir: docsDir,
  srcExclude: ["superpowers/**"],
  cleanUrls: true,
  lastUpdated: false,
  // 正文在 website 之外的 docs/，vue 统一从 website/node_modules 解析。
  vite: {
    publicDir: path.resolve(__dirname, "../public"),
    resolve: { alias: [{ find: /^vue($|\/)/, replacement: path.resolve(__dirname, "../node_modules/vue") + "$1" }] },
  },
  themeConfig: {
    socialLinks: [{ icon: "github", link: "https://github.com/VaalaCat/mygo-kit" }],
    search: {
      provider: "local",
      options: {
        locales: {
          zh: { translations: { button: { buttonText: "搜索", buttonAriaLabel: "搜索" }, modal: { noResultsText: "没有结果", resetButtonTitle: "清除", footer: { selectText: "选择", navigateText: "切换", closeText: "关闭" } } } },
        },
      },
    },
  },
  locales: {
    root: {
      label: "English",
      lang: "en",
      description: "Native Go components for MyGo: a Zed-style workspace, controls and readers",
      themeConfig: {
        nav: [
          { text: "Docs", link: "/" },
          { text: "Components", link: "/components/button" },
          { text: "Gallery", link: "/preview/?page=Overview&full", target: "_blank" },
        ],
        sidebar: [
          {
            text: "Getting Started",
            items: [
              { text: "Introduction", link: "/" },
              { text: "Installation", link: "/installation" },
              { text: "Theming", link: "/theming" },
              { text: "Dark Mode", link: "/dark-mode" },
              { text: "Typography & Icons", link: "/typography" },
              { text: "Localization", link: "/i18n" },
            ],
          },
          { text: "Components", items: componentPages("") },
        ],
        outline: { label: "On this page", level: [2, 3] },
      },
    },
    zh: {
      label: "简体中文",
      lang: "zh-CN",
      link: "/zh/",
      description: "MyGo 原生 Go 组件库：Zed 风格的工作区、控件和阅读组件",
      themeConfig: {
        nav: [
          { text: "文档", link: "/zh/" },
          { text: "组件", link: "/zh/components/button" },
          { text: "画廊", link: "/preview/?page=Overview&full&lang=zh-Hans", target: "_blank" },
        ],
        sidebar: [
          {
            text: "入门",
            items: [
              { text: "介绍", link: "/zh/" },
              { text: "安装", link: "/zh/installation" },
              { text: "主题", link: "/zh/theming" },
              { text: "暗色模式", link: "/zh/dark-mode" },
              { text: "字体与图标", link: "/zh/typography" },
              { text: "多语言", link: "/zh/i18n" },
            ],
          },
          { text: "组件", items: componentPages("/zh") },
        ],
        outline: { label: "本页内容", level: [2, 3] },
        docFooter: { prev: "上一页", next: "下一页" },
        darkModeSwitchLabel: "外观",
        lightModeSwitchTitle: "切换到浅色模式",
        darkModeSwitchTitle: "切换到深色模式",
        sidebarMenuLabel: "菜单",
        returnToTopLabel: "回到顶部",
        langMenuLabel: "切换语言",
      },
    },
  },
})
