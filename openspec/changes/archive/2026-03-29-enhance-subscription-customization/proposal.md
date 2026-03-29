## Why

当前订阅自定义能力过于受限，仅支持规则的自定义，无法对代理、代理组、全局配置进行调整。用户需要完整的订阅配置自定义能力，以满足个性化分流需求。

## What Changes

- 新增 `SubscriptionCustomization` 数据模型，统一存储订阅的自定义配置
- 支持代理、代理组、规则的自定义操作：插入、追加、移除
- 支持全局配置覆盖
- 支持后处理脚本（支持 `main` 函数或直接返回对象）
- **BREAKING** 移除原有的 `Rule` 和 `Script` 独立表及相关 API
- 前端交互改为单一弹窗 + Tab 切换，每个分类使用 YAML 编辑器
- 前后端配置验证，错误时显示提示且不关闭弹窗
- API 重命名：`/subscriptions/:id/merged-config` → `/subscriptions/:id/merged`

## Capabilities

### New Capabilities
- `subscription-customization`: 订阅自定义配置能力，包括代理、代理组、规则、全局配置、脚本的自定义

### Modified Capabilities
- `rule-management`: 规则管理能力合并到 subscription-customization，原 spec 作废
- `script-management`: 脚本管理能力合并到 subscription-customization，原 spec 作废

## Impact

**后端**:
- 新增 `SubscriptionCustomization` 模型和 repository
- 新增 `/subscriptions/:id/customization` API (GET/PUT)
- 修改 `MergerService` 合并逻辑，增加验证和错误处理
- 修改脚本引擎支持 `main` 函数
- 删除 `Rule`、`Script` 相关代码
- API 重命名：`GetMergedConfig` → `GetMerged`

**前端**:
- 新增 `CustomizationDialog.vue` 组件（Tab 切换，YAML 编辑器，错误提示）
- 新增 `js-yaml` 依赖用于 YAML 格式验证
- 删除 `RuleEditorDialog.vue`、`ScriptEditorDialog.vue`、`CustomizationSection.vue`
- 修改 `SubscriptionCard.vue` 入口，调整菜单顺序
- 删除"重置"按钮
- 对话框点击外部不关闭，操作按钮使用小尺寸

**数据库**:
- 新增 `subscription_customizations` 表
- 删除 `rules`、`scripts` 表
