## ADDED Requirements

### Requirement: 运行模式切换
系统 SHALL 支持三种运行模式的切换。

#### Scenario: 显示模式切换器
- **WHEN** 用户进入代理页面
- **THEN** 页面顶部显示模式切换器
- **AND** 包含"规则"、"全局"、"直连"三个选项

#### Scenario: 规则模式
- **WHEN** 用户选择"规则"模式
- **THEN** 调用 `PATCH /api/proxies/mode` 设置模式为 rule
- **AND** 显示非预置代理组列表（排除 GLOBAL）
- **AND** 代理组按 Selector > URLTest > Fallback 顺序排列

#### Scenario: 全局模式
- **WHEN** 用户选择"全局"模式
- **THEN** 调用 `PATCH /api/proxies/mode` 设置模式为 global
- **AND** 显示所有代理（包含 DIRECT、REJECT 等预置代理）
- **AND** 用户选择的代理作为全局代理

#### Scenario: 直连模式
- **WHEN** 用户选择"直连"模式
- **THEN** 调用 `PATCH /api/proxies/mode` 设置模式为 direct
- **AND** 显示提示信息"当前为直连模式"

### Requirement: 可折叠代理组
系统 SHALL 支持代理组的展开和收起功能。

#### Scenario: 点击代理组头切换展开状态
- **WHEN** 用户点击代理组头
- **THEN** 代理组在展开和收起状态之间切换
- **AND** 展开状态保存到 localStorage

#### Scenario: 默认展开状态
- **WHEN** 用户首次访问代理页面
- **THEN** 第一个代理组默认展开
- **AND** 其他代理组默认收起

### Requirement: 多列网格布局
系统 SHALL 以网格形式展示代理节点。

#### Scenario: 响应式列数
- **WHEN** 用户调整浏览器窗口大小
- **THEN** 节点网格自动调整列数
- **AND** 宽度 > 1400px 时 4 列
- **AND** 宽度 > 1000px 时 3 列
- **AND** 宽度 > 700px 时 2 列
- **AND** 宽度 <= 700px 时 1 列

#### Scenario: 节点卡片展示
- **WHEN** 节点在网格中展示
- **THEN** 显示节点名称、延迟
- **AND** 选中的节点左边框高亮

### Requirement: 筛选功能
系统 SHALL 支持按名称筛选代理节点。

#### Scenario: 输入筛选条件
- **WHEN** 用户在筛选框输入文本
- **THEN** 当前代理组仅显示名称包含该文本的节点
- **AND** 筛选不区分大小写

### Requirement: 排序功能
系统 SHALL 支持代理节点排序。

#### Scenario: 按延迟升序排序
- **WHEN** 用户选择"延迟升序"
- **THEN** 节点按延迟从低到高排列

#### Scenario: 按延迟降序排序
- **WHEN** 用户选择"延迟降序"
- **THEN** 节点按延迟从高到低排列

#### Scenario: 按名称排序
- **WHEN** 用户选择"按名称"
- **THEN** 节点按名称字母顺序排列

### Requirement: 定位当前节点
系统 SHALL 支持定位到当前选中的节点。

#### Scenario: 点击定位按钮
- **WHEN** 用户点击"定位当前"按钮
- **THEN** 页面滚动到当前选中的节点位置
- **AND** 节点高亮显示

### Requirement: 测试延迟
系统 SHALL 支持测试节点延迟。

#### Scenario: 测试单个节点延迟
- **WHEN** 用户点击节点的延迟区域
- **THEN** 调用 API 测试该节点延迟
- **AND** 更新显示的延迟值

#### Scenario: 测试组内所有节点延迟
- **WHEN** 用户点击"测全部延迟"按钮
- **THEN** 调用 API 测试组内所有节点延迟
- **AND** 显示加载状态
- **AND** 完成后更新所有延迟值

### Requirement: 选中代理不刷新
系统 SHALL 在选中代理后不重新拉取数据。

#### Scenario: 选中代理
- **WHEN** 用户点击选择一个代理
- **THEN** 调用 API 切换代理
- **AND** 仅更新本地状态中的 `now` 值
- **AND** 页面不重新渲染或闪烁

### Requirement: 代理组排序
系统 SHALL 在规则模式下按类型排序代理组。

#### Scenario: 显示排序后的代理组
- **WHEN** 用户处于规则模式
- **THEN** Selector 类型代理组排在第一位
- **AND** URLTest 类型代理组排在第二位
- **AND** Fallback 类型代理组排在第三位