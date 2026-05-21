Status: done

# 打通同步接收与幂等写入

## Type

AFK

## What to build

交付 `POST /api/v1/sync` 的首个完整纵向切片，支持接收 client 发来的**上报批次**，完成请求校验、来源白名单校验、鉴权、批次单事务和幂等写入。该切片先不负责复杂定价命中，只需要把请求合法性、重复处理、兼容 client 的成功响应以及非重复数据库错误整批回滚这条核心链路跑通。完成后，现有 client 可以与服务端完成一次真实联调，并验证 `accepted_count` / `duplicate_count` 语义。

## Acceptance criteria

- [x] `POST /api/v1/sync` 在鉴权通过后可接收包含 `client_id` 和 `reports` 的请求体，并对空批次、非法来源和缺失字段返回 `400`
- [x] 单个请求批次使用单事务处理，`(client_id, request_id)` 重复会计入 `duplicate_count`，其他数据库错误会回滚整批并返回 `500`
- [x] 同一请求体内重复的 `request_id` 按 duplicate 处理，不会导致整批失败
- [x] 成功响应兼容现有 client 协议，返回 `code`、`message`、`accepted_count`、`duplicate_count`

## Blocked by

- `01-bootstrap-server-runtime-and-health.md`
- `02-persist-usage-reports-and-model-pricing.md`
