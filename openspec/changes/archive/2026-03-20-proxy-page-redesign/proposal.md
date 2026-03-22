## Why

当前代理界面将所有节点平铺展示，当节点数量多时界面过长、难以管理。需要参考 Clash Verge Rev 的设计，实现可折叠代理组、多列网格布局、筛选排序、运行模式切换等功能。

## What Changes

- **运行模式切换**: 支持规则/全局/直连三种模式切换
  - 规则模式: 显示非预置代理组，按规则分流
  - 全局模式: 显示所有代理供选择
  - 直连模式: 所有流量直连
- **可折叠代理组**: 点击代理组头展开/收起，减少界面占用
- **多列网格布局**: 节点以网格形式展示，响应式列数 (1-4列)
- **筛选功能**: 按名称筛选节点
- **排序功能**: 按延迟/名称排序
- **定位当前节点**: 一键滚动到当前选中的节点
- **测试延迟**: 支持测试单个节点或组内所有节点延迟
- **选中代理不刷新**: 选中代理后仅更新本地状态，不重新拉取数据

## Capabilities

### New Capabilities

- `proxy-ui`: 代理界面增强，支持模式切换、折叠、筛选、排序、多列布局
- `proxy-mode`: 运行模式切换能力

### Modified Capabilities

无（新增能力）

## Impact

**前端改动**:
- `web/src/pages/Proxies.vue` - 重构页面结构，支持模式切换
- `web/src/components/proxy/` - 新增组件目录
  - `ModeSwitcher.vue` - 模式切换组件
  - `ProxyGroup.vue` - 可折叠代理组
  - `ProxyHead.vue` - 操作栏
  - `ProxyItem.vue` - 节点卡片
  - `ProxyGrid.vue` - 节点网格
- `web/src/stores/proxy.ts` - 扩展状态管理，添加 mode 状态
- `web/src/composables/useProxyFilter.ts` - 筛选排序逻辑
- `web/src/api/proxy.ts` - 新增延迟测试、模式切换 API
- `web/src/types/api.ts` - 新增 Mihomo 类型定义

**后端改动**:
- `server/internal/service/core.go` - 新增 GetMode/SetMode 方法
- `server/internal/handler/proxy.go` - 新增延迟测试、模式切换接口

**用户体验改进**:
- 支持三种运行模式切换
- 大量节点时界面更清晰
- 快速定位和筛选节点
- 选中代理无页面闪烁