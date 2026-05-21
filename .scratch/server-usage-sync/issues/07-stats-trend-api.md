Status: done

# 提供趋势统计 API

## Type

AFK

## What to build

交付 `GET /api/v1/stats/trend` 的完整纵向切片，支持按 `hour|day` 两种粒度、基于 `Asia/Shanghai` 时区做趋势聚合，并补零返回完整时间轴。该切片需要把业务时间 `created_at` 的使用、时间范围过滤、空桶补零和统一响应风格全部打通。完成后，调用方可以直接用该接口驱动折线图或趋势监控，而无需自行拼装缺失时间桶。

## Acceptance criteria

- [x] `GET /api/v1/stats/trend` 仅接受 `interval=hour|day`
- [x] 趋势聚合按 `Asia/Shanghai` 时区基于业务时间 `created_at` 计算
- [x] 响应会补零并返回完整时间轴，而不仅是有数据的时间桶
- [x] 成功响应使用 `{ "data": ... }` 风格，未授权或参数错误时返回正确 HTTP 状态码

## Blocked by

- `03-sync-ingest-and-idempotent-write.md`
- `04-pricing-match-and-cost-calculation.md`
