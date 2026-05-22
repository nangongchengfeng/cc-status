Status: done

# 实现模型排行与展示名聚合

## What to build

在 dashboard 统计接口中补齐模型排行能力，并把模型展示名一并聚合出来。该切片需要在指定时间范围内按总 token 数计算模型使用强度，固定返回 TOP 10，且结果同时包含原始 `model` 和可选 `display_name`。`display_name` 需要来源于现有模型定价数据，缺失时仍保留原始模型标识，避免页面展示和排障语义割裂。

## Acceptance criteria

- [x] dashboard 统计结果返回模型排行数组，固定为 TOP 10。
- [x] 模型排行按 `input_tokens + output_tokens + cache_read_tokens + cache_creation_tokens` 的总 token 数降序排序，并在并列时保持稳定排序。
- [x] 每个模型排行项同时包含原始 `model` 和可选 `display_name`。
- [x] `display_name` 来源于现有模型定价数据，未配置展示名时不会丢失原始模型标识。
- [x] 聚合只统计指定时间范围内的**使用记录**，不引入新的模型元数据来源或表结构变更。

## Blocked by

- `.scratch/server-dashboard/issues/01-define-dashboard-contract-and-route.md`
