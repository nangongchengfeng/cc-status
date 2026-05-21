Status: done

# 打通定价匹配与费用计算

## Type

AFK

## What to build

把模型名小写规范化、精确匹配、最长前缀匹配、默认定价回退和费用计算真正接入 `POST /api/v1/sync` 写入链路。该切片需要为每条**使用记录**落库 `input_cost_usd`、`output_cost_usd`、缓存费用字段、`total_cost_usd` 与 `pricing_source`，并保证默认定价命中和 placeholder 唯一性规则被统一执行。完成后，服务端接收到的同步数据已经是可直接统计和对账的完整计费结果。

## Acceptance criteria

- [x] 上报 `model` 与 `model_id` 在匹配前统一小写规范化
- [x] 模型定价命中顺序固定为精确匹配、最长前缀匹配、全局默认定价，并正确写入 `pricing_source`
- [x] 每条新接收的**使用记录**都会计算并落库完整 USD 费用字段与 `total_cost_usd`
- [x] 当默认定价被命中时，系统行为可通过测试和查询明确验证，不会静默丢失匹配来源

## Blocked by

- `02-persist-usage-reports-and-model-pricing.md`
- `03-sync-ingest-and-idempotent-write.md`
