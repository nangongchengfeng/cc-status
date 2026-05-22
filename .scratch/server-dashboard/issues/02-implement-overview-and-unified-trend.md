Status: done

# 实现总览卡片与统一时间桶趋势

## What to build

基于新的 dashboard 统计接口，实现大屏所需的总览卡片和统一时间桶趋势数据。该切片需要从 `usage_reports` 中按显式时间范围聚合出总 token、总费用、总请求数、活跃客户端数，并生成按 `Asia/Shanghai` 对齐和补零的统一时间桶数组。每个时间桶都应携带输入、输出、缓存读、缓存创建 token 以及请求数和总费用，供多个图表共享同一条时间轴。

## Acceptance criteria

- [x] dashboard 统计结果返回 `total_tokens`、`total_cost_usd`、`total_requests`、`active_clients` 四个核心总览字段。
- [x] dashboard 统计结果返回统一时间桶数组，支持按 `hour` 或 `day` 聚合，并按 `Asia/Shanghai` 补零。
- [x] 每个时间桶至少包含 `bucket`、输入/输出/缓存读/缓存创建 token、请求数、总费用等原始业务字段。
- [x] 显式 `start_at`、`end_at` 会影响总览卡片和趋势数据，超出范围的记录不会被聚合。
- [x] 当指定时间范围内没有数据时，接口返回结构仍然稳定，不因空结果破坏前端消费。

## Blocked by

- `.scratch/server-dashboard/issues/01-define-dashboard-contract-and-route.md`
