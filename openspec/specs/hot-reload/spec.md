# hot-reload 规格说明

## Purpose

定义配置变更在 Mihomo 核心上的应用方式：热重载与软重启的选择依据、失败降级路径，以及合并失败时的处理约定。

## Requirements

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

#### Scenario: 热重载失败时的错误信息透传
- **WHEN** 热重载 API 调用失败且核心返回了错误响应
- **THEN** 返回的错误信息包含核心返回的原始错误内容
- **AND** 不得以泛化文案替代核心的具体错误

### Requirement: 智能判断配置变更类型
系统 SHALL 自动判断配置变更是否需要软重启。

#### Scenario: 需要软重启的配置变更
- **WHEN** 变更涉及 mixed-port, external-controller, secret, ipv6 等配置项
- **THEN** 调用 `SoftRestart()` 方法，通过 POST /restart API 软重启

#### Scenario: 可热重载的配置变更
- **WHEN** 变更仅涉及 proxies, rules, dns 等配置项
- **THEN** 调用 `HotReload()` 方法热重载

#### Scenario: 判断基于变更前后的配置差异
- **WHEN** 判断配置变更是否需要软重启
- **THEN** 比较变更前的已应用配置与本次生成的配置
- **AND** 不得以本次配置与自身比较作为判断依据

#### Scenario: 端口类变更确实触发重启
- **WHEN** 用户修改了 mixed-port、secret、allow-lan 等需重启生效的配置项
- **THEN** 系统执行软重启而非热重载
- **AND** 该变更在核心侧实际生效

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

#### Scenario: 降级重启也失败
- **WHEN** 软重启失败且降级为杀进程重启后核心仍无法启动
- **THEN** 返回错误信息说明核心未能启动
- **AND** 该错误对调用方可见，不得以成功响应掩盖

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

### Requirement: 配置合并失败处理
系统 SHALL 在配置合并失败时显式暴露错误，不得静默以其他配置继续下发。

#### Scenario: 合并失败
- **WHEN** 合并订阅与自定义配置的过程失败
- **THEN** 返回错误信息说明合并失败的原因
- **AND** 不以最小配置或其他替代配置静默下发

#### Scenario: 合并失败不改变已持久化状态
- **WHEN** 合并失败
- **THEN** 此前已持久化的配置保持原样
- **AND** 核心继续运行其原有配置
