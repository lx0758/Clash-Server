## 1. 类型定义更新

- [x] 1.1 `types/api.ts` - `Connection` 新增 `disconnected?: boolean` 字段
- [x] 1.2 `types/api.ts` - `ConnectionsData` 新增 `history: Connection[]` 字段

## 2. useWebSocket.ts 重构

- [x] 2.1 默认订阅所有类型：`['traffic', 'connections', 'logs', 'core_status', 'memory']`
- [x] 2.2 维护 `previousConnectionIds` 追踪上一轮活跃连接
- [x] 2.3 实现断开连接检测逻辑
- [x] 2.4 历史连接去重逻辑
- [x] 2.5 history 列表 500 条上限

## 3. 页面订阅调用清理

- [x] 3.1 `Logs.vue` - 移除 `subscribe(['logs'])` / `unsubscribe(['logs'])` / `connect()` / `disconnect()`
- [x] 3.2 `Connections.vue` - 移除 `subscribe(['connections'])` / `unsubscribe(['connections'])` / `connect()` / `disconnect()`
- [x] 3.3 `Dashboard.vue` - 移除 `subscribe(['connections'])`
- [x] 3.4 `SidebarTraffic.vue` - 移除 `subscribe(['traffic', 'memory'])`
- [x] 3.5 `RuntimeInfoCard.vue` - 移除 `subscribe(['memory'])`
- [x] 3.6 `useTrafficGraph.ts` - 移除 `subscribe(['traffic'])`

## 4. Connections.vue UI 改造

- [x] 4.1 合并活跃连接 + 历史连接显示
- [x] 4.2 历史连接显示"已断开"标签
- [x] 4.3 历史连接样式变暗 (opacity: 0.7)
- [x] 4.4 历史连接隐藏"关闭"按钮

## 5. 验证

- [x] 5.1 构建通过 `npm run build`
- [x] 5.2 无 TypeScript 类型错误
