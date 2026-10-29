# Tasks

## 1. 前置：修复 main spec 结构

delta 归档要求目标 spec 结构合法。当前 `openspec/specs/` 下 15 个 main spec 误用 delta 头 `## ADDED Requirements`，其需求对 validate/list/archive 不可见。

- [x] 1.1 修复 `openspec/specs/web-api/spec.md` —— 首个 delta 头改为 `## Requirements`，验证 `openspec show web-api --type spec --json --no-scenarios` 不再报 "Main spec contains delta header" 且能返回 requirement 列表
- [x] 1.2 修复 `openspec/specs/hot-reload/spec.md` —— 同上，验证 `openspec show hot-reload --type spec --json --no-scenarios` 正常返回
- [x] 1.3 修复 `openspec/specs/proxy-ui/spec.md` —— 同上，验证 `openspec show proxy-ui --type spec --json --no-scenarios` 正常返回
- [x] 1.4 修复 `openspec/specs/clash-core-integration/spec.md` —— 同上，注意该文件含 `## REMOVED Requirements` 段落需一并归入 `## Requirements`，验证 `openspec show clash-core-integration --type spec --json --no-scenarios` 正常返回
- [x] 1.5 修复 `openspec/specs/traffic-monitoring/spec.md` —— 同上，验证 `openspec show traffic-monitoring --type spec --json --no-scenarios` 正常返回
- [x] 1.6 运行 `openspec validate improve-operation-feedback --strict`，验证 5 个 delta 的 "Archive would refuse" 提示全部消失且变更仍为 valid

## 2. 后端：错误语义与状态码

- [x] 2.1 在 `server/pkg/response/response.go` 增加 `partial_success` 字段支持 —— 新增构造函数承载"已持久化但未生效"语义（code 1007 + data 内含 core_error 与 partial_success=true），验证单元测试断言返回体的 code 与两个字段
- [x] 2.2 迁移 `handler/subscription.go` 的 6 处 `SuccessWithCoreError` 调用 —— 核心重载失败时改用新构造函数返回 code 1007，验证订阅创建/更新/删除/激活/刷新/内容写入共 6 个操作在核心失败时均返回非零 code 且 data 保留已保存对象
- [x] 2.3 迁移 `handler/config.go` 的 `SuccessWithCoreError` 调用 —— `PUT /api/config` 核心重启失败时返回 code 1007 + partial_success，验证保存成功但核心失败时响应体符合新契约
- [x] 2.4 修改 `handler/proxy.go` 的 `List`/`Get`/`GetMode` —— 核心不可用时返回非零业务码而非空结果集 + code 0，验证核心停止时 `GET /api/proxies` 返回非零 code 且响应体说明核心未运行
- [x] 2.5 修改 `handler/rule.go` 的 `GetRules` 与 `handler/connection.go` 的查询接口 —— 同样返回非零业务码，验证核心停止时规则与连接查询不再返回空数组伪装成功
- [x] 2.6 修改 `handler/proxy.go` 的 `CheckDelay` —— 延迟测试失败时返回核心错误信息而非 `delay:0`，验证核心不可用时该接口返回非零 code 且含失败原因

## 3. 后端：核心应用逻辑修复

- [x] 3.1 修复 `server/internal/service/core.go` 的 `ApplyConfig` —— 先取 `lastConfig` 快照完成 `needsRestart` 比较再写入新配置，验证修改 mixed-port/secret/allow-lan/tun/ipv6 时走 `SoftRestart` 分支（单元测试覆盖 `needsRestart` 对端口类差异返回 true）
- [x] 3.2 实现 `SoftRestart` 失败降级为杀进程重启 —— 按 `hot-reload` spec 要求降级，降级也失败时返回错误而非 `HotReloaded: true`，验证两条失败路径均返回非空 Error
- [x] 3.3 移除 `ApplyConfig` 中合并失败静默回退最小配置的逻辑 —— 合并失败时返回错误，验证合并失败时核心继续运行原有配置且错误信息可见
- [x] 3.4 修改 `apiRequest` 的非 2xx 分支 —— 读取响应体内容并包含在返回错误中，验证核心返回 YAML 语法错误等具体信息时该错误出现在 `core_error` 中而非泛化的 "API request failed"
- [x] 3.5 在 `handler/customization.go` 的 `Update` 中接入 `restoreCustomization` 回滚 —— 保存前读取旧配置，核心应用失败时回滚并返回错误，验证核心应用失败后 DB 中的自定义配置恢复为保存前状态
- [x] 3.6 校验 `ApplyConfigResult.Restarted` 字段 —— 软重启路径正确置位而非恒为 `HotReloaded: true`，验证软重启与热重载两条路径的返回字段各自正确

## 4. 前端：反馈基础设施

- [x] 4.1 新增 `web/src/stores/pending.ts` —— 基于 `Set<string>` 的对象级进行中状态，提供 `isPending(key)`/`start(key)`/`finish(key)`/`withPending(key, fn)`，其中 `withPending` 用 `finally` 保证清理，验证连续两次 `withPending` 后 `isPending` 为 false 且并发不同 key 互不影响
- [x] 4.2 新增 `web/src/composables/useDelayedPending.ts` —— 基于 pending store 的延迟可见封装（默认 300ms 阈值，可按操作覆盖），验证快速完成的操作不渲染加载指示而超时操作在阈值后显示
- [x] 4.3 新增 `web/src/components/ToastHost.vue` —— 读取 `toastStore.toasts` 渲染通知列表，挂在 `App.vue` 中，验证 `request.ts` 拦截器触发的错误提示在界面上可见
- [x] 4.4 新增 `web/src/components/LoadingBar.vue` —— 读取 `loadingStore` 在页面顶部渲染细进度条（不实现全屏遮罩），挂在 `App.vue` 中，验证慢请求期间顶部出现进度条且不遮挡页面内容
- [x] 4.5 统一 `web/src/api/` 下的取值方式 —— 全部改用 `request.get/put/post/delete` 返回完整响应，消除 `get/put/patch` 助手的已解包形态，验证各 api 方法的调用方取值层级一致且 `res.data.data` 能取到业务数据

## 5. 后端：逐节点并发测速与增量推送

- [x] 5.1 在 `server/internal/ws/message.go` 新增 `delay` 消息类型与数据体 —— 定义携带 group、name、delay、done 字段的结构，验证 `ws.MessageType` 常量与结构体可被 `BroadcastToType` 使用
- [x] 5.2 在 `server/internal/service/delay.go` 新增逐节点并发测速方法 —— 取目标组节点列表，用带缓冲 channel 信号量限流（默认并发 10）并发调用单节点测速端点，每完成一个即回调，验证 3 个响应耗时不同的节点（0.1s/1s/2s）并发总耗时接近最慢者而非三者之和
- [x] 5.3 在 `server/internal/handler/proxy.go` 新增测延迟触发接口 —— 接收组名后立即返回并启动异步测速，验证接口在测速进行中即返回而不阻塞至最慢节点完成
- [x] 5.4 将测速结果经 `hub.BroadcastToType("delay", ...)` 推送 —— 每节点结果与组完成信号分别推送，验证订阅 delay 类型的客户端能逐条收到结果且最终收到 done 信号
- [x] 5.5 实现测速任务的取消与生命周期管理 —— 核心停止或服务关闭时可取消进行中的测速，验证服务关闭后无 goroutine 泄漏且不向已关闭的 hub 广播
- [x] 5.6 限制单节点测速的并发数以保护核心 —— 并发上限可配置，验证 100 节点场景下并发请求数不超过配置值
- [x] 5.7 处理单节点测速失败 —— 转为该节点的失败结果并继续测试同组其他节点，验证单节点失败不影响其余节点且失败节点以超时语义呈现

## 6. 前端：代理页操作级反馈

- [x] 6.1 在 `web/src/composables/useWebSocket.ts` 订阅并处理 `delay` 类型消息 —— 按 group/name 更新对应节点延迟，收到 done 时清理该组全部节点状态，验证快节点的结果先于慢节点显示
- [x] 6.2 改造 `web/src/stores/proxy.ts` —— 用逐节点测速替换 `checkGroupDelay`，移除测延迟路径对 `fetchProxies()` 的调用，验证测延迟期间 `proxyStore.loading` 保持 false
- [x] 6.3 为 `fetchProxies` 增加静默模式 —— 测延迟完成后同步核心侧延迟历史但不触发整页骨架，验证该刷新不导致页面塌陷与滚动跳动
- [x] 6.4 改造 `web/src/components/proxy/ProxyHead.vue` —— 接收目标组的 pending 状态，测延迟按钮显示 loading 并禁用，验证该组测延迟时按钮不可重复触发且其他组的按钮不受影响
- [x] 6.5 改造 `web/src/components/proxy/ProxyGrid.vue` 与 `ProxyItem.vue` —— 支持节点级 pending，每个节点各自显示进行中状态，验证组 A 测延迟时组 B 的节点仍可点击且组内各节点转圈圈独立
- [x] 6.6 改造 `web/src/pages/Proxies.vue` —— 用组级 pending 替代 `proxyStore.loading` 对整页 `<el-skeleton>` 的绑定，验证测延迟时页面其余部分不塌陷且滚动位置保持不变
- [x] 6.7 实现测延迟中断的超时兜底 —— 结果通道断开且未收到 done 时自动清除该组节点状态，验证节点不会永久转圈
- [x] 6.8 改造 `web/src/stores/proxy.ts` 的 `selectProxy` —— 实现乐观更新与失败回滚，先写本地 `now` 再请求，失败时恢复原值，验证请求失败后界面回到切换前的选中节点
- [x] 6.9 替换 `Proxies.vue` 中测延迟相关处理中的 `console.error` 吞异常 —— 改为可见的错误提示，验证测延迟失败时界面出现错误语义提示且该组恢复可操作

## 7. 前端：订阅页操作级反馈

- [x] 7.1 修复 `web/src/pages/Subscriptions.vue` 的 `checkCoreError` 读取层级 —— 改为从 `res.data.data.core_error` 读取，验证订阅操作导致核心重载失败时页面顶部出现可关闭的 Core 错误提示
- [x] 7.2 拆分 `Subscriptions.vue` 的全局 `loading` 布尔 —— 拆为按操作标识的 pending 状态（`refresh:{id}`/`activate:{id}`/`delete:{id}`/`save`），验证保存订阅时删除对话框的按钮不再错误地一同转圈
- [x] 7.3 改造 `web/src/components/SubscriptionCard.vue` —— 在 meta 行的更新时间位置显示"⟳ 更新中"/"⟳ 激活中"状态，启用已声明但从未使用的 `loading` prop，验证刷新进行中时该卡片显示旋转图标且其余卡片保持可交互
- [x] 7.4 实现订阅操作的三态结果提示 —— 依据响应中的 `partial_success` 区分黄色警告"已保存但未生效"与红色错误"操作失败"，验证核心未生效时提示为警告语义且说明未生效原因
- [x] 7.5 在 `Subscriptions.vue` 实现同一订阅的操作互斥 —— 该订阅正在刷新时其其他操作不可重复提交，验证并发触发被阻止且订阅恢复就绪后操作重新可用

## 8. 前端：设置页与全局错误处理

- [x] 8.1 改造 `web/src/pages/Settings.vue` 的保存结果处理 —— 区分完全成功、半成功（配置已保存但核心未生效）与完全失败三种结果，验证核心重启失败时显示警告语义并说明后续动作
- [x] 8.2 在端口类配置变更前提示重启影响 —— 当变更涉及 mixed-port/secret/allow-lan/tun/ipv6 时提示"此操作将重启核心并断开当前连接"，验证保存此类配置前出现确认提示
- [x] 8.3 校验 `web/src/utils/request.ts` 的错误通道 —— 拦截器写入 `toastStore` 的错误经 `ToastHost` 可见，且 401 跳转登录页时不额外弹出冗余提示，验证断网/超时/业务错误三类失败均有可见提示

## 9. 集成验证

- [x] 9.1 后端测试 —— 为 `needsRestart`、`SoftRestart` 降级路径、`ApplyConfig` 的合并失败路径、逐节点并发测速的限流与结果顺序补充单元测试，验证 `go test ./...` 全部通过
- [x] 9.2 端到端手动验证清单 —— 逐项核对：测延迟（快节点先显示、每节点独立转圈、其他组不受影响）、刷新订阅（卡片级"更新中"）、切换节点（乐观更新+失败回滚）、修改端口（提示+真实重启+连接中断可见）、核心崩溃时规则页（显示"核心未运行"而非"暂无规则"），验证每一项均符合 spec 场景
- [x] 9.3 更新 `openspec/specs/` 中受影响的 5 个能力并归档本变更 —— 验证 `openspec validate --strict` 通过且归档后 `openspec list --specs` 显示这 5 个能力的 requirement 数量大于 0

## Workflow follow-up

- 剩余 10 个格式错误的 main spec（core-rules-viewer、dashboard-responsiveness、navigation-traffic-indicator、proxy-management、proxy-mode-control、subscription-content-viewer、subscription-customization、traffic-visualization、user-authentication、version-display）另开变更修复
- 暗色主题与 CSS 变量化、虚拟滚动与列表分页另开变更处理
