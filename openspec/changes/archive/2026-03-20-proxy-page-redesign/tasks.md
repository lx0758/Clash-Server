## 1. 依赖与基础设施

- [x] 1.1 创建 `web/src/components/proxy/` 目录结构

## 2. 状态管理扩展

- [x] 2.1 扩展 `ProxyStore`，新增 `groupStates` 状态
- [x] 2.2 新增 `useProxyFilter.ts` composable，实现筛选/排序逻辑
- [x] 2.3 实现响应式列数计算
- [x] 2.4 新增 `mode` 状态，支持运行模式
- [x] 2.5 新增 `globalGroup` 计算属性，获取全局代理组
- [x] 2.6 适配 Mihomo API 类型（`all` 为字符串数组）

## 3. 组件开发

- [x] 3.1 开发 `ModeSwitcher.vue` - 模式切换器
- [x] 3.2 开发 `ProxyGroup.vue` - 可折叠代理组容器
- [x] 3.3 开发 `ProxyHead.vue` - 操作栏（筛选、排序、定位、测延迟）
- [x] 3.4 开发 `ProxyItem.vue` - 节点卡片
- [x] 3.5 开发 `ProxyGrid.vue` - 节点网格

## 4. 页面重构

- [x] 4.1 重构 `Proxies.vue`，使用新组件
- [x] 4.2 实现三种模式的 UI 展示
- [x] 4.3 规则模式：过滤 GLOBAL 组，按类型排序
- [x] 4.4 全局模式：显示所有代理供选择
- [x] 4.5 直连模式：显示提示信息
- [x] 4.6 实现选中代理不刷新页面

## 5. 后端 API 扩展

- [x] 5.1 新增 `GET /api/proxies/:name/delay` 接口
- [x] 5.2 新增 `GET /api/proxies/group/:group/delay` 接口
- [x] 5.3 新增 `GET /api/proxies/mode` 接口
- [x] 5.4 新增 `PATCH /api/proxies/mode` 接口
- [x] 5.5 修复路由参数名错误（`c.Param("group")` → `c.Param("name")`）

## 6. 测试与优化

- [x] 6.1 测试响应式布局在不同窗口大小下的表现
- [x] 6.2 测试延迟测试功能
- [x] 6.3 测试模式切换功能
- [x] 6.4 测试选中代理功能