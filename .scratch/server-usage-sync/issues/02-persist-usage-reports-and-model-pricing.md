Status: done

# 建立使用记录与模型定价持久化骨架

## Type

AFK

## What to build

交付 `usage_reports` 和 `model_pricing` 的持久化骨架，覆盖实体定义、最小字段集、唯一键、自动建表和基础 repository。该切片需要把**使用记录**、**模型定价**、**默认定价**、`pricing_source`、`inserted_at` 等核心领域概念真正落到数据库中，并在启动流程中幂等初始化常见 Claude 模型定价及唯一的全局 placeholder 默认定价。完成后，服务端具备后续同步接收与查询 API 所需的底层存储结构。

## Acceptance criteria

- [x] `usage_reports` 与 `model_pricing` 已通过自动建表落地，包含 PRD 约定的最小字段集与约束
- [x] `(client_id, request_id)` 唯一约束、`pricing_source` 字段和 `inserted_at` 到达时间字段已落库
- [x] 启动时会幂等初始化常见 Claude 模型定价和一条全局 placeholder 默认定价
- [x] 全局 placeholder 默认定价在持久化层具备唯一性保护，重复初始化或重复写入不会产生多条生效默认价

## Blocked by

- `01-bootstrap-server-runtime-and-health.md`
