# Design: WebSocket 全局订阅与数据持久化

## 技术决策

### 1. 模块级单例模式

**决策**：使用模块级变量存储 WebSocket 状态

```typescript
// composables/useWebSocket.ts
const ws = ref<WebSocket | null>(null)
const connected = ref(false)
const connections = ref<ConnectionsState>({ connections: [], history: [], ... })
const logs = ref<LogData[]>([])
```

**理由**：
- Vue composable 通常返回响应式引用，多处调用会各自复制状态
- 模块级变量确保所有调用方共享同一份数据
- `ref()` 包装保证 Vue 的响应式追踪

**权衡**：
- 测试时需要 mock 模块级变量
- SSR 场景可能有问题（本项目为纯 SPA，无此担忧）

---

### 2. 连接状态追踪策略

**问题**：如何判断连接已断开？

**方案**：通过 ID 对比检测

```typescript
let previousConnectionIds = new Set<string>()

// WS 新推送时
const newConnectionIds = new Set(data.connections.map(c => c.id))

// 上一轮有，但现在没有 → 已断开
for (const id of previousConnectionIds) {
  if (!newConnectionIds.has(id)) {
    // 找到断开的连接
  }
}
previousConnectionIds = newConnectionIds
```

**替代方案对比**：

| 方案 | 优点 | 缺点 |
|------|------|------|
| ID 对比（采用） | 简单，不依赖服务端 | 瞬时断开重连可能丢失 |
| 服务端推送断开事件 | 精确 | 需要改服务端 |
| 心跳检测 | 能发现静默断开 | 增加复杂度 |

---

### 3. 历史连接去重策略

**问题**：同一连接可能反复断开/重连

**方案**：基于 ID 的去重

```typescript
const existingHistoryIds = new Set(connections.value.history.map(c => c.id))
const newHistory = disconnectedConnections.filter(c => !existingHistoryIds.has(c.id))
```

**场景**：
1. 连接 A 断开 → 进入 history
2. 连接 A 重连 → 从 history 移除
3. 连接 A 再次断开 → 视为新的断开事件

---

### 4. 数据上限策略

| 数据类型 | 上限 | 策略 |
|----------|------|------|
| logs | 500 | 新日志在前，超出从尾部丢弃 |
| connections (活跃) | 无 | 服务端推送多少显示多少 |
| connections (历史) | 500 | 新断开在前，超出从尾部丢弃 |

**理由**：
- 活跃连接不应限制，确保信息完整性
- 历史记录限制防止内存溢出
- 使用 `slice(0, 500)` 是简洁的实现方式

---

### 5. 连接数据结构分离

**原方案**：单一 connections 数组，内部标记 disconnected

**现方案**：分离为 connections + history 两个数组

```typescript
interface ConnectionsState {
  connections: Connection[]  // 活跃
  history: Connection[]      // 已断开
}
```

**理由**：
- UI 层只需简单拼接：`[...active, ...history]`
- 无需在渲染时过滤
- disconnected 标记仅用于 UI 样式控制

---

## 架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                         WS 服务端 (Golang)                          │
│                                                                     │
│  Hub.broadcastToType()                                              │
│       │                                                             │
└───────┼─────────────────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Browser WebSocket                               │
│                                                                     │
│  ws.onmessage() ──────────────────────────────────────────────────┐│
│        │                                                           ││
└────────┼───────────────────────────────────────────────────────────┘│
         │                                                             │
         ▼                                                             │
┌─────────────────────────────────────────────────────────────────────┐
│                    useWebSocket (模块级单例)                          │
│                                                                     │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐              │
│  │ connections │   │    logs    │   │   traffic   │              │
│  │  ref.value  │   │  ref.value  │   │  ref.value  │              │
│  └──────┬──────┘   └──────┬──────┘   └──────┬──────┘              │
│         │                 │                 │                       │
│  ┌──────▼──────┐   ┌──────▼──────┐   (直接使用)                    │
│  │ connections │   │ logs.value │                                │
│  │  + history  │   │  (500上限)  │                                │
│  └─────────────┘   └─────────────┘                                │
│                                                                     │
│  所有 Vue 组件直接读取模块级变量，共享同一份数据                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 关键代码

### 连接断开检测

```typescript
case 'connections': {
  const data = msg.data as { connections: Connection[], downloadTotal: number, uploadTotal: number }
  const newConnectionIds = new Set(data.connections.map(c => c.id))
  
  const disconnectedConnections: Connection[] = []
  for (const id of previousConnectionIds) {
    if (!newConnectionIds.has(id)) {
      const lastState = connections.value.connections.find(c => c.id === id) ||
                       connections.value.history.find(c => c.id === id)
      if (lastState) {
        disconnectedConnections.push({ ...lastState, disconnected: true })
      }
    }
  }
  previousConnectionIds = newConnectionIds
  
  const existingHistoryIds = new Set(connections.value.history.map(c => c.id))
  const newHistory = disconnectedConnections.filter(c => !existingHistoryIds.has(c.id))
  const history = [...newHistory, ...connections.value.history].slice(0, 500)
  
  connections.value = { connections: data.connections, history, ... }
  break
}
```

### UI 状态区分

```vue
<div class="connection-item" :class="{ 'is-disconnected': conn.disconnected }">
  <div class="conn-host">
    {{ host }}
    <el-tag v-if="conn.disconnected" type="info" size="small">已断开</el-tag>
  </div>
  <el-button v-if="!conn.disconnected" type="danger">关闭</el-button>
</div>
```

---

## 性能考量

### 内存
- connections 活跃列表：无上限，但通常服务端会限制（几千条以内）
- history 列表：固定 500 条上限
- logs 列表：固定 500 条上限
- 每次 WS 消息触发 O(n) 的 ID 集合操作，n=活跃连接数

### 渲染
- Vue 的 `ref` 自动追踪依赖
- filteredConnections 是 computed，WS 更新时自动重算
- 大列表建议加虚拟滚动（本项目暂未实现）

---

## 未来可能的优化

1. **虚拟滚动**：连接数超过 100 时考虑
2. **Web Worker**：大数据量时将数据处理移至 Worker
3. **服务端缓存**：服务端维护历史连接，客户端按需拉取
4. **增量更新**：服务端推送 delta 而非全量列表
