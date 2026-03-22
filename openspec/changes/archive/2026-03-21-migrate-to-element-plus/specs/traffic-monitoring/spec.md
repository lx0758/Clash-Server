## ADDED Requirements

### Requirement: 仪表盘卡片响应式
仪表盘卡片 SHALL 支持响应式布局。

#### Scenario: 桌面端仪表盘
- **WHEN** 屏幕宽度 ≥ 1200px
- **THEN** 状态卡片以 4 列网格排列
- **AND** 使用 el-row/el-col 栅格
- **AND** 每个卡片使用 el-card 组件

#### Scenario: 平板仪表盘
- **WHEN** 屏幕宽度在 768px-1199px
- **THEN** 状态卡片以 2 列网格排列

#### Scenario: 移动端仪表盘
- **WHEN** 屏幕宽度 < 768px
- **THEN** 状态卡片单列排列
- **AND** 卡片全宽

### Requirement: 流量显示组件
流量数据 SHALL 使用 Element Plus 组件展示。

#### Scenario: 流量卡片展示
- **WHEN** 显示流量数据
- **THEN** 使用 el-statistic 组件展示数值
- **AND** 使用 el-progress 组件展示进度
- **AND** 数字格式化 (KB/MB/GB)

### Requirement: 核心状态显示
核心状态 SHALL 使用 Element Plus 标签组件显示。

#### Scenario: 运行状态标签
- **WHEN** 核心运行中
- **THEN** 使用 el-tag 显示绿色 "运行中"

#### Scenario: 停止状态标签
- **WHEN** 核心已停止
- **THEN** 使用 el-tag 显示红色 "已停止"

#### Scenario: 错误状态
- **WHEN** 核心出错
- **THEN** 使用 el-alert 显示错误信息