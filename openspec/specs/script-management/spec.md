# Script Management Specification

## Requirements

### Requirement: 脚本关联订阅
脚本 SHALL 关联到具体订阅，不再支持全局脚本。

#### Scenario: 创建脚本
- **WHEN** 调用 POST /api/subscriptions/:subscriptionId/scripts
- **THEN** 创建脚本并关联到指定订阅
- **AND** 脚本自动获取 `subscription_id`

#### Scenario: 脚本列表
- **WHEN** 调用 GET /api/subscriptions/:subscriptionId/scripts
- **THEN** 返回该订阅的所有脚本
- **AND** 不包含其他订阅的脚本

#### Scenario: 更新脚本
- **WHEN** 调用 PUT /api/subscriptions/:subscriptionId/scripts/:id
- **THEN** 更新指定脚本
- **AND** 脚本必须属于该订阅

#### Scenario: 删除脚本
- **WHEN** 调用 DELETE /api/subscriptions/:subscriptionId/scripts/:id
- **THEN** 删除指定脚本
- **AND** 脚本必须属于该订阅

#### Scenario: 测试脚本
- **WHEN** 调用 POST /api/subscriptions/:subscriptionId/scripts/:id/test
- **THEN** 执行脚本测试
- **AND** 脚本必须属于该订阅

### Requirement: 扩展脚本管理
系统 SHALL 支持用户创建、编辑、删除扩展脚本，用于改写传入 Mihomo 的配置。

#### Scenario: 脚本语言
- **WHEN** 编写脚本时
- **THEN** 仅支持 JavaScript 语言

### Requirement: 脚本执行
系统 SHALL 在配置合并的最后阶段执行启用的脚本。

#### Scenario: 执行时机
- **WHEN** 生成最终 Mihomo 配置时
- **THEN** 按顺序执行所有 enabled=true 的脚本
- **AND** 脚本可以修改配置对象

#### Scenario: 脚本输入
- **WHEN** 执行脚本时
- **THEN** 传入当前配置对象 (YAML 解析后的 JavaScript 对象)

#### Scenario: 脚本输出
- **WHEN** 脚本执行完成
- **THEN** 返回修改后的配置对象
- **AND** 继续传递给下一个脚本或输出为最终配置

#### Scenario: 脚本示例
```javascript
// 配置对象作为全局变量 config 传入
function main(config) {
  // 修改 DNS 配置
  config.dns = {
    "enable": true,
    "enhanced-mode": "fake-ip",
    "nameserver": ["8.8.8.8", "1.1.1.1"]
  };
  
  // 添加自定义规则
  config.rules.unshift("DOMAIN-SUFFIX,google.com,PROXY");
  
  return config;
}
```

### Requirement: 脚本执行范围
配置合并时 SHALL 只执行激活订阅的脚本。

#### Scenario: 合并配置
- **WHEN** Merger.Merge() 执行时
- **THEN** 只执行 `subscription_id = 激活订阅ID` 且 `enabled = true` 的脚本
- **AND** 不再执行全局脚本

### Requirement: 脚本测试
系统 SHALL 支持在不应用的情况下测试脚本执行。

#### Scenario: 测试脚本
- **WHEN** 调用 POST /api/subscriptions/:subscriptionId/scripts/:id/test
- **THEN** 使用当前配置执行脚本
- **AND** 返回执行结果和修改后的配置预览

#### Scenario: 脚本错误处理
- **WHEN** 脚本执行出错
- **THEN** 返回错误信息和堆栈
- **AND** 不影响当前运行的配置

### Requirement: 脚本安全
系统 SHALL 限制脚本的执行权限。

#### Scenario: 执行超时
- **WHEN** 脚本执行超过限制时间 (默认 5 秒)
- **THEN** 终止执行并返回超时错误

#### Scenario: 资源限制
- **WHEN** 脚本消耗过多资源
- **THEN** 限制内存和 CPU 使用