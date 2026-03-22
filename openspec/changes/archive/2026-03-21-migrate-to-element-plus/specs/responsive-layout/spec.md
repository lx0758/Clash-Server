## ADDED Requirements

### Requirement: 断点检测
系统 SHALL 提供响应式断点检测能力。

#### Scenario: 检测当前断点
- **WHEN** 组件需要响应式适配
- **THEN** 可通过 `useBreakpoint` hook 获取当前断点
- **AND** 断点类型包括 xs/sm/md/lg/xl

#### Scenario: 断点值定义
- **WHEN** 检测屏幕宽度
- **THEN** xs < 768px
- **AND** sm ≥ 768px
- **AND** md ≥ 992px
- **AND** lg ≥ 1200px
- **AND** xl ≥ 1920px

### Requirement: 导航响应式
系统 SHALL 支持导航的响应式布局。

#### Scenario: 桌面端导航
- **WHEN** 屏幕宽度 ≥ 992px
- **THEN** 显示固定侧边栏导航
- **AND** 侧边栏可通过按钮折叠为图标模式

#### Scenario: 移动端导航
- **WHEN** 屏幕宽度 < 992px
- **THEN** 侧边栏隐藏
- **AND** 显示汉堡菜单按钮
- **AND** 点击按钮打开抽屉导航

#### Scenario: 导航折叠状态记忆
- **WHEN** 用户折叠或展开导航
- **THEN** 状态保存到 localStorage
- **AND** 刷新页面后保持该状态

### Requirement: 表单响应式
系统 SHALL 支持表单的响应式布局。

#### Scenario: 桌面端表单标签
- **WHEN** 屏幕宽度 ≥ 992px
- **AND** 表单标签位置为右侧对齐

#### Scenario: 移动端表单标签
- **WHEN** 屏幕宽度 < 992px
- **THEN** 表单标签位置为顶部对齐
- **AND** 输入框宽度 100%

#### Scenario: 表单验证响应式
- **WHEN** 表单验证失败
- **THEN** 错误提示在所有设备上清晰可见

### Requirement: 对话框响应式
系统 SHALL 支持对话框的响应式布局。

#### Scenario: 桌面端对话框
- **WHEN** 屏幕宽度 ≥ 992px
- **THEN** 对话框居中显示
- **AND** 宽度为固定值或自适应

#### Scenario: 移动端对话框
- **WHEN** 屏幕宽度 < 992px
- **THEN** 对话框全屏显示
- **AND** 关闭按钮在顶部

### Requirement: 表格响应式
系统 SHALL 支持表格的响应式布局。

#### Scenario: 桌面端表格
- **WHEN** 屏幕宽度 ≥ 992px
- **THEN** 显示所有列
- **AND** 表格宽度自适应

#### Scenario: 移动端表格
- **WHEN** 屏幕宽度 < 992px
- **THEN** 表格支持横向滚动
- **AND** 或隐藏次要列
- **AND** 关键信息优先显示

### Requirement: 卡片网格响应式
系统 SHALL 支持卡片网格的响应式布局。

#### Scenario: 桌面端卡片网格
- **WHEN** 屏幕宽度 ≥ 1200px
- **THEN** 卡片网格显示 4 列

#### Scenario: 平板卡片网格
- **WHEN** 屏幕宽度在 768px-1199px
- **THEN** 卡片网格显示 2-3 列

#### Scenario: 移动端卡片网格
- **WHEN** 屏幕宽度 < 768px
- **THEN** 卡片网格显示 1 列
- **AND** 卡片宽度 100%

### Requirement: 移动端触摸优化
系统 SHALL 优化移动端触摸交互。

#### Scenario: 按钮触摸区域
- **WHEN** 在移动设备上
- **THEN** 按钮最小高度 44px
- **AND** 最小宽度 44px

#### Scenario: 列表项触摸
- **WHEN** 在移动设备上点击列表项
- **THEN** 触摸区域足够大
- **AND** 提供视觉反馈