Status: done

# 实现缓存效益分析口径

## What to build

在 dashboard 统计接口中补齐缓存效益分析能力，并落实已经确认的口径：`cache_read` 带来的节省金额等于“如果没有缓存按输入单价本应支付的理论成本”减去“实际缓存读取成本”；`cache_creation` 仅作为缓存建设成本单独展示，不计入节省金额。该切片需要让前端能够直接展示累计节省、缓存读取成本和缓存建设成本等核心数据。

## Acceptance criteria

- [x] dashboard 统计结果返回缓存效益分析数据域，至少包含节省金额、缓存读取成本、缓存建设成本所需的原始业务字段。
- [x] `saved_cost_usd` 的计算严格基于 `cache_read_tokens` 的理论输入成本减去实际缓存读取成本。
- [x] `cache_creation_cost_usd` 单独累计返回，不并入节省金额。
- [x] 当记录没有 `cache_read_tokens` 或没有缓存相关成本时，缓存效益字段仍返回稳定结果且节省金额为零或正确值。
- [x] 缓存效益分析沿用当前命中的模型定价口径，不新增新的计价来源或数据库字段。

## Blocked by

- `.scratch/server-dashboard/issues/01-define-dashboard-contract-and-route.md`
