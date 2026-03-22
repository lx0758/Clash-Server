## ADDED Requirements

### Requirement: 用户登录
系统 SHALL 支持单用户登录认证。

#### Scenario: 首次登录（设置密码）
- **WHEN** 系统无用户，调用 GET /api/init
- **THEN** 返回 {"initialized": false} 引导设置密码
- **AND** 调用 POST /api/init 设置初始密码

#### Scenario: 用户登录
- **WHEN** 调用 POST /api/session，传入 username 和 password
- **AND** 凭证正确
- **THEN** 创建 Session，返回成功
- **AND** 设置 Cookie (session_id)

#### Scenario: 登录失败
- **WHEN** 调用 POST /api/session，凭证错误
- **THEN** 返回 code 401，提示用户名或密码错误
- **AND** 不创建 Session

### Requirement: 用户登出
系统 SHALL 支持用户登出。

#### Scenario: 登出
- **WHEN** 调用 DELETE /api/session
- **THEN** 销毁 Session，返回成功

### Requirement: 获取当前用户
系统 SHALL 提供当前登录用户信息。

#### Scenario: 获取当前用户
- **WHEN** 调用 GET /api/users/me，带有效 Session
- **AND** 返回用户信息（不含密码）

### Requirement: 修改密码
系统 SHALL 支持修改密码。

#### Scenario: 修改密码
- **WHEN** 调用 PUT /api/users/me/password，传入 old_password 和 password
- **AND** 旧密码正确
- **THEN** 更新密码为新密码
- **AND** 保持 Session 有效

#### Scenario: 旧密码错误
- **WHEN** 调用 PUT /api/users/me/password，旧密码错误
- **THEN** 返回 code 422，提示旧密码错误

### Requirement: 初始密码设置
系统 SHALL 能在首次运行时设置管理员密码。

#### Scenario: 首次设置密码
- **WHEN** 系统无用户，调用设置密码接口
- **THEN** 创建第一个用户，密码 bcrypt 存储
- **AND** 禁止再创建新用户（单用户模式）