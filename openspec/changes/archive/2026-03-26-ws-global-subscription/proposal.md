# Proposal: WebSocket 全局订阅与数据持久化

## 背景问题

### 现状分析

当前各页面独立订阅 WebSocket 事件，导致以下问题：

```
┌────────────────────────────────────────────────────────────────────┐
│                           Browser                                   │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐        │
│  │  Dashboard   │    │    Logs      │    │  Connections  │        │
│  │  (mounted)   │    │   (unmounted)│    │   (mounted)   │        │
│  │              │    │              │    │               │        │
│  │ subscribe    │    │ unsubscribe  │    │  subscribe    │        │
│  │ [traffic,    │    │   [logs]  ✗  │    │  [connections] │      │
│  │  connections] │    │              │    │               │        │
│  └──────────────┘    └──────────────┘    └──────────────┘        │
│          │                  │                     │                 │
│          └──────────────────┼─────────────────────┘                 │
│                             │                                       │
│                    ┌────────▼────────┐                             │
│                    │   useWebSocket   │                             │
│                    │   (模块级单例)    │                             │
│                    └──────────────────┘                            │
└────────────────────────────────────────────────────────────────────┘
```

### 问题 1：页面切换时消息丢失

```
用户操作:  Dashboard → Logs

Dashboard unmount:
  → unsubscribe(['connections', 'traffic'])
  
此时 WS 服务器推送 connections 更新
  → 客户端已取消订阅，收不到！
  
Logs mount:
  → subscribe(['logs'])
  → 开始接收新的 logs，但 connections 已经丢失了更新
```

### 问题 2：日志连续性中断

用户在 Logs 页面时，切到其他页面再切回来：
- 切走期间的新日志永远丢失
- 切回后从当前位置继续，但之前的数据已不在

### 问题 3：订阅泄漏（设计不一致）

| 页面/组件 | subscribe | unsubscribe |
|-----------|-----------|-------------|
| MainLayout | `[traffic, core_status]` | 无 (disconnect整体) |
| Dashboard | `[connections]` | **缺失！** |
| Logs | `[logs]` | ✓ 有 |
| Connections | `[connections]` | ✓ 有 |
| RuntimeInfoCard | `[memory]` | **缺失！** |

---

## 解决方案

### 方案 A：全局订阅（推荐）✓ 已采用

**核心思路**：连接时统一订阅所有类型，各页面按需读取共享状态

```
┌────────────────────────────────────────────────────────────────────┐
│                           Browser                                   │
│                                                                     │
│  MainLayout 在 connect() 时订阅所有类型:                           │
│  subscribe(['traffic', 'connections', 'logs', 'memory', ...])   │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐        │
│  │  Dashboard   │    │    Logs      │    │  Connections  │        │
│  │              │    │              │    │               │        │
│  │ 直接读取      │    │ 直接读取      │    │ 直接读取       │        │
│  │ logs.value   │    │ logs.value   │    │ logs.value    │        │
│  │ connections  │    │ connections  │    │ connections   │        │
│  │              │    │              │    │               │        │
│  │ 无 subscribe │    │ 无 subscribe │    │ 无 subscribe  │        │
│  │ 无 unsubscribe│    │ 无 unsubscribe│    │ 无 unsubscribe │        │
│  └──────────────┘    └──────────────┘    └──────────────┘        │
└────────────────────────────────────────────────────────────────────┘
```

---

## 数据管理策略

### 日志数据
- 保留最新的 **500 条历史日志**
- 超出时从旧的一端丢弃：`[newLog, ...oldLogs].slice(0, 500)`

### 连接数据

**数据结构**：
```typescript
interface ConnectionsData {
  connections: Connection[]        // 当前激活的连接，全部保留，无上限
  history: Connection[]           // 历史连接（已断开的），最多 500 条
  downloadTotal: number
  uploadTotal: number
}
```

**处理逻辑**：
1. WS 推送新的 connections 列表时
2. 对比新旧列表，找出"消失"的连接（已断开的）
3. 将断开的连接加入 history 列表头部
4. history 列表超出 500 条时，从尾部丢弃
5. 通过 `id` 字段去重，避免重复

**排序规则**：
- 激活的连接排在前面
- 历史连接按时间倒序排列（新的在前）

**连接状态展示**：
- 激活的连接：正常显示，有"关闭"按钮
- 已断开的连接：显示灰色 + "已断开"标签，无"关闭"按钮

**Connection 数据结构扩展**：
```typescript
interface Connection {
  id: string
  metadata?: { ... }
  upload?: number
  download?: number
  start?: string
  chains?: string[]
  rule?: string
  rulePayload?: string
  disconnected?: boolean  // 新增：标识是否为历史连接
}
```

---

## 改动范围

### 1. types/api.ts
- `Connection` 新增 `disconnected?: boolean` 字段
- `ConnectionsData` 新增 `history: Connection[]` 字段

### 2. composables/useWebSocket.ts
- 默认订阅所有类型：`['traffic', 'connections', 'logs', 'core_status', 'memory']`
- connections 处理逻辑改为：
  - 维护 `previousConnectionIds` 追踪上一轮活跃连接
  - 对比找出断开的连接，标记 `disconnected: true` 后加入 history
  - history 最多 500 条，超出从尾部丢弃
  - 通过 ID 去重
- 日志维持 500 条上限

### 3. 移除订阅调用
| 文件 | 移除内容 |
|------|----------|
| Logs.vue | `subscribe(['logs'])` / `unsubscribe(['logs'])` |
| Connections.vue | `subscribe(['connections'])` / `unsubscribe(['connections'])` |
| Dashboard.vue | `subscribe(['connections'])` |
| SidebarTraffic.vue | `subscribe(['traffic', 'memory'])` |
| RuntimeInfoCard.vue | `subscribe(['memory'])` |
| useTrafficGraph.ts | `subscribe(['traffic'])` |

### 4. pages/Connections.vue
- 合并活跃连接 + 历史连接显示
- 激活连接在前，历史连接在后
- 历史连接显示"已断开"标签
- 历史连接样式变暗
- 历史连接隐藏"关闭"按钮

---

## 验收标准

1. **日志**：始终显示最新的 500 条日志
2. **连接**：
   - 激活的连接全部保留，无数量限制
   - 历史连接最多保留 500 条
   - 激活连接排在前面，历史连接按时间倒序
   - 通过 ID 去重
   - 断开连接的 ITEM 有明确的状态标识
3. 页面切换后，原有数据保持连续不断
4. 构建通过，无类型错误
5. 双滚动条问题已修复

---

## 实现细节

### useWebSocket.ts 核心逻辑

```typescript
// 维护上一轮活跃连接的 ID 集合
let previousConnectionIds = new Set<string>()

// connections 处理
case 'connections': {
  const data = msg.data as { connections: Connection[], downloadTotal: number, uploadTotal: number }
  const newConnectionIds = new Set(data.connections.map(c => c.id))
  
  // 找出断开的连接
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
  
  // 去重并合并到 history
  const existingHistoryIds = new Set(connections.value.history.map(c => c.id))
  const newHistory = disconnectedConnections.filter(c => !existingHistoryIds.has(c.id))
  const history = [...newHistory, ...connections.value.history].slice(0, 500)
  
  connections.value = {
    connections: data.connections,
    history,
    downloadTotal: data.downloadTotal,
    uploadTotal: data.uploadTotal,
  }
  break
}
```

### Connections.vue UI 差异

```vue
<!-- 激活连接 -->
<div class="connection-item">
  <div class="conn-host">{{ host }}</div>
  <el-button type="danger" class="close-btn">关闭</el-button>
</div>

<!-- 历史连接 -->
<div class="connection-item is-disconnected">
  <div class="conn-host">
    {{ host }}
    <el-tag type="info" size="small">已断开</el-tag>
  </div>
  <!-- 无关闭按钮 -->
</div>
```

---

## 状态

- **Proposed**: 2026-03-26
- **Status**: ✓ implemented
- **Completed**: 2026-03-26

---

## 更新记录

### 2026-03-26 新需求澄清

**原需求理解**：
- connections 总数不超过 500 条

**新需求**：
- 激活连接不受限制，全部保留
- 历史连接（已断开的）最多 500 条
- 排列：激活在前，历史按时间倒序
- 通过 ID 去重
- 历史连接 ITEM 需展示"已断开"状态

### 2026-03-26 实现完成

所有需求已实现并验证通过：
- [x] 日志 500 条上限
- [x] 连接分类：活跃 + 历史
- [x] 历史连接 500 条上限
- [x] ID 去重
- [x] 已断开状态展示
- [x] 构建通过
