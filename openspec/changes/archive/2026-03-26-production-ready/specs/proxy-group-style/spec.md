# 代理组样式优化

## 概述

优化代理组的视觉呈现，添加边框框架包裹节点列表。

## 样式要求

- 代理组头部：折叠面板，显示组名、类型、节点数、当前选中
- 内容区域：灰色背景 + 边框框架
- 选中节点：左上角蓝色三角形角标
- 全局模式：独立的卡片式布局

## 文件

- `web/src/pages/Proxies.vue`
- `web/src/components/proxy/ProxyGroup.vue`
- `web/src/components/proxy/ProxyItem.vue`
