# Spec Delta

## MODIFIED Requirements

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
- **THEN** 返回 {"code": 1007, "data": {"core_error": "...", "partial_success": true}}
- **AND** 客户端据此提示"已保存但未生效"而非"操作失败"

### Requirement: 配置变更智能应用
系统 SHALL 在配置变更时智能选择热重载或软重启。

#### Scenario: 订阅变更触发智能应用
- **WHEN** 创建/更新/删除/刷新订阅
- **THEN** 保存成功后智能应用配置（热重载或软重启）
- **AND** 返回操作结果

#### Scenario: 规则变更触发智能应用
- **WHEN** 创建/更新/删除自定义规则
- **THEN** 保存成功后智能应用配置
- **AND** 返回操作结果

#### Scenario: 脚本变更触发智能应用
- **WHEN** 创建/更新/删除扩展脚本
- **THEN** 保存成功后智能应用配置
- **AND** 返回操作结果

#### Scenario: 面板配置变更触发智能应用
- **WHEN** 更新面板配置 (PUT /api/config)
- **THEN** 保存成功后智能应用配置
- **AND** 返回操作结果

#### Scenario: 应用失败时的结果表达
- **WHEN** 配置已保存但智能应用失败
- **THEN** 返回 code 为 1007 且携带 partial_success 标记
- **AND** 返回体说明未生效原因

#### Scenario: 自定义配置应用失败时回滚
- **WHEN** 自定义配置保存后核心应用失败
- **THEN** 已保存的自定义配置回滚为保存前的状态
- **AND** 返回错误说明应用失败

## ADDED Requirements

### Requirement: 核心数据查询的失败语义
系统 SHALL 在核心不可用时明确返回失败，而非以空结果伪装成功。

#### Scenario: 核心未运行时的代理查询
- **WHEN** 核心未运行且请求代理列表
- **THEN** 返回非零业务状态码
- **AND** 返回体说明核心未运行

#### Scenario: 核心未运行时的规则与连接查询
- **WHEN** 核心未运行且请求规则列表或连接列表
- **THEN** 返回非零业务状态码
- **AND** 返回体说明核心未运行

#### Scenario: 核心错误信息透传
- **WHEN** 核心 API 调用失败并返回了错误响应
- **THEN** 返回给客户端的错误信息包含核心的原始错误内容
- **AND** 不得以泛化文案替代核心的具体错误
