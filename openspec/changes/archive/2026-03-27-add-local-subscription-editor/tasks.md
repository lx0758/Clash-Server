## 1. 后端 API

- [x] 1.1 新增 `GET /api/subscriptions/:id/content` API 返回订阅内容

## 2. 内容查看/编辑组件

- [x] 2.1 创建 `web/src/components/SubscriptionContentEditor.vue`
- [x] 2.2 实现对话框布局和 textarea 编辑器
- [x] 2.3 支持 readonly 模式（远程订阅）和 edit 模式（本地订阅）
- [x] 2.4 调用 `PUT /api/subscriptions/:id/content` 保存内容
- [x] 2.5 添加复制和下载功能
- [x] 2.6 YAML 内容不自动折行
- [x] 2.7 优化小屏幕显示效果

## 3. 集成到订阅卡片

- [x] 3.1 修改 `SubscriptionCard.vue`，在操作菜单添加"查看源文件"选项
- [x] 3.2 优化本地订阅列表展示，显示节点数量和编辑时间

## 4. 页面集成

- [x] 4.1 修改 `Subscriptions.vue`，集成内容查看/编辑组件
- [x] 4.2 处理编辑完成后的刷新逻辑

## 5. 后端优化

- [x] 5.1 更新内容时计算节点数量

## 6. 验证

- [x] 6.1 测试所有订阅显示"查看源文件"选项
- [x] 6.2 测试远程订阅内容只读
- [x] 6.3 测试本地订阅内容可编辑
- [x] 6.4 测试内容编辑和保存
- [x] 6.5 测试本地订阅列表展示
