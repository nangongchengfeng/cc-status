Status: done

# 落地总览卡片与趋势图区域

## What to build

在大屏主页面中实现总览卡片和趋势图的主要展示区域。该切片需要把 dashboard 接口返回的总览字段和统一时间桶趋势渲染成可演示的卡片与图表，并让时间范围切换后图表和卡片联动刷新。重点是让页面中部的核心视觉区域先真正成型，能直观展示成本和 token 变化。

## Acceptance criteria

- [x] 页面顶部展示四个总览卡片：总 Token、总费用、总请求数、活跃客户端数。
- [x] 页面展示至少一个费用趋势图和一个 Token 细分趋势图，并共享同一组时间桶数据。
- [x] 图表支持基础 Tooltip 或等效交互，能查看时间点对应的详细数值。
- [x] 时间范围切换后，总览卡片和趋势图会一起刷新并保持口径一致。
- [x] 当接口无数据时，该区域展示稳定空态而不是报错或错位布局。

## Blocked by

- `.scratch/web-dashboard/issues/02-connect-dashboard-query-and-time-range-presets.md`
