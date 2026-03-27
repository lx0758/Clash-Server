## Why

本地订阅（`source_type=local`）没有界面可以编辑订阅内容。用户创建本地订阅后，无法通过 UI 输入或修改 Clash 配置内容，导致本地订阅功能无法正常使用。同时，远程订阅用户也希望能查看订阅内容，而不仅仅是查看合并后的配置。

此外，本地订阅在列表中显示信息不足，缺少节点数量和编辑时间等关键信息。

## What Changes

- 在订阅卡片的操作菜单中添加"查看源文件"选项
- 创建内容查看/编辑对话框，提供 YAML 格式的代码查看和编辑
- 本地订阅支持编辑，远程订阅仅支持只读查看
- 添加复制和下载功能
- 优化本地订阅列表展示，显示节点数量和编辑时间
- 优化小屏幕下的显示效果
- YAML 内容不自动折行

## Capabilities

### New Capabilities

- `subscription-content-viewer`: 订阅内容查看/编辑功能，本地订阅可编辑，远程订阅只读

### Modified Capabilities

无

## Impact

- 新增 `web/src/components/SubscriptionContentEditor.vue`
- 修改 `web/src/components/SubscriptionCard.vue`（添加菜单项、优化本地订阅展示）
- 修改 `web/src/pages/Subscriptions.vue`（集成编辑器）
- 新增 `GET /api/subscriptions/:id/content` API（后端）
- 修改 `server/internal/service/subscription.go`（更新内容时计算节点数量）
