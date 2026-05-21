Status: done

# 补强服务端回归测试与联调样例

## Type

AFK

## What to build

在主要 API 与核心领域规则打通后，补上一组覆盖鉴权、幂等、事务、定价、统计和分页行为的回归测试，并整理最小联调样例。该切片不是新能力扩展，而是把已经落地的服务端能力变成可持续回归验证的资产。完成后，后续对同步接收、定价管理和查询 API 的修改都能更早暴露破坏性回归。

## Acceptance criteria

- [x] 自动化测试覆盖 `sync` 鉴权、参数校验、重复处理、事务回滚和兼容响应
- [x] 自动化测试覆盖定价匹配、placeholder 唯一性、总览排行、趋势补零和日志分页
- [x] 提供最小联调样例或说明，能够指导使用现有 client 与 server 做本地验证
- [x] 测试与说明聚焦外部行为和业务契约，不依赖内部实现细节

## Blocked by

- `03-sync-ingest-and-idempotent-write.md`
- `04-pricing-match-and-cost-calculation.md`
- `05-manage-model-pricings.md`
- `06-stats-overview-api.md`
- `07-stats-trend-api.md`
- `08-query-usage-logs.md`
