## ADDED Requirements

### Requirement: 登录页面响应式
登录页面 SHALL 支持响应式布局。

#### Scenario: 桌面端登录
- **WHEN** 屏幕宽度 ≥ 992px
- **THEN** 登录表单居中显示
- **AND** 使用 el-card 包裹表单
- **AND** 表单宽度固定 400px

#### Scenario: 移动端登录
- **WHEN** 屏幕宽度 < 992px
- **THEN** 登录表单全宽显示
- **AND** 表单标签在输入框上方
- **AND** 卡片无边距

#### Scenario: 登录表单组件
- **WHEN** 渲染登录表单
- **THEN** 使用 el-form 组件
- **AND** 使用 el-input 组件
- **AND** 使用 el-button 组件
- **AND** 使用 Element Plus 表单验证

### Requirement: 登录错误提示
系统 SHALL 使用 Element Plus 消息组件显示登录错误。

#### Scenario: 登录失败提示
- **WHEN** 登录凭证错误
- **THEN** 使用 ElMessage.error 显示错误信息
- **AND** 消息自动消失