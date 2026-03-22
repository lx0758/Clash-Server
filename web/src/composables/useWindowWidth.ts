import { ref, onMounted, onUnmounted } from 'vue'

export function useWindowWidth() {
  const width = ref(window.innerWidth)

  const updateWidth = () => {
    width.value = window.innerWidth
  }

  onMounted(() => {
    window.addEventListener('resize', updateWidth)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', updateWidth)
  })

  const calculateColumns = (w: number): number => {
    if (w > 1400) return 4
    if (w > 1000) return 3
    if (w > 700) return 2
    return 1
  }

  const columns = computed(() => calculateColumns(width.value))

  return {
    width,
    columns,
    calculateColumns,
  }
}

import { computed } from 'vue'