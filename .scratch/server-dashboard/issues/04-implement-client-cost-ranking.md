Status: done

# 实现客户端费用排行聚合

## What to build

在 dashboard 统计接口中补齐客户端费用排行能力。该切片需要按指定时间范围聚合各 `client_id` 的总费用，固定返回 TOP 10，并维持现有领域对 `client_id` 的语义约束：它是稳定技术标识，而不是展示别名。结果应直接供前端排行模块使用，不额外引入客户端资料表或别名管理。

## Acceptance criteria

- [x] dashboard 统计结果返回客户端费用排行数组，固定为 TOP 10。
- [x] 客户端排行按 `total_cost_usd DESC, client_id ASC` 排序，保证输出稳定。
- [x] 每个排行项直接返回 `client_id` 和聚合后的总费用，不引入客户端别名体系。
- [x] 聚合结果只统计指定时间范围内的**使用记录**。
- [x] 新增能力不影响现有 `GET /api/v1/stats/overview` 中前 5 客户端排行的既有口径。

## Blocked by

- `.scratch/server-dashboard/issues/01-define-dashboard-contract-and-route.md`
