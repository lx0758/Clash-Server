## ADDED Requirements

### Requirement: 脚本管理 UI 位置
脚本管理功能 SHALL 位于"订阅"页面的"脚本"Tab 中。

#### Scenario: Tab 布局
- **WHEN** 用户访问"订阅"页面
- **THEN** 显示三个 Tab：订阅、规则、脚本
- **AND** "脚本"Tab 展示扩展脚本管理功能

#### Scenario: 功能完整
- **WHEN** 在"脚本"Tab 操作脚本
- **THEN** 支持所有原有功能：添加、编辑、删除、启用/禁用、测试

## REMOVED Requirements

### Requirement: 独立脚本页面
**Reason**: 脚本管理迁移到"订阅"页面的 Tab 中
**Migration**: 访问 /subscriptions 页面的"脚本"Tab