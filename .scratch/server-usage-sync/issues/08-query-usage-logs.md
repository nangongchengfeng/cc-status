Status: done

# 提供原始日志查询 API

## Type

AFK

## What to build

交付 `GET /api/v1/logs` 的完整纵向切片，支持按客户端、模型、请求标识和时间范围过滤**使用记录**，并使用 `offset` 分页返回最近记录。该切片需要把默认分页规则、默认排序、完整 token 字段、完整费用字段和 `pricing_source` 一并暴露出来，以便后台排障和计费核对。完成后，调用方可以直接查看原始**使用记录**列表，而无需手查数据库。

## Acceptance criteria

- [x] `GET /api/v1/logs` 支持 `client_id`、`model`、`request_id`、`start_time`、`end_time` 等过滤条件
- [x] 接口采用 `offset` 分页，默认 `limit=20`、最大 `100`
- [x] 默认按 `created_at DESC, id DESC` 排序，并返回完整 token、完整费用字段与 `pricing_source`
- [x] 成功响应使用 `{ "data": ... }` 风格，同时返回分页所需的 `total`、`page/offset` 类信息

## Blocked by

- `03-sync-ingest-and-idempotent-write.md`
- `04-pricing-match-and-cost-calculation.md`
