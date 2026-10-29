# traffic-monitoring 规格说明

## Purpose

定义面板通过 WebSocket 向客户端推送实时数据的约定，涵盖流量、连接、日志、内存使用与延迟测试结果的推送格式和订阅机制。

## Requirements

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

### Requirement: 延迟结果推送
系统 SHALL 通过 WebSocket 推送延迟测试过程中单个节点的延迟结果。

#### Scenario: 延迟消息推送
- **WHEN** 客户端订阅 delay 类型
- **THEN** 服务端在每个节点测试完成时推送该节点的延迟结果
- **AND** 推送格式: {"type": "delay", "data": {"group": "G1", "name": "HK-01", "delay": 120}}

#### Scenario: 测速完成信号
- **WHEN** 目标组的全部节点测试完成
- **THEN** 服务端推送该组的完成信号
- **AND** 推送格式: {"type": "delay", "data": {"group": "G1", "done": true}}

#### Scenario: 节点测试失败
- **WHEN** 某个节点的延迟测试失败
- **THEN** 该节点的结果中体现失败状态
- **AND** 不影响同组其他节点的测试继续进行

#### Scenario: 逐节点结果与完成信号的关系
- **WHEN** 客户端接收 delay 类型消息
- **THEN** 在收到该组的完成信号前，该组节点处于测试中状态
- **AND** 收到完成信号后，该组全部节点退出测试中状态
