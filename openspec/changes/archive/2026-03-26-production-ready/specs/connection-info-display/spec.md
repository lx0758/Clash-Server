# 连接信息显示优化

## 概述

优化连接列表的信息展示，使用两行框架布局。

## 显示格式

**第一行**：网络类型 / 连接类型 / 目标IP:端口
- 示例：`TCP / HTTP / 142.250.185.78:443`

**第二行**：规则 / 规则参数 => 链路
- 示例：`GEOIP / CN => 节点A → 节点B`

**右侧**：上行/下行流量 + 关闭按钮

## 文件

- `web/src/pages/Connections.vue`
- `web/src/types/api.ts` - 新增 `destinationIP`, `rulePayload` 字段
