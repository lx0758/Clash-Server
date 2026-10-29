import { computed, ref, watch, onUnmounted } from 'vue'
import { usePendingStore } from '@/stores/pending'

/**
 * 延迟显示进行中状态的 composable。
 *
 * 立即禁用是必要的(防止重复提交), 但视觉上延迟显示可避免快操作闪烁。
 * 默认阈值 300ms, 可按操作覆盖。
 */
export function useDelayedPending(key: string | (() => string), threshold = 300) {
  const pendingStore = usePendingStore()
  const visible = ref(false)
  let timer: ReturnType<typeof setTimeout> | null = null

  const resolveKey = () => (typeof key === 'function' ? key() : key)
  const pending = computed(() => pendingStore.isPending(resolveKey()))

  watch(pending, (isPending) => {
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
    if (isPending) {
      // 超过阈值仍未完成才显示加载指示
      timer = setTimeout(() => {
        visible.value = true
        timer = null
      }, threshold)
    } else {
      visible.value = false
    }
  })

  onUnmounted(() => {
    if (timer) clearTimeout(timer)
  })

  return { pending, visible }
}
