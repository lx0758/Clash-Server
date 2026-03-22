# 代理延迟颜色编码

## 概述

根据代理节点的延迟值显示不同颜色，直观反映节点响应速度。

## 颜色规则

| 延迟范围 | 颜色 | CSS Class |
|---------|------|-----------|
| ≤400ms | 绿色 #22c55e | `fast` |
| 400-800ms | 黄绿色 #84cc16 | `medium` |
| 800-1500ms | 黄色 #eab308 | `slow` |
| 1500-3000ms | 橙色 #f97316 | `very-slow` |
| >3000ms | 红色 #ef4444 | `error` |
| 超时(0) | 灰色 #94a3b8 | `timeout` |

## 文件

- `web/src/components/proxy/ProxyItem.vue`
