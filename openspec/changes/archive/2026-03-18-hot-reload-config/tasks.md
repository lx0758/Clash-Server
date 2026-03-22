## 1. CoreService 重构

- [x] 1.1 新增 `RestartRequiredKeys` 常量定义需要软重启的配置项列表
- [x] 1.2 新增 `ApplyConfigResult` 结构体定义返回结果
- [x] 1.3 新增 `needsRestart()` 方法判断配置变更是否需要软重启
- [x] 1.4 新增 `HotReload()` 方法，调用 Mihomo PUT /configs API 热重载
- [x] 1.5 新增 `SoftRestart()` 方法，调用 Mihomo POST /restart API 软重启
- [x] 1.6 新增 `ApplyConfig()` 方法，智能选择热重载或软重启
- [x] 1.7 重构 `start()` 方法，从 CoreConfig 获取 API 地址，生成 Base64 编码的最小配置
- [x] 1.8 移除临时配置文件相关逻辑（tmpDir, tmpFile 字段）
- [x] 1.9 重构 `stop()` 方法，仅在降级场景使用

## 2. 响应结构改造

- [x] 2.1 新增 `SuccessWithCoreResult()` 响应方法
- [x] 2.2 新增 `NeedConfirmRestart()` 响应方法

## 3. Handler 适配

- [x] 3.1 改造 `SubscriptionHandler.applyConfig()` 使用 ApplyConfig
- [x] 3.2 改造 `ConfigHandler.applyConfig()` 使用 ApplyConfig
- [x] 3.3 检查并改造其他 Handler 中的核心重启调用

## 4. 测试验证

- [x] 4.1 编写单元测试：needsRestart 判断逻辑
- [x] 4.2 编写单元测试：HotReload 热重载
- [x] 4.3 编写单元测试：SoftRestart 软重启
- [x] 4.4 编写集成测试：完整启动和配置应用流程