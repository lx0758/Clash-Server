import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 对象级进行中状态。
 *
 * 用 key 标识一个具体操作, 使不同操作、以及同一列表中不同对象的
 * 执行状态互不影响。
 *
 * key 命名约定: `<动作>:<对象标识>`
 *   delay:GROUP_A      测某代理组的延迟
 *   delay-node:HK-01   测某个节点的延迟
 *   refresh:12         刷新 id 为 12 的订阅
 *   activate:12        激活 id 为 12 的订阅
 *   select:HK-01       切换到某个节点
 *   mode:global        切换运行模式
 */
export const usePendingStore = defineStore('pending', () => {
  const keys = ref<Set<string>>(new Set())

  const isPending = (key: string) => keys.value.has(key)

  const start = (key: string) => {
    // Set 的插入会触发响应式更新
    keys.value = new Set(keys.value).add(key)
  }

  const finish = (key: string) => {
    if (!keys.value.has(key)) return
    const next = new Set(keys.value)
    next.delete(key)
    keys.value = next
  }

  /**
   * 包装一个异步操作, 自动管理进行中状态。
   * 使用 finally 保证异常路径下状态也会被清理, 避免永久卡在加载态。
   */
  const withPending = async <T>(key: string, fn: () => Promise<T>): Promise<T> => {
    start(key)
    try {
      return await fn()
    } finally {
      finish(key)
    }
  }

  /** 是否存在任意进行中的操作 */
  const hasAny = () => keys.value.size > 0

  /** 清除所有进行中状态 */
  const clear = () => {
    keys.value = new Set()
  }

  return { keys, isPending, start, finish, withPending, hasAny, clear }
})
