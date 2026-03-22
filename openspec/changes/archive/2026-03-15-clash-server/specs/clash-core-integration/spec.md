## ADDED Requirements

### Requirement: Core 自动启动
系统 SHALL 在 Server 启动时自动启动 Mihomo 核心，无论是否有订阅配置。

#### Scenario: Server 启动时自动启动 Core
- **WHEN** Server 启动完成
- **THEN** Core 自动启动
- **AND** 如果有订阅，使用合并后的配置
- **AND** 如果无订阅，使用最小配置 (mode: direct)

#### Scenario: 无订阅时的最小配置
- **WHEN** Server 启动时没有任何订阅
- **THEN** Core 使用最小配置运行
- **AND** mixed-port 可用
- **AND** 所有流量直连 (MATCH,DIRECT)
- **AND** 所有 API 可用，返回空数据

### Requirement: 配置变更自动重启
系统 SHALL 在配置变更时自动重启 Core，并返回重启结果。

#### Scenario: 订阅变更触发重启
- **WHEN** 创建/更新/删除/刷新订阅
- **THEN** 保存成功后立即重启 Core
- **AND** 返回操作结果 (包含 Core 启动错误，如果有)

#### Scenario: 规则变更触发重启
- **WHEN** 创建/更新/删除自定义规则
- **THEN** 保存成功后立即重启 Core
- **AND** 返回操作结果

#### Scenario: 脚本变更触发重启
- **WHEN** 创建/更新/删除扩展脚本
- **THEN** 保存成功后立即重启 Core
- **AND** 返回操作结果

#### Scenario: 面板配置变更触发重启
- **WHEN** 更新面板配置 (PUT /api/config)
- **THEN** 保存成功后立即重启 Core
- **AND** 返回操作结果

### Requirement: Core 启动错误暴露
系统 SHALL 显式暴露 Core 启动失败的错误信息。

#### Scenario: Core 启动失败
- **WHEN** Core 启动失败 (配置错误等)
- **THEN** 记录错误信息
- **AND** GET /api/system/info 返回 core.error 字段
- **AND** WebSocket 推送 core_status 包含错误
- **AND** 保持失败状态，等待下次配置变更

#### Scenario: 配置变更 API 返回 Core 错误
- **WHEN** 配置保存成功但 Core 启动失败
- **THEN** 返回 {"code": 0, "data": {"core_error": "..."}}

### Requirement: Core 崩溃处理
系统 SHALL 检测 Core 崩溃并更新状态，不自动重启。

#### Scenario: Core 崩溃
- **WHEN** Core 进程异常退出
- **THEN** 更新运行状态为 stopped
- **AND** WebSocket 推送 core_status
- **AND** 等待用户修改配置触发重启

### Requirement: 系统信息整合
系统 SHALL 通过统一接口提供 Core 状态、版本和看板数据。

#### Scenario: 获取系统信息
- **WHEN** 调用 GET /api/system/info
- **THEN** 返回整合的系统状态
```json
{
  "core": {
    "running": true,
    "version": "v1.18.1",
    "error": null
  },
  "subscription": {
    "count": 2,
    "proxy_count": 156
  },
  "traffic": {
    "up": 1024,
    "down": 20480
  }
}
```

### Requirement: API 降级行为
系统 SHALL 在 Core 未运行时返回空数据而非错误。

#### Scenario: 代理列表降级
- **WHEN** Core 未运行时调用 GET /api/proxies
- **THEN** 返回 {"code": 0, "data": {"proxies": []}}

#### Scenario: 连接列表降级
- **WHEN** Core 未运行时调用 GET /api/connections
- **THEN** 返回 {"code": 0, "data": {"connections": []}}

#### Scenario: 流量统计降级
- **WHEN** Core 未运行时通过 GET /api/system/info 获取流量
- **THEN** 返回 traffic: null

## REMOVED Requirements

### Requirement: 手动启停 API (已移除)
~~系统 SHALL 提供手动启动/停止/重启 Core 的 API。~~

已移除的 API:
- GET /api/core/version
- GET /api/core/status
- POST /api/core/start
- POST /api/core/stop
- POST /api/core/restart
- POST /api/core/reload

**移除原因:** Core 现在完全自动管理，用户无需手动操作。
