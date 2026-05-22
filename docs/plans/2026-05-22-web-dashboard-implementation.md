# Web Dashboard Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在 `web/` 中实现可运行、可构建、可演示的单页数据大屏，并按 issue 顺序逐步接入查询、图表、排行、缓存效益与最近请求。

**Architecture:** 前端采用单页结构，`/` 直接渲染 `DashboardPage`。页面持有时间范围状态，通过工具函数转换成服务端查询参数，再用 TanStack Query 驱动 dashboard 与 logs 两组数据。页面专属展示组件优先保留在 `pages/Dashboard/` 目录内，等复用清晰后再抽离。

**Tech Stack:** React 18、TypeScript 5、Vite、Tailwind CSS、React Router v6、TanStack Query、Axios、Recharts、Vitest、React Testing Library

---

### Task 1: Bootstrap App And Home Route

**Files:**
- Create: `web/package.json`
- Create: `web/index.html`
- Create: `web/tsconfig.json`
- Create: `web/tsconfig.app.json`
- Create: `web/vite.config.ts`
- Create: `web/tailwind.config.ts`
- Create: `web/postcss.config.js`
- Create: `web/src/main.tsx`
- Create: `web/src/router.tsx`
- Create: `web/src/App.tsx`
- Create: `web/src/styles/index.css`
- Create: `web/src/pages/Dashboard/DashboardPage.tsx`
- Create: `web/src/pages/Dashboard/DashboardPage.test.tsx`

**Step 1: Write the failing test**

- 编写页面测试，验证根路由 `/` 会渲染仪表盘页头和卡片骨架标题。

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- DashboardPage`
- 预期：测试失败，因为应用和页面文件还不存在。

**Step 3: Write minimal implementation**

- 初始化 Vite React TS 工程文件。
- 接入 React Router 和 QueryClientProvider。
- 创建 `DashboardPage` 骨架，包含页头、主体区和卡片式容器。
- 添加中文注释说明主题布局意图。

**Step 4: Run test to verify it passes**

- 运行：`npm run test -- DashboardPage`
- 预期：测试通过。

**Step 5: Run minimal build verification**

- 运行：`npm run build`
- 预期：构建通过。

**Step 6: Commit**

- 提交：`feat(web): 初始化dashboard前端骨架`

### Task 2: Connect Queries And Time Presets

**Files:**
- Create: `web/src/api/http.ts`
- Create: `web/src/api/dashboard.ts`
- Create: `web/src/api/logs.ts`
- Create: `web/src/hooks/useDashboardQuery.ts`
- Create: `web/src/hooks/useRecentLogsQuery.ts`
- Create: `web/src/constants/timeRanges.ts`
- Create: `web/src/types/dashboard.ts`
- Create: `web/src/types/logs.ts`
- Create: `web/src/utils/timeRange.ts`
- Create: `web/src/utils/timeRange.test.ts`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: Write the failing test**

- 为时间预设转换工具编写测试，覆盖“今天、最近 7 天、最近 30 天、本月、全部”的参数转换。

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- timeRange`
- 预期：测试失败，因为工具函数未实现。

**Step 3: Write minimal implementation**

- 定义 dashboard 和 logs 的请求类型。
- 建立 axios 客户端和静态 Bearer token 注入。
- 实现时间预设常量与参数转换工具。
- 实现 query hooks，并在页面接入时间范围状态。

**Step 4: Run tests**

- 运行：`npm run test -- timeRange`
- 预期：测试通过。

**Step 5: Run minimal build verification**

- 运行：`npm run build`
- 预期：构建通过。

**Step 6: Commit**

- 提交：`feat(web): 接入dashboard查询与时间预设`

### Task 3: Build Overview And Trend Section

**Files:**
- Create: `web/src/pages/Dashboard/components/OverviewCards.tsx`
- Create: `web/src/pages/Dashboard/components/CostTrendChart.tsx`
- Create: `web/src/pages/Dashboard/components/TokenTrendChart.tsx`
- Create: `web/src/utils/format.ts`
- Create: `web/src/utils/format.test.ts`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: Write the failing test**

- 为格式化函数编写测试，覆盖金额与时间轴标签格式化。

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- format`
- 预期：测试失败。

**Step 3: Write minimal implementation**

- 渲染四个总览卡片。
- 使用共享时间桶渲染费用趋势和 token 细分趋势。
- 加入 Tooltip 映射和空态。

**Step 4: Run tests**

- 运行：`npm run test -- format`
- 预期：测试通过。

**Step 5: Build verification**

- 运行：`npm run build`
- 预期：构建通过。

**Step 6: Commit**

- 提交：`feat(web): 落地dashboard总览与趋势区域`

### Task 4: Build Ranking Modules

**Files:**
- Create: `web/src/pages/Dashboard/components/ModelRanking.tsx`
- Create: `web/src/pages/Dashboard/components/ClientRanking.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Modify: `web/src/utils/format.ts`
- Modify: `web/src/utils/format.test.ts`

**Step 1: Write the failing test**

- 为模型展示名回退函数编写测试。

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- format`
- 预期：测试失败。

**Step 3: Write minimal implementation**

- 渲染模型排行与客户端排行模块。
- 模型名优先展示 `display_name`，缺失时回退 `model`。
- 排行为空时渲染稳定标题与空态。

**Step 4: Run tests**

- 运行：`npm run test -- format`
- 预期：测试通过。

**Step 5: Build verification**

- 运行：`npm run build`
- 预期：构建通过。

**Step 6: Commit**

- 提交：`feat(web): 落地dashboard排行模块`

### Task 5: Build Cache Benefit And Recent Requests

**Files:**
- Create: `web/src/pages/Dashboard/components/CacheAnalysis.tsx`
- Create: `web/src/pages/Dashboard/components/RecentRequestsTable.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Modify: `web/src/utils/format.ts`
- Modify: `web/src/utils/format.test.ts`

**Step 1: Write the failing test**

- 为最近请求表格列映射或关键格式化逻辑编写测试。

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- format`
- 预期：测试失败。

**Step 3: Write minimal implementation**

- 渲染缓存节省、缓存读取成本、缓存建设成本。
- 渲染最近请求表格，展示时间、模型、输入/输出 token、总费用、客户端标识。
- 保持空态和时间范围联动。

**Step 4: Run tests**

- 运行：`npm run test -- format`
- 预期：测试通过。

**Step 5: Build verification**

- 运行：`npm run build`
- 预期：构建通过。

**Step 6: Commit**

- 提交：`feat(web): 落地缓存效益与最近请求`

### Task 6: Polish States And Verify Formatting

**Files:**
- Create: `web/src/pages/Dashboard/DashboardPage.states.test.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Modify: `web/src/utils/format.ts`
- Modify: `web/src/utils/format.test.ts`

**Step 1: Write the failing test**

- 编写页面状态测试，覆盖 loading、error、empty、success 四种主状态。

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- DashboardPage.states`
- 预期：测试失败。

**Step 3: Write minimal implementation**

- 完善统一 loading/error/empty 展示。
- 收口金额、时间、长标签格式化。
- 修正边界数据下的布局稳定性。

**Step 4: Run tests**

- 运行：`npm run test`
- 预期：关键测试通过。

**Step 5: Final verification**

- 运行：`npm run build`
- 预期：构建通过。

**Step 6: Commit**

- 提交：`test(web): 完成dashboard状态与格式化验证`
