## Why

项目上线需要单一二进制部署，避免外部依赖静态资源文件。同时优化代码结构，使用社区推荐的中间件替代手写实现，提升可维护性。此外，发现若干 UI 问题需要修复，提升用户体验。

## What Changes

### 后端改进

- **静态资源嵌入**: 使用 Go `embed` 将 Vue 构建产物嵌入二进制文件，实现单一可执行文件部署
- **Session 管理重构**: 从 `gorilla/sessions` 迁移到 `gin-contrib/sessions` + `memstore`，简化使用方式
- **CORS 中间件改进**: 使用 `gin-contrib/cors` 替代手写 CORS 逻辑
- **模块名简化**: `github.com/clash-server/server` → `clash-server`，缩短导入路径
- **资源目录重组**: 静态资源从 `server/web/` 移至 `server/res/web/`

### 前端 UI 改进

- **PC 端登出入口**: 在侧边栏底部添加用户信息和退出登录按钮，与移动端保持一致
- **仪表盘核心状态**: 运行中的核心只显示版本号，简化界面
- **订阅卡片激活样式**: 优化激活订阅的视觉表现，移除绿色边框，改用更柔和的指示方式
- **规则编辑删除按钮**: 小屏幕下将删除按钮移至规则项右上角，节省空间
- **侧边栏流量组件**: 底部显示迷你流量图表和上行/下行/内存数据
- **代理延迟颜色编码**: 根据延迟值显示不同颜色（绿/黄绿/黄/橙/红/灰）
- **代理组样式优化**: 代理组添加边框框架，选中代理显示蓝色角标
- **连接信息显示优化**: 时间显示在主机名前（`HH:mm:ss.SSS - 主机名`）
- **日志格式优化**: 统一时间格式 `HH:mm:ss.SSS`，小屏采用两行布局（时间和 level 在第一行）
- **连接数统计**: 仅显示激活连接数，使用 Link 图标
- **流量图优化**: 时间选择器居中，1分钟范围显示秒数

## Capabilities

### New Capabilities

- `static-asset-embedding`: 静态资源嵌入二进制，支持 SPA 路由回退
- `pc-logout-entry`: PC 端侧边栏登出入口

### Modified Capabilities

- `session-management`: Session 管理方式变更，从 gorilla 迁移到 gin-contrib
- `dashboard-core-status`: 仪表盘核心状态显示优化
- `subscription-active-style`: 激活订阅视觉样式优化
- `rule-editor-mobile`: 规则编辑器小屏幕布局优化

## Impact

### 后端

- **代码影响**: 约 25+ 个 Go 文件导入路径更新
- **新增文件**: `server/res/res.go` (embed 定义)
- **目录变更**: 静态资源输出路径 `server/web/` → `server/res/web/`
- **依赖变更**: 
  - 新增: `gin-contrib/cors`, `gin-contrib/sessions`, `quasoft/memstore`
  - 移除: 手写 session 初始化逻辑
- **部署影响**: 二进制文件增大（包含静态资源），但部署简化为单一文件

### 前端

- **修改文件**: 
  - `MainLayout.vue` - 添加 PC 端登出入口
  - `SidebarTraffic.vue` - 新增侧边栏流量组件
  - `RuntimeInfoCard.vue` - 核心状态显示优化
  - `SubscriptionCard.vue` - 激活样式优化
  - `RuleEditorDialog.vue` - 删除按钮位置调整
  - `Settings.vue` - 卡片宽度优化
  - `ProxyItem.vue` - 延迟颜色编码、选中样式
  - `Proxies.vue` - 代理组边框框架
  - `Connections.vue` - 连接信息两行布局
  - `Logs.vue` - 日志格式优化
  - `TrafficGraph.vue` - 时间选择器布局优化
