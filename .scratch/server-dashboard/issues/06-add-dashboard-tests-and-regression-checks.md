Status: done

# 补齐仪表盘统计测试与回归验证

## What to build

为新的 dashboard 统计能力补齐自动化测试和最小回归验证，覆盖 handler、service 或 repository 层的关键统计契约。目标不是测试内部实现细节，而是锁定外部可观察行为，包括参数校验、统一时间桶补零、排序稳定性、模型展示名回退以及缓存节省公式，确保后续前端联调不会因为统计口径漂移而反复返工。

## Acceptance criteria

- [x] 新增 dashboard 接口的 handler 测试，覆盖鉴权、非法参数和成功响应包装。
- [x] 新增聚合逻辑测试，覆盖总览字段、统一时间桶补零和显式时间范围过滤。
- [x] 新增排行测试，覆盖模型排行与客户端排行的主要排序规则和并列时稳定性。
- [x] 新增缓存效益测试，覆盖节省金额公式和缓存建设成本单独累计的行为。
- [x] 相关服务端测试命令可在本地运行并作为该功能的最小回归证据。

## Blocked by

- `.scratch/server-dashboard/issues/02-implement-overview-and-unified-trend.md`
- `.scratch/server-dashboard/issues/03-implement-model-ranking-and-display-name.md`
- `.scratch/server-dashboard/issues/04-implement-client-cost-ranking.md`
- `.scratch/server-dashboard/issues/05-implement-cache-benefit-analysis.md`
