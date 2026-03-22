# 侧边栏流量组件

## 概述

在 PC 端侧边栏底部显示实时流量数据和迷你图表。

## 功能要求

### 迷你图表

- 使用 Canvas 绘制贝塞尔曲线
- 数据右对齐（新数据在右边，与仪表盘一致）
- 保存 30 秒历史数据（30 个数据点）
- 显示上行（橙色 #f59e0b）和下行（蓝色 #3b82f6）两条曲线
- 曲线张力参数 0.3，线宽 1.5
- 每秒采集数据，使用 `setInterval(1000)`

### 数据显示

- 上行速度：使用 Upload 图标，橙色
- 下行速度：使用 Download 图标，蓝色
- 内存使用：使用 Cpu 图标，紫色
- 图标固定宽度 16px，右对齐数值

## 技术实现

- 使用 `useWebSocket` 获取 traffic 和 memory 数据
- 使用 Element Plus 图标组件（Upload, Download, Cpu）
- 使用 `requestAnimationFrame` 实现绘制循环
- 使用 `setInterval(1000)` 采集数据
- 支持高分屏（devicePixelRatio）

## 文件

- `web/src/components/Layout/SidebarTraffic.vue`
