## ADDED Requirements

### Requirement: Core 规则列表展示
系统 SHALL 在"规则"页面展示 Clash Core 实时生效的规则列表。

#### Scenario: 获取规则列表
- **WHEN** 用户访问"规则"页面
- **THEN** 调用 GET /api/core/rules 获取 Core 规则
- **AND** 展示规则列表

#### Scenario: 规则数据字段
- **WHEN** 展示规则时
- **THEN** 显示以下字段: index, type, payload, proxy
- **AND** 可选显示 extra.hitCount, extra.hitAt

#### Scenario: Core 未运行
- **WHEN** Core 未运行时访问规则页面
- **THEN** 显示"核心未运行"提示

### Requirement: 规则列表性能优化
系统 SHALL 支持大量规则的高效渲染。

#### Scenario: 虚拟滚动
- **WHEN** 规则数量超过 100 条
- **THEN** 使用虚拟滚动仅渲染可视区域

#### Scenario: 规则搜索
- **WHEN** 用户输入搜索关键词
- **THEN** 实时过滤匹配的规则

### Requirement: 规则详情查看
系统 SHALL 支持查看单条规则的详细信息。

#### Scenario: 查看详情
- **WHEN** 用户点击某条规则
- **THEN** 展开显示完整信息包括 extra 字段