import { onMounted, onUnmounted } from 'vue'

/**
 * 维护 --app-height CSS 变量：动态视口高度的 JS 兜底方案。
 *
 * 三级降级链：100vh（极老内核）→ 100dvh（现代内核）→ var(--app-height)（本模块）。
 * dvh 可用的浏览器由纯 CSS 覆盖；不支持 dvh 的老内核（如部分 X5/XWeb）
 * 通过本模块感知地址栏收起/展开，避免 100vh 遮挡底部内容。
 *
 * 在 App.vue setup 阶段同步执行首帧赋值，使 var() 在首次渲染前即已就绪。
 */
export function useAppHeight(): void {
  const update = () => {
    const height = window.visualViewport?.height ?? window.innerHeight
    document.documentElement.style.setProperty('--app-height', `${height}px`)
  }

  // setup 阶段同步首帧赋值（非 mounted），避免 dvh 不可用内核的首帧回退窗口
  update()

  onMounted(() => {
    window.addEventListener('resize', update)
    window.addEventListener('orientationchange', update)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', update)
    window.removeEventListener('orientationchange', update)
  })
}
