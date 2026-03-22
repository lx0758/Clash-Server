import { ref } from 'vue'
import type { Traffic, CoreStatus, ConnectionsData, LogData, MemoryData } from '@/types/api'
import { useSystemStore } from '@/stores/system'

type MessageType = 'traffic' | 'connections' | 'logs' | 'core_status' | 'memory'

interface WebSocketMessage {
  type: MessageType
  data: unknown
}

const ws = ref<WebSocket | null>(null)
const connected = ref(false)
const traffic = ref<Traffic>({ up: 0, down: 0 })
const coreStatus = ref<CoreStatus>({ running: false })
const connections = ref<ConnectionsData>({ connections: [], downloadTotal: 0, uploadTotal: 0 })
const logs = ref<LogData[]>([])
const memory = ref<MemoryData>({ inuse: 0, oslimit: 0 })
const subscribedTypes = ref<Set<MessageType>>(new Set(['traffic', 'core_status']))
let reconnectTimer: number | null = null

const connect = () => {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    return
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws.value = new WebSocket(`${protocol}//${window.location.host}/api/ws`)

  ws.value.onopen = () => {
    connected.value = true
    if (reconnectTimer) {
      window.clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (subscribedTypes.value.size > 0) {
      ws.value?.send(JSON.stringify({ action: 'subscribe', types: Array.from(subscribedTypes.value) }))
    }
  }

  ws.value.onclose = () => {
    connected.value = false
    ws.value = null
    if (!reconnectTimer) {
      reconnectTimer = window.setTimeout(() => {
        reconnectTimer = null
        connect()
      }, 3000)
    }
  }

  ws.value.onerror = () => {
    ws.value?.close()
  }

  ws.value.onmessage = (event) => {
    const msg: WebSocketMessage = JSON.parse(event.data)
    switch (msg.type) {
      case 'traffic':
        traffic.value = msg.data as Traffic
        break
      case 'connections':
        connections.value = msg.data as ConnectionsData
        break
      case 'logs':
        const logData = msg.data as LogData
        logs.value = [logData, ...logs.value].slice(0, 500)
        break
      case 'memory':
        memory.value = msg.data as MemoryData
        break
      case 'core_status':
        const status = msg.data as CoreStatus
        coreStatus.value = status
        try {
          const systemStore = useSystemStore()
          systemStore.updateCoreStatus(status)
        } catch {
          // Store not initialized yet
        }
        break
    }
  }
}

const disconnect = () => {
  if (reconnectTimer) {
    window.clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (ws.value) {
    ws.value.onopen = null
    ws.value.onclose = null
    ws.value.onerror = null
    ws.value.onmessage = null
    ws.value.close()
    ws.value = null
  }
  connected.value = false
}

const subscribe = (types: MessageType[]) => {
  types.forEach(t => subscribedTypes.value.add(t))
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({ action: 'subscribe', types }))
  } else {
    setTimeout(() => subscribe(types), 100)
  }
}

const unsubscribe = (types: MessageType[]) => {
  types.forEach(t => subscribedTypes.value.delete(t))
  ws.value?.send(JSON.stringify({ action: 'unsubscribe', types }))
}

const clearLogs = () => {
  logs.value = []
}

export function useWebSocket() {
  return {
    connect,
    disconnect,
    subscribe,
    unsubscribe,
    clearLogs,
    connected,
    traffic,
    coreStatus,
    connections,
    logs,
    memory,
  }
}
