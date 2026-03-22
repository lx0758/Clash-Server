## Context

当前代理界面 (`web/src/pages/Proxies.vue`) 将所有节点平铺展示：
- 所有代理组全部展开
- 节点单行排列
- 无筛选/排序功能
- 无运行模式切换
- 节点多时界面过长

参考 Clash Verge Rev 的设计，需要实现更现代的代理管理界面。

## Goals / Non-Goals

**Goals:**
- 运行模式切换（规则/全局/直连）
- 可折叠代理组，减少界面占用
- 多列网格布局，展示更多节点
- 筛选/排序功能，快速找到节点
- 代理组导航，快速跳转
- 选中代理不刷新页面

**Non-Goals:**
- 不实现链式代理模式（Clash Verge 特有功能）
- 不实现虚拟滚动（节点数量不多时非必需）

## Decisions

### 决策 1：组件架构

```
web/src/
├── pages/
│   └── Proxies.vue
├── components/proxy/
│   ├── ModeSwitcher.vue        ← 模式切换器
│   ├── ProxyGroup.vue          ← 可折叠代理组
│   ├── ProxyHead.vue           ← 操作栏
│   ├── ProxyItem.vue           ← 节点卡片
│   └── ProxyGrid.vue           ← 节点网格
├── stores/proxy.ts             ← 扩展状态
└── composables/useProxyFilter.ts ← 筛选排序
```

**理由**:
- 组件拆分清晰，职责单一
- 便于维护和测试
- 参照 Clash Verge Rev 的组件结构

### 决策 2：运行模式

| 模式 | 显示内容 | API 调用 |
|------|----------|----------|
| 规则 | 非预置代理组（排除 GLOBAL） | `GET /proxies` |
| 全局 | 所有代理（含 DIRECT/REJECT） | `GET /proxies/GLOBAL` |
| 直连 | 空白/提示文字 | 无需额外 API |

**规则模式代理组排序**:
1. Selector - 手动选择组
2. URLTest - 自动选择延迟最低
3. Fallback - 故障转移

**理由**:
- 与 Clash Verge Rev 行为一致
- 用户习惯的交互方式

### 决策 3：状态管理扩展

```typescript
interface ProxyGroupState {
  open: boolean       // 是否展开
  filterText: string  // 筛选文本
  sortType: 'default' | 'delay-asc' | 'delay-desc' | 'name'
}

interface ProxyState {
  rawProxies: Record<string, MihomoProxy | MihomoProxyGroup>
  groupStates: Record<string, ProxyGroupState>
  mode: 'rule' | 'global' | 'direct'
  loading: boolean
}
```

**理由**:
- 每个代理组独立管理展开/筛选/排序状态
- 模式状态持久化到后端
- 展开状态持久化到 localStorage

### 决策 4：Mihomo API 类型适配

**问题**: Mihomo API 返回的 `all` 字段是字符串数组（代理名称），不是完整对象。

**解决方案**:
```typescript
interface MihomoProxyGroup {
  name: string
  type: string
  all: string[]        // 代理名称列表
  now: string
}

interface MihomoProxy {
  name: string
  type: string
  alive: boolean
  history: Array<{ time: string; delay: number }>
}

// 从名称查找代理详情
const getProxyInfo = (name: string): Proxy => {
  const p = rawProxies.value[name]
  const history = p?.history || []
  const delay = history.length > 0 ? history[history.length - 1]?.delay ?? 0 : 0
  return { name, type: p?.type, alive: p?.alive, delay }
}
```

### 决策 5：响应式列数

```typescript
function calculateColumns(width: number): number {
  if (width > 1400) return 4
  if (width > 1000) return 3
  if (width > 700) return 2
  return 1
}
```

**理由**:
- 根据窗口宽度自动调整
- 提供良好的移动端体验

### 决策 6：选中代理不刷新

**问题**: 选中代理后调用 `fetchProxies()` 导致页面闪烁。

**解决方案**:
```typescript
const selectProxy = async (group: string, name: string) => {
  await proxyApi.select(group, name)
  // 仅更新本地状态，不重新拉取
  if (rawProxies.value[group]) {
    rawProxies.value = {
      ...rawProxies.value,
      [group]: { ...rawProxies.value[group], now: name }
    }
  }
}
```

### 决策 7：API 扩展

新增后端接口：

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/proxies/:name/delay` | GET | 测试单个节点延迟 |
| `/api/proxies/group/:group/delay` | GET | 测试组内所有节点延迟 |
| `/api/proxies/mode` | GET | 获取当前运行模式 |
| `/api/proxies/mode` | PATCH | 切换运行模式 |

**理由**:
- 前端需要调用 mihomo API 测试延迟和切换模式
- 后端代理调用，统一错误处理

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|---------|
| 组件重构范围大 | 分步骤实现，先核心功能 |
| 延迟测试并发量大 | 后端限流，前端队列控制 |
| Mihomo API 版本差异 | 类型定义兼容处理 |