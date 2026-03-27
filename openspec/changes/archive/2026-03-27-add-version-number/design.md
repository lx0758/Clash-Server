## Context

项目前后端分离，需要统一的版本号机制：
- 后端：Go，通过 `ldflags` 在构建时注入
- 前端：Vue 3 + Vite，通过环境变量注入
- 版本号来源：Git tag（未来关联 GitHub Release）

版本号格式采用日历版本：`vYYYY.M.NUM`

## Goals / Non-Goals

**Goals:**
- 定义统一的版本号，格式 `vYYYY.M.NUM`
- 构建时从 Git tag 读取版本号并注入前后端
- 提供 API 端点查询版本号
- 在 UI 侧边栏显示版本号

**Non-Goals:**
- 版本更新检查功能
- 自动发布到 GitHub Release
- 版本号持久化到数据库

## Decisions

### 1. 版本号格式：`vYYYY.M.NUM`

**决定**: 使用日历版本格式，如 `v2026.3.1`

**理由**:
- 一眼看出发布时间
- 不需要纠结 major/minor/patch
- 适合小团队频繁发布
- `v` 前缀符合 Git tag 惯例

### 2. 版本号存储：Git tag

**决定**: 版本号来源于 Git tag，构建时读取

**理由**:
- 与 GitHub Release 天然关联
- 无需维护额外的 VERSION 文件
- CI/CD 友好

**备选**: 独立 VERSION 文件
- 优点：无 Git 依赖
- 缺点：需要额外维护，可能与 tag 不同步

### 3. 后端注入方式：ldflags

**决定**: 使用 Go ldflags 在构建时注入

```bash
go build -ldflags "-X clash-server/pkg/version.Version=v2026.3.1"
```

**理由**:
- Go 标准做法
- 无需修改代码逻辑
- 支持开发模式默认值

### 4. 前端注入方式：Vite 环境变量

**决定**: 通过环境变量 `VITE_APP_VERSION` 注入

```bash
VITE_APP_VERSION=v2026.3.1 npm run build
```

**理由**:
- Vite 原生支持
- 通过 `import.meta.env.VITE_APP_VERSION` 访问
- 构建时静态替换

### 5. API 端点：`/api/version`

**决定**: 新增独立端点

```json
GET /api/version
{ "version": "v2026.3.1" }
```

**理由**:
- 简洁明了
- 便于扩展（未来可加 build time、commit 等）
- 前端可直接调用

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 无 Git tag 时构建失败 | Makefile 提供默认值 `dev` |
| 忘记打 tag | 文档说明发布流程 |
| 版本号格式不一致 | Makefile 验证 tag 格式 |
