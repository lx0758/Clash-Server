## 1. Backend Data Model

- [x] 1.1 Create `SubscriptionCustomization` model in `model/customization.go`
- [x] 1.2 Create `CustomizationRepository` in `repository/customization.go`
- [x] 1.3 Add database migration for `subscription_customizations` table

## 2. Backend Service Layer

- [x] 2.1 Create `CustomizationService` in `service/customization.go`
- [x] 2.2 Update `MergerService` to apply customization during merge
- [x] 2.3 Implement remove/insert/append logic for proxies
- [x] 2.4 Implement remove/insert/append logic for proxy groups
- [x] 2.5 Implement remove/insert/append logic for rules
- [x] 2.6 Implement global override merge logic with error handling
- [x] 2.7 Update script execution to support `main` function
- [x] 2.8 Add comprehensive config validation (proxies/proxy-groups/rules)
- [x] 2.9 Use names instead of indices in validation error messages
- [x] 2.10 Add `extractName` helper function using reflection

## 3. Backend API

- [x] 3.1 Create `CustomizationHandler` in `handler/customization.go`
- [x] 3.2 Add `GET /subscriptions/:id/customization` endpoint
- [x] 3.3 Add `PUT /subscriptions/:id/customization` endpoint with validation
- [x] 3.4 Register routes in `main.go`
- [x] 3.5 Rename API: `GET /subscriptions/:id/merged-config` → `GET /subscriptions/:id/merged`
- [x] 3.6 Rename function: `GetMergedConfig` → `GetMerged`
- [x] 3.7 Add rollback mechanism on validation/core error

## 4. Backend Cleanup

- [x] 4.1 Delete `model/rule.go`
- [x] 4.2 Delete `repository/rule.go`
- [x] 4.3 Delete `service/rule.go`
- [x] 4.4 Delete `handler/sub_rule.go`
- [x] 4.5 Delete `model/script.go`
- [x] 4.6 Delete `repository/script.go`
- [x] 4.7 Delete `service/script.go`
- [x] 4.8 Delete `handler/sub_script.go`
- [x] 4.9 Remove old rule/script routes from `main.go`
- [x] 4.10 Update `db.go` to remove Rule/Script from autoMigrate

## 5. Frontend API Layer

- [x] 5.1 Create `api/customization.ts`
- [x] 5.2 Delete `api/rule.ts`
- [x] 5.3 Delete `api/script.ts`
- [x] 5.4 Rename `getMergedConfig` → `getMerged` in subscription API

## 6. Frontend Components

- [x] 6.1 Create `CustomizationDialog.vue` component
- [x] 6.2 Implement tab switching (代理/代理组/规则/全局/脚本)
- [x] 6.3 Implement YAML editors with autosize and placeholders
- [x] 6.4 Add error alert display at top of dialog
- [x] 6.5 Install `js-yaml` and `@types/js-yaml` dependencies
- [x] 6.6 Implement frontend YAML validation
- [x] 6.7 Remove "重置" button
- [x] 6.8 Set `close-on-click-modal="false"` for all dialogs
- [x] 6.9 Change buttons to `size="small"` for copy/download
- [x] 6.10 Simplify button labels: "复制 YAML" → "复制", "下载文件" → "下载"
- [x] 6.11 Remove "关闭" button from content editor footer

## 7. Frontend Integration

- [x] 7.1 Update `SubscriptionCard.vue` to add "自定义配置" menu item
- [x] 7.2 Update menu order: 编辑 → 刷新 → 查看源文件 → 自定义配置 → 查看合并配置 → 删除
- [x] 7.3 Remove `v-if="subscription.active"` from "查看合并配置" menu item
- [x] 7.4 Update `Subscriptions.vue` to use new dialog
- [x] 7.5 Delete `RuleEditorDialog.vue`
- [x] 7.6 Delete `ScriptEditorDialog.vue`
- [x] 7.7 Delete `CustomizationSection.vue`
- [x] 7.8 Update `types/api.ts` to add Customization type

## 8. Testing & Verification

- [x] 8.1 Test API endpoints with curl/Postman
- [x] 8.2 Test frontend dialog interactions
- [x] 8.3 Test merge logic with various customization combinations
- [x] 8.4 Test error handling (YAML format, script error, core error)
- [x] 8.5 Test mobile responsive layout
