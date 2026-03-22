## ADDED Requirements

### Requirement: 代理模式切换
系统 SHALL 在仪表盘页面提供代理模式切换功能。

#### Scenario: 三种代理模式显示
- **WHEN** 用户访问仪表盘页面
- **THEN** 系统显示 3 个模式按钮：Rule（规则模式）、Global（全局模式）、Direct（直连模式）
- **AND** 每个按钮包含图标和文字标签
- **AND** Rule 按钮使用 MultipleStop 图标
- **AND** Global 按钮使用 Language 图标
- **AND** Direct 按钮使用 Directions 图标

#### Scenario: 当前模式高亮显示
- **WHEN** 代理模式显示
- **THEN** 系统高亮显示当前激活的模式按钮
- **AND** 激活模式按钮使用主题色（primary.main）背景
- **AND** 激活模式按钮文字加粗
- **AND** 激活模式按钮底部显示指示线

#### Scenario: 代理模式切换操作
- **WHEN** 用户点击非激活的模式按钮
- **THEN** 系统调用代理 API 切换模式
- **AND** 切换成功后更新按钮高亮状态
- **AND** 切换过程中显示加载状态（按钮禁用、显示旋转图标）
- **AND** 切换失败时显示错误提示并恢复原模式

#### Scenario: 代理模式切换动画
- **WHEN** 用户点击模式按钮
- **THEN** 按钮显示悬停效果（向上平移 1px，显示阴影）
- **AND** 按钮显示按下效果（向下平移 1px）
- **AND** 过渡动画时长为 200ms（ease-in-out）

#### Scenario: 移动端代理模式显示
- **WHEN** 屏幕宽度小于 768px（xs/sm 断点）
- **THEN** 系统保持 3 个模式按钮的完整显示
- **AND** 按钮宽度自适应（平均分配可用空间）
- **AND** 按钮文字可能改为简化版本或图标-only（如果空间不足）

#### Scenario: 代理模式说明文字
- **WHEN** 代理模式显示
- **THEN** 系统在模式按钮下方显示当前模式的说明文字
- **AND** Rule 模式说明："根据规则自动选择代理"
- **AND** Global 模式说明："所有流量使用同一代理"
- **AND** Direct 模式说明："直连网络，不使用代理"
- **AND** 说明文字居中显示，使用小号字体（caption）
- **AND** 说明文字外框与激活模式按钮颜色一致

### Requirement: 代理模式配置
系统 SHALL 根据用户配置切换代理模式时的行为。

#### Scenario: 切换模式时关闭现有连接
- **WHEN** 用户配置了"切换模式时关闭连接"选项
- **AND** 用户切换代理模式
- **THEN** 系统在切换模式前关闭所有现有连接
- **AND** 然后调用代理 API 切换模式

#### Scenario: 切换模式时保持连接
- **WHEN** 用户未配置"切换模式时关闭连接"选项
- **AND** 用户切换代理模式
- **THEN** 系统直接调用代理 API 切换模式
- **AND** 不关闭现有连接（让它们自然过期）

### Requirement: 代理模式状态同步
系统 SHALL 确保代理模式状态在多个组件间同步。

#### Scenario: 代理模式状态管理
- **WHEN** 代理模式状态变化
- **THEN** 系统通过 Pinia store 代理模式状态（proxyStore.mode）
- **AND** 所有使用代理模式的组件自动获取最新状态
- **AND** 代理模式状态持久化到 localStorage（如果需要）

#### Scenario: 代理模式初始化
- **WHEN** 用户访问仪表盘页面
- **THEN** 系统从 API 获取当前代理模式
- **AND** 更新 proxyStore.mode 状态
- **AND** 高亮显示对应的模式按钮

#### Scenario: 代理模式变化通知
- **WHEN** 代理模式通过其他页面（如代理页面）切换
- **THEN** 系统通过 store 通知仪表盘页面
- **AND** 仪表盘页面的代理模式卡片自动更新高亮状态

### Requirement: 代理模式错误处理
系统 SHALL 处理代理模式切换过程中的错误。

#### Scenario: API 调用失败
- **WHEN** 代理模式切换 API 调用失败
- **THEN** 系统显示错误提示（使用 el-message）
- **AND** 错误提示内容："切换模式失败：[错误信息]"
- **AND** 恢复原模式按钮的高亮状态
- **AND** 记录错误日志到控制台

#### Scenario: 网络连接问题
- **WHEN** 检测到网络连接问题（WebSocket 断开、API 超时）
- **THEN** 系统禁用代理模式切换按钮
- **AND** 按钮显示警告图标和提示文本
- **AND** 网络恢复后自动重新启用按钮

#### Scenario: 核心（Clash）未运行
- **WHEN** 检测到 Clash 核心未运行
- **THEN** 系统禁用代理模式切换按钮
- **AND** 按钮显示灰色的禁用状态
- **AND** 提示用户先启动核心
