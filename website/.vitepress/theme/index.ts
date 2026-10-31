import DefaultTheme from "vitepress/theme"
import type { Theme } from "vitepress"
import Preview from "./Preview.vue"
import "./custom.css"

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component("Preview", Preview)
  },
} satisfies Theme
