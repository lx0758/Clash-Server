## MODIFIED Requirements

### Requirement: 脚本关联订阅
脚本 SHALL 关联到具体订阅，不再支持全局脚本。

#### Scenario: 创建脚本
- **WHEN** 调用 POST /api/subscriptions/:subscriptionId/scripts
- **THEN** 创建脚本并关联到指定订阅
- **AND** 脚本自动获取 `subscription_id`

#### Scenario: 脚本列表
- **WHEN** 调用 GET /api/subscriptions/:subscriptionId/scripts
- **THEN** 返回该订阅的所有脚本
- **AND** 不包含其他订阅的脚本

#### Scenario: 更新脚本
- **WHEN** 调用 PUT /api/subscriptions/:subscriptionId/scripts/:id
- **THEN** 更新指定脚本
- **AND** 脚本必须属于该订阅

#### Scenario: 删除脚本
- **WHEN** 调用 DELETE /api/subscriptions/:subscriptionId/scripts/:id
- **THEN** 删除指定脚本
- **AND** 脚本必须属于该订阅

#### Scenario: 测试脚本
- **WHEN** 调用 POST /api/subscriptions/:subscriptionId/scripts/:id/test
- **THEN** 执行脚本测试
- **AND** 脚本必须属于该订阅

### Requirement: 脚本执行范围
配置合并时 SHALL 只执行激活订阅的脚本。

#### Scenario: 合并配置
- **WHEN** Merger.Merge() 执行时
- **THEN** 只执行 `subscription_id = 激活订阅ID` 且 `enabled = true` 的脚本
- **AND** 不再执行全局脚本

## REMOVED Requirements

### Requirement: 全局脚本 API
**Reason**: 脚本已关联到订阅，不再支持全局脚本

**Migration**:
- 旧 API: GET /api/scripts
- 新 API: GET /api/subscriptions/:id/scripts