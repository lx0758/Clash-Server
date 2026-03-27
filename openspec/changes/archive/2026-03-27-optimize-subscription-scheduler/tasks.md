## 1. 调度器启动时计算刷新时间

- [x] 1.1 在 `scheduleRefresh` 中根据 `last_refresh` 计算下次刷新时间
- [x] 1.2 当刷新时间已过时，设置 5 分钟最小重试间隔

## 2. 排除本地订阅

- [x] 2.1 在 `scheduleRefresh` 开头检查 `SourceType != remote` 并提前返回

## 3. 修复 Stop/Start 重入

- [x] 3.1 在 `Start()` 方法中重新创建 `stopChan`
- [x] 3.2 在 `Start()` 方法中重新创建 `timers` map

## 4. 添加日志

- [x] 4.1 刷新失败时记录日志
- [x] 4.2 获取订阅失败时记录日志

## 5. 验证

- [ ] 5.1 编写单元测试验证调度逻辑
- [ ] 5.2 手动测试服务重启后刷新时间计算
