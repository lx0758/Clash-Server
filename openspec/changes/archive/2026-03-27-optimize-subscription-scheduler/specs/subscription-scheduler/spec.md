## Note

本变更为 Bug 修复，不涉及需求变更。

现有 `subscription-management` 规格中的"订阅定时刷新"需求保持不变：

> ### Requirement: 订阅定时刷新
> 系统 SHALL 支持定时自动刷新订阅。
>
> #### Scenario: 定时刷新
> - **WHEN** 订阅配置了 interval（分钟）
> - **THEN** 后台定时任务按间隔自动拉取新数据
> - **AND** 更新数据库中的内容

本变更修复该需求的实现缺陷，确保定时刷新功能正常工作。
