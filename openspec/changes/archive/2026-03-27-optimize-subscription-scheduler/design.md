## Context

订阅调度器 `SubscriptionScheduler` 负责定时刷新远程订阅。当前实现存在以下问题：

1. **启动时忽略 `last_refresh`**：调度器从当前时间计算下次刷新，导致服务重启后刷新时间错乱
2. **本地订阅被调度**：`source_type=local` 的订阅不应该参与定时刷新
3. **Stop/Start 重入失败**：`stopChan` 关闭后不会重新创建
4. **刷新失败立即重试**：当 `last_refresh + interval` 已过，`nextRefresh=0` 导致立即重试

## Goals / Non-Goals

**Goals:**
- 调度器启动时正确计算下次刷新时间（基于 `last_refresh`）
- 排除本地订阅参与定时调度
- 支持 Stop/Start 正确重入
- 刷新失败后设置最小重试间隔，避免频繁请求
- 添加日志便于问题排查

**Non-Goals:**
- 不改变订阅刷新的 API 接口
- 不改变用户可见的需求行为
- 不引入新的外部依赖

## Decisions

### 1. 使用 `last_refresh` 计算下次刷新时间

**决定**: 启动时检查 `last_refresh`，计算 `nextRefreshTime = lastRefresh + interval`

**理由**:
- 保证服务重启后刷新时间一致性
- 如果已过刷新时间，设置 5 分钟最小间隔后重试

**替代方案**:
- 直接立即刷新：可能导致服务启动时大量请求机场 API
- 记录下次刷新时间到数据库：增加复杂度，当前方案足够

### 2. 排除本地订阅

**决定**: 在 `scheduleRefresh` 开头检查 `sub.SourceType != model.SourceTypeRemote`

**理由**:
- 本地订阅没有 URL，无法刷新
- `interval` 字段对本地订阅无意义

### 3. Start() 中重新创建 stopChan

**决定**: 在 `Start()` 方法中重新创建 `stopChan` 和 `timers`

**理由**:
- 允许 Stop 后重新启动
- 保持资源一致性

### 4. 最小重试间隔

**决定**: 当刷新时间已过，设置 `nextRefresh = 5 * time.Minute`

**理由**:
- 避免刷新失败后立即重试导致的频繁请求
- 5 分钟足够让网络问题恢复，又不会等待太久

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 5 分钟最小间隔可能延迟刷新 | 用户可手动刷新；权衡频繁请求的风险更高 |
| 服务长时间停机后启动会触发刷新 | 已实现：超过刷新时间后等待 5 分钟再刷新，而非立即 |
