## MODIFIED Requirements

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

## REMOVED Requirements

### Requirement: 配置变更自动重启
**Reason**: 替换为智能应用配置（热重载优先，软重启兜底）
**Migration**: 使用 SmartApply 逻辑，通过 Mihomo API 管理配置

### Requirement: 使用临时配置文件启动
**Reason**: 改为最小配置启动 + API 应用全量配置
**Migration**: 启动时不指定 -f 参数，启动后调用 PUT /configs