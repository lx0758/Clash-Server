## Context

ClashServer 是一个自托管的代理管理 Web 面板，后端使用 Go + Gin，前端使用 Vue3，核心使用 Mihomo。用户通过浏览器访问面板，可以管理机场订阅、切换代理节点、查看实时流量。

**实现状态**: 已完成全部功能开发 (119/119 任务)。

## Goals / Non-Goals

**Goals:**
- 提供完整的 Web 管理界面，功能和 UI 模仿 Clash Verge
- 支持接入机场 Clash 订阅，本地缓存，定时刷新
- 通过 HTTP API 控制 Mihomo 核心（启动、停止、重载配置）
- 实时流量和连接数据通过 WebSocket 推送
- Session-based 单用户认证
- 支持规则手动编辑和自定义规则组
- 支持自定义 Clash 扩展脚本改写配置

**Non-Goals:**
- 多用户支持
- 节点数据持久化存储（只读订阅数据）
- 移动端原生应用
- 日志/流量统计持久化存储

## Decisions

### 1. 核心通信模式: 子进程自动管理 + HTTP API

```
┌─────────────────────────────────────────────────────────────┐
│                    Server 启动流程                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Server 启动                                                │
│      ├── 初始化数据库、Session、WS Hub                      │
│      ├── 加载配置                                           │
│      ├── 合并配置 (有订阅用订阅，无订阅用最小配置)           │
│      ├── 启动 Core (成功/失败都继续)                        │
│      └── 启动 HTTP 服务 (所有 API 立即可用)                 │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────┐      HTTP       ┌─────────────────┐
│   Backend   │◀───────────────▶│  Mihomo Child   │
│   (Go)      │    :9090/ui     │   Process       │
└─────────────┘                 └─────────────────┘
配置传递: mihomo -d {tmpdir} -f -
```

**设计原则:**
- Server 就绪时所有 API 可用，不依赖 Core 状态
- 无订阅时 Core 运行最小配置 (mode: direct, rules: [MATCH,DIRECT])
- Core 是实现细节，用户不关心启停，只关心"有代理"或"无代理"

**备选方案:**
- 手动启停 API - 增加用户心智负担，已弃用
- 独立进程 (Systemd/Docker) - 需要用户自行管理 Mihomo 生命周期
- 嵌入式 (cgo 绑定) - 部署简单但灵活性差

**选择理由:** 自动管理模式由后端统一管理，部署简单，用户体验最佳。

### 2. API 设计: RESTful + 业务状态码

**HTTP 动词规范:**
- GET: 获取资源
- POST: 创建资源 / 执行操作
- PUT: 完整更新资源
- PATCH: 部分更新资源
- DELETE: 删除资源

**响应格式:**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

业务状态码与 HTTP 状态码含义一致，避免与网络层状态码混淆。

**核心 API 列表:**
```
# 认证
/api/session          GET(登录) / DELETE(登出)
/api/users/me         GET / PUT(password)

# 订阅
/api/subscriptions    GET / POST
/api/subscriptions/:id  GET / PUT / DELETE
/api/subscriptions/:id/refresh  POST

# 代理
/api/proxies          GET
/api/proxies/:name    GET / PUT(select)
/api/proxies/:name/delay  POST

# 连接
/api/connections      GET / DELETE
/api/connections/:id  DELETE

# 规则 (自定义)
/api/rules            GET / POST
/api/rules/:id        GET / PUT / DELETE

# 扩展脚本
/api/scripts          GET / POST
/api/scripts/:id      GET / PUT / DELETE

# 配置
/api/config           GET / PUT

# 系统信息 (整合 core version/status + 看板数据)
/api/system/info      GET

# 统计 (实时，不持久化)
/api/statistics/traffic  GET
```

**已移除的 API (Core 自动管理):**
```
/api/core/version     → 整合到 /api/system/info
/api/core/status      → 整合到 /api/system/info
/api/core/start       → 移除 (自动启动)
/api/core/stop        → 移除 (自动管理)
/api/core/restart     → 移除 (配置变更自动触发)
/api/core/reload      → 移除 (改为重启)
```

### 3. WebSocket: 实时数据推送

**连接:** `ws://host:port/api/ws`

**消息格式:**
```json
// Server → Client
{"type": "traffic", "data": {"up": 100, "down": 200}}
{"type": "connections", "data": [...]}
{"type": "logs", "data": [...]}
{"type": "core_status", "data": {"running": true, "version": "v1.18.1", "error": null}}

// Client → Server
{"action": "subscribe", "types": ["traffic", "connections"]}
{"action": "unsubscribe", "types": ["logs"]}
```

**core_status 推送时机:**
- Server 启动后 (无论成功失败)
- 配置变更触发重启后 (无论成功失败)
- Core 崩溃时

### 4. 数据库: SQLite + GORM

**核心表:**
- `users`: id, username, password(bcrypt), created_at
- `subscriptions`: id, name, url, interval, content, active, updated_at
- `config`: key, value(JSON)
- `rules`: id, name, type, payload, proxy, enabled, mode(insert/append), priority, created_at
- `scripts`: id, name, description, content, enabled, created_at

**数据关系:**
```
┌─────────────────┐     ┌─────────────────┐
│ subscriptions   │     │ proxy_groups    │
│ (机场订阅)       │     │ (自定义规则组)   │
└────────┬────────┘     └────────┬────────┘
         │                       │
         ▼                       ▼
┌─────────────────────────────────────────┐
│              合并引擎                    │
│  订阅配置 + 面板配置 + 自定义规则/规则组  │
└────────────────┬────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│              扩展脚本                    │
│  JavaScript/Lua 改写最终配置            │
└────────────────┬────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│            Mihomo Core                  │
└─────────────────────────────────────────┘
```

### 5. 前端架构: Vue3 + 响应式

- 使用 Vue3 Composition API
- 布局和交互完全模仿 Clash Verge
- 页面: 仪表盘、代理、规则、设置
- 响应式布局，支持桌面和移动端

### 6. 配置合并策略

激活订阅时，按以下顺序合并配置：

```
┌─────────────────────────────────────────────────────────────┐
│                     配置合并流程                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 机场订阅 (基础)                                         │
│     └─ proxies, proxy-providers, rule-providers            │
│                                                             │
│  2. 自定义规则 (插入/追加)                                  │
│     └─ mode=insert: 插入到订阅规则之前                      │
│     └─ mode=append: 追加到订阅规则之后                      │
│                                                             │
│  3. 面板配置 (覆盖)                                         │
│     └─ 端口、DNS、日志级别等系统设置                        │
│                                                             │
│  4. 扩展脚本 (改写)                                         │
│     └─ 执行启用的 JavaScript 脚本，对配置进行最终改写       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**合并规则:**
- 面板配置优先覆盖订阅的端口、DNS、日志级别等设置
- 自定义规则支持两种模式：
  - `insert`: 插入到订阅规则之前 (高优先级，先匹配)
  - `append`: 追加到订阅规则之后 (低优先级，兜底)
- 扩展脚本最后执行，可任意改写配置

### 7. 配置系统 (环境变量 + 数据库)

配置分为两类：Server 配置和 Core 配置。Server 配置通过环境变量设置，Core 配置存储在数据库中。

**Server 配置 (环境变量, 不可编辑)**

| 环境变量 | 默认值 | 说明 |
|---------|-------|------|
| CS_SERVER_HOST | 0.0.0.0 | 服务监听地址 |
| CS_SERVER_PORT | 7000 | 服务监听端口 |
| CS_SERVER_DATABASE | data.db | 数据库文件路径 |

**Core 配置 (数据库, 可编辑)**

| 数据库键 | 环境变量 (初始化用) | 默认值 | 说明 |
|---------|-------------------|-------|------|
| core:api_host | CS_CORE_API_HOST | 127.0.0.1 | Mihomo API 地址 |
| core:api_port | CS_CORE_API_PORT | 9090 | Mihomo API 端口 |
| core:api_secret | CS_CORE_API_SECRET | (空) | Mihomo API 密钥 |
| core:mixed_port | CS_CORE_MIXED_PORT | 7890 | HTTP/SOCKS5 混合端口 |
| core:allow_lan | CS_CORE_ALLOW_LAN | true | 允许局域网连接 |
| core:mode | CS_CORE_MODE | rule | 运行模式 (rule/global/direct) |
| core:log_level | CS_CORE_LOG_LEVEL | info | 日志级别 (silent/error/warning/info/debug) |
| core:ipv6 | CS_CORE_IPV6 | false | 启用 IPv6 |

**硬编码配置:**
- Mihomo 二进制名称: `clash/mihomo` (通过 PATH 查找)

**Session 配置**
- 会话有效期: 浏览器关闭时失效 (session cookie)
- 密钥: 每次启动随机生成

**Core 配置读取优先级:**
```
数据库 > 环境变量 > 默认值
```

**API 设计:**
```
GET  /api/config          获取 Server + Core 配置
PUT  /api/config          批量更新 Core 配置，触发 Core 重启
```

**GET /api/config 响应:**
```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 7000,
    "database": "data.db"
  },
  "core": {
    "api_host": "127.0.0.1",
    "api_port": 9090,
    "api_secret": "",
    "mixed_port": 7890,
    "allow_lan": true,
    "mode": "rule",
    "log_level": "info",
    "ipv6": false
  }
}
```

**PUT /api/config 请求:**
```json
{
  "core": {
    "api_host": "127.0.0.1",
    "api_port": 9091,
    "api_secret": "my-secret",
    "mixed_port": 7890,
    "allow_lan": true,
    "mode": "rule",
    "log_level": "info",
    "ipv6": false
  }
}
```

**PUT /api/config 响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "core": { ... },
    "core_error": ""
  }
}
```

### 9. Core 自动管理架构

**生命周期:**
```
┌─────────────────────────────────────────────────────────────┐
│                Core 生命周期                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Server 启动 → Core 启动 (有订阅/无订阅都启动)              │
│                                                             │
│  配置变更 (订阅/规则/脚本/面板配置) → 立即重启 Core        │
│      ├── 生成新配置                                         │
│      ├── 停止旧 Core                                        │
│      ├── 启动新 Core                                        │
│      └── 返回结果 (成功/失败+错误信息)                      │
│                                                             │
│  Core 崩溃 → 保持停止，记录错误，等待配置变更触发重启      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**无订阅时的最小配置:**
```yaml
mixed-port: 7890
allow-lan: true
bind-address: "*"
mode: direct
log-level: info
ipv6: false
proxies: []
proxy-groups:
  - name: DIRECT
    type: select
    proxies: []
rules:
  - MATCH,DIRECT
```

**系统信息 API (GET /api/system/info):**
```json
{
  "core": {
    "running": true,
    "version": "v1.18.1",
    "error": null
  },
  "subscription": {
    "count": 2,
    "proxy_count": 156
  },
  "traffic": {
    "up": 1024,
    "down": 20480
  }
}
```

**配置变更 API 返回格式:**
```json
// 成功
{"code": 0, "message": "success", "data": null}

// 保存成功但 Core 启动失败
{"code": 0, "message": "success", "data": {"core_error": "yaml: line 45: ..."}}
```

### 10. 项目代码结构

```
clash-server/
├── server/                    # 后端 (Go)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # 入口，路由定义，CORS 内联
│   ├── internal/
│   │   ├── config/            # 配置加载
│   │   │   ├── config.go      # 配置结构体
│   │   │   └── loader.go      # 环境变量加载
│   │   ├── handler/           # HTTP handlers
│   │   │   ├── auth.go        # 认证相关
│   │   │   ├── subscription.go
│   │   │   ├── proxy.go
│   │   │   ├── connection.go
│   │   │   ├── rule.go        # 自定义规则
│   │   │   ├── script.go      # 扩展脚本
│   │   │   ├── config.go
│   │   │   ├── system.go      # 系统信息 (整合 core + 看板数据)
│   │   │   ├── user.go        # 用户管理
│   │   │   ├── scheduler.go   # 调度器集成
│   │   │   └── websocket.go
│   │   ├── middleware/        # 中间件
│   │   │   └── auth.go        # Session 验证
│   │   ├── model/             # 数据模型
│   │   │   ├── db.go          # 数据库初始化
│   │   │   ├── model.go       # 通用模型
│   │   │   ├── user.go
│   │   │   ├── subscription.go
│   │   │   ├── rule.go        # 自定义规则
│   │   │   ├── script.go      # 扩展脚本
│   │   │   └── config.go
│   │   ├── repository/        # 数据访问层
│   │   │   ├── user.go
│   │   │   ├── subscription.go
│   │   │   ├── rule.go
│   │   │   ├── script.go
│   │   │   └── config.go
│   │   ├── service/           # 业务逻辑
│   │   │   ├── auth.go
│   │   │   ├── subscription.go
│   │   │   ├── rule.go        # 规则管理
│   │   │   ├── script.go      # 脚本执行引擎
│   │   │   ├── merger.go      # 配置合并引擎
│   │   │   ├── core.go        # Core 自动管理 (整合 mihomo 进程管理)
│   │   │   ├── core_unix.go   # Unix 平台支持
│   │   │   └── core_windows.go # Windows 平台支持
│   │   ├── scheduler/         # 定时任务
│   │   │   └── subscription.go
│   │   └── ws/                # WebSocket
│   │       ├── hub.go         # 连接管理
│   │       ├── client.go      # 客户端
│   │       └── message.go     # 消息类型
│   ├── pkg/
│   │   ├── response/          # 统一响应
│   │   │   └── response.go
│   │   ├── crypto/            # 加密工具
│   │   │   └── password.go
│   │   └── script/            # JavaScript 脚本引擎
│   │       └── engine.go      # goja 执行器
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
│
├── web/                       # 前端 (Vue3)
│   ├── src/
│   │   ├── api/               # API 请求
│   │   │   ├── auth.ts
│   │   │   ├── subscription.ts
│   │   │   ├── proxy.ts
│   │   │   ├── connection.ts
│   │   │   ├── rule.ts        # 规则 API
│   │   │   ├── script.ts      # 脚本 API
│   │   │   ├── config.ts
│   │   │   └── system.ts      # 系统信息
│   │   ├── components/        # 通用组件
│   │   │   ├── Layout/
│   │   │   │   └── MainLayout.vue
│   │   │   ├── ProxyCard/
│   │   │   ├── TrafficGraph/
│   │   │   ├── ConnectionTable/
│   │   │   ├── RuleEditor/    # 规则编辑器
│   │   │   ├── ScriptEditor/  # 脚本编辑器
│   │   │   ├── Loading.vue    # 加载组件
│   │   │   └── Toast.vue      # 通知组件
│   │   ├── composables/       # 组合式函数
│   │   │   └── useWebSocket.ts
│   │   ├── pages/             # 页面
│   │   │   ├── Login.vue
│   │   │   ├── Dashboard.vue
│   │   │   ├── Proxies.vue
│   │   │   ├── Rules.vue
│   │   │   ├── Connections.vue
│   │   │   ├── Settings.vue
│   │   │   ├── Subscriptions.vue
│   │   │   └── Scripts.vue
│   │   ├── router/
│   │   │   └── index.ts
│   │   ├── stores/            # Pinia stores
│   │   │   ├── auth.ts
│   │   │   ├── config.ts
│   │   │   ├── loading.ts     # 加载状态
│   │   │   ├── proxy.ts
│   │   │   ├── system.ts      # 系统状态
│   │   │   └── toast.ts       # 通知状态
│   │   ├── types/             # TypeScript 类型 (统一)
│   │   │   └── api.ts         # 所有类型定义
│   │   ├── utils/
│   │   │   └── request.ts     # Axios 封装
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── style.css          # 全局样式
│   ├── public/
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   └── Makefile
│
├── Makefile                   # 根 Makefile
└── README.md
```

**实现说明:**
- CORS 中间件内联在 `main.go` 中，未单独抽文件
- `core.go` 整合了 Core 管理和 Mihomo 进程封装
- 前端类型定义统一在 `types/api.ts`
- 前端样式使用 `style.css` (根目录)
- 新增 `loading.ts`, `toast.ts`, `system.ts` stores 用于状态管理

## Risks / Trade-offs

| 风险 | 描述 | 缓解措施 |
|------|------|----------|
| 配置变更导致连接中断 | 每次配置变更都会重启 Core | 接受此 trade-off，重启速度快 |
| Core 启动失败 | 错误配置导致 Core 无法启动 | 显式暴露错误，用户修正后自动恢复 |
| 订阅 URL 安全 | 订阅 URL 包含敏感 token | 存储时加密，API 返回时脱敏 |
| 日志量过大 | Mihomo 日志可能很大 | 配置日志轮转，默认 INFO 级别 |
| 子进程崩溃 | Mihomo 崩溃后服务中断 | 更新状态，等待用户修正配置 |

## Migration Plan

1. 编译后端二进制和前端静态资源
2. 初始化 SQLite 数据库
3. 启动后端服务，配置端口和认证密码
4. 前端构建产物通过嵌入或反向代理访问
5. 用户首次访问时设置管理员密码

## Open Questions

- [x] ~~是否需要支持 Docker 一键部署？~~ 暂不支持，用户可自行 Docker 化