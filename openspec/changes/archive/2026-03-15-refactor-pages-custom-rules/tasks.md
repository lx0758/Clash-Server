## 1. 后端 API

- [x] 1.1 CoreService 新增 GetRules() 方法
- [x] 1.2 新增 GET /api/core/rules 路由端点
- [x] 1.3 新增 CoreRule 类型定义

## 2. 前端类型定义

- [x] 2.1 新增 CoreRule 接口定义
- [x] 2.2 新增 CoreRulesResponse 接口定义
- [x] 2.3 新增 coreRules API 调用函数

## 3. 新规则页面

- [x] 3.1 重写 Rules.vue，展示 Core 实时规则
- [x] 3.2 实现虚拟滚动 (useVirtualList)
- [x] 3.3 实现规则搜索过滤功能
- [x] 3.4 实现 Core 未运行时的提示状态

## 4. 订阅页面 Tab 改造

- [x] 4.1 Subscriptions.vue 改为 Tab 布局（订阅/规则/脚本）
- [x] 4.2 将原 Rules.vue 的规则管理逻辑迁移到"规则"Tab
- [x] 4.3 将 Scripts.vue 的脚本管理逻辑迁移到"脚本"Tab
- [x] 4.4 删除 Scripts.vue 文件

## 5. 路由和导航

- [x] 5.1 更新 router/index.ts 路由配置（删除 /scripts 路由）
- [x] 5.2 更新 MainLayout.vue 导航菜单
- [x] 5.3 验证导航顺序: 仪表盘 → 代理 → 规则 → 连接 → 订阅 → 设置