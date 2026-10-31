# Proposal

## Why

移动端浏览器存在两个高频使用场景的视口 Bug：

1. **"电脑窗口"模式下查看器塌陷**：手机浏览器开启"桌面版网站"后，布局视口宽度约为 980px，恰好落在 JS 断点（`isMobile < 992px`）与 CSS 媒体查询（`max-width: 768px`）之间的盲区。订阅"查看源文件"弹窗按手机逻辑走了全屏 + `rows=undefined`，但拉伸样式因媒体查询不命中而失效，textarea 退化为默认 2~3 行，内容几乎不可读。
2. **"手机窗口"模式下内容被浏览器工具栏遮挡**：全站使用裸 `100vh` 计算页面高度，而移动端 `100vh` 等于地址栏收起时的大视口高度。地址栏可见时，页面底部内容及 `.el-main` 滚动条末端被工具栏挡住约 60~100px。

两个问题同源：响应式适配逻辑分散在 JS 断点、CSS 媒体查询、视口单位三处且互不同步。现在修复它，避免后续新页面继续复制错误模式。

## What Changes

- 修正订阅"查看源文件"与"查看合并配置"弹窗的全屏布局样式：以 `.el-dialog.is-fullscreen` 为唯一布局条件，移除 `@media (max-width: 768px)` 包裹，消灭 768~991px 盲区
- 全站视口高度从裸 `100vh` 迁移到 `100vh → 100dvh → var(--app-height)` 三级降级链
- 新增 `--app-height` CSS 变量：由 JS 通过 `visualViewport`/`innerHeight` 计算并在 `resize`、`orientationchange` 时更新，兜底不支持 `dvh` 的老旧国产浏览器内核
- `index.html` 的 viewport meta 追加 `interactive-widget=resizes-content`

## Capabilities

### New Capabilities

（无）

### Modified Capabilities

- `dashboard-responsiveness`:
  - **新增"CSS 与 JS 断点一致性"需求**：CSS 媒体查询的布局切换阈值必须与 JS 断点系统（`useBreakpoint`，`isMobile < 992px`）同源，禁止出现两套不一致的断点造成布局盲区
  - **新增"移动端视口高度"需求**：应用根容器及全屏弹窗的高度必须使用动态视口高度（跟随地址栏收起/展开），浏览器工具栏可见时页面内容与滚动条不得被遮挡

## Impact

- **前端代码**（仅 `web/`，不涉及后端）：
  - `web/src/composables/useBreakpoint.ts`（断点同源的参照物，可能不需要改动）
  - `web/src/components/SubscriptionContentEditor.vue`（Bug 1）
  - `web/src/components/MergedConfigDialog.vue`（Bug 1 + Bug 2）
  - `web/src/components/CustomizationDialog.vue`（Bug 2，`calc(100vh - 54px)`）
  - `web/src/style.css`、`web/src/components/Layout/MainLayout.vue`、`web/src/pages/Login.vue`（Bug 2）
  - `web/index.html`（viewport meta）
  - 新增小型 composable（如 `useAppHeight`）用于维护 `--app-height`
- **依赖**：无新增依赖（`dvh` 为原生 CSS 能力，JS 兜底仅用标准 API）
- **兼容性**：老内核浏览器走 `--app-height` 兜底路径；桌面端行为不变
