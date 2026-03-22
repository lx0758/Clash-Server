## MODIFIED Requirements

### Requirement: 规则关联订阅
规则 SHALL 关联到具体订阅，不再支持全局规则。

#### Scenario: 创建规则
- **WHEN** 调用 POST /api/subscriptions/:subscriptionId/rules
- **THEN** 创建规则并关联到指定订阅
- **AND** 规则自动获取 `subscription_id`

#### Scenario: 规则列表
- **WHEN** 调用 GET /api/subscriptions/:subscriptionId/rules
- **THEN** 返回该订阅的所有规则
- **AND** 不包含其他订阅的规则

#### Scenario: 更新规则
- **WHEN** 调用 PUT /api/subscriptions/:subscriptionId/rules/:id
- **THEN** 更新指定规则
- **AND** 规则必须属于该订阅

#### Scenario: 删除规则
- **WHEN** 调用 DELETE /api/subscriptions/:subscriptionId/rules/:id
- **THEN** 删除指定规则
- **AND** 规则必须属于该订阅

### Requirement: 规则应用到订阅配置
配置合并时 SHALL 只使用激活订阅的规则。

#### Scenario: 合并配置
- **WHEN** Merger.Merge() 执行时
- **THEN** 只获取 `subscription_id = 激活订阅ID` 的规则
- **AND** 不再使用全局规则

## REMOVED Requirements

### Requirement: 全局规则 API
**Reason**: 规则已关联到订阅，不再支持全局规则

**Migration**: 
- 旧 API: GET /api/rules
- 新 API: GET /api/subscriptions/:id/rules