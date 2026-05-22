# Web Dashboard Design

**目标**

在 `web/` 中落地一个面向桌面演示的数据大屏单页应用，首页直接使用 `/`，对接服务端的 dashboard 聚合接口与日志接口，展示总览、趋势、排行、缓存效益和最近请求。

**约束**

- 采用 `React + TypeScript + Vite + Tailwind CSS + Recharts + TanStack Query`。
- 保持单页，不扩展为后台管理台。
- 遵循 `web/CLAUDE.md` 的目录结构与前端工程约定。
- 遵循 `web/AGENTS.md` 的视觉禁区，避免紫色 SaaS 模板感，使用深色、非对称、卡片化的大屏布局。
- 前端负责时间预设、金额与时间格式化，不让后端返回展示文案。

**总体方案**

页面只保留一个根路由 `/`，渲染 `DashboardPage`。页面持有时间范围预设状态，并把预设转换为 `start_at`、`end_at`、`interval`，同时驱动 dashboard 和最近请求两组查询。展示层以页面目录为主，先在 `pages/Dashboard/` 内完成首版模块，再把明显复用三次以上的通用块提取到 `components/`。

**数据流**

1. 页面初始化时默认选择“最近 7 天”。
2. 工具函数把预设转换为后端查询参数。
3. `useDashboardQuery()` 请求 `/stats/dashboard`。
4. `useRecentLogsQuery()` 请求 `/logs`。
5. 页面根据 query 状态渲染 `loading / error / empty / success`。
6. 页面内模块只消费已格式化或已映射的数据，不直接操作原始请求细节。

**布局**

- 顶部为页头，包含标题、说明和时间范围切换器。
- 中部为主视觉区，左侧偏宽放总览卡片和趋势图，右侧放辅助信息或摘要卡片。
- 下部为排行、缓存效益和最近请求表格。
- 使用不对称网格，避免平均分栏。

**模块拆分**

- `pages/Dashboard/DashboardPage.tsx`：页面编排与状态分发。
- `pages/Dashboard/components/`：页面专属卡片、图表、排行、表格模块。
- `api/dashboard.ts` / `api/logs.ts`：资源级 API。
- `hooks/useDashboardQuery.ts` / `hooks/useRecentLogsQuery.ts`：查询 hooks。
- `constants/timeRanges.ts`：前端时间预设。
- `utils/dashboard.ts` / `utils/format.ts`：参数转换、格式化、展示名回退。
- `types/dashboard.ts` / `types/logs.ts`：类型定义。

**测试策略**

- issue 01：用最小页面测试和构建验证锁住根路由与骨架。
- issue 02：重点测试时间预设转换工具和查询 hooks。
- issue 03-05：优先测试关键状态映射，不测试图表库内部。
- issue 06：统一补齐格式化函数和页面状态测试，保留 `npm run build` 作为最小回归证据。

**执行方式**

按 `.scratch/web-dashboard/issues/01-06` 顺序逐个执行。每完成一个 issue，更新 issue 状态与复选框，并使用中文约定式提交。
