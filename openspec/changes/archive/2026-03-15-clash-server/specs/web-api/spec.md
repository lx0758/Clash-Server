## ADDED Requirements

### Requirement: RESTful API 设计
系统 SHALL 提供遵循 RESTful 规范的 API。

#### Scenario: 资源命名
- **WHEN** 设计 API 端点
- **THEN** 使用名词而非动词，如 /api/subscriptions 而非 /api/getSubscriptions
- **AND** 使用 HTTP 方法表达操作语义

#### Scenario: HTTP 方法规范
- **WHEN** 设计 API
- **THEN** GET 用于查询，POST 用于创建，PUT 用于完整更新，DELETE 用于删除
- **AND** PATCH 用于部分更新

### Requirement: 响应格式
系统 SHALL 提供统一的响应格式。

#### Scenario: 成功响应
- **WHEN** 请求成功
- **THEN** 返回 {"code": 0, "message": "success", "data": ...}

#### Scenario: 错误响应
- **WHEN** 请求失败
- **THEN** 返回 {"code": <业务状态码>, "message": "...", "detail": "..."}
- **AND** 不使用 HTTP 状态码表达业务状态

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

### Requirement: 配置管理
系统 SHALL 提供面板配置的获取和更新。

#### Scenario: 获取配置
- **WHEN** 调用 GET /api/config
- **AND** 返回面板配置（代理端口、模式、DNS 等）

#### Scenario: 更新配置
- **WHEN** 调用 PUT /api/config
- **AND** 更新面板配置
- **AND** 自动重启 Core
- **AND** 返回 Core 启动结果 (成功/失败+错误)

### Requirement: 系统信息
系统 SHALL 提供整合的系统信息接口。

#### Scenario: 获取系统信息
- **WHEN** 调用 GET /api/system/info
- **THEN** 返回整合的系统状态（包含 Core 状态、订阅统计、流量数据）
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

> **注意**: 流量数据已整合到 /api/system/info，无独立的 /api/statistics/traffic 接口。

### Requirement: 配置变更响应
系统 SHALL 在配置变更 API 中返回 Core 启动结果。

#### Scenario: 配置保存成功 + Core 启动成功
- **WHEN** 配置变更成功且 Core 启动成功
- **THEN** 返回 {"code": 0, "message": "success", "data": null}

#### Scenario: 配置保存成功 + Core 启动失败
- **WHEN** 配置变更成功但 Core 启动失败
- **THEN** 返回 {"code": 0, "message": "success", "data": {"core_error": "..."}}

#### Scenario: 配置保存失败
- **WHEN** 配置验证失败
- **THEN** 返回 {"code": 400, "message": "Invalid config", "data": null}