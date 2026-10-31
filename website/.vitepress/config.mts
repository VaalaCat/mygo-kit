import fs from "node:fs"
import path from "node:path"
import { defineConfig } from "vitepress"

const docsDir = path.resolve(__dirname, "../../docs")

// componentPages 读取 docs/components 下每个组件页的 title，按名称排序，作为侧栏「组件」一节。
function componentPages() {
  const dir = path.join(docsDir, "components")
  return fs
    .readdirSync(dir)
    .filter((file) => file.endsWith(".md"))
    .map((file) => {
      const source = fs.readFileSync(path.join(dir, file), "utf8")
      const title = /^title:\s*(.+)$/m.exec(source)?.[1].trim() ?? file.replace(/\.md$/, "")
      return { text: title, link: `/components/${file.replace(/\.md$/, "")}` }
    })
    .sort((a, b) => a.text.localeCompare(b.text))
}

export default defineConfig({
  title: "mygo-kit",
  description: "MyGo 原生 Go 组件库：Zed 风格的工作区、控件和阅读组件",
  lang: "zh-CN",
  // GitHub Pages 发布在 /mygo-kit/ 下，由 DOCS_BASE 指定；本地开发用根路径。
  base: process.env.DOCS_BASE ?? "/",
  srcDir: docsDir,
  srcExclude: ["superpowers/**"],
  cleanUrls: true,
  lastUpdated: false,
  // 正文在 website 之外的 docs/，依赖统一从 website/node_modules 解析。
  vite: { publicDir: path.resolve(__dirname, "../public"), resolve: { alias: [{ find: /^vue($|\/)/, replacement: path.resolve(__dirname, "../node_modules/vue") + "$1" }] } },
  themeConfig: {
    nav: [
      { text: "文档", link: "/" },
      { text: "组件", link: "/components/button" },
      { text: "画廊", link: "/preview/?page=Overview&full", target: "_blank" },
    ],
    sidebar: [
      {
        text: "入门",
        items: [
          { text: "介绍", link: "/" },
          { text: "安装", link: "/installation" },
          { text: "主题", link: "/theming" },
          { text: "暗色模式", link: "/dark-mode" },
          { text: "字体与图标", link: "/typography" },
          { text: "多语言", link: "/i18n" },
        ],
      },
      { text: "组件", items: componentPages() },
    ],
    socialLinks: [{ icon: "github", link: "https://github.com/VaalaCat/mygo-kit" }],
    search: { provider: "local" },
    outline: { label: "本页内容", level: [2, 3] },
    docFooter: { prev: "上一页", next: "下一页" },
    darkModeSwitchLabel: "外观",
    sidebarMenuLabel: "菜单",
    returnToTopLabel: "回到顶部",
  },
})
