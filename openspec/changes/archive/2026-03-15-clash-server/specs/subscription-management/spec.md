## ADDED Requirements

### Requirement: 订阅接入
系统 SHALL 支持接入机场的 Clash 订阅 URL。

#### Scenario: 添加订阅
- **WHEN** 调用 POST /api/subscriptions，传入 name 和 url
- **THEN** 创建订阅记录并立即拉取一次数据
- **AND** 返回订阅详情包含节点数量

#### Scenario: 订阅 URL 格式
- **WHEN** 订阅 URL 是标准 Clash 订阅格式
- **THEN** 系统能够解析并提取代理节点和规则

### Requirement: 订阅本地缓存
系统 SHALL 本地化存储订阅内容，避免每次请求机场 API。

#### Scenario: 订阅内容存储
- **WHEN** 订阅首次拉取成功
- **THEN** 将原始 YAML/Base64 内容存入数据库
- **AND** 后续读取优先返回缓存

### Requirement: 订阅定时刷新
系统 SHALL 支持定时自动刷新订阅。

#### Scenario: 定时刷新
- **WHEN** 订阅配置了 interval（分钟）
- **THEN** 后台定时任务按间隔自动拉取新数据
- **AND** 更新数据库中的内容

#### Scenario: 手动刷新
- **WHEN** 调用 POST /api/subscriptions/:id/refresh
- **THEN** 立即拉取最新订阅内容
- **AND** 返回更新后的节点数量

### Requirement: 订阅激活状态
系统 SHALL 支持激活订阅，同一时间只能有一个订阅处于激活状态。

#### Scenario: 数据模型
- **WHEN** 存储订阅时
- **THEN** 订阅包含 `active` 布尔字段标识是否激活
- **AND** 同一时间只有一个订阅的 active 为 true

#### Scenario: 激活订阅
- **WHEN** 调用 PUT /api/subscriptions/:id/activate
- **THEN** 将该订阅的 active 设为 true
- **AND** 将其他所有订阅的 active 设为 false
- **AND** 合并订阅配置和面板配置
- **AND** 重启 Mihomo 核心
- **AND** 返回 Core 启动结果

#### Scenario: 创建首个订阅
- **WHEN** 系统中无订阅时创建第一个订阅
- **THEN** 自动将该订阅设为激活状态

#### Scenario: 删除激活订阅
- **WHEN** 删除当前激活的订阅
- **THEN** 自动激活剩余订阅中的第一个（如果有）
- **AND** 重启 Mihomo 核心

#### Scenario: 获取激活订阅
- **WHEN** 合并配置时
- **THEN** 使用当前激活的订阅配置
- **AND** 如果无激活订阅但有订阅存在，自动激活第一个

### Requirement: 订阅管理
系统 SHALL 支持订阅的增删改查。

#### Scenario: 订阅列表
- **WHEN** 调用 GET /api/subscriptions
- **THEN** 返回所有订阅的摘要（不含 content）
- **AND** 支持分页

#### Scenario: 订阅详情
- **WHEN** 调用 GET /api/subscriptions/:id
- **AND** 返回完整订阅信息包括节点数量

#### Scenario: 更新订阅
- **WHEN** 调用 PUT /api/subscriptions/:id
- **THEN** 更新订阅的 name、url、interval
- **AND** 触发一次刷新

#### Scenario: 删除订阅
- **WHEN** 调用 DELETE /api/subscriptions/:id
- **THEN** 删除订阅记录和缓存内容