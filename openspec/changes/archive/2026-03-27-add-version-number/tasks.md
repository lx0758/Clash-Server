## 1. 后端版本号模块

- [x] 1.1 创建 `server/pkg/version/version.go`，定义 `Version` 变量
- [x] 1.2 创建 `server/internal/handler/version.go`，实现 `/api/version` 端点
- [x] 1.3 在 `main.go` 中注册版本路由

## 2. 构建脚本修改

- [x] 2.1 修改 `server/Makefile`，构建时通过 ldflags 注入版本号
- [x] 2.2 修改根目录 `Makefile`，统一版本号来源

## 3. 前端版本号集成

- [x] 3.1 创建 `web/src/api/version.ts`，调用版本 API
- [x] 3.2 修改 `web/src/components/Layout/MainLayout.vue`，在底部显示版本号
- [x] 3.3 修改前端构建脚本，注入 `VITE_APP_VERSION` 环境变量

## 4. 验证

- [x] 4.1 测试 API `/api/version` 返回正确版本号
- [x] 4.2 测试侧边栏显示版本号
- [x] 4.3 测试无 Git tag 时默认值为 `dev`
