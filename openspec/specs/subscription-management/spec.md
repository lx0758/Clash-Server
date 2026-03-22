# Subscription Management Specification

## Requirements

### Requirement: 订阅来源类型
订阅 SHALL 支持远程和本地两种来源类型。

#### Scenario: 远程订阅
- **WHEN** 创建订阅时指定 `source_type = "remote"`
- **THEN** 必须提供 `url` 字段
- **AND** 可选提供 `interval` 字段
- **AND** 支持定时刷新和手动刷新

#### Scenario: 本地配置
- **WHEN** 创建订阅时指定 `source_type = "local"`
- **THEN** `url` 和 `interval` 字段为空
- **AND** 用户直接编辑 `content` 字段
- **AND** 不支持刷新操作

#### Scenario: 订阅列表展示
- **WHEN** 获取订阅列表时
- **THEN** 显示 `source_type` 字段
- **AND** 远程订阅显示 URL 和刷新状态
- **AND** 本地配置显示"本地"标签

### Requirement: 订阅详情包含规则和脚本
获取订阅详情 SHALL 返回关联的规则和脚本。

#### Scenario: 获取详情
- **WHEN** 调用 GET /api/subscriptions/:id
- **THEN** 返回订阅基本信息
- **AND** 返回 `rules` 数组 (该订阅的所有规则)
- **AND** 返回 `scripts` 数组 (该订阅的所有脚本)

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
- **AND** 使用该订阅的规则和脚本合并配置
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

### Requirement: 删除订阅级联删除规则和脚本
删除订阅 SHALL 同时删除关联的规则和脚本。

#### Scenario: 级联删除
- **WHEN** 调用 DELETE /api/subscriptions/:id
- **THEN** 删除订阅记录
- **AND** 删除该订阅关联的所有规则
- **AND** 删除该订阅关联的所有脚本

### Requirement: 代理更新订阅
远程订阅 SHALL 支持通过代理服务器更新。

#### Scenario: 启用代理更新
- **WHEN** 创建或更新订阅时设置 `use_proxy = true`
- **THEN** 刷新订阅时通过 127.0.0.1:7890 代理请求

### Requirement: 自定义 User-Agent
远程订阅 SHALL 支持自定义请求 User-Agent。

#### Scenario: 自定义 UA
- **WHEN** 创建或更新订阅时提供 `user_agent` 字段
- **THEN** 刷新订阅时使用该 User-Agent
- **AND** 为空时使用默认 UA

### Requirement: 证书校验控制
远程订阅 SHALL 支持跳过 SSL 证书校验。

#### Scenario: 跳过证书校验
- **WHEN** 创建或更新订阅时设置 `skip_cert = true`
- **THEN** 刷新订阅时跳过 SSL 证书校验
- **AND** 用于解决某些机场的证书问题

### Requirement: 查看合并后配置
系统 SHALL 支持查看订阅合并后的最终配置。

#### Scenario: 获取合并配置
- **WHEN** 调用 GET /api/subscriptions/:id/merged-config
- **THEN** 返回该订阅的合并后配置 (JSON 和 YAML 格式)
- **AND** 包含订阅基础配置 + 规则 + 脚本处理结果

### Requirement: 直接编辑订阅内容
订阅 SHALL 支持直接编辑订阅内容 (主要用于本地配置)。

#### Scenario: 更新内容
- **WHEN** 调用 PUT /api/subscriptions/:id/content
- **THEN** 直接更新订阅的 `content` 字段
- **AND** 如果是激活订阅，触发 Core 重启

### Requirement: 订阅元数据展示
订阅列表 SHALL 展示完整的元数据信息。

#### Scenario: 时间信息
- **WHEN** 获取订阅列表或详情
- **THEN** 显示 `created_at` 创建时间
- **AND** 远程订阅显示 `last_refresh` 上次刷新时间

#### Scenario: 节点数量
- **WHEN** 获取订阅列表或详情
- **THEN** 显示 `node_count` 节点数量
- **AND** 每次刷新后更新

### Requirement: 订阅流量信息
远程订阅 SHALL 展示机场提供的流量信息。

#### Scenario: 解析 subscription-userinfo
- **WHEN** 刷新远程订阅时
- **AND** 响应头包含 `subscription-userinfo`
- **THEN** 解析并存储 `upload_used`、`download_used`、`total_transfer`、`expire_at`
- **AND** 更新到订阅记录

#### Scenario: 流量信息格式
- **WHEN** 响应头为 `subscription-userinfo: upload=240935845; download=15758204665; total=1099511627776; expire=1796263142`
- **THEN** `upload_used` = 240935845 (字节)
- **AND** `download_used` = 15758204665 (字节)
- **AND** `total_transfer` = 1099511627776 (字节)
- **AND** `expire_at` = 2026-12-31 对应的日期时间

#### Scenario: 无流量信息
- **WHEN** 响应头不包含 `subscription-userinfo`
- **THEN** 流量相关字段保持为空或 0
- **AND** UI 不显示流量信息

#### Scenario: 流量进度展示
- **WHEN** 订阅有流量信息 (`total_transfer > 0`)
- **THEN** 显示已用流量 / 总流量
- **AND** 显示流量使用百分比
- **AND** 显示流量进度条
- **AND** 显示上传和下载分项

#### Scenario: 过期时间展示
- **WHEN** 订阅有过期时间 (`expire_at` 不为空)
- **THEN** 显示订阅过期日期
- **AND** 显示剩余天数

#### Scenario: 即将过期警告
- **WHEN** 订阅剩余天数 <= 7 天
- **THEN** 过期时间显示为警告颜色
- **AND** 提示用户续费

#### Scenario: 已过期提示
- **WHEN** 订阅已过期
- **THEN** 过期时间显示为错误颜色
- **AND** 提示用户订阅已过期

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