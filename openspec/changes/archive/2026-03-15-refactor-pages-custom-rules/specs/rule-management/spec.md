## ADDED Requirements

### Requirement: 规则管理 UI 位置
规则管理功能 SHALL 位于"订阅"页面的"规则"Tab 中。

#### Scenario: Tab 布局
- **WHEN** 用户访问"订阅"页面
- **THEN** 显示三个 Tab：订阅、规则、脚本
- **AND** "规则"Tab 展示自定义规则管理功能

#### Scenario: 功能完整
- **WHEN** 在"规则"Tab 操作规则
- **THEN** 支持所有原有功能：添加、编辑、删除、启用/禁用

## REMOVED Requirements

### Requirement: 独立规则页面
**Reason**: 规则管理迁移到"订阅"页面的 Tab 中
**Migration**: 访问 /subscriptions 页面的"规则"Tab