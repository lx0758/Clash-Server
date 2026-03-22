<template>
  <div class="traffic-graph-wrapper">
    <div
      ref="containerRef"
      class="traffic-graph"
      :class="{ 'compact': compact }"
      @click="toggleStyle"
    >
      <canvas ref="canvasRef"></canvas>
    </div>
    <div class="graph-footer">
      <el-radio-group v-model="selectedTimeRange" size="small" @change="handleTimeRangeChange">
        <el-radio-button value="1m">1分钟</el-radio-button>
        <el-radio-button value="5m">5分钟</el-radio-button>
        <el-radio-button value="10m">10分钟</el-radio-button>
        <el-radio-button value="30m">30分钟</el-radio-button>
      </el-radio-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue'
import { useTrafficGraph, type TimeRange } from '@/composables/useTrafficGraph'
import { formatBytes } from '@/utils/format'

interface Props {
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compact: false,
})

const { dataPoints, upData, downData, maxUp, maxDown, isAreaMode, timeRange, maxPoints, toggleStyle, setTimeRange } = useTrafficGraph()

const selectedTimeRange = ref<TimeRange>(timeRange.value)

const canvasRef = ref<HTMLCanvasElement | null>(null)
const containerRef = ref<HTMLDivElement | null>(null)

let animationFrameId: number | null = null
let resizeObserver: ResizeObserver | null = null

const currentUpSpeed = computed(() => {
  if (upData.value.length === 0) return 0
  return upData.value[upData.value.length - 1] ?? 0
})

const currentDownSpeed = computed(() => {
  if (downData.value.length === 0) return 0
  return downData.value[downData.value.length - 1] ?? 0
})

const formatSpeed = (bytesPerSec: number): string => {
  if (!bytesPerSec || !isFinite(bytesPerSec)) return '0 B/s'
  return formatBytes(bytesPerSec) + '/s'
}

const formatAxisLabel = (bytes: number): string => {
  if (bytes <= 0 || !isFinite(bytes)) return '0'
  const units = ['', 'K', 'M', 'G', 'T']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const safeI = Math.max(0, Math.min(i, units.length - 1))
  const value = bytes / Math.pow(k, safeI)
  const formatted = value < 10 ? value.toFixed(1) : Math.round(value).toString()
  return formatted + units[safeI]
}

const handleTimeRangeChange = (range: TimeRange) => {
  setTimeRange(range)
}

const draw = () => {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const container = containerRef.value
  if (!container) return

  const width = container.clientWidth
  const height = container.clientHeight

  canvas.width = width * window.devicePixelRatio
  canvas.height = height * window.devicePixelRatio

  ctx.scale(window.devicePixelRatio, window.devicePixelRatio)
  ctx.clearRect(0, 0, width, height)

  if (dataPoints.value.length < 2) {
    drawEmptyState(ctx, width, height)
    return
  }

  const padding = props.compact
    ? { top: 20, right: 20, bottom: 55, left: 40 }
    : { top: 20, right: 20, bottom: 60, left: 45 }
  const graphWidth = width - padding.left - padding.right
  const graphHeight = height - padding.top - padding.bottom

  const rawMaxUp = maxUp.value ?? 0
  const rawMaxDown = maxDown.value ?? 0
  const maxVal = Math.max(rawMaxUp, rawMaxDown, 1024)
  const minVal = 0
  const currentMaxPoints = maxPoints.value

  const getX = (index: number) => {
    const dataLength = dataPoints.value.length
    const startIndex = Math.max(0, currentMaxPoints - dataLength)
    const actualIndex = startIndex + index
    return padding.left + (actualIndex / currentMaxPoints) * graphWidth
  }

  const getY = (value: number) => {
    const safeMaxVal = maxVal || 1
    const safeValue = Math.max(0, value)
    return height - padding.bottom - (safeValue / safeMaxVal) * graphHeight
  }

  ctx.save()
  ctx.beginPath()
  ctx.rect(padding.left - 2, padding.top - 2, graphWidth + 4, graphHeight + 4)
  ctx.clip()

  drawCurve(ctx, upData.value, '#f59e0b', isAreaMode.value, getX, getY, height - padding.bottom)
  drawCurve(ctx, downData.value, '#3b82f6', isAreaMode.value, getX, getY, height - padding.bottom)

  ctx.restore()

  drawGrid(ctx, width, height, padding, maxVal, minVal, currentMaxPoints)
  drawLegend(ctx, width, height, padding)
}

const drawEmptyState = (ctx: CanvasRenderingContext2D, width: number, height: number) => {
  ctx.fillStyle = '#94a3b8'
  ctx.font = '14px sans-serif'
  ctx.textAlign = 'center'
  ctx.fillText('等待数据...', width / 2, height / 2)
}

const drawGrid = (
  ctx: CanvasRenderingContext2D,
  width: number,
  height: number,
  padding: { top: number; right: number; bottom: number; left: number },
  maxVal: number,
  minVal: number,
  totalSeconds: number
) => {
  const graphHeight = height - padding.top - padding.bottom
  const graphWidth = width - padding.left - padding.right

  ctx.strokeStyle = '#e2e8f0'
  ctx.lineWidth = 1
  ctx.setLineDash([4, 4])

  for (let i = 0; i <= 4; i++) {
    const y = height - padding.bottom - (i / 4) * graphHeight
    ctx.beginPath()
    ctx.moveTo(padding.left, y)
    ctx.lineTo(width - padding.right, y)
    ctx.stroke()

    const value = (maxVal - minVal) * (i / 4) + minVal
    const safeValue = isFinite(value) && value >= 0 ? value : 0
    const label = formatAxisLabel(safeValue)
    ctx.fillStyle = '#94a3b8'
    ctx.font = '10px sans-serif'
    ctx.textAlign = 'right'
    ctx.fillText(label, padding.left - 5, y + 3)
  }

  const tickCount = 5
  const tickInterval = totalSeconds / tickCount

  for (let i = 0; i <= tickCount; i++) {
    const x = padding.left + (i / tickCount) * graphWidth
    ctx.beginPath()
    ctx.moveTo(x, padding.top)
    ctx.lineTo(x, height - padding.bottom)
    ctx.stroke()

    const secondsAgo = Math.round((tickCount - i) * tickInterval)
    let label: string
    if (secondsAgo === 0) {
      label = '现在'
    } else if (totalSeconds <= 60) {
      label = `-${secondsAgo}s`
    } else {
      const minutes = secondsAgo / 60
      label = `-${minutes}m`
    }
    ctx.fillStyle = '#94a3b8'
    ctx.font = '10px sans-serif'
    ctx.textAlign = 'center'
    ctx.fillText(label, x, height - padding.bottom + 15)
  }

  ctx.setLineDash([])
}

const drawLegend = (
  ctx: CanvasRenderingContext2D,
  width: number,
  height: number,
  _padding: { top: number; right: number; bottom: number; left: number }
) => {
  const legendY = height - 18
  const centerX = width / 2
  const legendSpacing = 100

  ctx.font = '11px sans-serif'

  ctx.fillStyle = '#f59e0b'
  ctx.beginPath()
  ctx.arc(centerX - legendSpacing, legendY, 5, 0, Math.PI * 2)
  ctx.fill()
  ctx.fillStyle = '#374151'
  ctx.textAlign = 'left'
  ctx.fillText(`上行 ${formatSpeed(currentUpSpeed.value)}`, centerX - legendSpacing + 10, legendY + 4)

  ctx.fillStyle = '#3b82f6'
  ctx.beginPath()
  ctx.arc(centerX + legendSpacing - 80, legendY, 5, 0, Math.PI * 2)
  ctx.fill()
  ctx.fillStyle = '#374151'
  ctx.fillText(`下行 ${formatSpeed(currentDownSpeed.value)}`, centerX + legendSpacing - 70, legendY + 4)
}

const drawCurve = (
  ctx: CanvasRenderingContext2D,
  data: number[],
  color: string,
  isArea: boolean,
  getX: (index: number) => number,
  getY: (value: number) => number,
  baselineY: number
) => {
  if (data.length < 2) return

  const points: { x: number; y: number }[] = []
  for (let i = 0; i < data.length; i++) {
    points.push({ x: getX(i), y: getY(data[i] ?? 0) })
  }

  const firstPoint = points[0]
  if (!firstPoint) return

  ctx.beginPath()
  ctx.moveTo(firstPoint.x, firstPoint.y)

  for (let i = 0; i < points.length - 1; i++) {
    const p0 = points[Math.max(0, i - 1)]
    const p1 = points[i]
    const p2 = points[i + 1]
    const p3 = points[Math.min(points.length - 1, i + 2)]

    if (!p0 || !p1 || !p2 || !p3) continue

    const cp1x = p1.x + (p2.x - p0.x) / 6
    const cp1y = p1.y + (p2.y - p0.y) / 6
    const cp2x = p2.x - (p3.x - p1.x) / 6
    const cp2y = p2.y - (p3.y - p1.y) / 6

    ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, p2.x, p2.y)
  }

  ctx.strokeStyle = color
  ctx.lineWidth = 2
  ctx.stroke()

  if (isArea && points.length > 0) {
    const lastPoint = points[points.length - 1]
    if (lastPoint) {
      ctx.lineTo(lastPoint.x, baselineY)
      ctx.lineTo(firstPoint.x, baselineY)
      ctx.closePath()
      ctx.fillStyle = color + '20'
      ctx.fill()
    }
  }
}

const animate = () => {
  draw()
  animationFrameId = requestAnimationFrame(animate)
}

onMounted(() => {
  nextTick(() => {
    animate()

    if (containerRef.value && window.ResizeObserver) {
      resizeObserver = new ResizeObserver(() => {
        draw()
      })
      resizeObserver.observe(containerRef.value)
    }
  })
})

onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
  if (resizeObserver && containerRef.value) {
    resizeObserver.disconnect()
  }
})

watch([dataPoints, isAreaMode], () => {
  draw()
})
</script>

<style scoped>
.traffic-graph-wrapper {
  width: 100%;
}

.graph-footer {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

.traffic-graph {
  width: 100%;
  height: 200px;
  cursor: pointer;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  overflow: hidden;
  background: #ffffff;
  transition: all 0.2s ease-in-out;
}

.traffic-graph:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
  border-color: #cbd5e1;
}

.traffic-graph.compact {
  height: 140px;
}

canvas {
  display: block;
  width: 100%;
  height: 100%;
}

@media (max-width: 576px) {
  .graph-footer :deep(.el-radio-button__inner) {
    padding: 6px 10px;
    font-size: 12px;
  }
}
</style>
