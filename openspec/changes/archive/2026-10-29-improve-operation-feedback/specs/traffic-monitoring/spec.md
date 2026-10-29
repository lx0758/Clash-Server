# Spec Delta

## ADDED Requirements

### Requirement: 延迟结果推送
系统 SHALL 通过 WebSocket 推送延迟测试过程中单个节点的延迟结果。

#### Scenario: 延迟消息推送
- **WHEN** 客户端订阅 delay 类型
- **THEN** 服务端在每个节点测试完成时推送该节点的延迟结果
- **AND** 推送格式: {"type": "delay", "data": {"group": "G1", "name": "HK-01", "delay": 120}}

#### Scenario: 测速完成信号
- **WHEN** 目标组的全部节点测试完成
- **THEN** 服务端推送该组的完成信号
- **AND** 推送格式: {"type": "delay", "data": {"group": "G1", "done": true}}

#### Scenario: 节点测试失败
- **WHEN** 某个节点的延迟测试失败
- **THEN** 该节点的结果中体现失败状态
- **AND** 不影响同组其他节点的测试继续进行

#### Scenario: 逐节点结果与完成信号的关系
- **WHEN** 客户端接收 delay 类型消息
- **THEN** 在收到该组的完成信号前，该组节点处于测试中状态
- **AND** 收到完成信号后，该组全部节点退出测试中状态
