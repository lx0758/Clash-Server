# Spec Delta

## MODIFIED Requirements

### Requirement: 响应格式
系统 SHALL 提供统一的响应格式。

#### Scenario: 成功响应
- **WHEN** 请求成功
- **THEN** 返回 {"code": 0, "message": "success", "data": ...}

#### Scenario: 错误响应
- **WHEN** 请求失败
- **THEN** 返回 {"code": <业务状态码>, "message": "...", "detail": "..."}
- **AND** 不使用 HTTP 状态码表达业务状态

#### Scenario: 部分成功响应
- **WHEN** 变更已持久化但未在运行中的核心生效
- **THEN** 返回 {"code": 1007, "message": "...", "data": {..., "partial_success": true}}
- **AND** data 中携带未生效的原因
- **AND** 客户端可据此区分"未生效"与"操作失败"

#### Scenario: 列表响应
- **WHEN** 返回列表资源
- **THEN** 包含 meta 字段如 {"total": 100, "page": 1, "page_size": 20}

### Requirement: 业务状态码
系统 SHALL 使用独立的业务状态码。

#### Scenario: 状态码规范
- **WHEN** 返回业务错误
- **AND** code 与 HTTP 状态码含义一致
- **AND** 避免与网络层状态码混淆

| code | 含义 |
|------|------|
| 0 | 成功 |
| 400 | 请求格式错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 422 | 业务逻辑错误 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用 |

#### Scenario: 核心未生效状态码
- **WHEN** 变更已持久化但核心重载失败
- **THEN** 返回 code 为 1007
- **AND** 响应携带 partial_success 标记

#### Scenario: 核心不可用时的数据查询
- **WHEN** 查询依赖运行中核心的数据（代理、规则、连接）且核心不可用
- **THEN** 返回非零业务状态码
- **AND** 不得以空结果集伪装为成功的查询

### Requirement: 配置变更响应
系统 SHALL 在配置变更 API 中返回 Core 启动结果。

#### Scenario: 配置保存成功 + Core 启动成功
- **WHEN** 配置变更成功且 Core 启动成功
- **THEN** 返回 {"code": 0, "message": "success", "data": null}

#### Scenario: 配置保存成功 + Core 启动失败
- **WHEN** 配置变更成功但 Core 启动失败
- **THEN** 返回 {"code": 1007, "message": "...", "data": {"core_error": "...", "partial_success": true}}
- **AND** 客户端据此以警告语义而非错误语义提示用户

#### Scenario: 配置保存失败
- **WHEN** 配置验证失败
- **THEN** 返回 {"code": 400, "message": "Invalid config", "data": null}

#### Scenario: 核心错误原因透传
- **WHEN** 核心返回失败响应
- **THEN** 返回的 core_error 包含核心返回的原始错误信息
- **AND** 不得以泛化文案替代核心的具体错误
