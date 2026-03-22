## Why

当前仪表盘功能过于简单，仅显示 4 个基础统计卡片（核心状态、上行流量、下行流量、订阅数量），缺少 Clash Verge 等现代代理客户端的丰富功能。用户体验不够友好，无法快速查看实时流量趋势、切换代理模式、管理订阅配置等核心操作。移动端体验尤其需要改进，缺少导航栏流量指示器等便捷功能。

## What Changes

- 添加实时流量图表组件（TrafficGraph），使用 Canvas 绘制双曲线图（上行橙色/下行蓝色），支持响应式显示（桌面端完整显示，移动端精简显示）
- 添加增强型流量统计组件（EnhancedTrafficStats），包含 6 个统计卡片：上行速度、下行速度、活跃连接数、总上传量、总下载量、内存使用
- 添加代理模式切换组件（ClashModeCard），支持 Rule/Global/Direct 三种模式切换，带动画效果和模式说明
- 添加当前代理组件（CurrentProxyCard），显示当前选择的代理，支持快速切换和延迟测试
- 添加网络设置组件（ProxyTunCard），支持系统代理和 TUN 模式开关
- 添加订阅配置组件（HomeProfileCard），显示当前订阅信息、流量使用进度、过期时间
- 在移动端导航栏添加流量指示器（TrafficIndicator），实时显示上行/下行速度
- 所有新增组件完全支持响应式布局，适配移动端（<768px）、平板端（768-1200px）、桌面端（>1200px）
- 优化现有 Dashboard 布局，使用 Element Plus Grid 系统实现弹性响应式

## Capabilities

### New Capabilities

- `traffic-visualization`: 流量数据可视化能力，包括实时流量图表和统计卡片展示
- `proxy-mode-control`: 代理模式交互控制，支持模式切换和状态显示
- `dashboard-responsiveness`: 仪表盘响应式布局，确保在所有设备尺寸下都有良好体验
- `navigation-traffic-indicator`: 导航栏流量实时指示器，为移动端提供便捷的流量查看

### Modified Capabilities

无。现有 traffic-monitoring 规范已满足后端数据需求，无需修改。

## Impact

**前端代码**
- 新增 7 个 Vue 组件：TrafficGraph.vue, TrafficIndicator.vue, EnhancedTrafficStats.vue, StatCard.vue, ClashModeCard.vue, CurrentProxyCard.vue, NetworkSettingsCard.vue, HomeProfileCard.vue
- 修改 Dashboard.vue，集成所有新组件并实现响应式布局
- 修改 MainLayout.vue，在移动端 Header 添加 TrafficIndicator

**API 和数据**
- 复用现有的 WebSocket 流量数据推送（traffic, connections, memory）
- 复用现有的代理 API（proxies/mode, proxies/{name}/delay）
- 无需后端 API 修改

**依赖**
- 使用原生 Canvas API，无需额外绘图库依赖
- 复用现有 Vue 生态：Vue 3, Pinia, Vue Router, Element Plus, @vueuse/core

**兼容性**
- 保持现有功能不变，仅为增强体验
- 所有新增组件可选显示，不影响现有布局
- 响应式设计确保在所有设备上可访问核心功能
