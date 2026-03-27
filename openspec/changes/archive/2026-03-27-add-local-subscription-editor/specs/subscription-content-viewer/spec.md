## ADDED Requirements

### Requirement: 订阅内容查看
系统 SHALL 允许用户查看订阅的配置内容。

#### Scenario: 显示查看入口
- **WHEN** 用户查看订阅的操作菜单
- **THEN** 显示"查看源文件"选项

#### Scenario: 打开查看器
- **WHEN** 用户点击"查看源文件"
- **THEN** 显示内容对话框
- **AND** 显示当前订阅的 YAML 配置内容
- **AND** YAML 内容不自动折行

#### Scenario: 远程订阅只读
- **WHEN** 用户查看远程订阅内容
- **THEN** 内容为只读模式
- **AND** 不显示保存按钮

#### Scenario: 复制内容
- **WHEN** 用户点击复制按钮
- **THEN** 内容复制到剪贴板

#### Scenario: 下载内容
- **WHEN** 用户点击下载按钮
- **THEN** 下载 YAML 文件

### Requirement: 本地订阅内容编辑
系统 SHALL 允许用户编辑本地订阅的配置内容。

#### Scenario: 本地订阅可编辑
- **WHEN** 用户查看本地订阅内容
- **THEN** 内容为可编辑模式
- **AND** 显示保存按钮

#### Scenario: 保存内容
- **WHEN** 用户修改本地订阅内容并点击保存
- **THEN** 调用 API 更新订阅内容
- **AND** 更新节点数量
- **AND** 如果订阅处于激活状态，重新加载配置

#### Scenario: 内容为空
- **WHEN** 本地订阅内容为空
- **THEN** 编辑器显示空白内容
- **AND** 允许用户输入新内容

### Requirement: 本地订阅列表展示
本地订阅 SHALL 在收缩状态下显示节点数量和编辑时间。

#### Scenario: 显示节点数量
- **WHEN** 用户查看本地订阅卡片
- **THEN** 显示节点数量

#### Scenario: 显示编辑时间
- **WHEN** 本地订阅曾被编辑
- **THEN** 显示"XXX前编辑"
