## 1. 项目设置和基础设施

- [x] 1.1 创建 dashboard 组件目录结构（web/src/components/dashboard/）
- [x] 1.2 创建流量图表相关 composables（useTrafficGraph.ts）
- [x] 1.3 添加流量数据格式化工具函数（formatTraffic, formatBytes）
- [x] 1.4 配置响应式 Grid 系统断点（确认 useBreakpoint 与 Element Plus Grid 一致）

## 2. 流量指示器组件（TrafficIndicator）

- [x] 2.1 创建 TrafficIndicator.vue 组件（web/src/components/dashboard/TrafficIndicator.vue）
- [x] 2.2 实现流量指示器基本布局（上行/下行两个流量项）
- [x] 2.3 集成 WebSocket 流量数据（使用 useWebSocket）
- [x] 2.4 实现流量数据格式化（简化单位：MB/s → M，KB/s → K）
- [x] 2.5 添加流量指示器样式（橙色/蓝色、圆角、紧凑布局）
- [x] 2.6 实现点击跳转到仪表盘功能
- [x] 2.7 添加悬停效果和高对比度模式支持
- [x] 2.8 实现响应式适配（小屏设备字体和间距调整）
- [x] 2.9 添加无障碍支持（aria-label、键盘导航）

## 3. 流量图表组件（TrafficGraph）

- [x] 3.1 创建 TrafficGraph.vue 组件（web/src/components/dashboard/TrafficGraph.vue）
- [x] 3.2 设置 Canvas 元素和响应式尺寸
- [x] 3.3 实现流量数据缓冲区（支持 1m/5m/10m/30m 时间范围）
- [x] 3.4 实现双曲线绘制（上行橙色、下行蓝色）
- [x] 3.5 添加 Catmull-Rom 样条曲线平滑处理（bezierCurveTo）
- [x] 3.6 实现图表样式切换（折线图 ↔ 面积图）
- [x] 3.7 添加网格线和图例（底部居中显示上行/下行速度）
- [x] 3.8 实现响应式高度调整（桌面端 200px，移动端 140px）
- [x] 3.9 实现时间范围选择器（1分钟、5分钟、10分钟、30分钟）
- [x] 3.10 添加 requestAnimationFrame 优化和性能降级
- [x] 3.11 实现组件卸载时清理资源（取消动画、清理 Canvas）
- [x] 3.12 优化图表边界裁剪（防止曲线越界）
- [x] 3.13 实现数据对齐时间轴（数据少时从右侧开始绘制）
- [x] 3.14 优化纵轴标签格式（精简显示：1.5M 格式）

## 4. 统计卡片组件（StatCard）

- [x] 4.1 创建 StatCard.vue 通用组件（web/src/components/dashboard/StatCard.vue）
- [x] 4.2 实现卡片布局（图标 + 数值 + 单位）
- [x] 4.3 添加彩色图标背景和样式（44px 图标、22px 数值）
- [x] 4.4 实现悬停效果（背景加深、阴影）
- [x] 4.5 添加 raw 属性支持纯数字显示（如活跃连接数）
- [x] 4.6 实现响应式适配（移动端紧凑布局、垂直居中）
- [x] 4.7 添加数据格式化（智能单位转换、字节单位合并）

## 5. 增强型流量统计组件（EnhancedTrafficStats）

- [x] 5.1 创建流量信息卡片布局（流量图 + 5个统计子卡片）
- [x] 5.2 集成 TrafficGraph 组件（支持时间范围选择）
- [x] 5.3 创建 5 个 StatCard 实例（上行速度、下行速度、上传流量、下载流量、活跃连接）
- [x] 5.4 集成 WebSocket 数据（traffic、connections）
- [x] 5.5 实现响应式 Grid 布局（桌面端 4 列，平板端 2 列，移动端 2 列）
- [x] 5.6 实现统计卡片最小高度（72px）和间距（12px）
- [x] 5.7 修复小屏幕滚动问题和布局自适应

## 6. 代理模式卡片组件（ClashModeCard）

- [x] 6.1 创建 ClashModeCard.vue 组件（web/src/components/dashboard/ClashModeCard.vue）
- [x] 6.2 实现 3 个模式按钮（规则、全局、直连）
- [x] 6.3 添加模式图标（Operation、Promotion、Aim）
- [x] 6.4 实现当前模式高亮显示（主题色背景、加粗）
- [x] 6.5 集成代理 API（proxyApi.setMode）
- [x] 6.6 实现模式切换逻辑（调用 API、更新状态、错误处理）
- [x] 6.7 添加模式切换动画（200ms、ease-in-out）
- [x] 6.8 实现模式说明文字显示（规则/全局/直连的说明）
- [x] 6.9 实现响应式适配（按钮宽度自适应、均分布局）

## 7. 当前代理卡片组件（CurrentProxyCard）

- [ ] 7.1 创建 CurrentProxyCard.vue 组件（web/src/components/dashboard/CurrentProxyCard.vue）
- [ ] 7.2 实现当前代理显示（代理名称、类型、延迟）
- [ ] 7.3 集成 proxyStore 获取当前代理状态
- [ ] 7.4 添加代理选择下拉菜单或快速切换功能
- [ ] 7.5 实现延迟测试按钮和结果显示
- [ ] 7.6 添加代理状态指示器（在线/离线、延迟颜色）
- [ ] 7.7 实现响应式适配（移动端简化布局）

## 8. 网络设置卡片组件（NetworkSettingsCard）

- [ ] 8.1 创建 NetworkSettingsCard.vue 组件（web/src/components/dashboard/NetworkSettingsCard.vue）
- [ ] 8.2 实现系统代理开关（Toggle）
- [ ] 8.3 实现系统代理状态显示（启用/禁用、端口信息）
- [ ] 8.4 实现系统代理 API 调用（获取状态、设置开关）
- [ ] 8.5 添加 TUN 模式开关（如果后端支持）
- [ ] 8.6 实现 TUN 模式状态显示和 API 调用
- [ ] 8.7 添加 Tab 切换（系统代理 ↔ TUN 模式）
- [ ] 8.8 实现响应式适配（移动端简化布局）

## 9. 订阅配置卡片组件（ActiveSubscriptionCard）

- [x] 9.1 创建 ActiveSubscriptionCard.vue 组件（web/src/components/dashboard/ActiveSubscriptionCard.vue）
- [x] 9.2 实现激活订阅信息显示（订阅名称、状态图标）
- [x] 9.3 集成 subscriptionApi 获取订阅数据
- [x] 9.4 实现流量使用进度条（上传 + 下载 / 总流量）
- [x] 9.5 添加流量百分比计算和颜色指示（正常蓝色、警告黄色、异常红色）
- [x] 9.6 实现过期时间显示和倒计时（已过期、N天后过期、日期）
- [x] 9.7 实现节点数量显示
- [x] 9.8 添加加载状态和空状态提示
- [x] 9.9 添加响应式适配（移动端简化布局）

## 10. 运行信息卡片组件（RuntimeInfoCard）

- [x] 10.1 创建 RuntimeInfoCard.vue 组件（web/src/components/dashboard/RuntimeInfoCard.vue）
- [x] 10.2 实现核心状态显示（运行中/已停止、版本号）
- [x] 10.3 实现核心内存显示（通过 WebSocket 实时更新）
- [x] 10.4 实现混合端口和日志级别显示
- [x] 10.5 集成配置 API 获取核心配置
- [x] 10.6 添加加载状态和错误处理
- [x] 10.7 实现响应式网格布局（桌面端 4 列、平板端 2 列、移动端 2 列）

## 11. 仪表盘页面集成（Dashboard.vue）

- [x] 11.1 修改 Dashboard.vue，移除旧的 4 个基础卡片
- [x] 11.2 集成所有新组件（TrafficGraph、StatCard、ClashModeCard、ActiveSubscriptionCard、RuntimeInfoCard）
- [x] 11.3 实现响应式 Grid 布局（使用 el-row 和 el-col）
- [x] 11.4 配置各个组件的响应式断点（:xs/:sm/:md/:lg/:xl）
- [x] 11.5 设置 Grid 间距响应式（桌面端 16px、移动端 12px）
- [x] 11.6 实现三行布局（代理模式+订阅信息、流量信息、运行信息）
- [x] 11.7 添加卡片图标和标题样式
- [x] 11.8 修复小屏幕布局问题（卡片间距、滚动问题）

## 12. 布局集成（MainLayout.vue）

- [x] 12.1 修改 MainLayout.vue，导入 TrafficIndicator 组件
- [x] 12.2 在移动端 Header（el-header v-if="isMobile"）添加 TrafficIndicator
- [x] 12.3 配置流量指示器位置（居中显示）
- [x] 12.4 添加 header-right 布局（用户名居右）
- [x] 12.5 移除移动端 Header 标题（避免臃肿）
- [x] 12.6 初始化 WebSocket 连接（onMounted 中调用 connect）
- [x] 12.7 修复小屏幕滚动问题（flex: 1、min-height: 0、overflow-y: auto）

## 13. WebSocket 数据管理优化

- [x] 13.1 将 useWebSocket 改为单例模式（共享数据）
- [x] 13.2 将 useTrafficGraph 改为单例模式（共享数据点）
- [x] 13.3 实现订阅类型管理（避免重复订阅）
- [x] 13.4 实现引用计数管理 watcher（防止重复创建/销毁）

## 14. 样式优化和主题适配

- [x] 14.1 统一全局样式（字体、颜色、间距）
- [x] 14.2 添加暗色模式支持（如果需要）
- [x] 14.3 实现高对比度模式（流量指示器和卡片）
- [x] 14.4 优化动画和过渡效果（200-300ms、ease-in-out）
- [x] 14.5 添加阴影和模糊效果（桌面端）
- [x] 14.6 优化滚动条样式（移动端、平板端、桌面端）
