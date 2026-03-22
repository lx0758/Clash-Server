## ADDED Requirements

### Requirement: 配置热重载
系统 SHALL 支持通过 Mihomo API 热重载配置，无需重启进程。

#### Scenario: 热重载成功
- **WHEN** 配置变更且无需重启（proxies, rules 等）
- **THEN** 调用 `HotReload()` 方法，通过 Mihomo PUT /configs API 热重载
- **AND** 返回 hot_reloaded: true
- **AND** 现有连接保持不中断

#### Scenario: 热重载失败需用户确认
- **WHEN** 热重载 API 调用失败
- **THEN** 返回 need_confirm: true 和错误信息
- **AND** 内核保持旧配置运行

### Requirement: 智能判断配置变更类型
系统 SHALL 自动判断配置变更是否需要软重启。

#### Scenario: 需要软重启的配置变更
- **WHEN** 变更涉及 mixed-port, external-controller, secret, ipv6 等配置项
- **THEN** 调用 `SoftRestart()` 方法，通过 POST /restart API 软重启

#### Scenario: 可热重载的配置变更
- **WHEN** 变更仅涉及 proxies, rules, dns 等配置项
- **THEN** 调用 `HotReload()` 方法热重载

### Requirement: 软重启
系统 SHALL 通过 Mihomo POST /restart API 实现软重启，而非杀进程。

#### Scenario: 软重启成功
- **WHEN** 调用 POST /restart API
- **THEN** Mihomo 内核自行重启
- **AND** 返回 restarted: true

#### Scenario: 软重启失败降级
- **WHEN** 软重启 API 失败
- **THEN** 降级为杀进程重启
- **AND** 记录错误日志

### Requirement: 启动流程重构
系统 SHALL 在 Mihomo 启动时通过 `-config` 参数传入 Base64 编码的最小配置，启动后通过 API 应用全量配置。

#### Scenario: Base64 最小配置启动
- **WHEN** 启动 Mihomo
- **THEN** 从 CoreConfig 获取 APIHost 和 APIPort
- **AND** 生成最小配置 `external-controller: {APIHost}:{APIPort}`
- **AND** Base64 编码后通过 `-config` 参数传入
- **AND** 不生成配置文件

#### Scenario: 启动后应用配置
- **WHEN** Mihomo API 就绪
- **THEN** 调用 PUT /configs 应用全量配置

#### Scenario: 启动时无订阅
- **WHEN** 启动时没有任何订阅
- **THEN** 应用最小运行配置（mode: direct）

### Requirement: 移除临时配置文件
系统 SHALL 不再使用临时配置文件管理 Mihomo 配置。

#### Scenario: 配置变更不生成临时文件
- **WHEN** 配置变更
- **THEN** 直接通过 API 传递配置 YAML
- **AND** 不生成临时配置文件