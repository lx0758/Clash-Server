<template>
  <div class="sidebar-traffic">
    <canvas ref="canvasRef" class="mini-chart"></canvas>
    <div class="traffic-stats">
      <div class="stat-row">
        <el-icon class="stat-icon up"><Upload /></el-icon>
        <span class="stat-value">{{ formatSpeed(traffic.up) }}</span>
      </div>
      <div class="stat-row">
        <el-icon class="stat-icon down"><Download /></el-icon>
        <span class="stat-value">{{ formatSpeed(traffic.down) }}</span>
      </div>
      <div class="stat-row">
        <el-icon class="stat-icon mem"><Cpu /></el-icon>
        <span class="stat-value">{{ formatBytes(memory.inuse) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Upload, Download, Cpu } from '@element-plus/icons-vue'
import { useWebSocket } from '@/composables/useWebSocket'
import { formatBytes } from '@/utils/format'

const { traffic, memory } = useWebSocket()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const historyUp = ref<number[]>([])
const historyDown = ref<number[]>([])
const maxHistory = 30

let animationFrameId: number | null = null
let dataIntervalId: number | null = null

const formatSpeed = (bytes: number) => {
  if (!bytes || bytes < 0) return '0 B/s'
  return formatBytes(bytes) + '/s'
}

const draw = () => {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const width = canvas.width / (window.devicePixelRatio || 1)
  const height = canvas.height / (window.devicePixelRatio || 1)

  ctx.clearRect(0, 0, width, height)

  if (historyUp.value.length < 2) return

  const upMax = Math.max(...historyUp.value, 0)
  const downMax = Math.max(...historyDown.value, 0)
  const maxVal = Math.max(upMax, downMax, 1024)

  const drawLine = (data: number[], color: string) => {
    if (data.length < 2) return

    const points: { x: number; y: number }[] = []
    const step = width / (maxHistory - 1)
    const startIndex = maxHistory - data.length

    for (let i = 0; i < data.length; i++) {
      const val = data[i] ?? 0
      points.push({
        x: (startIndex + i) * step,
        y: height - (val / maxVal) * height * 0.9
      })
    }

    ctx.beginPath()
    ctx.strokeStyle = color
    ctx.lineWidth = 1.5

    const first = points[0]
    if (!first) return
    ctx.moveTo(first.x, first.y)

    for (let i = 0; i < points.length - 1; i++) {
      const p0 = points[Math.max(0, i - 1)]
      const p1 = points[i]
      const p2 = points[i + 1]
      const p3 = points[Math.min(points.length - 1, i + 2)]
      
      if (!p0 || !p1 || !p2 || !p3) continue

      const tension = 0.3
      const cp1x = p1.x + (p2.x - p0.x) * tension
      const cp1y = p1.y + (p2.y - p0.y) * tension
      const cp2x = p2.x - (p3.x - p1.x) * tension
      const cp2y = p2.y - (p3.y - p1.y) * tension

      ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, p2.x, p2.y)
    }
    ctx.stroke()
  }

  drawLine(historyDown.value, '#3b82f6')
  drawLine(historyUp.value, '#f59e0b')
}

const collectData = () => {
  const t = traffic.value
  const up = t?.up || 0
  const down = t?.down || 0
  historyUp.value.push(up)
  historyDown.value.push(down)

  if (historyUp.value.length > maxHistory) {
    historyUp.value.shift()
    historyDown.value.shift()
  }
}

const animate = () => {
  draw()
  animationFrameId = requestAnimationFrame(animate)
}

onMounted(() => {
  const canvas = canvasRef.value
  if (canvas) {
    const dpr = window.devicePixelRatio || 1
    const rect = canvas.getBoundingClientRect()
    canvas.width = rect.width * dpr
    canvas.height = rect.height * dpr
    const ctx = canvas.getContext('2d')
    if (ctx) {
      ctx.scale(dpr, dpr)
    }
  }
  
  dataIntervalId = window.setInterval(collectData, 1000)
  animate()
})

onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
  if (dataIntervalId) {
    clearInterval(dataIntervalId)
  }
})
</script>

<style scoped>
.sidebar-traffic {
  padding: 10px 12px;
}

.mini-chart {
  width: 100%;
  height: 40px;
  display: block;
  margin-bottom: 8px;
}

.traffic-stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-icon {
  font-size: 14px;
  width: 16px;
  height: 16px;
}

.stat-icon.up {
  color: #f59e0b;
}

.stat-icon.down {
  color: #3b82f6;
}

.stat-icon.mem {
  color: #8b5cf6;
}

.stat-value {
  font-size: 12px;
  color: #1e293b;
  font-family: 'SF Mono', monospace;
  flex: 1;
  text-align: right;
}
</style>
