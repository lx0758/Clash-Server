## ADDED Requirements

### Requirement: 自定义规则管理
系统 SHALL 支持用户手动添加、编辑、删除自定义规则。

#### Scenario: 添加规则
- **WHEN** 调用 POST /api/rules，传入 type, payload, proxy, enabled, mode, priority
- **THEN** 创建规则记录

#### Scenario: 规则类型
- **WHEN** 创建规则时
- **THEN** 支持 Clash 标准规则类型: DOMAIN, DOMAIN-SUFFIX, DOMAIN-KEYWORD, IP-CIDR, IP-CIDR6, SRC-IP-CIDR, GEOIP, DST-PORT, SRC-PORT, PROCESS-NAME, RULE-SET, MATCH

#### Scenario: 更新规则
- **WHEN** 调用 PUT /api/rules/:id
- **THEN** 更新规则的所有属性
- **AND** 需要重新加载配置生效

#### Scenario: 删除规则
- **WHEN** 调用 DELETE /api/rules/:id
- **THEN** 删除规则记录
- **AND** 需要重新加载配置生效

### Requirement: 规则插入/追加模式
系统 SHALL 支持两种规则应用模式。

#### Scenario: 插入模式 (insert)
- **WHEN** 规则的 mode 为 insert
- **THEN** 该规则插入到订阅规则之前
- **AND** 优先于订阅规则匹配

#### Scenario: 追加模式 (append)
- **WHEN** 规则的 mode 为 append
- **THEN** 该规则追加到订阅规则之后
- **AND** 作为兜底规则

#### Scenario: 同模式内排序
- **WHEN** 同一 mode 下有多条规则
- **THEN** 按 priority 从小到大排序

### Requirement: 规则优先级
系统 SHALL 支持规则按优先级排序。

#### Scenario: 规则排序
- **WHEN** 合并配置时
- **THEN** insert 规则按 priority 排序后放在订阅规则之前
- **AND** append 规则按 priority 排序后放在订阅规则之后

### Requirement: 规则启用/禁用
系统 SHALL 支持单独启用或禁用规则。

#### Scenario: 禁用规则
- **WHEN** 规则的 enabled 为 false
- **THEN** 合并配置时跳过该规则