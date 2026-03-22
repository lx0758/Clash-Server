## ADDED Requirements

### Requirement: 代理节点列表
系统 SHALL 提供当前可用的代理节点列表。

#### Scenario: 获取节点列表
- **WHEN** 调用 GET /api/proxies
- **THEN** 返回所有代理节点，包含名称、类型、延迟
- **AND** 按分组和类型组织

#### Scenario: 获取单个节点
- **WHEN** 调用 GET /api/proxies/:name
- **AND** 返回节点详细信息，包括所有可用选项

### Requirement: 代理节点选择
系统 SHALL 支持手动选择代理节点或使用自动选择。

#### Scenario: 选择代理
- **WHEN** 调用 PUT /api/proxies/:name，指定目标节点
- **THEN** Mihomo 配置更新，使用该节点作为主代理
- **AND** 返回当前选中的节点

#### Scenario: 自动选择
- **WHEN** 选择自动模式 (url-test/fallback/load-balance)
- **THEN** Mihomo 使用对应策略自动选择节点
- **AND** 返回当前策略类型

### Requirement: 延迟测试
系统 SHALL 支持测试代理节点的延迟。

#### Scenario: 测试节点延迟
- **WHEN** 调用 POST /api/proxies/:name/delay
- **THEN** 对节点进行 HTTP/HTTPS 延迟测试
- **AND** 返回延迟值（毫秒）

### Requirement: 连接管理
系统 SHALL 提供当前活动连接列表和管理能力。

#### Scenario: 获取连接列表
- **WHEN** 调用 GET /api/connections
- **AND** 返回所有活动连接，包含源、目的、流量

#### Scenario: 断开连接
- **WHEN** 调用 DELETE /api/connections/:id
- **AND** 断开指定连接

#### Scenario: 断开所有连接
- **WHEN** 调用 DELETE /api/connections
- **AND** 断开所有活动连接