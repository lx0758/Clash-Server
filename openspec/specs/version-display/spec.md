## ADDED Requirements

### Requirement: 版本号定义
系统 SHALL 使用统一的版本号格式 `vYYYY.M.NUM`。

#### Scenario: 版本号格式
- **WHEN** 发布新版本
- **THEN** 版本号格式为 `vYYYY.M.NUM`
- **AND** YYYY 为四位年份
- **AND** M 为月份（无前导零）
- **AND** NUM 为当月发布序号

#### Scenario: 开发模式默认值
- **WHEN** 构建 时无 Git tag
- **THEN** 版本号默认为 `dev`

### Requirement: 版本号 API
系统 SHALL 提供 API 端点查询版本号。

#### Scenario: 查询版本号
- **WHEN** 调用 `GET /api/version`
- **THEN** 返回 `{ "version": "<version>" }`
- **AND** version 为当前版本号

### Requirement: UI 显示版本号
系统 SHALL 在侧边栏底部显示版本号。

#### Scenario: 显示版本号
- **WHEN** 用户查看侧边栏
- **THEN** 底部显示版本号

#### Scenario: 点击版本号
- **WHEN** 用户点击版本号
- **THEN** 复制版本号到剪贴板
