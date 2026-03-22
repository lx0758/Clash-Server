# 实现任务

## 1. 后端 - 数据模型

- [x] 1.1 修改 `model/subscription.go`
- [x] 1.2 修改 `model/rule.go`
- [x] 1.3 修改 `model/script.go`
- [x] 1.4 创建数据库迁移 (GORM AutoMigrate + migrateData)
- [x] 1.5 创建 `pkg/userinfo/parser.go`

## 2. 后端 - 服务层

- [x] 2.1 修改 `service/subscription.go`
- [x] 2.2 修改 `service/rule.go`
- [x] 2.3 修改 `service/script.go`
- [x] 2.4 修改 `service/merger.go`

## 3. 后端 - 路由和处理器

- [x] 3.1 修改路由配置
- [x] 3.2 修改 `handler/rule.go`
- [x] 3.3 修改 `handler/script.go`
- [x] 3.4 修改 `handler/subscription.go`

## 4. 前端 - 类型定义

- [x] 4.1 修改 `types/api.ts`

## 5. 前端 - 工具函数

- [x] 5.0 创建 `utils/format.ts`

## 6. 前端 - API 模块

- [x] 6.1 修改 `api/rule.ts`
- [x] 6.2 修改 `api/script.ts`
- [x] 6.3 修改 `api/subscription.ts`

## 7. 前端 - 页面重构

- [x] 7.1 重构 `pages/Subscriptions.vue`
- [x] 7.2 创建 `components/SubscriptionCard.vue`
- [x] 7.3 创建 `components/TrafficInfo.vue`
- [x] 7.4 创建 `components/SubscriptionEditDialog.vue`
- [x] 7.5 创建 `components/MergedConfigDialog.vue`
- [x] 7.6 创建 `components/RuleEditorDialog.vue`
- [x] 7.7 创建 `components/ScriptEditorDialog.vue`

## 8. 数据迁移

- [x] 8.1 编写迁移逻辑
  - [x] 现有订阅设为 source_type = "remote"
  - [x] 现有订阅设为 use_proxy = false, skip_cert = false
  - [x] 流量字段初始化为 0
  - [x] 创建默认本地配置
  - [x] 将现有规则关联到激活订阅或默认配置
  - [x] 将现有脚本关联到激活订阅或默认配置

- [x] 8.2 测试迁移 (编译验证通过)