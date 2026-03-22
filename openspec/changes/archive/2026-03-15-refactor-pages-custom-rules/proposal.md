## Why

当前前端存在概念混淆：名为"规则"的页面实际管理的是用户自定义规则（插入/追加到配置），而非 Clash Core 实时生效的规则列表。用户无法查看当前 Core 实际使用的规则，这影响了对代理行为的理解和调试。

同时，自定义规则和脚本都是为订阅服务的配置扩展，分散在独立页面不便管理。

## What Changes

- **重写**："规则"页面改为展示 Clash Core 实时规则列表（只读）
- **合并**：将"规则管理"和"脚本"功能合并到"订阅"页面，采用 Tab 布局
- **删除**：独立的"脚本"页面
- **导航调整**：顺序变为 仪表盘 → 代理 → 规则 → 连接 → 订阅 → 设置

### 具体变更

1. 重写 `Rules.vue`，从数据库规则管理改为 Core 规则展示
2. 修改 `Subscriptions.vue`，新增"规则"和"脚本"两个 Tab
3. 删除 `Scripts.vue`
4. 后端新增 `GetRules()` 方法获取 Core 实时规则
5. 更新路由和导航菜单

## Capabilities

### New Capabilities

- `core-rules-viewer`: 展示 Clash Core 实时规则列表，支持查看规则类型、payload、代理、命中统计等

### Modified Capabilities

- `subscription-management`: 新增"规则"和"脚本"两个 Tab，整合自定义配置管理
- `rule-management`: 从独立页面迁移到订阅页面的 Tab 中
- `script-management`: 从独立页面迁移到订阅页面的 Tab 中