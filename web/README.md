# ClashServer Web Frontend

基于 Vue 3 + TypeScript + Vite + Element Plus 构建的 Clash 代理服务器管理界面。

## 技术栈

- **框架**: Vue 3.5
- **UI 组件库**: Element Plus 2.x
- **状态管理**: Pinia 3.x
- **路由**: Vue Router 4.x
- **构建工具**: Vite 7.x
- **类型系统**: TypeScript 5.x

## 功能特性

- 响应式布局，支持桌面/平板/手机三端适配
- 深色主题
- 代理节点管理与切换
- 订阅管理
- 规则与脚本编辑
- 实时流量监控
- 连接管理
- 日志查看

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

## 响应式断点

| 断点 | 宽度范围 | 布局特点 |
|------|---------|---------|
| xs | < 768px | 手机布局，抽屉导航 |
| sm | ≥ 768px | 小平板布局 |
| md | ≥ 992px | 平板布局 |
| lg | ≥ 1200px | 桌面布局，折叠侧边栏 |
| xl | ≥ 1920px | 大屏布局 |