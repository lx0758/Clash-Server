# Design

## Context

见 proposal.md - Why。此处只记录影响技术方案的现状约束。

**前端现状约束**

| 现状 | 位置 | 对方案的约束 |
|---|---|---|
| `toast` store 的 `toasts` 数组无任何消费方 | `web/src/stores/toast.ts:12` | store 本身可直接复用，只需补渲染层 |
| `loading` store 的 `count`/`message` 无任何消费方 | `web/src/stores/loading.ts` | 同上；但其"全局遮罩"语义与本次"对象级"目标冲突 |
| `request.ts` 拦截器已把错误写入 `toastStore.error()` | `web/src/utils/request.ts:46,74` | 通道修好后错误提示自动生效，无需改拦截器错误分支 |
| `Proxies.vue` 整页 `<el-skeleton :loading="proxyStore.loading">` | `web/src/pages/Proxies.vue:8` | `proxyStore.loading` 被"初始加载"与"测延迟后刷新"共用，需解耦 |
| `CheckGroupDelay` 调 Mihomo `/group/{name}/delay`，整组端点只返回一个时间点 | `server/internal/service/core.go:226` | 需改为逐节点并发测速 + 增量推送 |
| `proxyStore.checkDelay`（单节点测速）已实现但零 UI 调用方 | `web/src/stores/proxy.ts:160` | 逐节点方案正好可复用 |
| `ws.MessageType` 已定义 5 种类型，`BroadcastToType` 支持按类型订阅 | `server/internal/ws/message.go`、`hub.go` | 新增 `delay` 类型即可复用现成机制 |
| `server/go.mod` 无 `golang.org/x/sync` | `server/go.mod` | 限流器需自行实现（带缓冲 channel 信号量），不引入新依赖 |
| 后端 `http.Client` 超时 10s，无自定义 Transport | `server/internal/service/core.go:70` | 并发测速请求共用该 client，需确认 Transport 的连接池配置 |
| `SubscriptionCard.vue:135` 声明了 `loading?: boolean` prop 但父组件从未传、模板也未用 | `web/src/components/SubscriptionCard.vue` | 死 prop 可直接启用，但需支持"卡片级"而非"整卡遮罩" |
| 刷新动作通过 `el-dropdown-item` 触发 | `SubscriptionCard.vue` 模板 | `el-dropdown-item` 无 loading 能力，需改为在卡片头部独立呈现状态 |
| `handleSaveSub`/`handleActivate`/`handleRefresh`/`handleDelete` 共用一个 `loading` 布尔 | `web/src/pages/Subscriptions.vue:88` | 必须拆分为按操作标识的独立状态 |
| `checkCoreError(res.data)` 读取层级错误 | `web/src/pages/Subscriptions.vue:119` | 后端把 `core_error` 放在 `data` 内，前端应读 `res.data.data.core_error` |
| `api/proxy.ts` 用已解包的 `get/put/patch` 助手，其余 api 文件用 `request.get` 返回完整响应 | `web/src/api/` | 两套取值方式并存，是 `checkCoreError` 读错层级的根源，需统一 |

**后端现状约束**

| 现状 | 位置 | 对方案的约束 |
|---|---|---|
| `ApplyConfig` 先写 `cs.lastConfig = mergedConfig` 再调 `needsRestart(cs.getLastConfig(), mergedConfig)` | `server/internal/service/core.go:604-607` | 自比恒真，`SoftRestart` 不可达 |
| `SuccessWithCoreError` 把 `core_error` 塞进 `data` 但 `code` 恒为 0 | `server/pkg/response/response.go:13` | 现有 7 处调用点需一并迁移语义 |
| `response.CodeCoreError`(1007) 与 `response.NeedConfirmRestart` 已定义但零调用方 | `server/pkg/response/response.go:32,59` | 可直接启用，无需新增常量 |
| `apiRequest` 非 2xx 时返回 `fmt.Errorf("API request failed")` | `server/internal/service/core.go:298` | 丢弃了核心返回的错误体 |
| 代理/规则/连接查询失败时返回空结果集 + `code:0` | `handler/proxy.go:23`、`handler/rule.go:24` | 需改为非零码，但会影响前端"核心未运行"分支 |
| `merger.Merge()` 失败时静默回退 `GetMinimalConfig()` | `server/internal/service/core.go:595` | 改为返回错误 |
| `restoreCustomization` 已定义但零调用方 | `server/internal/handler/customization.go:130` | 可直接接入回滚 |
| 所有 handler 返回 HTTP 200，业务状态在 body 的 `code` | 全部 handler | `web-api` spec 明确此约定，本次不改 HTTP 状态码语义 |
| `openspec/specs/` 下 15 个 main spec 误用 delta 头 `## ADDED Requirements` | `openspec/specs/*/spec.md` | 归档本变更的 delta 时会被拒绝，需先修复这 5 个相关 spec 的结构 |

## Goals / Non-Goals

**Goals:**

- 建立对象级 pending 机制，使任意操作的状态可按操作标识独立查询与清除
- 建立统一的三层反馈通道（对象级 / 页面级 / 全局），每层职责不重叠
- 后端操作结果能区分完全成功、部分成功、完全失败三类
- 修复三个导致"成功说谎"的缺陷，使界面结论可信
- 前端 api 层取值方式统一，消除层级误读

**Non-Goals:**

- 不做暗色主题与 CSS 变量化（视觉一致性，另行处理）
- 不做虚拟滚动与列表分页（性能，另行处理）
- 不引入全局状态管理库（Pinia 已足够）
- 不改动 HTTP 状态码语义（保持全 200 + body code 的既有约定）
- 不重构 WebSocket 协议

## Decisions

### 决策 1：pending 状态用 key 集合而非嵌套对象

**选择**：`Set<string>` 语义的 `ref`，key 形如 `delay:GROUP_A`、`refresh:12`、`select:HK-01`、`activate:12`、`mode:global`。

```
pending store
   |
   +-- keys: Ref<Set<string>>
   +-- isPending(key)      -> boolean
   +-- start(key)          -> 加入集合
   +-- finish(key)         -> 从集合移除
   +-- withPending(key, fn)  -> 包装异步操作，自动 start/finish
```

**理由**：

1. **粒度天然由 key 表达**，不需要额外的 `loadingGroups: string[]` + `loadingSubs: number[]` 这类平行结构
2. **`withPending` 包装器保证清理**。裸 `start()`/`finish()` 极易在 `catch` 分支漏掉 `finish`，导致状态永久卡死；包装器用 `finally` 兜底
3. **Set 的增删是 O(1) 且天然去重**，同一 key 的重复触发不会重复计数

**考虑过的替代方案**：

| 方案 | 否决理由 |
|---|---|
| 每种操作一个独立的 `ref<boolean>` | 操作一多就是十几个平行布尔，且无法表达"同一个列表里两个对象各自进行中" |
| 嵌套对象 `Record<string, boolean>` | 语义等价但每次增删触发整对象响应式更新；Set 更省 |
| 组件内各自维护 `loading` ref | 状态无法跨组件共享（如卡片知道自己 loading，但父组件不知道），重复代码多 |
| 引入 VueUse 的 `useAsyncState` | 只能管理单个异步状态，无法覆盖"同类型多实例"场景 |

**与现有 `proxyStore.loading` 的关系**：`proxyStore.loading` 保留，但语义收窄为**仅表示代理列表的初始加载**。测延迟不再经过 `fetchProxies()` 的全局 loading 路径，改为局部 pending + 局部数据合并。

### 决策 2：pending 状态与 toast 分属不同层

```
用户点击「测延迟」
      |
      +-- 对象级: 该组按钮转圈 + 该组节点卡片显示测速中   <- pending store
      |
      +-- 全局结果: 测速完成 toast「已测试 12 个节点」   <- toast store
```

**理由**：pending 是"进行中"的连续状态，toast 是"已完成"的瞬时事件。混在一起会导致 toast 频繁弹出打断操作节奏。两者用不同 store、不同生命周期，避免互相污染。

**考虑过的替代方案**：用 toast 承载进行中状态（如"测延迟中..."的 info toast）。否决——每个并发操作都要一条 toast，多个组同时测延迟时屏幕被 toast 淹没，且 toast 会自动消失无法表达"仍在进行"。

### 决策 3：慢操作阈值 300ms

**选择**：pending 状态**立即**写入 store（保证按钮立刻禁用、防重复提交），但**视觉呈现**延迟 300ms 才显示。300ms 内完成则不渲染任何加载指示。

**理由**：
- 立即禁用是必要的 —— 否则用户可以连点"测延迟"触发多个并发请求
- 视觉延迟是必要的 —— 多数操作（切换节点、切换模式）在本地网络下 <100ms，立刻显示骨架会造成刺眼闪烁

**实现**：pending store 额外维护 `visibleKeys`（延迟后才对 UI 暴露）。或由组件侧用 `useTimeoutFn` 处理。倾向**组件侧**处理：pending store 保持纯粹的状态容器，延迟逻辑作为通用 composable `useDelayedPending(key)` 提供。

**考虑过的替代方案**：全局硬编码 300ms 阈值。否决——阈值应该是可配置的（不同操作耗时差异大，测延迟可能 30 秒，切换节点 50ms）。但**先按全局默认 300ms 实现**，保留按操作覆盖的能力。

### 决策 4：测延迟改为逐节点增量反馈（经实测验证）

**实测依据**（真实 mihomo v1.19.22，3 个响应耗时不同的节点 0.1s / 1s / 2s）：

```
实验A: GET /group/G1/delay          <- 当前项目的方式
  [ 2.040s] <<< 唯一一次响应 >>> {"fast":190,"mid":990,"slow":1980}
  只有 1 个时间点。快节点被压着等最慢节点，数据一起到达

实验B: GET /proxies/{name}/delay     <- 每节点一个独立请求
  [ 0.105s] <<< fast-node 独立到达 >>>
  [ 1.006s] <<< mid-node  独立到达 >>>
  [ 2.006s] <<< slow-node 独立到达 >>>
  3 个独立时间点，快节点先显示

并发验证: 3 节点并发请求总耗时 2.047s (最慢节点单独 2.0s, 串行则 3.1s)
  => Mihomo 核心并行处理测速请求
```

**选择**：后端对目标组内节点**并发逐个调用单节点测速端点**，每完成一个即通过 WebSocket 推送该节点结果；前端订阅 `delay` 类型消息，逐个更新节点卡片。

```
用户点击「测延迟」(组 G1, 12 个节点)
       |
       +-> POST /api/proxies/group/G1/delay-test   立即返回 202
       |     |
       |     +-> 后端起 goroutine, 遍历组内节点
       |           |
       |           +-> 并发(限流 N) 调用 core /proxies/{name}/delay
       |           |      |
       |           |      +-> 每完成一个 -> hub.BroadcastToType("delay", {group, name, delay})
       |           |         -> WS 推送 -> 前端更新该节点卡片
       |           |
       |           +-> 全部完成 -> 推送 {group, done: true}
       |
       +-> 前端收到 delay 消息
             |
             +-> 该节点: pending 移除, 延迟值刷新
             +-> 收到 done: 该组 pending 移除, toast「已测试 12 个节点」
```

**理由**：
1. **快节点不必等慢节点**。当前实现下 0.1s 的节点要等 2s 的节点才显示，这个延迟是纯粹的浪费
2. **每个节点有独立的进行中状态**。这正是客户端转圈圈效果的来源，也正好落在"细化到单个对象"的既定目标上
3. **消除整页骨架屏**。当前 `checkGroupDelay` 内部调 `fetchProxies()`，`fetchProxies` 置 `proxyStore.loading = true`，触发 `Proxies.vue:8` 的整页 `<el-skeleton>`——7 个代理组同时塌陷、滚动归零、DOM 全量重建
4. **消除整表替换**。逐节点推送只改单个节点的 delay 字段，不触碰 `rawProxies` 整体

**关键实现点**：
- 逐节点推送**不替代**最终的 `GET /proxies` 刷新。测速结果同时要写回 Mihomo 的代理历史（供 `history` 字段读取），前端在收到 `done` 后仍需拉一次列表以同步核心侧状态
- 但这次刷新**不能**触发整页骨架——`fetchProxies` 需要一个 `silent` 模式，或拆分为「静默刷新数据」与「显示整页骨架」两条路径

### 决策 4b：并发编排放在后端而非前端

**选择**：并发测速由**后端**编排（goroutine + 限流），结果经 WebSocket 推送；前端只负责发一次触发请求并订阅结果。

**理由**：

1. **浏览器同源并发上限 6**。前端逐节点发请求时，HTTP/1.1 下每域名并发上限为 6，100 个节点会自动分批 —— 虽然能工作，但批次边界会造成不均匀的刷新节奏，且浏览器会为每个节点建立/复用连接
2. **WS 连接是复用的**。项目已有 WebSocket 通道，新增 `delay` 类型只需在 `ws/message.go` 加一个常量 + 复用现成的 `BroadcastToType`（已支持按类型订阅）
3. **限流可控**。后端可以用信号量把并发压到合理值（如 10），避免节点多时瞬时压垮核心；前端无法做这种控制
4. **错误语义集中**。单节点测速失败时后端统一转成"该节点超时"，前端不需要处理各种异常形态

**考虑过的替代方案**：

| 方案 | 否决理由 |
|---|---|
| 前端逐节点并发请求 | 受浏览器同源并发上限 6 约束，批次化刷新不均匀；无法对核心做限流；连接开销大 |
| 后端并发 + 轮询查询进度 | 需要额外的进度存储与查询端点，比 WS 推送复杂且有延迟 |
| 保持整组端点，仅修前端整页骨架 | 治标不治本——快节点仍被慢节点阻塞，且拿不到"每节点独立转圈"的效果 |
| 继续用整组端点但让核心分批返回 | Mihomo 的 `/group/{name}/delay` 不支持分批返回，无此能力 |

### 决策 4c：并发限流阈值

**选择**：后端并发测速的并发上限设为 **10**，通过带缓冲 channel 实现的信号量控制。

**理由**：
- Mihomo 核心本身对并发测速是并行的（实测已验证），10 并发足以让快节点迅速返回
- 100 个节点在 10 并发下约需 10 批，比 1 并发快一个数量级，又不至于瞬时产生 100 个请求
- 需为限流器设置**可配置项**，因为不同部署规模（10 节点 vs 500 节点）差异很大

### 决策 5：三态结果模型与 `partial_success`

**选择**：后端在"已持久化但核心未生效"时返回：

```json
{
  "code": 1007,
  "message": "订阅已刷新，但核心重载失败",
  "data": {
    "subscription": { ... },
    "core_error": "API request failed: 403 Forbidden",
    "partial_success": true
  }
}
```

前端按 `partial_success` 区分：

| `partial_success` | 视觉 | 文案 |
|---|---|---|
| 无 / `false` + `code:0` | 绿色成功 | "订阅已刷新" |
| `true` + `code:1007` | 黄色警告 | "订阅已刷新，但核心未生效：xxx" |
| `false` + `code≠0` | 红色错误 | "刷新失败：xxx" |

**理由**：核心未生效时**数据已落库**，用户不需要重试，重试反而可能造成重复副作用（重复拉取订阅、重复写配置）。用红色错误会误导用户重复操作。黄色警告 + 明确的后续动作指引才是准确的表达。

**考虑过的替代方案**：

| 方案 | 否决理由 |
|---|---|
| 保持 `code:0` + `core_error` 在 data 内 | 这是当前做法，问题在于前端读错层级导致提示从未生效；即使修好层级，客户端仍需自行区分三态，不如后端明确表达 |
| 直接返回 `code:1007` 不带 `partial_success` | 丢失"数据已保存"这一事实，客户端无法区分"要不要重试" |
| 引入 `status` 字段（`success`/`partial`/`failed`）三值枚举 | 语义更明确，但与项目既有的 `code` 数字体系并存会造成两套真相；`partial_success` 布尔 + `code` 组合已足够表达 |

**与既有 spec 的冲突**：`web-api` 与 `clash-core-integration` 的 main spec 明确要求"配置保存成功 + Core 启动失败 → 返回 `{"code": 0, ...}`"。本次变更修订该要求（已在 delta spec 中体现）。

### 决策 6：`core_error` 的层级契约

**选择**：后端保持 `core_error` 位于 `data` 对象内（`SuccessWithCoreError` 的既有行为不变），前端统一从 `res.data.data.core_error` 读取。

**理由**：把 `core_error` 提到响应顶层需要改动 `ApiResponse` 的公共结构，且 `data` 内已有 `subscription` 等同级字段时顶层放置会造成语义层级混乱。**真正的 bug 是前端读错层级**，修复点在调用方而非数据结构。

**考虑过的替代方案**：把 `core_error` 提升到响应顶层。否决——`data` 是"本次操作的结果载荷"，`core_error` 描述的是"载荷的生效状态"，放在 `data` 内更内聚。

**配套决策：统一 api 层取值方式**。当前 `api/proxy.ts` 用已解包的 `get/put/patch`（返回 `data`），其余 api 文件用 `request.get`（返回完整响应）。两套并存正是层级误读的温床。**全部 api 文件统一使用 `request.get` 返回完整响应**，业务代码显式取 `res.data.data`。这会让 `checkCoreError` 这类代码的层级变得显式可读。

### 决策 7：`needsRestart` 的修复方式

**选择**：`ApplyConfig` 中先取旧配置快照并完成比较，再写入新配置：

```go
oldConfig := cs.getLastConfig()          // 先取快照
needsRestart := cs.needsRestart(oldConfig, mergedConfig)
cs.mu.Lock()
cs.lastConfig = mergedConfig              // 后写入
cs.mu.Unlock()
```

**理由**：当前实现先赋值后比较，`getLastConfig()` 返回的就是刚写入的 `mergedConfig`，`DeepEqual` 恒真，导致 `SoftRestart` 分支永不可达。修改 mixed-port / secret / allow-lan / tun / ipv6 时核心不重启，变更不生效，但接口返回"配置已保存"。

**副作用（必须告知用户）**：修复后修改端口类配置**会真正触发 Mihomo 软重启**，当前代理连接会中断。这是使配置生效的必要代价 —— 但界面必须明确提示，不能让用户在不知情的情况下丢失连接。

**附带修复**：`SoftRestart` 失败时按 `hot-reload` spec 要求降级为杀进程重启（当前未实现）。降级也失败时必须返回错误，不得以 `code:0` 掩盖。

### 决策 8：核心错误信息透传

**选择**：`apiRequest` 在非 2xx 时读取并返回响应体内容：

```go
if resp.StatusCode < 200 || resp.StatusCode >= 300 {
    body, _ := io.ReadAll(resp.Body)
    return fmt.Errorf("core API %s %s failed (%d): %s", method, path, resp.StatusCode, strings.TrimSpace(string(body)))
}
```

**理由**：当前返回 `"API request failed"`，用户看到的是无信息量的文案。核心返回的错误（如 YAML 语法错误的具体行号）才是用户能据此行动的依据。`hot-reload` spec 已要求"内核保持旧配置运行"且需暴露错误信息。

### 决策 9：自定义配置失败回滚

**选择**：在 `CustomizationHandler.Update` 中，保存前读取旧的自定义配置，`ApplyConfig` 失败时调用既有的 `restoreCustomization` 回滚，并返回错误。

**理由**：`restoreCustomization` 已实现但零调用方。当前行为：自定义配置已写入 DB，但核心应用失败 —— DB 状态与运行态不一致，用户下次重启面板时配置才突然生效，期间无法解释"为什么改了没反应"。

**考虑过的替代方案**：不回滚，仅返回错误。否决——留下"DB 里有但核心没生效"的隐性不一致状态，用户会陷入无法理解的困惑。

### 决策 10：反馈通道的渲染位置

```
+----------------------------------------------------------+
| App.vue                                                  |
|   <ToastHost />        右上角堆叠, z-index 最高          |
|   <LoadingBar />       顶部 2px 细进度条, 页面级          |
|   <router-view />                                        |
+----------------------------------------------------------+
```

- `ToastHost` 读取 `toastStore.toasts`，渲染 `el-notification` 或自绘
- `LoadingBar` 读取 `loadingStore`，但**只对慢操作显示**（配合决策 3 的阈值）
- `loadingStore` 当前的"全局遮罩"语义**不实现** —— 全屏遮罩在桌面浏览器管理场景下过于粗暴，与"对象级"目标冲突。只做顶部细条。

### 决策 11：代理切换的乐观更新与回滚

**选择**：`selectProxy` 先写入本地 `now` 状态，请求失败时回滚到原值。

```
点击 JP-02
   |
   +-> 立即更新本地 now = "JP-02"    <- 界面立即响应
   +-> await API
        |
        +-- 成功: 保持
        +-- 失败: 回滚 now = 原值 + 错误提示
```

**理由**：当前实现 `await` 完才更新 `now`，网络慢时点击完全无响应。代理切换是高频低延迟操作，乐观更新体感提升明显。

**风险**：乐观期间界面显示"已切换"但实际未切换，若用户此时刷新页面可能困惑。通过失败时的回滚 + 提示缓解。

### 决策 12：核心不可用时的查询失败语义

**选择**：`handler/proxy.go`、`handler/rule.go`、`handler/connection.go` 在核心不可用时返回非零业务码（如 `CodeCoreError`），而非空结果集 + `code:0`。

**理由**：当前核心崩溃时规则页显示"暂无规则"，用户无法区分"真的没有规则"和"核心挂了"。

**影响**：前端需要新增"核心不可用"分支处理。`Rules.vue` 已有 `coreRunning` 状态和"核心未运行"的空状态分支，可复用该模式。

## Risks / Trade-offs

**[修复 needsRestart 后改端口会真的重启核心，连接中断]** → 界面上在触发端口类变更前明确提示"此操作将重启核心并断开当前连接"。这是使配置生效的必要代价，但必须让用户知情。

**[`code` 从 0 变为 1007 是接口契约变更]** → 前后端必须同批次发布。需在后端响应中保留 `core_error` 字段与 `partial_success` 标记，使旧版前端至少能显示错误而非崩溃。

**[主 spec 结构损坏，delta 无法归档]** → `openspec/specs/` 下 15 个 main spec 误用了 delta 头 `## ADDED Requirements`，导致其需求对 validate/list/archive 不可见。本变更涉及的 5 个（`web-api`、`hot-reload`、`proxy-ui`、`clash-core-integration`、`traffic-monitoring`）必须先修正为 `## Requirements` 结构，否则归档被拒。此修复独立于本变更，需单独执行。

**[统一 api 层取值方式会触及所有 api 文件与调用点]** → 改动面比看起来大，但属于一次性机械替换。替代方案是保持两套并存、仅修复 `checkCoreError` 一处 —— 治标不治本，同类层级 bug 会再次发生。

**[pending key 命名不统一会导致状态串扰]** → 在 design 中固化 key 命名规则，并在 `withPending` 中对 key 做类型前缀校验。

**[顶部进度条的可见阈值暂按全局默认 300ms 实现]** → 不做按操作类型的差异化配置，保留后续调整空间。当前所有操作共用一个阈值，简单可预测。

**[合并失败不再静默回退最小配置，可能让原本"能用的面板"变成报错]** → 这是正确行为（之前用户拿到的是一份与预期完全不同的配置且毫不知情）。但需确保错误信息清晰，说明是合并失败而非其他原因。

**[逐节点测速对核心的请求数上升]** → 组内 N 个节点产生 N 次测速调用（当前整组端点为 1 次）。后端以并发上限 10 限流，且该上限可配置。实测证明核心侧并行处理，10 并发下 100 节点约需 10 批次完成。

**[测速触发接口变为异步]** → 新接口触发后立即返回，测速结果经 WS 推送。需要定义"触发成功但测速中途中断"的状态：前端在收到 `done` 前若连接断开，应将未完成的节点标记为超时而非永久转圈。

**[WS 推送丢失导致节点永久转圈]** → `BroadcastToType` 在客户端缓冲满时**静默丢弃**消息（`hub.go:76-77`）。需要在前端为每组测速设置超时兜底（如 15 秒后自动清除 pending），以及在收到 `done` 时统一收敛所有节点状态。

**[测速结果需写回核心历史]** → 逐节点调用 `/proxies/{name}/delay` 时 Mihomo 会自行记录延迟历史，因此最终仍需一次 `GET /proxies` 同步核心侧状态。但该刷新必须是静默的，不得触发整页骨架。

**[后端需新增异步任务的生命周期管理]** → 测速 goroutine 需在核心停止或服务关闭时可取消，避免 goroutine 泄漏与向已关闭的 hub 广播。

**[`loadingStore` 的全局遮罩语义被废弃]** → 现有无消费方，实际影响为零。但若未来有代码依赖全局遮罩，需重新设计。

## Migration Plan

1. **前置**：修复 5 个相关 main spec 的结构（`## ADDED Requirements` → `## Requirements`），使 delta 可归档
2. **后端先行**：部署后端（`code:1007` + `partial_success`），此时旧前端会把半成功显示为普通错误 —— 有提示但语义不准，可接受
3. **前端跟进**：部署前端（pending 机制、toast 渲染层、三态处理），与后端对齐
4. **回滚策略**：前后端均为单一二进制/静态包，可独立回滚。后端回滚后前端仍能工作（前端把 `code:1007` 当普通错误处理）；前端回滚后后端返回 `code:1007` 会被旧前端显示为错误 —— 可接受但不精确

## 刷新状态的位置（已决策）

订阅卡片的 meta 行已有一个时间元素（`{{ formatRelativeTime(subscription.last_refresh) }}更新`）。刷新进行中时，**在该位置替换为旋转图标 + "更新中"**，而非改动 dropdown 菜单结构。

```
刷新前                              刷新中
+------------------------------+    +------------------------------+
| 远程订阅  156 节点            |    | 远程订阅  156 节点            |
| 3 分钟前更新                  |    | ⟳ 更新中                     |
+------------------------------+    +------------------------------+
```

**理由**：
- 该位置本就是"此订阅上次变更时间"的语义槽位，"更新中"是它最自然的替代态
- 完全不改动现有 dropdown 交互结构，规避了 `el-dropdown-item` 无 loading 能力的问题
- 旋转图标比文字更能表达"进行中"而非"已更新"

同一槽位也用于表达激活进行中（"⟳ 激活中"），保持卡片反馈位置的一致性。

## 主 spec 修复范围（已决策）

`openspec/specs/` 下 15 个 main spec 误用 delta 头 `## ADDED Requirements`，导致其需求对 validate/list/archive 不可见。本变更**只修复 delta 直接涉及的 5 个**（`web-api`、`hot-reload`、`proxy-ui`、`clash-core-integration`、`traffic-monitoring`），其余 10 个另开变更处理。

修复方式：将首个 delta 头替换为 `## Requirements`，并把误置其中的 `## REMOVED Requirements` 段落内容归入 `## Requirements`（保留其"已移除"标注的 requirement 名称）。修复后 `openspec list --specs` 应能正确显示这 5 个能力的 requirement 数量。
