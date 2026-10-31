# Design

## Context

动机见 proposal.md（Why）。此处只列约束现状：

- JS 断点唯一真源是 `web/src/composables/useBreakpoint.ts`（`isMobile = width < 992px`），组件通过它控制弹窗 `fullscreen` 属性
- 两个订阅弹窗（`SubscriptionContentEditor.vue`、`MergedConfigDialog.vue`）的全屏布局样式写在**非 scoped** `<style>` 中并挂在 `@media (max-width: 768px)` 下——与 JS 断点不一致，形成 768~991px 盲区
- 全站高度硬编码 `100vh`：`style.css`（body/#app）、`MainLayout.vue`（.layout，主滚动容器）、`Login.vue`、`MergedConfigDialog.vue`、`CustomizationDialog.vue`（calc）
- 目标用户使用国内手机浏览器（大量较老 X5/XWeb 内核），不能假设 `dvh` 可用
- Element Plus 的 `.el-dialog.is-fullscreen` 自带 `height: 100%`，是全屏样式的可靠锚点

## Goals / Non-Goals

**Goals:**

- 消灭 JS/CSS 断点错配导致的全屏弹窗布局失效
- 全站高度感知浏览器工具栏动态变化，内容与滚动条不被遮挡
- 对不支持 `dvh` 的浏览器有可验证的降级路径

**Non-Goals:**

- 不重写断点系统本身（`useBreakpoint` 的阈值、API 不变）
- 不改桌面端任何视觉表现
- 不迁移现有 `<style>` 为 scoped（teleport 到 body 的弹窗必须用全局样式，维持现状）
- 不处理主 spec `openspec/specs/dashboard-responsiveness/spec.md` 的结构性损坏（归档阻碍问题另案处理）

## Decisions

### D1: 全屏布局以 `.el-dialog.is-fullscreen` 为唯一条件，删除媒体查询包裹

全屏样式条件改为 `:global(.el-dialog.is-fullscreen) .content-editor` 等选择器，**不再包 `@media`**。

- 依据：`is-fullscreen` 只在 `isMobile=true` 时出现，媒体查询是冗余条件；删掉后 CSS 与 JS 天然同源，盲区消失
- 备选（否决）：把 JS 断点改成 768 对齐 CSS——会改变 768~991px 区间的整体布局行为（侧边栏变抽屉），影响面远超本 bug
- 备选（否决）：把 CSS 媒体查询改成 992px 对齐 JS——留下两套断点，未来仍会漂移

### D2: 视口高度用三级降级链 `100vh → 100dvh → var(--app-height)`

```css
height: 100vh;                          /* 兜底：极老内核 / 变量未就绪 */
height: 100dvh;                         /* 支持则优先：纯 CSS 动态跟随，无 JS 抖动 */
height: var(--app-height, 100vh);       /* JS 兜底：变量就绪后覆盖；必须带 fallback */
```

- `--app-height` 由 JS 在 `document.documentElement.style` 上维护
- **注意**：`var()` 不带 fallback 时，变量未定义会导致整条声明在计算期失效（属性变 unset），**不会**回退到前一条声明——所以第三级必须写 `var(--app-height, 100vh)`；`calc()` 内同理
- 备选（否决）：纯 JS `--app-height`——放弃 dvh 的原生平滑性，工具栏收起过程会有阶梯感
- 备选（否决）：纯 dvh——老内核国产浏览器 bug 依旧，不满足降级需求
- 依据：三级链让每类浏览器都取到最优可用值，且互不冲突（后者覆盖前者）

### D3: 新增 `useAppHeight` composable 统一维护 `--app-height`

- 计算取 `window.visualViewport?.height ?? window.innerHeight`
- 监听 `resize` 与 `orientationchange`（含 100ms 防抖外的立即首帧赋值），在 `App.vue` 挂载一次全局生效
- 桌面端该值恒等于视口高度，行为等价于原 `100vh`
- 备选（否决）：`window.innerHeight` 单用——iOS Safari 上它不随地址栏展开实时收缩，visualViewport 更准

### D4: 所有裸 `100vh` 统一替换为降级链

涉及：`style.css`（body min-height、#app）、`MainLayout.vue`（.layout）、`Login.vue`、`MergedConfigDialog.vue`（fullscreen 高度）、`CustomizationDialog.vue`（`calc(100vh - 54px)` → `calc(var(--app-height, 100dvh) - 54px)`，注意保持降级次序一致）。

### D5: viewport meta 追加 `interactive-widget=resizes-content`

让 Chrome 系在软键盘弹出时收缩布局视口，编辑弹窗场景直接受益；老内核忽略该属性，无副作用。

## Risks / Trade-offs

- [首帧 `--app-height` 未就绪，dvh 不支持的老内核退回 `100vh`] → composable 在 `App.vue` 同步首帧赋值（setup 阶段即写入一次），仅极短窗口内为 `100vh`
- [`visualViewport` resize 在 iOS 上滚动时也触发（地址栏收起动画期间高频）] → 更新仅是写一个 CSS 变量，代价极低；不做防抖以免动画期间高度滞后
- [弹窗样式是非 scoped 全局样式，`.el-dialog.is-fullscreen` 选择器可能误伤其他弹窗] → 现有写法已如此（`.el-dialog.is-fullscreen .config-viewer`），按组件内类名（.content-editor/.config-viewer）限定作用域，保持现状边界
- [`calc()` 内 `var()` 无 fallback 时整条声明失效] → CustomizationDialog 的 calc 必须写 `var(--app-height, 100dvh)` 形式，写入 tasks 时明确

## Migration Plan

纯前端样式/脚本变更，无数据、无 API、无部署顺序问题。回滚 = revert 单个 commit。发布后真机验证清单写入 tasks。

## Open Questions

（无——主 spec 结构损坏的修复已明确另案处理，不影响本设计。）
