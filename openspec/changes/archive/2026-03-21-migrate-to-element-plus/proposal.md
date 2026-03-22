## Why

当前前端使用纯手写 CSS 和原生 Vue 组件，缺乏 UI 组件库支持。这导致：

1. **开发效率低** - 每个表单、对话框、导航都需要手写，维护成本高
2. **无响应式支持** - 当前布局仅适配桌面端，移动端体验差
3. **样式不统一** - 各组件样式分散，缺乏统一设计规范

引入 Element Plus 可以获得成熟的组件库、内置响应式系统和设计规范，显著提升开发效率和用户体验。

## What Changes

- **引入 Element Plus** - Vue 3 官方推荐的 UI 组件库
- **改造所有现有页面** - 使用 Element Plus 组件重写
- **实现响应式布局** - 支持桌面、平板、手机三端适配
- **配置浅色主题** - 采用浅色风格，统一配色方案
- **移除冗余自定义组件** - Toast、Loading 等由 Element Plus 替代

### 页面改造清单

| 页面 | 改造内容 |
|------|---------|
| MainLayout | el-container + el-menu (折叠导航 + 移动端抽屉) |
| Dashboard | el-card + el-col 响应式卡片 + el-statistic |
| Proxies | el-collapse + ProxyGrid 响应式网格 + 模式切换器 |
| Rules | 多列网格布局 (4/3/2/1列响应式) |
| Connections | 卡片列表 (替代 el-table，更适合移动端) |
| Logs | 卡片列表 + 时间戳 + 错误/警告突出显示 |
| Subscriptions | 两行布局卡片 + el-dialog (响应式对话框) |
| Settings | el-form 左侧标签布局 + 分隔线按钮区 |
| Login | el-form + el-card 响应式登录页 |

### UI 细节优化

- **日志页面**: 时间戳显示、错误/警告左侧边框突出、等宽字体
- **代理Item**: 类型标签 + 延迟数值显示、移动端响应式适配
- **连接页面**: chains 顺序反转、渐变背景强化
- **规则页面**: payload 和 proxy 位置调整，多列网格布局
- **订阅页面**: 两行布局 (标题+按钮 / 类型+节点+更新时间)，规则/脚本编辑器滚动优化
- **设置页面**: 表单标签左侧对齐，保存按钮独立区域

## Capabilities

### New Capabilities

- `responsive-layout`: 响应式布局能力，支持桌面/平板/手机三端适配，包括导航折叠、表单布局切换、对话框全屏等

### Modified Capabilities

- `proxy-ui`: 代理选择界面使用 Element Plus 组件重写，增加响应式支持，优化代理卡片信息展示
- `user-authentication`: 登录页面使用 el-form 重写，增加响应式支持
- `subscription-management`: 订阅管理使用两行布局卡片，规则/脚本编辑器列表可滚动
- `traffic-monitoring`: 流量监控卡片使用 el-card 重写
- `log-viewer`: 日志页面改用卡片列表，优化错误/警告显示
- `connection-viewer`: 连接页面改用卡片列表，chains 顺序反转
- `rule-viewer`: 规则页面改用多列网格布局，信息层次优化
- `settings-page`: 设置页面表单布局优化，左侧标签对齐

## Impact

### 代码影响

- `web/src/components/` - 所有组件重写
- `web/src/pages/` - 所有页面改造
- `web/src/style.css` - 替换为 Element Plus 主题配置
- `web/src/styles/theme.css` - 浅色主题 CSS 变量
- `web/src/composables/useBreakpoint.ts` - 响应式断点 hook
- `web/package.json` - 新增 element-plus 依赖

### 依赖变更

```diff
+ element-plus: ^2.x
+ @element-plus/icons-vue: ^2.x
- (移除手写 Toast/Loading 组件)
```

### 兼容性

- **BREAKING**: 移除自定义 Toast/Loading 组件，改用 Element Plus 全局 API
- 现有组件 API 保持不变，仅内部实现变化
- 主题从深色改为浅色

### 测试验证

- ✅ 桌面端响应式测试通过
- ✅ iPhone 12 Pro 移动端测试通过
- ✅ 所有页面构建成功，无 TypeScript 错误