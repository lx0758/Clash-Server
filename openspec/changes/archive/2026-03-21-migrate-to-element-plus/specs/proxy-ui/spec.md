## MODIFIED Requirements

### Requirement: 多列网格布局
系统 SHALL 以网格形式展示代理节点，并支持响应式列数。

#### Scenario: 响应式列数
- **WHEN** 用户调整浏览器窗口大小
- **THEN** 节点网格自动调整列数
- **AND** 宽度 ≥ 1200px 时 4 列
- **AND** 宽度 ≥ 992px 时 3 列
- **AND** 宽度 ≥ 768px 时 2 列
- **AND** 宽度 < 768px 时 1 列

#### Scenario: 节点卡片展示
- **WHEN** 节点在网格中展示
- **THEN** 显示节点名称、延迟
- **AND** 选中的节点左边框高亮
- **AND** 使用 Element Plus el-card 样式

## ADDED Requirements

### Requirement: 代理页面响应式
代理页面 SHALL 支持响应式布局。

#### Scenario: 移动端代理组列表
- **WHEN** 屏幕宽度 < 768px
- **THEN** 代理组全宽显示
- **AND** 筛选和排序控件堆叠排列

#### Scenario: 移动端节点选择
- **WHEN** 用户在移动设备上选择代理节点
- **THEN** 点击区域足够大
- **AND** 提供明确的选中反馈