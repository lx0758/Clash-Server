import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWebSocket } from '@/composables/useWebSocket'

export interface TrafficDataPoint {
  up: number
  down: number
}

export type TimeRange = '1m' | '5m' | '10m' | '30m'

const TIME_RANGE_POINTS: Record<TimeRange, number> = {
  '1m': 60,
  '5m': 300,
  '10m': 600,
  '30m': 1800,
}

const dataPoints = ref<TrafficDataPoint[]>([])
const timeRange = ref<TimeRange>('10m')
const isAreaMode = ref(false)
const initialized = ref(false)
const refCount = ref(0)

let watcherId: number | null = null
let animationFrameId: number | null = null

const maxPoints = computed(() => TIME_RANGE_POINTS[timeRange.value])

const startWatcher = (traffic: { value: { up: number; down: number } }) => {
  if (watcherId) return
  watcherId = window.setInterval(() => {
    addDataPoint(traffic.value.up, traffic.value.down)
  }, 1000)
}

const stopWatcher = () => {
  if (watcherId) {
    clearInterval(watcherId)
    watcherId = null
  }
}

const addDataPoint = (up: number, down: number) => {
  const safeUp = up ?? 0
  const safeDown = down ?? 0
  dataPoints.value.push({ up: safeUp, down: safeDown })
  const max = maxPoints.value
  while (dataPoints.value.length > max) {
    dataPoints.value.shift()
  }
}

const clearData = () => {
  dataPoints.value = []
}

const toggleStyle = () => {
  isAreaMode.value = !isAreaMode.value
}

const setTimeRange = (range: TimeRange) => {
  timeRange.value = range
  const max = TIME_RANGE_POINTS[range]
  while (dataPoints.value.length > max) {
    dataPoints.value.shift()
  }
}

const upData = computed(() => dataPoints.value.map(p => p.up ?? 0))
const downData = computed(() => dataPoints.value.map(p => p.down ?? 0))

const maxUp = computed(() => {
  const values = upData.value.filter(v => isFinite(v))
  return values.length > 0 ? Math.max(...values, 1) : 1
})
const maxDown = computed(() => {
  const values = downData.value.filter(v => isFinite(v))
  return values.length > 0 ? Math.max(...values, 1) : 1
})

export function useTrafficGraph() {
  const { traffic, subscribe } = useWebSocket()

  onMounted(() => {
    refCount.value++
    if (!initialized.value) {
      initialized.value = true
      subscribe(['traffic'])
    }
    startWatcher(traffic)
  })

  onUnmounted(() => {
    refCount.value--
    if (refCount.value <= 0) {
      stopWatcher()
      if (animationFrameId) {
        cancelAnimationFrame(animationFrameId)
        animationFrameId = null
      }
    }
  })

  return {
    dataPoints,
    upData,
    downData,
    maxUp,
    maxDown,
    isAreaMode,
    timeRange,
    maxPoints,
    toggleStyle,
    clearData,
    setTimeRange,
  }
}
