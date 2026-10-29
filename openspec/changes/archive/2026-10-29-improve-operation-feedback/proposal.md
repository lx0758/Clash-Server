# Proposal

## Why

面板的**操作反馈链路完全断裂**：点击「测延迟」整页代理组塌成骨架屏、点击「刷新订阅」界面零变化、所有网络错误提示落进没有渲染层的 toast store 而用户永远看不到。更严重的是后端在核心重载失败时仍返回 `code:0`，前端因此弹出「订阅已刷新」——系统会在操作失败时说谎，用户对面板的信任被反复消耗。

## What Changes

### 反馈机制重建

- **对象级 pending 状态**: 新增按 key 索引的 pending 状态集合（`delay:GROUP_A`、`refresh:sub_12`、`select:HK-01`），取代当前 1 个全局 `loading` 布尔撑起 7 种操作的模式
- **反馈通道复活**: `toast` store 挂载 `ToastHost` 渲染组件，使 `request.ts` 拦截器已有的错误提示真正可见；`loading` store 挂载顶部细进度条，仅慢操作（>300ms）触发，避免快请求闪烁
- **反馈分层**: 全局结果用 toast，页面级用进度条，操作级用按钮转圈 + 局部骨架；局部操作不得影响其他对象的可交互性
- **乐观更新**: 切换代理节点先本地高亮，请求失败时回滚并提示

### 测延迟改造为逐节点增量反馈

当前「测延迟」调用 Mihomo 的整组测速端点 `GET /group/{group}/delay`，该端点**只返回一个时间点**——必须等最慢的节点测完，全部数据同时返回。经实测确认：

- 整组端点返回全量 map，**只有 1 个响应时间点**，节点数据一起到达
- 单节点端点 `GET /proxies/{name}/delay` 每次调用**独立返回一个结果**，多个调用并发时各自按完成顺序到达
- Mihomo 核心对并发测速请求**并行处理**：3 个节点（响应耗时 0.1s / 1s / 2s）并发请求总耗时 2.047s，而非串行的 3.1s

因此改造为：后端对目标组内的节点**并发逐个调用单节点测速**，每测完一个节点就通过 WebSocket 推送该节点的延迟结果；前端收到即刷新该节点卡片，未完成的节点各自显示独立的进行中状态。

- 快节点的延迟不再被慢节点阻塞
- 每个节点有独立的进行中状态，与其他节点、其他组互不影响
- 该组测延迟不再触发整页骨架屏，也不再整表替换代理数据

### 后端语义修正

- **三态结果模型**: 引入 `partial_success` 标记区分「完全成功」「已保存但核心未生效」「完全失败」，`code:1007` 承载半成功，前端据此选择黄色警告而非红色错误
- **修复 `core_error` 读取路径**: 前端 `checkCoreError` 读取层级错误，导致订阅页的 Core 错误提示从未生效过
- **修复 `needsRestart` 判断**: `ApplyConfig` 先写入 `lastConfig` 再拿它与新配置比较，`DeepEqual` 恒真，`SoftRestart` 分支不可达——修改 mixed-port / secret / allow-lan 等需重启项时核心不重启却返回「配置已保存」
- **修复失败静默吞没**: 合并失败不再静默回退最小配置；`apiRequest` 保留 Mihomo 错误响应体；代理/规则/连接查询失败不再伪装成空数组 + `code:0`
- **接通 `need_confirm`**: `response.NeedConfirmRestart` 已定义但零调用方，自定义配置保存后核心应用失败时接入回滚逻辑

### 破坏性变更

- **BREAKING**: 配置变更类接口（`PUT /api/config`、订阅增删改/激活/刷新、自定义配置保存）在核心重载失败时，`code` 从 `0` 变为 `1007`，响应体新增 `partial_success` 字段。所有前端调用点必须同步处理

## Capabilities

### New Capabilities

- `operation-feedback`: 操作反馈的通用契约——pending 状态的对象级粒度、三层反馈通道的职责边界、操作结果的三态语义

### Modified Capabilities

- `web-api`: 错误响应增加 `detail` 字段；业务状态码表补充 `1007` 与半成功语义
- `hot-reload`: 修正「智能判断配置变更类型」的判断逻辑；接通 `need_confirm`；实现「软重启失败降级为杀进程重启」
- `clash-core-integration`: 「配置变更 API 返回 Core 错误」场景从 `code:0` 改为 `code:1007` + `partial_success`
- `proxy-ui`: 测延迟改为逐节点增量反馈（每个节点独立进行中状态、延迟结果陆续到达）；切换节点的乐观更新与失败回滚
- `subscription-management`: 刷新/激活的对象级反馈；核心未生效必须在界面上可见
- `traffic-monitoring`: 新增 `delay` 类型的 WebSocket 推送，承载测延迟过程中单个节点的延迟结果

## Impact

**前端**:

- `web/src/stores/` — 新增 `pending.ts`（对象级状态集合）；`toast.ts`/`loading.ts` 接线渲染层
- `web/src/components/ToastHost.vue`、`web/src/components/LoadingBar.vue` — 新增
- `web/src/utils/request.ts` — 拦截器接入真实错误通道，处理 `partial_success`
- `web/src/pages/Proxies.vue`、`web/src/pages/Subscriptions.vue`、`web/src/pages/Settings.vue` — 操作改为对象级 pending
- `web/src/components/proxy/ProxyHead.vue`、`ProxyGrid.vue`、`ProxyItem.vue` — 接收 pending 状态，局部骨架替代整页骨架
- `web/src/components/SubscriptionCard.vue` — 启用已声明但从未使用的 `loading` prop
- `web/src/composables/useWebSocket.ts` — 订阅并处理 `delay` 类型消息，逐个更新节点延迟
- `web/src/stores/proxy.ts` — 启用已存在但零调用方的 `checkDelay`，改为逐节点并发测速

**后端**:

- `server/pkg/response/response.go` — `partial_success` 字段、`SuccessWithCoreError` 语义调整
- `server/internal/service/core.go` — `ApplyConfig` 的 `needsRestart` 判断顺序；`apiRequest` 错误信息透传；合并失败不再静默回退；新增逐节点并发测速并逐个推送结果
- `server/internal/ws/message.go` — 新增 `delay` 消息类型
- `server/internal/handler/` — `subscription.go`、`config.go`、`customization.go`、`proxy.go`、`rule.go` 的失败语义
- `server/internal/service/customization.go` — 接入 `restoreCustomization` 回滚

**兼容性**: 前端与后端需同批次发布。后端先发布会使旧前端把 `code:1007` 当失败处理（误报但不至于崩溃）；前端先发布会漏读 `partial_success`（半成功被显示为普通错误）。`delay` 类型的 WebSocket 消息为新增，旧前端不订阅该类型时不受影响。

**已知副作用**: 修复 `needsRestart` 后，修改端口类配置将真正触发 Mihomo 软重启，**当前代理连接会中断**。此为使配置生效的必要代价，需在界面上明确提示。

**测延迟改造的代价**: 后端需对组内每个节点各发起一次测速调用，节点数量大时对核心的请求数相应增加（核心侧并行处理，但请求总数上升）。后端需限制并发数以避免瞬时压垮核心。
