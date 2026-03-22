## Why

当前系统在所有配置变更（订阅、规则、脚本等）时都会重启 Mihomo 进程，导致连接中断和用户体验差。Mihomo 内核支持热重载（PUT /configs）和软重启（POST /restart）API，但项目未利用此能力。

## What Changes

- **启动流程重构**: Mihomo 启动时通过 `-config` 参数传入 Base64 编码的最小配置（`external-controller: {APIHost}:{APIPort}`），启动后通过 `HotReload()` 应用全量配置
- **统一 API 管理**: 所有配置变更通过 Mihomo API 管理，不再使用临时配置文件
- **热重载优先**: 大多数配置变更通过 `HotReload()` 热重载
- **软重启替代硬重启**: 需要重启时调用 `SoftRestart()`，而非杀进程
- **热重载失败处理**: 失败时询问用户是否软重启

## Capabilities

### New Capabilities

- `hot-reload`: 配置热重载能力，支持不重启进程的情况下更新配置

### Modified Capabilities

- `clash-core-integration`: 修改核心启动和配置管理流程

## Impact

**代码改动**:
- `server/internal/service/core.go` - 重构启动流程，新增 API 管理方法
- `server/pkg/response/response.go` - 新增响应方法
- `server/internal/handler/*.go` - 适配新响应结构

**移除**:
- 不再需要临时配置文件（temp/config-*.yaml）
- 不再需要 `mihomo -f <config>` 指定配置文件

**用户体验改进**:
- 订阅/规则/脚本变更时，大多数情况下不中断连接
- 变更生效时间从 3-5 秒降至毫秒级
- 软重启比硬重启更稳定