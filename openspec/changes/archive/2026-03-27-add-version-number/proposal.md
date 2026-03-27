## Why

项目没有统一的版本号，无法区分不同版本，也难以追踪发布状态。未来需要关联 GitHub Release 进行版本管理。

## What Changes

- 新增版本号机制，格式为 `vYYYY.M.NUM`（如 `v2026.3.1`）
- 后端通过 ldflags 在构建时注入版本号
- 前端通过 Vite 环境变量注入版本号
- 新增 API 端点 `/api/version` 返回版本信息
- 侧边栏底部显示版本号

## Capabilities

### New Capabilities

- `version-display`: 版本号显示功能，在 UI 和 API 中展示统一的版本号

### Modified Capabilities

无

## Impact

- 新增 `server/pkg/version/version.go`
- 新增 `server/internal/handler/version.go`
- 修改 `server/cmd/server/main.go`（注册路由）
- 修改 `server/Makefile`（构建时注入版本号）
- 修改根目录 `Makefile`
- 新增 `web/src/api/version.ts`
- 修改 `web/src/components/Layout/SidebarTraffic.vue`（显示版本号）
- 修改 `web/vite.config.ts` 或构建脚本
