## Why

订阅调度器存在多个问题导致订阅无法按时刷新：服务重启后不根据 `last_refresh` 计算下次刷新时间、本地订阅被错误调度、Stop/Start 重入失败、刷新失败后立即重试导致频繁请求。这些问题使得用户配置的定时刷新功能完全失效。

## What Changes

- 修复调度器启动时根据 `last_refresh` 计算下次刷新时间，而非从当前时间开始
- 排除本地订阅（`source_type=local`）参与定时调度
- 修复 `Stop()` 后再 `Start()` 失败的问题（重新创建 `stopChan`）
- 当刷新时间已过时，设置最小 5 分钟重试间隔，避免立即重试
- 添加刷新失败日志记录，便于问题排查

## Capabilities

### New Capabilities

无新功能引入，本变更为 Bug 修复

### Modified Capabilities

无需求变更，现有 `subscription-management` 规格中的"订阅定时刷新"需求保持不变，本变更修复其实现缺陷

## Impact

- 代码：`server/internal/scheduler/subscription.go`
- 依赖：无新增依赖
- 系统：订阅刷新行为变化，本地订阅不再参与定时调度
