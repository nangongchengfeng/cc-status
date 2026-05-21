Status: done

# 提供模型定价管理 API

## Type

AFK

## What to build

交付 `GET /api/v1/model-pricings`、`POST /api/v1/model-pricings`、`PUT /api/v1/model-pricings/:id` 三个管理 API，并把 placeholder 唯一性和全量更新语义真正对外暴露。该切片需要保证管理员可以查看全部定价记录、创建新模型价格、修正错误价格，同时在尝试创建或更新第二条全局默认定价时收到 `409`。完成后，服务端不再依赖手工改库即可维护定价规则。

## Acceptance criteria

- [x] `GET /api/v1/model-pricings` 返回全部定价记录，并包含 `is_placeholder`
- [x] `POST /api/v1/model-pricings` 支持创建普通模型定价，并对第二条 placeholder 默认定价返回 `409`
- [x] `PUT /api/v1/model-pricings/:id` 采用全量更新语义，不支持部分更新
- [x] 更新普通定价为 placeholder 时若会造成冲突，同样返回 `409`

## Blocked by

- `02-persist-usage-reports-and-model-pricing.md`
- `04-pricing-match-and-cost-calculation.md`
