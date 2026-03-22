## ADDED Requirements

### Requirement: 实时流量推送
系统 SHALL 通过 WebSocket 推送实时流量数据。

#### Scenario: 流量数据推送
- **WHEN** 客户端订阅 traffic 类型
- **THEN** 服务端定期推送上行/下行速度
- **AND** 推送格式: {"type": "traffic", "data": {"up": 100, "down": 200}}

#### Scenario: 流量数据格式
- **WHEN** 推送流量数据
- **AND** 包含 up (上行 bytes/s) 和 down (下行 bytes/s)

### Requirement: 连接状态推送
系统 SHALL 通过 WebSocket 推送连接状态变化。

#### Scenario: 连接列表推送
- **WHEN** 客户端订阅 connections 类型
- **THEN** 服务端定期推送当前连接列表
- **AND** 推送格式: {"type": "connections", "data": {"connections": [...], "downloadTotal": 0, "uploadTotal": 0}}

#### Scenario: 核心状态推送
- **WHEN** 客户端订阅 core_status 类型
- **AND** 核心状态变化时推送 (start/stop/crash)

### Requirement: 日志推送
系统 SHALL 通过 WebSocket 推送 Mihomo 运行日志。

#### Scenario: 日志推送
- **WHEN** 客户端订阅 logs 类型
- **THEN** 服务端实时推送 Mihomo 日志
- **AND** 推送格式: {"type": "logs", "data": {"type": "info", "payload": "..."}}

### Requirement: 内存使用推送
系统 SHALL 通过 WebSocket 推送 Mihomo 内存使用情况。

#### Scenario: 内存推送
- **WHEN** 客户端订阅 memory 类型
- **THEN** 服务端定期推送内存使用情况
- **AND** 推送格式: {"type": "memory", "data": {"inuse": 12345, "oslimit": 0}}

### Requirement: 订阅管理
系统 SHALL 支持客户端管理 WebSocket 订阅。

#### Scenario: 订阅消息
- **WHEN** 客户端发送 {"action": "subscribe", "types": ["traffic", "connections"]}
- **AND** 服务端开始推送指定类型的消息

#### Scenario: 取消订阅
- **WHEN** 客户端发送 {"action": "unsubscribe", "types": ["logs"]}
- **AND** 服务端停止推送指定类型的消息