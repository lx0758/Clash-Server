# Spec Delta

## MODIFIED Requirements

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
- **AND** 使用该订阅的规则和脚本合并配置
- **AND** 重启 Mihomo 核心
- **AND** 返回 Core 启动结果

#### Scenario: 激活后核心未生效
- **WHEN** 订阅激活已持久化但核心重载失败
- **THEN** 返回 code 为 1007 且携带 partial_success 标记
- **AND** 界面以警告语义提示"已激活但未在核心生效"，而非提示失败

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

### Requirement: 订阅刷新结果反馈
刷新订阅 SHALL 返回详细的结果信息。

#### Scenario: 刷新成功
- **WHEN** 调用 POST /api/subscriptions/:id/refresh 成功
- **THEN** 返回更新后的订阅信息
- **AND** 返回 `node_count` 节点数量

#### Scenario: 刷新失败
- **WHEN** 调用 POST /api/subscriptions/:id/refresh 失败
- **THEN** 返回 `refresh_error` 错误信息
- **AND** 保持原有订阅内容不变

#### Scenario: 刷新成功但核心未生效
- **WHEN** 订阅内容已更新但核心重载失败
- **THEN** 返回 code 为 1007 且携带 partial_success 标记
- **AND** 返回体包含 core_error 说明未生效原因

## ADDED Requirements

### Requirement: 订阅操作的对象级反馈
系统 SHALL 将刷新与激活操作的进行中状态限制在被操作的订阅卡片上。

#### Scenario: 刷新期间的状态表达
- **WHEN** 用户刷新某个订阅
- **THEN** 该订阅卡片显示进行中状态
- **AND** 列表中其他订阅保持可交互

#### Scenario: 激活期间的状态表达
- **WHEN** 用户激活某个订阅
- **THEN** 该订阅的操作入口显示进行中状态
- **AND** 列表中其他订阅保持可交互

#### Scenario: 同一订阅的操作互斥
- **WHEN** 某个订阅正在刷新
- **THEN** 该订阅的其他操作不被重复提交
- **AND** 该订阅恢复就绪后操作重新可用

#### Scenario: 核心错误提示可见
- **WHEN** 订阅操作导致核心重载失败
- **THEN** 界面明确显示未生效的错误信息
- **AND** 该提示可被用户关闭
