## 实现状态

✅ **全部完成** (133/133 任务)

---

## 1. 项目初始化

- [x] 1.1 创建项目根目录结构 (server/, web/, Makefile)
- [x] 1.2 初始化 server/ Go 项目 (go mod init)
- [x] 1.3 创建 server/ 目录结构 (cmd/, internal/, pkg/)
- [x] 1.4 添加 Go 依赖 (gin, gorm, sqlite driver, gorilla/sessions, gorilla/websocket, gopkg.in/yaml.v3, github.com/dop251/goja)
- [x] 1.5 创建 server/config.yaml 配置文件模板
- [x] 1.6 初始化 web/ Vue3 项目 (Vite + Vue3 + TypeScript)
- [x] 1.7 创建 web/ 目录结构 (api/, components/, pages/, stores/, etc.)
- [x] 1.8 配置根 Makefile (build, run, clean)

## 2. 数据库与模型

- [x] 2.1 创建 internal/model 基础结构
- [x] 2.2 实现 internal/model/user.go (id, username, password, created_at)
- [x] 2.3 实现 internal/model/subscription.go (id, name, url, interval, content, updated_at)
- [x] 2.4 实现 internal/model/config.go (key, value JSON)
- [x] 2.5 实现 internal/model/rule.go (id, name, type, payload, proxy, enabled, mode, priority, created_at)
- [x] 2.6 实现 internal/model/script.go (id, name, description, content, enabled, created_at)
- [x] 2.7 创建数据库初始化和迁移逻辑
- [x] 2.8 实现 internal/repository 层 (user, subscription, config, rule, script)

## 3. 认证系统

- [x] 3.1 实现 pkg/crypto/password.go (bcrypt 加密)
- [x] 3.2 实现 internal/service/auth.go (登录/登出/密码修改逻辑)
- [x] 3.3 实现 internal/middleware/auth.go (Session 验证中间件)
- [x] 3.4 实现 internal/handler/auth.go
  - [x] GET /api/session 登录
  - [x] DELETE /api/session 登出
- [x] 3.5 实现 internal/handler/user.go
  - [x] GET /api/users/me 当前用户
  - [x] PUT /api/users/me/password 修改密码
- [x] 3.6 实现首次启动密码初始化流程
- [x] 3.7 配置 Session 存储 (gorilla/sessions + cookie 或 file)

## 4. 订阅管理

- [x] 4.1 实现 internal/service/subscription.go (拉取/解析)
- [x] 4.2 实现 internal/handler/subscription.go
  - [x] GET /api/subscriptions 列表
  - [x] POST /api/subscriptions 创建
  - [x] GET /api/subscriptions/:id 详情
  - [x] PUT /api/subscriptions/:id 更新
  - [x] DELETE /api/subscriptions/:id 删除
  - [x] POST /api/subscriptions/:id/refresh 刷新
- [x] 4.3 实现 internal/scheduler/subscription.go (定时刷新)

## 5. 规则管理 (自定义)

- [x] 5.1 实现 internal/service/rule.go (CRUD + 排序 + 插入/追加模式)
- [x] 5.2 实现 internal/handler/rule.go
  - [x] GET /api/rules 列表
  - [x] POST /api/rules 创建
  - [x] GET /api/rules/:id 详情
  - [x] PUT /api/rules/:id 更新
  - [x] DELETE /api/rules/:id 删除
- [x] 5.3 实现规则类型验证 (DOMAIN, DOMAIN-SUFFIX, IP-CIDR, etc.)
- [x] 5.4 实现规则模式: insert(插入到订阅规则前) / append(追加到订阅规则后)

## 6. 扩展脚本管理

- [x] 6.1 实现 pkg/script/engine.go (JavaScript 执行引擎 - goja)
- [x] 6.2 实现 internal/service/script.go (CRUD + 执行)
- [x] 6.3 实现 internal/handler/script.go
  - [x] GET /api/scripts 列表
  - [x] POST /api/scripts 创建
  - [x] GET /api/scripts/:id 详情
  - [x] PUT /api/scripts/:id 更新
  - [x] DELETE /api/scripts/:id 删除
  - [x] POST /api/scripts/:id/test 测试脚本执行

## 7. 配置合并引擎

- [x] 7.1 实现 internal/service/merger.go
  - [x] 合并订阅配置
  - [x] 应用自定义规则 (insert 模式插入前，append 模式追加后)
  - [x] 应用面板配置
  - [x] 执行 JavaScript 扩展脚本改写
- [x] 7.2 实现配置验证和错误处理
- [x] 7.3 生成最终 Mihomo 配置 (YAML)

## 8. 配置管理

- [x] 8.1 实现 internal/config/config.go (YAML 配置结构体定义)
- [x] 8.2 实现 internal/config/loader.go (配置加载: 文件 → 环境变量 → 默认值)
- [x] 8.3 实现 internal/handler/config.go
  - [x] GET /api/config 获取
  - [x] PUT /api/config 更新
- [x] 8.4 实现配置验证逻辑
- [x] 8.5 实现配置热更新 (修改 config.yaml 后无需重启)

## 9. Mihomo 核心集成

- [x] 9.1 实现 internal/service/mihomo.go (子进程管理 + HTTP API 封装)
  - [x] Start() 启动子进程
  - [x] Stop() 停止子进程
  - [x] Restart() 重启子进程
  - [x] Reload() 配置热更新
  - [x] Status() 状态获取
  - [x] Version() 版本获取
- [x] 9.2 实现 internal/handler/core.go
  - [x] GET /api/core/version
  - [x] GET /api/core/status
  - [x] POST /api/core/start
  - [x] POST /api/core/stop
  - [x] POST /api/core/restart
  - [x] POST /api/core/reload
- [x] 9.3 实现配置文件生成和 base64 编码

## 10. 代理管理

- [x] 10.1 实现 internal/handler/proxy.go
  - [x] GET /api/proxies 节点列表
  - [x] GET /api/proxies/:name 节点详情
  - [x] PUT /api/proxies/:name 选中节点
  - [x] POST /api/proxies/:name/delay 延迟测试
- [x] 10.2 代理数据从 mihomo service 获取

## 11. 连接管理

- [x] 11.1 实现 internal/handler/connection.go
  - [x] GET /api/connections 连接列表
  - [x] DELETE /api/connections/:id 断开指定
  - [x] DELETE /api/connections 断开所有
- [x] 11.2 连接数据从 mihomo service 获取

## 12. WebSocket 实时推送

- [x] 12.1 实现 internal/ws/hub.go (连接管理)
- [x] 12.2 实现 internal/ws/client.go (客户端连接)
- [x] 12.3 实现 internal/ws/message.go (消息类型定义)
- [x] 12.4 实现 internal/handler/websocket.go (WS 入口 /api/ws)
- [x] 12.5 实现流量推送 (traffic)
- [x] 12.6 实现连接状态推送 (connections)
- [x] 12.7 实现日志推送 (logs)
- [x] 12.8 实现核心状态推送 (core_status)
- [x] 12.9 实现订阅管理 (subscribe/unsubscribe)

## 13. 统一响应格式

- [x] 13.1 实现 pkg/response/response.go (统一响应结构)
- [x] 13.2 定义业务状态码常量

## 14. 前端基础架构

- [x] 14.1 创建 web/src/types/ TypeScript 类型定义
- [x] 14.2 实现 web/src/utils/request.ts (Axios 封装)
- [x] 14.3 实现 web/src/api/ 所有 API 模块
  - [x] auth.ts
  - [x] subscription.ts
  - [x] rule.ts
  - [x] script.ts
  - [x] proxy.ts
  - [x] connection.ts
  - [x] config.ts
  - [x] core.ts
  - [x] statistics.ts
- [x] 14.4 实现 web/src/composables/useWebSocket.ts
- [x] 14.5 配置 web/src/router/index.ts 路由
- [x] 14.6 实现 web/src/stores/ Pinia stores

## 15. 前端页面开发

- [x] 15.1 实现 web/src/pages/Login.vue 登录页
- [x] 15.2 实现 web/src/components/Layout/ 布局组件
- [x] 15.3 实现 web/src/pages/Dashboard.vue 仪表盘
- [x] 15.4 实现 web/src/pages/Proxies.vue 代理页
- [x] 15.5 实现 web/src/pages/Rules.vue 规则页 (展示 + 编辑)
- [x] 15.6 实现 web/src/pages/Connections.vue 连接页
- [x] 15.7 实现 web/src/pages/Settings.vue 设置页
- [x] 15.8 实现 web/src/pages/Subscriptions.vue 订阅管理页
- [x] 15.9 实现 web/src/pages/Scripts.vue 扩展脚本管理页

## 16. 前端组件开发

- [x] 16.1 实现 web/src/components/ProxyCard/ 代理卡片
- [x] 16.2 实现 web/src/components/TrafficGraph/ 流量图表
- [x] 16.3 实现 web/src/components/ConnectionTable/ 连接表格
- [x] 16.4 实现 web/src/components/RuleEditor/ 规则编辑器
- [x] 16.5 实现 web/src/components/ScriptEditor/ 脚本编辑器 (JavaScript 代码高亮)
- [x] 16.6 实现响应式布局样式

## 17. 集成与测试

- [x] 17.1 前后端联调
- [x] 17.2 后端单元测试
- [x] 17.3 构建脚本优化 (前端嵌入后端二进制)
- [x] 17.4 Docker Compose 配置 (可选)

## 18. Core 自动管理重构

### 18.1 后端 - Core 生命周期管理

- [x] 18.1.1 实现 internal/service/core_manager.go
  - [x] NewCoreManager() 构造函数
  - [x] Init() Server 启动时调用
  - [x] Restart() 配置变更时调用
  - [x] GetStatus() 返回运行状态和错误信息
  - [x] GetVersion() 获取版本
  - [x] 生成最小配置 (无订阅时)
- [x] 18.1.2 修改 internal/service/mihomo.go
  - [x] 成为 CoreManager 的内部实现
  - [x] 移除公开的 Start/Stop/Restart 方法
- [x] 18.1.3 修改 internal/service/merger.go
  - [x] 支持无订阅时返回最小配置
  - [x] GetMinimalConfig() 方法

### 18.2 后端 - Handler 变更

- [x] 18.2.1 删除 internal/handler/core.go
- [x] 18.2.2 新增 internal/handler/system.go
  - [x] GET /api/system/info (整合 core + 看板数据)
- [x] 18.2.3 修改 internal/handler/proxy.go
  - [x] Core 未运行时返回空数组
- [x] 18.2.4 修改 internal/handler/connection.go
  - [x] Core 未运行时返回空数组
- [x] 18.2.5 修改 internal/handler/subscription.go
  - [x] CRUD 操作后触发 Core 重启
  - [x] 返回值包含 Core 启动错误 (如果有)
- [x] 18.2.6 修改 internal/handler/rule.go
  - [x] CRUD 操作后触发 Core 重启
- [x] 18.2.7 修改 internal/handler/script.go
  - [x] CRUD 操作后触发 Core 重启
- [x] 18.2.8 修改 internal/handler/config.go
  - [x] 更新配置后触发 Core 重启
- [x] 18.2.9 修改 internal/handler/websocket.go
  - [x] 推送 core_status 消息 (Server 启动后、配置变更后、Core 崩溃时)
- [x] 18.2.10 修改 cmd/server/main.go
  - [x] 启动时初始化 CoreManager
  - [x] 注册 /api/system/info 路由
  - [x] 移除 /api/core/* 路由

### 18.3 前端 - API 层变更

- [x] 18.3.1 删除 web/src/api/core.ts
- [x] 18.3.2 新增 web/src/api/system.ts
  - [x] getSystemInfo() 方法
- [x] 18.3.3 修改 web/src/api/proxy.ts
  - [x] 处理空数据响应
- [x] 18.3.4 修改 web/src/api/connection.ts
  - [x] 处理空数据响应

### 18.4 前端 - 页面/组件变更

- [x] 18.4.1 修改 web/src/composables/useWebSocket.ts
  - [x] 处理 core_status 消息
- [x] 18.4.2 修改 web/src/pages/Dashboard.vue
  - [x] 使用 /api/system/info 获取数据
  - [x] 显示 Core 启动错误 (如果有)
- [x] 18.4.3 修改 web/src/pages/Settings.vue
  - [x] 移除 Core 启停按钮
  - [x] 配置保存后显示 Core 启动结果
- [x] 18.4.4 修改 web/src/pages/Subscriptions.vue
  - [x] 订阅操作后显示 Core 启动结果
- [x] 18.4.5 修改 web/src/pages/Rules.vue
  - [x] 规则操作后显示 Core 启动结果
- [x] 18.4.6 修改 web/src/pages/Scripts.vue
  - [x] 脚本操作后显示 Core 启动结果
- [x] 18.4.7 修改 web/src/components/Layout/
  - [x] 更新状态指示器 (使用 system info)
  - [x] 显示 Core 错误提示

## 19. 配置系统重构

### 19.1 后端 - 配置模块重构

- [x] 19.1.1 删除 config.yaml 文件
- [x] 19.1.2 重写 server/internal/config/config.go
  - [x] ServerConfig 结构体 (Host, Port, Database)
  - [x] CoreConfig 结构体 (APIHost, APIPort, APISecret, MixedPort, AllowLan, Mode, LogLevel, IPv6)
  - [x] 默认值定义
- [x] 19.1.3 重写 server/internal/config/loader.go
  - [x] LoadServerConfig() 从环境变量加载
  - [x] GetCoreConfig() 从数据库读取，fallback 到环境变量/默认值
  - [x] 移除 YAML 加载/保存逻辑
- [x] 19.1.4 修改 server/internal/repository/config.go
  - [x] 支持批量获取/设置 core:* 键
- [x] 19.1.5 修改 server/internal/handler/config.go
  - [x] GET /api/config 返回 Server + Core 配置
  - [x] PUT /api/config/core 批量更新 Core 配置
- [x] 19.1.6 修改 server/cmd/server/main.go
  - [x] 使用新配置加载逻辑
  - [x] Session 密钥随机生成
  - [x] Session 有效期改为浏览器关闭失效
- [x] 19.1.7 修改 server/internal/service/*.go
  - [x] 更新配置引用方式 (config.GetCoreConfig())
  - [x] 移除 config.Global 相关引用

### 19.2 前端 - 配置页面重构

- [x] 19.2.1 修改 web/src/types/api.ts
  - [x] 新 ServerConfig 类型
  - [x] 新 CoreConfig 类型 (含 mixed_port, allow_lan, bind_address, mode, log_level, ipv6)
- [x] 19.2.2 修改 web/src/api/config.ts
  - [x] getConfig() 返回新格式
  - [x] updateConfig() 使用 PUT /api/config
- [x] 19.2.3 修改 web/src/stores/config.ts
  - [x] 适配新配置结构
- [x] 19.2.4 修改 web/src/pages/Settings.vue
  - [x] Server 配置只读展示
  - [x] Core 配置可编辑 (含所有字段)
  - [x] 批量保存按钮

## 20. 订阅激活状态

### 20.1 后端 - 模型与仓库层

- [x] 20.1.1 修改 server/internal/model/subscription.go
  - [x] 添加 Active 字段
- [x] 20.1.2 修改 server/internal/repository/subscription.go
  - [x] GetActive() 获取激活订阅
  - [x] SetActive(id) 设置激活订阅 (事务)
  - [x] Count() 统计订阅数量

### 20.2 后端 - 服务层

- [x] 20.2.1 修改 server/internal/service/subscription.go
  - [x] Create() 首个订阅自动激活
  - [x] Activate(id) 激活指定订阅
  - [x] GetActive() 获取激活订阅 (无激活时自动激活第一个)
- [x] 20.2.2 修改 server/internal/service/merger.go
  - [x] Merge() 使用激活的订阅 (移除 subscriptionID 参数)
- [x] 20.2.3 修改 server/internal/service/core.go
  - [x] Start()/Restart() 调用 Merge() 不再传参

### 20.3 后端 - Handler 层

- [x] 20.3.1 修改 server/internal/handler/subscription.go
  - [x] Delete() 删除激活订阅时自动激活下一个
  - [x] 新增 Activate() 激活 API
- [x] 20.3.2 修改 server/cmd/server/main.go
  - [x] 添加 PUT /api/subscriptions/:id/activate 路由

### 20.4 前端 - API 与类型

- [x] 20.4.1 修改 web/src/types/api.ts
  - [x] Subscription 添加 active 字段
- [x] 20.4.2 修改 web/src/api/subscription.ts
  - [x] 添加 activate(id) 方法

### 20.5 前端 - 页面

- [x] 20.5.1 修改 web/src/pages/Subscriptions.vue
  - [x] 显示激活状态徽章
  - [x] 激活订阅边框高亮
  - [x] 添加激活按钮 (非激活订阅显示)
  - [x] 激活操作触发 Core 重启

## 21. WebSocket 实时消息增强

### 21.1 后端 - CoreService 整合 WebSocket 代理

- [x] 21.1.1 修改 server/internal/service/core.go
  - [x] 添加 Hub 引用和 SetHub 方法
  - [x] 添加 WebSocket 连接管理 (wsConns, wsStopChan, wsDialer)
  - [x] startWSProxies / stopWSProxies 方法
  - [x] runWSEndpoint / connectWSEndpoint 方法
  - [x] Core 启动后自动启动 WebSocket 代理
  - [x] Core 停止后自动停止 WebSocket 代理
- [x] 21.1.2 删除 server/internal/ws/mihomo_proxy.go
- [x] 21.1.3 修改 server/internal/ws/message.go
  - [x] 添加 ConnectionsData, LogData, MemoryData 类型
  - [x] 添加 TypeMemory 常量
- [x] 21.1.4 修改 server/internal/ws/hub.go
  - [x] 添加 BroadcastToTypeRaw 方法
- [x] 21.1.5 修改 server/internal/handler/websocket.go
  - [x] InitWebSocketHub 返回 *ws.Hub
- [x] 21.1.6 修改 server/cmd/server/main.go
  - [x] 调用 coreService.SetHub(hub)
  - [x] 移除 mihomoProxy 相关代码

### 21.2 前端 - 类型定义

- [x] 21.2.1 修改 web/src/types/api.ts
  - [x] 添加 ConnectionsData, LogData, MemoryData 类型
  - [x] 清理重复的 Connection 定义

### 21.3 前端 - WebSocket 处理

- [x] 21.3.1 修改 web/src/composables/useWebSocket.ts
  - [x] 添加 connections, logs, memory 响应式数据
  - [x] 处理 connections, logs, memory 消息类型
  - [x] 添加 clearLogs 方法
  - [x] 日志自动限制 500 条

### 21.4 前端 - 页面更新

- [x] 21.4.1 修改 web/src/pages/Connections.vue
  - [x] 使用 WebSocket 实时获取连接数据
  - [x] 显示连接数和流量统计
  - [x] 移除轮询逻辑