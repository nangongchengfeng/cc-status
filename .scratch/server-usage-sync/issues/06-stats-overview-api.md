Status: done

# 提供总览统计 API

## Type

AFK

## What to build

交付 `GET /api/v1/stats/overview` 的完整纵向切片，从 `usage_reports` 聚合生成总 token、总成本、总请求数、活跃客户端数，以及固定规则的模型和客户端排行。该切片需要把 `top_models` 与 `top_clients` 的排序与数量规则真正落地，并保持查询响应风格符合服务端统一的 `{ "data": ... }` 约定。完成后，调用方可以直接通过单个总览接口拿到 dashboard 和监控最常用的概览数据。

## Acceptance criteria

- [x] `GET /api/v1/stats/overview` 返回总 token、总成本、总请求数和活跃客户端数
- [x] `top_models` 固定返回前 5 条，按 `tokens DESC, model ASC` 排序
- [x] `top_clients` 固定返回前 5 条，按 `total_cost_usd DESC, client_id ASC` 排序
- [x] 成功响应使用 `{ "data": ... }` 风格，鉴权和错误状态码行为与服务端统一约定一致

## Blocked by

- `03-sync-ingest-and-idempotent-write.md`
- `04-pricing-match-and-cost-calculation.md`
