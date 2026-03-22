## 1. 基础设施配置

- [x] 1.1 安装 Element Plus 和图标库 (`npm install element-plus @element-plus/icons-vue`)
- [x] 1.2 配置 Element Plus 按需引入 (使用 unplugin-vue-components)
- [x] 1.3 创建浅色主题 CSS 变量文件 (`src/styles/theme.css`)
- [x] 1.4 在 main.ts 中注册 Element Plus 和主题
- [x] 1.5 创建 `useBreakpoint` 响应式断点 hook

## 2. 布局改造

- [x] 2.1 创建 `MainLayout.vue` 使用 el-container 和 el-menu
- [x] 2.2 实现侧边栏折叠功能 (el-menu :collapse)
- [x] 2.3 创建移动端抽屉导航组件
- [x] 2.4 实现导航折叠状态 localStorage 持久化
- [x] 2.5 测试布局在桌面端和移动端的响应式

## 3. Login 页面改造

- [x] 3.1 使用 el-card 重写登录卡片
- [x] 3.2 使用 el-form 重写登录表单
- [x] 3.3 使用 ElMessage 替代现有错误提示
- [x] 3.4 实现登录页面响应式布局 (桌面居中/移动全宽)
- [x] 3.5 测试首次登录密码设置流程

## 4. Dashboard 页面改造

- [x] 4.1 使用 el-row/el-col 重写状态卡片网格
- [x] 4.2 使用 el-card 重写状态卡片
- [x] 4.3 使用 el-statistic 展示流量数据
- [x] 4.4 使用 el-tag 显示核心运行状态
- [x] 4.5 使用 el-alert 显示核心错误
- [x] 4.6 实现响应式列数 (xl:4, lg:3, md:2, sm:1)

## 5. Settings 页面改造

- [x] 5.1 使用 el-form 重写配置表单
- [x] 5.2 使用 el-input/el-select/el-switch 替换表单控件
- [x] 5.3 实现表单标签响应式 (桌面右侧/移动顶部)
- [x] 5.4 使用 ElMessage 显示保存结果
- [x] 5.5 测试配置保存和 Core 重启流程

## 6. Subscriptions 页面改造

- [x] 6.1 使用 el-row/el-col 重写订阅卡片网格
- [x] 6.2 使用 el-card 重写订阅卡片
- [x] 6.3 使用 el-progress 显示流量进度
- [x] 6.4 使用 el-tag 显示订阅状态标签
- [x] 6.5 实现响应式卡片布局

## 7. SubscriptionEditDialog 改造

- [x] 7.1 使用 el-dialog 重写对话框
- [x] 7.2 使用 el-form 重写表单
- [x] 7.3 使用 el-radio-group 替换类型选择器
- [x] 7.4 使用 el-switch 替换复选框
- [x] 7.5 实现移动端对话框全屏
- [x] 7.6 使用 el-collapse 实现高级设置折叠

## 8. Proxies 页面改造

- [x] 8.1 使用 el-radio-group 重写模式切换器
- [x] 8.2 使用 el-collapse 重写代理组折叠面板
- [x] 8.3 使用 el-input 重写筛选框
- [x] 8.4 使用 el-select 重写排序选择器
- [x] 8.5 保留 ProxyGrid 自定义组件，添加响应式列数
- [x] 8.6 实现移动端代理组全宽显示

## 9. Logs 页面改造

- [x] 9.1 使用 el-input 重写搜索框
- [x] 9.2 改用卡片列表替代 el-table (更适合移动端)
- [x] 9.3 使用 el-tag 显示日志级别
- [x] 9.4 使用 el-button 重写操作按钮
- [x] 9.5 实现响应式布局

## 10. Connections 页面改造

- [x] 10.1 改用卡片列表替代 el-table (更适合移动端)
- [x] 10.2 使用 el-tag 显示连接状态
- [x] 10.3 实现响应式布局
- [x] 10.4 添加搜索和筛选功能

## 11. Rules 页面改造

- [x] 11.1 改用卡片列表替代 el-table (更适合移动端)
- [x] 11.2 使用 el-tag 显示规则类型
- [x] 11.3 实现响应式布局
- [x] 11.4 添加搜索功能

## 12. 清理和优化

- [x] 12.1 删除废弃的 Toast.vue 组件 (使用 ElMessage)
- [x] 12.2 删除废弃的 Loading.vue 组件 (使用 ElLoading)
- [x] 12.3 清理 style.css 中冗余样式
- [x] 12.4 更新 README 文档
- [x] 12.5 执行全页面响应式测试
- [x] 12.6 执行浅色主题一致性检查

## 13. UI 细节优化

- [x] 13.1 日志页面：修复布局问题，添加时间戳显示
- [x] 13.2 日志页面：错误/警告日志添加左侧边框和背景色突出显示
- [x] 13.3 日志页面：使用等宽字体显示日志内容
- [x] 13.4 代理Item：添加类型标签显示
- [x] 13.5 代理Item：在线显示延迟数值，离线显示"离线"标签
- [x] 13.6 代理Item：移动端响应式适配（iPhone 12 Pro 测试通过）
- [x] 13.7 连接页面：chains 顺序反转（从代理到目标）
- [x] 13.8 连接页面：chains 添加渐变背景强化显示效果
- [x] 13.9 规则页面：payload 和 proxy 位置调整，proxy 更醒目
- [x] 13.10 订阅页面：规则编辑器对话框响应式适配
- [x] 13.11 订阅页面：脚本管理器对话框响应式适配

## 14. 布局优化

- [x] 14.1 订阅页面：移除最大宽度限制，卡片横向铺满
- [x] 14.2 订阅卡片：小屏幕两行布局 (标题+按钮 / 类型+节点+更新时间)
- [x] 14.3 订阅卡片：PC端按钮保持默认大小，移动端按钮缩小
- [x] 14.4 规则页面：宽屏多列网格布局 (1400px:4列, 1000px:3列, 640px:2列, <640px:1列)
- [x] 14.5 规则/脚本对话框：小屏幕列表区域可滚动，footer固定底部
- [x] 14.6 设置页面：表单标签左侧对齐，保存按钮独立区域
- [x] 14.7 设置页面：优化卡片布局和按钮位置

## 15. Bug 修复

- [x] 15.1 规则页面：修复直接打开显示"核心未运行"问题（添加 fetchSystemInfo 调用）
- [x] 15.2 规则页面：添加核心状态变化监听，自动获取规则
- [x] 15.3 订阅卡片：修复点击操作按钮触发展开/收缩问题