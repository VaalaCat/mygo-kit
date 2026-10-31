<!--
  Preview 是组件页里的「预览 | 代码」块：预览页签用 iframe 运行画廊的一页（WebAssembly），
  代码页签显示默认插槽里的源码，通常是 <<< 引入的画廊源文件。

  <Preview page="Button" :height="360">

  <<< @/../cmd/gallery/button.go

  </Preview>

  width 让画廊按这个宽度布局，再整体缩小到正文宽度，例如 Workspace 要 900 以上才会并排显示右 Dock：
  <Preview page="Workspace" :height="560" :width="1180">
-->
<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue"
import { withBase } from "vitepress"

const props = withDefaults(defineProps<{ page: string; height?: number; width?: number }>(), { height: 360 })
const tab = ref<"preview" | "code">("preview")

// 1. 给了 width 且正文比它窄时，iframe 按 width 布局，再按正文宽度等比缩小；height 是缩小后的高度。
const frame = ref<HTMLElement>()
const frameWidth = ref(0)
let resizeObserver: ResizeObserver | undefined
onMounted(() => {
  resizeObserver = new ResizeObserver(([entry]) => (frameWidth.value = entry.contentRect.width))
  if (frame.value) resizeObserver.observe(frame.value)
})
onBeforeUnmount(() => resizeObserver?.disconnect())
const scale = computed(() => (props.width && frameWidth.value && frameWidth.value < props.width ? frameWidth.value / props.width : 1))
const iframeStyle = computed(() =>
  scale.value === 1
    ? { height: `${props.height}px` }
    : { width: `${props.width}px`, height: `${props.height / scale.value}px`, transform: `scale(${scale.value})`, transformOrigin: "0 0" },
)
const src = withBase(`/preview/?page=${encodeURIComponent(props.page)}`)
const fullSrc = withBase(`/preview/?page=${encodeURIComponent(props.page)}&full`)
</script>

<template>
  <div class="kit-preview">
    <div class="kit-preview-tabs" role="tablist">
      <button role="tab" :aria-selected="tab === 'preview'" :class="{ active: tab === 'preview' }" @click="tab = 'preview'">预览</button>
      <button role="tab" :aria-selected="tab === 'code'" :class="{ active: tab === 'code' }" @click="tab = 'code'">代码</button>
      <a class="kit-preview-open" :href="fullSrc" target="_blank" rel="noopener">在画廊中打开 ↗</a>
    </div>
    <div v-show="tab === 'preview'" ref="frame" class="kit-preview-frame" :style="{ height: `${height}px` }">
      <iframe :src="src" :style="iframeStyle" loading="lazy" :title="`${page} 预览`" />
    </div>
    <div v-show="tab === 'code'" class="kit-preview-code">
      <slot />
    </div>
  </div>
</template>
