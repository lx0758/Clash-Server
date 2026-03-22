## Context

当前前端架构：
- Vue 3.5 + TypeScript + Vite 7
- 纯手写 CSS，无 UI 框架
- 固定布局，仅适配桌面端
- 自定义 Toast/Loading 组件

目标架构：
- 引入 Element Plus 2.x
- 响应式布局（桌面/平板/手机）
- 保持深色主题

## Goals / Non-Goals

**Goals:**
- 引入 Element Plus 组件库，替换手写组件
- 实现响应式布局，支持三端适配
- 配置深色主题，保持现有视觉风格
- 减少维护成本，提高开发效率

**Non-Goals:**
- 不改变现有 API 接口
- 不改变现有业务逻辑
- 不支持 IE 浏览器
- 不引入 SSR

## Decisions

### 1. 组件库选择：Element Plus

**选择**: Element Plus

**理由**:
- Vue 3 官方推荐，生态成熟
- 内置响应式栅格系统和断点
- 支持按需引入，包体积可控
- 深色模式官方支持 (CSS 变量)

**备选方案**:
- Ant Design Vue: 更偏企业级，组件更重
- Naive UI: 更现代，但生态较小
- Vuetify: Material Design 风格，与当前深色风格差异大

### 2. 响应式断点策略

**选择**: 采用 Element Plus 默认断点

```
xs: < 768px   (手机)
sm: ≥ 768px   (小平板)
md: ≥ 992px   (平板)
lg: ≥ 1200px  (桌面)
xl: ≥ 1920px  (大屏)
```

**理由**: 与 Element Plus 栅格系统一致，无需额外配置

### 3. 导航布局策略

**选择**: 
- 桌面端 (≥992px): 固定侧边栏 (可折叠)
- 移动端 (<992px): 抽屉导航 + 底部固定触发按钮

**理由**: 
- `el-menu` 内置折叠功能
- 移动端使用 `el-drawer` 是最佳实践
- 避免在小屏幕上占用过多空间

### 4. 深色主题实现

**选择**: CSS 变量覆盖 + Element Plus 暗黑模式

```css
/* 使用 Element Plus 暗黑模式作为基础 */
@use "element-plus/theme-chalk/dark/css-vars.css";

/* 自定义颜色覆盖 */
:root {
  --el-bg-color: #1a1a2e;
  --el-bg-color-page: #16213e;
  --el-text-color-primary: #ffffff;
  --el-color-primary: #e94560;
}
```

**理由**: 
- 官方支持，兼容性好
- 可精细控制每个组件的颜色
- 保持现有深色风格

### 5. 表单响应式策略

**选择**: 
- 标签位置: 桌面端右对齐，移动端顶部对齐
- 对话框: 移动端全屏
- 表格: 横向滚动或隐藏次要列

**实现**:
```vue
<el-form :label-position="isMobile ? 'top' : 'right'">
<el-dialog :fullscreen="isMobile">
```

### 6. 保留的自定义组件

**选择**: ProxyGrid 保持自定义实现

**理由**:
- 虚拟滚动需求 (vue-virtual-scroller)
- 特殊的卡片布局逻辑
- Element Plus 无完全匹配组件
- 可使用 `el-row`/`el-col` 配合实现响应式

## Risks / Trade-offs

### 风险 1: 样式冲突
- **风险**: Element Plus 样式可能与现有 CSS 冲突
- **缓解**: 使用 scoped 样式 + 按需引入 + 全局样式隔离

### 风险 2: 包体积增加
- **风险**: Element Plus 增加 ~500KB (gzip ~150KB)
- **缓解**: 按需引入 + Vite tree-shaking

### 风险 3: 深色主题不一致
- **风险**: 部分组件深色效果不理想
- **缓解**: 统一 CSS 变量 + 逐个组件验证

### 风险 4: 响应式测试覆盖
- **风险**: 多端测试工作量大
- **缓解**: 定义关键断点场景，优先验证高频页面

## Migration Plan

### 阶段 1: 基础设施 (1 天)
1. 安装 Element Plus + 图标库
2. 配置按需引入
3. 配置深色主题 CSS 变量
4. 创建响应式 hook (`useBreakpoint`)

### 阶段 2: 布局改造 (1 天)
1. MainLayout 使用 el-container/el-menu
2. 实现移动端抽屉导航
3. 测试导航响应式

### 阶段 3: 页面改造 (3-4 天)
1. Login 页面
2. Dashboard 页面
3. Settings 页面
4. Subscriptions 页面
5. Proxies 页面 (最复杂)
6. Logs/Connections/Rules 页面

### 阶段 4: 收尾 (1 天)
1. 移除废弃组件 (Toast/Loading)
2. 全局样式清理
3. 响应式测试
4. 文档更新

### 回滚策略
- 保留 `feature/element-plus-migration` 分支
- 每个阶段完成后创建 tag
- 出现重大问题时可回滚到上一阶段

## Open Questions

1. **是否需要支持平板横屏特殊布局?** — 暂不处理，跟随响应式自动适配
2. **是否需要底部导航栏?** — 移动端使用抽屉导航，暂不加底部导航