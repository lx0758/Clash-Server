## ADDED Requirements

### Requirement: 订阅页面 Tab 布局
订阅页面 SHALL 采用 Tab 布局整合订阅、规则、脚本三个功能。

#### Scenario: Tab 导航
- **WHEN** 用户访问"订阅"页面
- **THEN** 显示三个 Tab：订阅、规则、脚本
- **AND** 默认显示"订阅"Tab

#### Scenario: Tab 切换
- **WHEN** 用户点击不同 Tab
- **THEN** 切换显示对应内容
- **AND** 保持各 Tab 的状态（如编辑中的表单）