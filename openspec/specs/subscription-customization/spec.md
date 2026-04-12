## ADDED Requirements

### Requirement: 订阅自定义存储

系统 SHALL 将订阅自定义配置存储在专用的 `SubscriptionCustomization` 表中，包含以下字段：
- `subscription_id`: 关联订阅（唯一，1:1）
- `proxy_insert`, `proxy_append`, `proxy_remove`: 代理自定义 YAML
- `proxy_group_insert`, `proxy_group_append`, `proxy_group_remove`: 代理组自定义 YAML
- `rule_insert`, `rule_append`, `rule_remove`: 规则自定义 YAML
- `global_override`: 全局配置覆盖 YAML
- `script`: 后处理脚本

#### Scenario: 首次保存时创建自定义配置
- **WHEN** 用户为没有自定义配置的订阅保存自定义配置
- **THEN** 系统创建新的 `SubscriptionCustomization` 记录

#### Scenario: 更新已有自定义配置
- **WHEN** 用户为已有自定义配置的订阅保存自定义配置
- **THEN** 系统更新已有记录

### Requirement: 代理自定义

系统 SHALL 支持代理自定义，包含插入、追加和移除操作。

#### Scenario: 插入代理
- **WHEN** 用户在 `proxy_insert` 字段中提供 YAML
- **THEN** 系统在合并时将代理插入到订阅代理列表的开头

#### Scenario: 追加代理
- **WHEN** 用户在 `proxy_append` 字段中提供 YAML
- **THEN** 系统在合并时将代理追加到订阅代理列表的末尾

#### Scenario: 移除代理
- **WHEN** 用户在 `proxy_remove` 字段中提供 YAML 数组
- **THEN** 系统在合并时从订阅中移除匹配的代理名称

### Requirement: 代理组自定义

系统 SHALL 支持代理组自定义，包含插入、追加和移除操作。

#### Scenario: 插入代理组
- **WHEN** 用户在 `proxy_group_insert` 字段中提供 YAML
- **THEN** 系统在合并时将代理组插入到订阅代理组列表的开头

#### Scenario: 追加代理组
- **WHEN** 用户在 `proxy_group_append` 字段中提供 YAML
- **THEN** 系统在合并时将代理组追加到订阅代理组列表的末尾

#### Scenario: 移除代理组
- **WHEN** 用户在 `proxy_group_remove` 字段中提供 YAML 数组
- **THEN** 系统在合并时从订阅中移除匹配的代理组名称

### Requirement: 规则自定义

系统 SHALL 支持规则自定义，包含插入、追加和移除操作。

#### Scenario: 插入规则
- **WHEN** 用户在 `rule_insert` 字段中提供 YAML
- **THEN** 系统在合并时将规则插入到订阅规则列表的开头

#### Scenario: 追加规则
- **WHEN** 用户在 `rule_append` 字段中提供 YAML
- **THEN** 系统在合并时将规则追加到订阅规则列表的末尾

#### Scenario: 移除规则
- **WHEN** 用户在 `rule_remove` 字段中提供 YAML 数组
- **THEN** 系统在合并时从订阅中精确移除匹配的规则

### Requirement: 全局配置覆盖

系统 SHALL 支持全局配置覆盖。

#### Scenario: 覆盖全局配置
- **WHEN** 用户在 `global_override` 字段中提供 YAML
- **THEN** 系统在合并流程的最后阶段合并配置（脚本执行之前）

#### Scenario: 无效 YAML 格式
- **WHEN** 用户在 `global_override` 字段中提供无效的 YAML
- **THEN** 系统返回错误信息，不保存

### Requirement: 后处理脚本

系统 SHALL 在所有其他自定义操作完成后执行后处理脚本。

#### Scenario: 执行包含 main 函数的脚本
- **WHEN** 用户提供包含 `main(config)` 函数的 JavaScript
- **THEN** 系统以 `config` 参数调用 `main` 并使用返回的配置

#### Scenario: 执行直接返回的脚本
- **WHEN** 用户提供直接修改并返回 `config` 的 JavaScript
- **THEN** 系统使用返回的配置

#### Scenario: 脚本执行错误
- **WHEN** 脚本执行失败
- **THEN** 系统返回错误信息，不保存自定义配置

### Requirement: 合并顺序

系统 SHALL 按以下顺序应用自定义配置：
1. 移除操作（代理、代理组、规则）
2. 插入操作（代理、代理组、规则）
3. 追加操作（代理、代理组、规则）
4. 全局覆盖
5. 脚本后处理

#### Scenario: 应用完整自定义配置
- **WHEN** 订阅配置了所有类型的自定义
- **THEN** 系统按定义的顺序依次应用

### Requirement: 配置校验

系统 SHALL 在应用合并配置之前进行校验。

#### Scenario: 校验代理
- **WHEN** 校验合并后的配置
- **THEN** 系统检查代理是否为包含有效 name 字段的数组

#### Scenario: 校验代理组
- **WHEN** 校验合并后的配置
- **THEN** 系统检查代理组是否为包含有效 name 字段的数组

#### Scenario: 校验规则格式
- **WHEN** 校验合并后的配置
- **THEN** 系统检查规则是否为格式正确的字符串数组

#### Scenario: 校验规则目标存在
- **WHEN** 规则引用代理或代理组
- **THEN** 系统检查目标是否存在于代理、代理组中，或是内置目标（DIRECT、REJECT、REJECT-DROP、PASS、COMPATIBLE）

#### Scenario: 规则目标不存在
- **WHEN** 规则引用了不存在的代理或代理组
- **THEN** 系统返回错误信息，指出哪条规则的目标无效

#### Scenario: 校验代理组中的代理存在
- **WHEN** 代理组引用了代理
- **THEN** 系统检查每个代理是否存在于代理或其他代理组中

#### Scenario: 代理组中的代理不存在
- **WHEN** 代理组引用了不存在的代理或代理组
- **THEN** 系统返回错误信息，指出哪个代理组的代理引用无效

#### Scenario: 校验失败
- **WHEN** 校验失败
- **THEN** 系统恢复之前的自定义配置并返回错误信息

### Requirement: 错误处理与回滚

系统 SHALL 在发生任何错误时回滚自定义配置。

#### Scenario: YAML 格式错误
- **WHEN** 用户提供了无效的 YAML
- **THEN** 系统返回错误信息，不保存

#### Scenario: 合并错误
- **WHEN** 合并过程失败
- **THEN** 系统恢复之前的自定义配置并返回错误信息

#### Scenario: Core 应用错误
- **WHEN** 向 Core 应用配置失败
- **THEN** 系统恢复之前的自定义配置、重新应用旧配置并返回错误信息

### Requirement: 自定义 API

系统 SHALL 提供用于自定义配置管理的 REST API。

#### Scenario: 获取自定义配置
- **WHEN** 客户端请求 `GET /subscriptions/:id/customization`
- **THEN** 系统返回自定义配置数据，如不存在则返回空对象

#### Scenario: 更新自定义配置并进行校验
- **WHEN** 客户端请求 `PUT /subscriptions/:id/customization` 并附带自定义配置数据
- **THEN** 系统校验、保存、合并并应用配置；如任一步骤失败则返回错误

### Requirement: 前端自定义对话框

系统 SHALL 提供带标签页的统一自定义对话框。

#### Scenario: 打开自定义对话框
- **WHEN** 用户从订阅卡片菜单点击"自定义配置"
- **THEN** 系统显示带标签页的对话框：[代理] [代理组] [规则] [全局] [脚本]

#### Scenario: 带占位符示例的 YAML 编辑器
- **WHEN** 用户查看任意标签页
- **THEN** 系统显示带有 YAML 格式示例占位符的文本编辑区

#### Scenario: 错误显示
- **WHEN** 保存操作失败
- **THEN** 系统在对话框顶部显示红色错误提示，不关闭对话框

#### Scenario: 前端 YAML 校验
- **WHEN** 用户点击保存按钮
- **THEN** 系统在发送到后端之前校验所有 YAML 字段

## REMOVED Requirements

### Requirement: 规则管理（来自 rule-management）

**移除原因**: 替换为统一的订阅自定义配置。

**迁移方式**: 使用自定义配置中的 `rule_insert`、`rule_append`、`rule_remove` 字段。

### Requirement: 脚本管理（来自 script-management）

**移除原因**: 替换为统一的订阅自定义配置。

**迁移方式**: 使用自定义配置中的 `script` 字段。
