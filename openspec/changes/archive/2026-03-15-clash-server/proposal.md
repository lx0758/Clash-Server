## Why

当前缺乏一个轻量级的自托管代理管理面板。Clash Verge 是优秀的桌面客户端，但需要本地安装。用户需要一个可以通过浏览器访问的 Web 面板，能够接入机场订阅、管理 Mihomo 核心、查看实时流量和连接。

## What Changes

创建全新的 ClashServer 项目，包含：

- **后端服务**: Go + Gin + GORM + SQLite，提供 RESTful API 和 WebSocket
- **Web 前端**: Vue3 响应式布局，UI/UX 模仿 Clash Verge
- **核心集成**: Mihomo 子进程自动管理，随 Server 启停，配置变更自动重启
- **订阅管理**: 支持机场 Clash 订阅接入，本地缓存，定时刷新，激活状态管理
- **规则管理**: 支持手动添加、编辑、删除自定义规则 (insert/append 模式)
- **扩展脚本**: 支持 JavaScript 脚本改写传入 Mihomo 的配置 (goja 引擎)
- **认证系统**: Session-based 单用户认证，首次启动密码初始化

### 核心架构原则

- **Server 就绪 = 所有 API 可用**: Core 状态不影响 API 可用性
- **无订阅 = 空数据，不是错误**: 未配置订阅时返回空数组，Core 运行最小配置
- **Core 是实现细节**: 用户不关心 Core 启停，只关心"有代理"或"无代理"
- **配置变更 = 立即重启 Core**: 订阅/规则/脚本/配置变更后立即重启并返回结果

### 实现状态

✅ 已完成全部功能开发 (119/119 任务)

## Capabilities

### New Capabilities

- `clash-core-integration`: Mihomo 核心自动管理，随 Server 启停，配置变更自动重启
- `subscription-management`: 机场订阅的接入、解析、存储和定时刷新
- `proxy-management`: 代理节点列表、选中、延迟测试
- `traffic-monitoring`: 实时流量、连接数、统计数据的 WebSocket 推送
- `user-authentication`: Session-based 单用户登录认证
- `web-api`: 完整的 RESTful API 和 WebSocket 接口
- `rule-management`: 自定义规则的 CRUD、优先级管理和插入/追加模式
- `script-management`: JavaScript 扩展脚本的管理和执行引擎

### Modified Capabilities

- (无) 这是全新项目

## Impact

- 新增代码库: 后端 (Go) 和前端 (Vue3)
- 依赖: Mihomo 核心库，Go 生态 (Gin, GORM, Gorilla Sessions, goja)，Vue 生态 (Vue3, Vite, Pinia)
- 部署: 二进制分发，用户自行 Docker 化