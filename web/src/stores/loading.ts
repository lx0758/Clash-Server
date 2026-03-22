import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const count = ref(0)
  const message = ref('')

  function show(msg = '加载中...') {
    count.value++
    message.value = msg
  }

  function hide() {
    if (count.value > 0) {
      count.value--
    }
    if (count.value === 0) {
      message.value = ''
    }
  }

  const isLoading = () => count.value > 0

  return { count, message, show, hide, isLoading }
})
