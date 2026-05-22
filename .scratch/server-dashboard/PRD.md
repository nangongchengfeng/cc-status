Status: ready-for-agent

# Claude Usage Dashboard Server API

## Problem Statement

当前 `cc-status/server` 已经具备同步接收、总览统计、趋势查询、日志分页和模型定价管理能力，但这些 API 主要面向通用查询与排障，不足以直接支撑 Web 数据大屏。现有 `overview` 和 `trend` 无法一次性提供输入/输出/缓存读/缓存创建四类 token 趋势、缓存节省金额、Top 10 排行以及模型展示名等大屏所需聚合结果。如果前端自己用分页日志拼装这些统计，会带来数据失真、性能低下和职责错位的问题。

## Solution

在现有服务端上新增一个面向 Web 数据大屏的独立聚合统计接口，保持既有 `sync`、`overview`、`trend`、`logs` 和模型定价管理接口不破坏。新的 dashboard 统计接口负责返回大屏所需的总览卡片、统一时间桶趋势、模型排行、客户端排行和缓存效益分析；最近请求列表继续复用现有 `GET /api/v1/logs`。接口仅返回原始业务字段，不混入展示文案，前端根据返回结果做格式化和图表渲染。

## User Stories

1. 作为项目展示者，我希望服务端直接提供面向大屏的聚合接口，以便前端无需用分页日志拼装复杂统计。
2. 作为前端开发者，我希望通过一次 dashboard 请求拿到大部分核心图表数据，以便降低页面数据编排复杂度。
3. 作为前端开发者，我希望保留现有 `GET /api/v1/logs` 用于最近请求列表，以便继续复用既有分页和过滤能力。
4. 作为平台管理员，我希望现有 `GET /api/v1/stats/overview` 和 `GET /api/v1/stats/trend` 继续可用，以便不破坏已有调试和联调方式。
5. 作为调用方，我希望新的 dashboard 接口仍然受静态 Bearer token 鉴权保护，以便与现有查询接口的安全边界保持一致。
6. 作为调用方，我希望 dashboard 接口使用显式 `start_at`、`end_at` 和 `interval` 参数，以便前端可以自由映射“今天”“最近 7 天”“本月”等预设。
7. 作为前端开发者，我希望 dashboard 接口支持 `interval=hour|day`，以便在短时间范围和长时间范围之间切换图表粒度。
8. 作为前端开发者，我希望 dashboard 趋势数据使用统一时间桶数组，以便多个图表共享同一条时间轴。
9. 作为用户，我希望在趋势图里区分输入、输出、缓存读和缓存创建 token，以便理解成本构成。
10. 作为用户，我希望看到指定时间范围内的总 token 消耗，以便快速掌握整体用量。
11. 作为用户，我希望看到指定时间范围内的总费用，以便快速掌握花费规模。
12. 作为用户，我希望看到指定时间范围内的总请求数，以便了解总体调用次数。
13. 作为用户，我希望看到指定时间范围内的活跃客户端数，以便知道有多少机器在产生数据。
14. 作为用户，我希望看到模型排行固定返回 TOP 10，以便快速识别最常用模型。
15. 作为开发者，我希望模型排行按总 token 数排序，以便“模型使用排行”的语义清晰稳定。
16. 作为用户，我希望看到客户端排行固定返回 TOP 10，以便快速识别成本最高的客户端标识。
17. 作为开发者，我希望客户端排行按总费用排序，以便和“客户端费用排行”的页面文案一致。
18. 作为用户，我希望 dashboard 接口返回模型的 `display_name`，以便大屏展示更友好。
19. 作为开发者，我希望 dashboard 接口同时保留原始 `model`，以便排障和核对时不丢失技术标识。
20. 作为用户，我希望看到缓存节省金额，以便直观看到缓存带来的收益。
21. 作为开发者，我希望缓存效益分析有明确公式，以便前后端和后续维护者对“节省”定义一致。
22. 作为用户，我希望缓存建设成本被单独展示，以便知道为了后续命中缓存实际付出了多少。
23. 作为前端开发者，我希望接口只返回原始业务字段，不返回页面专用格式化字符串，以便前端保留展示灵活性。
24. 作为开发者，我希望 dashboard 统计口径完全基于 `usage_reports` 和 `model_pricing`，以便复用现有领域模型和 ADR。
25. 作为项目维护者，我希望 dashboard 聚合逻辑形成独立深模块，以便后续测试和复用更加稳定。
26. 作为测试工程师，我希望 dashboard 接口有针对排序、补零、时间过滤和缓存节省公式的自动化测试，以便避免大屏上线后统计口径漂移。

## Implementation Decisions

- 服务端新增一个独立的 dashboard 统计接口，语义上属于现有 stats 资源下的聚合能力，不替代或重写既有 `overview`、`trend`、`logs`。
- dashboard 接口返回的数据域包括：总览卡片、统一时间桶趋势、模型排行、客户端排行、缓存效益分析；不包含最近请求列表。
- 最近请求列表继续由现有 `GET /api/v1/logs` 提供，前端按固定 `limit` 获取最近若干条**使用记录**。
- dashboard 接口查询参数采用显式时间范围设计：前端维护时间预设，服务端只接收 `start_at`、`end_at` 和 `interval`。
- `interval` 首版只支持 `hour` 和 `day`，与现有趋势接口保持一致，并继续按 `Asia/Shanghai` 进行桶对齐和补零。
- dashboard 趋势数据采用统一时间桶数组，每个桶同时携带输入、输出、缓存读、缓存创建 token，以及请求数、总费用、缓存节省金额等可用于图表的数据字段。
- dashboard 总览卡片固定返回四个核心值：`total_tokens`、`total_cost_usd`、`total_requests`、`active_clients`。
- dashboard 模型排行固定返回 TOP 10，按 `input_tokens + output_tokens + cache_read_tokens + cache_creation_tokens` 的总和降序排序；若并列则按 `model` 升序保证稳定性。
- dashboard 客户端排行固定返回 TOP 10，按 `total_cost_usd DESC, client_id ASC` 排序，并直接展示 `client_id`，不引入新的客户端资料表或别名体系。
- dashboard 模型相关结果同时返回原始 `model` 和可选 `display_name`；`display_name` 来源于 `model_pricing.display_name`，未命中时为空。
- dashboard 缓存效益分析遵循统一口径：`saved_cost_usd = theoretical_input_cost_for_cache_read - actual_cache_read_cost`。其中 theoretical 成本按命中定价的输入单价计算，`cache_creation_cost_usd` 仅作为建设成本单独返回，不计入节省金额。
- dashboard 统计结果仅返回原始业务字段和原始金额字符串，不返回 `$12.34`、短标签、中文时间范围名等展示层格式化字段。
- dashboard 聚合严格基于现有**使用记录**和**模型定价**领域模型，不修改 `usage_reports` 和 `model_pricing` 的表结构。
- 实现上建议把 dashboard 聚合沉淀为独立 service/repository 模块，避免在 handler 中堆积聚合逻辑。
- dashboard 统计模块应复用现有金额解析、时间桶截断和模型定价读取逻辑，减少口径分叉。
- 新接口成功响应继续遵循查询类 API 的 `{ "data": ... }` 包装风格；错误响应继续遵循现有错误码和 HTTP 状态码约定。
- 首版不为了 dashboard 引入预计算表、物化视图、Redis 缓存或异步聚合任务，仍以 SQLite 实时聚合为主。

## Testing Decisions

- 好的测试只验证 dashboard 作为外部契约的统计口径和输出结构，不绑定内部 SQL 细节或具体函数调用顺序。
- 重点测试模块是 dashboard 聚合 service、dashboard handler，以及必要的 repository 聚合查询。
- dashboard handler 测试应覆盖：鉴权失败、非法 `interval`、非法时间范围、成功响应包装结构。
- dashboard service 或 repository 测试应覆盖：显式时间范围过滤、按 `Asia/Shanghai` 聚合、空桶补零、统一时间轴长度正确。
- 排行测试应覆盖：模型排行按总 token 排序、客户端排行按总费用排序、并列时的稳定二级排序。
- 缓存效益测试应覆盖：`cache_read` 节省金额公式、`cache_creation` 成本单独累计、无缓存读时节省金额为零。
- 模型展示名测试应覆盖：命中 `display_name` 时返回展示名，未配置展示名时保留原始 `model`。
- repository 层优先使用 SQLite 集成测试，保持与现有服务端测试风格一致，真实验证金额聚合和排序结果。
- handler 层继续使用 `httptest`，与现有 `stats`、`logs`、`model_pricing` handler 测试模式保持一致。

## Out of Scope

- 不把最近请求列表并入 dashboard 统计接口。
- 不修改现有 `sync` 协议、总览接口、趋势接口和日志分页协议。
- 不新增客户端别名管理、模型展示名编辑 UI 或新的元数据维护资源。
- 不引入服务端返回的页面格式化字符串、国际化文案或图表专用 label 字段。
- 不引入预计算表、缓存层、异步聚合或单独的数据仓库。
- 不为了 dashboard 需求修改现有数据库 schema。

## Further Notes

- 该 PRD 依赖现有 `usage_reports` 作为**使用记录库**，并遵循 ADR 中“服务端核心域使用 `usage_reports`”的术语约束。
- 该 PRD 是对现有服务端能力的增量增强，不替代已有 `server-usage-sync` PRD。
- 与 Web 大屏配合时，前端应同时调用 dashboard 接口和现有 `GET /api/v1/logs`，而不是回退到前端全量聚合。
