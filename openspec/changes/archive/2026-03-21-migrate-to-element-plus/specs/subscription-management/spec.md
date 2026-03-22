## ADDED Requirements

### Requirement: 订阅卡片响应式
订阅卡片 SHALL 支持响应式布局。

#### Scenario: 桌面端订阅列表
- **WHEN** 屏幕宽度 ≥ 992px
- **THEN** 订阅卡片以网格形式排列
- **AND** 使用 el-row/el-col 栅格
- **AND** 每行 2-3 个卡片

#### Scenario: 移动端订阅列表
- **WHEN** 屏幕宽度 < 768px
- **THEN** 订阅卡片单列排列
- **AND** 卡片全宽

#### Scenario: 订阅卡片组件
- **WHEN** 渲染订阅卡片
- **THEN** 使用 el-card 组件
- **AND** 使用 el-progress 显示流量进度
- **AND** 使用 el-tag 显示订阅状态

### Requirement: 订阅编辑对话框响应式
订阅编辑对话框 SHALL 支持响应式布局。

#### Scenario: 桌面端编辑对话框
- **WHEN** 屏幕宽度 ≥ 992px
- **THEN** 对话框宽度 500px
- **AND** 居中显示

#### Scenario: 移动端编辑对话框
- **WHEN** 屏幕宽度 < 992px
- **THEN** 对话框全屏显示
- **AND** 表单标签顶部对齐

### Requirement: 订阅表单组件
订阅表单 SHALL 使用 Element Plus 组件。

#### Scenario: 表单组件使用
- **WHEN** 渲染订阅表单
- **THEN** 使用 el-form 组件
- **AND** 使用 el-input 组件
- **AND** 使用 el-select 组件
- **AND** 使用 el-switch 组件
- **AND** 使用 Element Plus 表单验证