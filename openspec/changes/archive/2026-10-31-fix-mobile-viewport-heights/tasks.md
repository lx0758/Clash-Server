# Tasks

## 1. 视口高度基础设施（D2/D3/D5）

- [x] 1.1 新增 `web/src/composables/useAppHeight.ts`：以 `window.visualViewport?.height ?? window.innerHeight` 计算高度，在 setup 阶段同步首帧写入 `document.documentElement.style.setProperty('--app-height', ...)`，并监听 `resize`、`orientationchange` 持续更新；验证：组件挂载后 `document.documentElement.style.getPropertyValue('--app-height')` 非空且等于当前视口高度
- [x] 1.2 在 `App.vue` 挂载 `useAppHeight()`；验证：刷新页面后控制台读取 `--app-height` 已设置，窗口 resize 后值跟随变化
- [x] 1.3 `web/index.html` 的 viewport meta 追加 `interactive-widget=resizes-content`；验证：查看渲染后的 meta 标签内容正确

## 2. 全站 100vh 迁移（D1 盲区修复 + D4）

- [x] 2.1 `web/src/components/SubscriptionContentEditor.vue`：删除非 scoped 样式中的 `@media (max-width: 768px)` 包裹，以 `.el-dialog.is-fullscreen` 为唯一布局条件；验证：Chrome DevTools 设备模拟 980px 宽（或"电脑窗口"模式真机）打开"查看源文件"，textarea 拉伸占满弹窗而非 2~3 行
- [x] 2.2 `web/src/components/MergedConfigDialog.vue`：同样删除 `@media` 包裹改用 `.el-dialog.is-fullscreen` 条件，同时将其 `height/max-height: 100vh` 改为三级降级链（100vh → 100dvh → var(--app-height)）；验证：980px 宽下"查看合并配置"弹窗布局正常，且移动端地址栏可见时弹窗底不被遮挡
- [x] 2.3 `web/src/style.css`：body `min-height` 与 `#app` `height` 的 `100vh` 替换为三级降级链；验证：浏览器计算样式中最终生效值为 dvh 或 --app-height（视内核而定）
- [x] 2.4 `web/src/components/Layout/MainLayout.vue`：`.layout` 的 `height: 100vh` 替换为三级降级链；验证：移动端浏览器地址栏可见时，`.main` 滚动条可滚动到最底部且最后一屏内容完整可见
- [x] 2.5 `web/src/pages/Login.vue`：`min-height: 100vh` 替换为降级链；验证：移动端登录页底部按钮在地址栏可见时不被工具栏遮挡
- [x] 2.6 `web/src/components/CustomizationDialog.vue`：`calc(100vh - 54px)` 改为 `calc(var(--app-height, 100dvh) - 54px)`（var 必须带 100dvh fallback，否则 calc 整条失效）；验证：计算样式在 var 未定义时回退 dvh、定义后取 --app-height，弹窗内容区不溢出

## 3. 集成验证

- [x] 3.1 `npm run build`（web/ 目录）构建通过，无新增警告；验证：构建产物成功生成
- [x] 3.2 桌面端回归：宽屏下打开仪表盘、订阅页及三个弹窗，视觉与迁移前一致；验证：对照截图或人工确认无布局变化
- [x] 3.3 移动端场景验证（真机或 DevTools 模拟）：① 手机模式滚动到底部内容/滚动条不被遮挡；② "电脑窗口"（980px 视口）打开"查看源文件"与"查看合并配置"内容区占满可用高度；③ 横竖屏切换后高度正确

## Workflow follow-up

- 实现完成后归档 change（`/opsx-archive`）；注意归档前需先处理主 spec `openspec/specs/dashboard-responsiveness/spec.md` 的结构性损坏（delta 头残留导致无法合并，需另案修复）。
