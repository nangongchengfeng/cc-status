# Dashboard Redesign Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 将 Dashboard 首页重构为“浅蓝指挥舱”风格的大屏，在不改数据接口和查询逻辑的前提下，显著提升首屏层次、布局张力和汇报观感。

**Architecture:** 以 `DashboardPage.tsx` 为总编排入口，保留现有 query、类型和数据流，只重构页面网格、主舞台、控制塔和各模块容器样式。组件内部尽量复用既有数据契约，通过小步测试锁定页头文案、状态提示和核心空态，再逐块替换视觉与布局。

**Tech Stack:** React 19, TypeScript, Vite, Tailwind CSS v4, Recharts, TanStack Query, Vitest, React Testing Library

---

### Task 1: 锁住新页头与布局预期

**Files:**
- Modify: `web/src/pages/Dashboard/DashboardPage.test.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.states.test.tsx`
- Check: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: 写失败测试，改成新的页头文案和布局语义**

```tsx
expect(screen.getByRole('heading', { name: 'Claude 用量指挥舱' })).toBeInTheDocument();
expect(screen.getByText('先看花费，再找来源。')).toBeInTheDocument();
expect(screen.getByText('时间范围')).toBeInTheDocument();
```

**Step 2: 运行页面测试，确认旧实现失败**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx
```

Expected:

```text
FAIL
Unable to find role "heading" with name "Claude 用量指挥舱"
```

**Step 3: 只做最小页面文案调整让测试靠近目标**

```tsx
<h1>Claude 用量指挥舱</h1>
<p>先看花费，再找来源。</p>
```

**Step 4: 重新运行页面测试，确认断言通过**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx
```

Expected:

```text
PASS
```

**Step 5: 提交**

```bash
git add web/src/pages/Dashboard/DashboardPage.test.tsx web/src/pages/Dashboard/DashboardPage.states.test.tsx web/src/pages/Dashboard/DashboardPage.tsx
git commit -m "test(web-dashboard): 锁定新版页头预期"
```

提交前先征得用户确认。

### Task 2: 重建 Dashboard 页面骨架

**Files:**
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Check: `web/src/constants/timeRanges.ts`
- Check: `web/src/utils/timeRange.ts`

**Step 1: 写失败测试，补充主舞台和控制塔的关键结构**

```tsx
expect(screen.getByText('当前按天粒度观察')).toBeInTheDocument();
expect(screen.getByText('费用是这页的主线。')).toBeInTheDocument();
```

**Step 2: 运行页面状态测试，确认新结构尚未出现**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx
```

Expected:

```text
FAIL
Unable to find text "当前按天粒度观察"
```

**Step 3: 重写页面网格和首屏分区**

```tsx
<main className="min-h-screen px-6 py-6">
  <div className="mx-auto grid max-w-[1720px] gap-6 xl:grid-cols-12">
    <section className="xl:col-span-8">{/* 主舞台 */}</section>
    <aside className="xl:col-span-4">{/* 控制塔 */}</aside>
  </div>
</main>
```

**Step 4: 运行页面测试，确认新骨架稳定**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx
```

Expected:

```text
PASS
```

**Step 5: 提交**

```bash
git add web/src/pages/Dashboard/DashboardPage.tsx web/src/pages/Dashboard/DashboardPage.test.tsx web/src/pages/Dashboard/DashboardPage.states.test.tsx
git commit -m "feat(web-dashboard): 重建大屏首屏骨架"
```

提交前先征得用户确认。

### Task 3: 升级全局背景与主题基调

**Files:**
- Modify: `web/src/styles/index.css`
- Check: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: 写最小回归测试或人工检查点，锁定浅色科技背景目标**

```text
人工检查点：
1. 根背景不再是暖黑色
2. 页面出现浅蓝雾光
3. 页面仍可正常渲染
```

**Step 2: 运行构建，记录当前样式基线**

Run:

```bash
npm run build
```

Expected:

```text
build completed successfully
```

**Step 3: 改写根背景、网格纹理和基础字体色**

```css
:root {
  color: #12304d;
  background:
    radial-gradient(circle at top left, rgba(108, 184, 255, 0.28), transparent 32%),
    radial-gradient(circle at right center, rgba(185, 226, 255, 0.42), transparent 24%),
    linear-gradient(135deg, #f7fbff 0%, #edf5fb 48%, #f4f8fc 100%);
}
```

**Step 4: 重新构建，确认样式改动不破坏编译**

Run:

```bash
npm run build
```

Expected:

```text
build completed successfully
```

**Step 5: 提交**

```bash
git add web/src/styles/index.css
git commit -m "style(web-dashboard): 切换浅蓝指挥舱底色"
```

提交前先征得用户确认。

### Task 4: 重做主指标卡与主舞台信息层级

**Files:**
- Modify: `web/src/pages/Dashboard/components/OverviewCards.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Check: `web/src/utils/format.ts`

**Step 1: 写失败测试，锁定主指标区域的新标签或说明**

```tsx
expect(screen.getByText('总费用')).toBeInTheDocument();
expect(screen.getByText('费用是这页的主线。')).toBeInTheDocument();
```

**Step 2: 运行测试，确认当前页面没有新提示语**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx
```

Expected:

```text
FAIL
Unable to find text "费用是这页的主线。"
```

**Step 3: 重构概览卡为主次组合**

```tsx
<section className="grid gap-4 xl:grid-cols-[1.4fr_0.9fr_0.9fr]">
  <article>{/* 总费用主卡 */}</article>
  <article>{/* 总 Token */}</article>
  <article>{/* 次级数据堆叠 */}</article>
</section>
```

**Step 4: 运行页面测试，确认主指标区可渲染**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx
```

Expected:

```text
PASS
```

**Step 5: 提交**

```bash
git add web/src/pages/Dashboard/components/OverviewCards.tsx web/src/pages/Dashboard/DashboardPage.tsx web/src/pages/Dashboard/DashboardPage.test.tsx
git commit -m "feat(web-dashboard): 强化主舞台指标层级"
```

提交前先征得用户确认。

### Task 5: 重做趋势图容器与浅蓝配色

**Files:**
- Modify: `web/src/pages/Dashboard/components/CostTrendChart.tsx`
- Modify: `web/src/pages/Dashboard/components/TokenTrendChart.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: 写失败测试，锁定新的趋势区说明文字**

```tsx
expect(screen.getByText('费用趋势')).toBeInTheDocument();
expect(screen.getByText('Token 轨迹')).toBeInTheDocument();
```

**Step 2: 运行测试，确认副标题需要更新**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx
```

Expected:

```text
FAIL
```

**Step 3: 修改趋势图外层容器和图表主题**

```tsx
<div className="rounded-[28px] border border-white/70 bg-white/55 p-5 shadow-[0_24px_80px_rgba(84,145,204,0.14)] backdrop-blur-xl">
  <ResponsiveContainer>{/* Recharts */}</ResponsiveContainer>
</div>
```

```tsx
<Line stroke="#3f8cff" />
<Area fill="#83c8ff" />
```

**Step 4: 运行页面测试与构建，确认图表模块稳定**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx
npm run build
```

Expected:

```text
PASS
build completed successfully
```

**Step 5: 提交**

```bash
git add web/src/pages/Dashboard/components/CostTrendChart.tsx web/src/pages/Dashboard/components/TokenTrendChart.tsx web/src/pages/Dashboard/DashboardPage.tsx
git commit -m "style(web-dashboard): 升级趋势图科技感"
```

提交前先征得用户确认。

### Task 6: 统一重做排行、缓存效益和最近请求

**Files:**
- Modify: `web/src/pages/Dashboard/components/ModelRanking.tsx`
- Modify: `web/src/pages/Dashboard/components/ClientRanking.tsx`
- Modify: `web/src/pages/Dashboard/components/CacheAnalysis.tsx`
- Modify: `web/src/pages/Dashboard/components/RecentRequestsTable.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: 写失败测试，确认空态和成功态文案仍可读**

```tsx
expect(screen.getByText('当前时间范围还没有缓存数据。')).toBeInTheDocument();
expect(screen.getByText('当前时间范围还没有最近请求。')).toBeInTheDocument();
```

**Step 2: 运行状态测试，确保现有空态仍被覆盖**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx
```

Expected:

```text
PASS
```

**Step 3: 统一替换底部模块容器和色板**

```tsx
<section className="rounded-[28px] border border-white/70 bg-white/60 p-5 shadow-[0_18px_64px_rgba(84,145,204,0.1)]">
  {/* 排行 / 缓存 / 表格 */}
</section>
```

**Step 4: 运行状态测试与构建，确认底部区域稳定**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx
npm run build
```

Expected:

```text
PASS
build completed successfully
```

**Step 5: 提交**

```bash
git add web/src/pages/Dashboard/components/ModelRanking.tsx web/src/pages/Dashboard/components/ClientRanking.tsx web/src/pages/Dashboard/components/CacheAnalysis.tsx web/src/pages/Dashboard/components/RecentRequestsTable.tsx web/src/pages/Dashboard/DashboardPage.tsx
git commit -m "style(web-dashboard): 统一底部模块层次"
```

提交前先征得用户确认。

### Task 7: 收口验证与诊断

**Files:**
- Check: `web/src/pages/Dashboard/DashboardPage.tsx`
- Check: `web/src/pages/Dashboard/components/OverviewCards.tsx`
- Check: `web/src/pages/Dashboard/components/CostTrendChart.tsx`
- Check: `web/src/pages/Dashboard/components/TokenTrendChart.tsx`
- Check: `web/src/pages/Dashboard/components/ModelRanking.tsx`
- Check: `web/src/pages/Dashboard/components/ClientRanking.tsx`
- Check: `web/src/pages/Dashboard/components/CacheAnalysis.tsx`
- Check: `web/src/pages/Dashboard/components/RecentRequestsTable.tsx`
- Check: `web/src/styles/index.css`

**Step 1: 获取诊断，清理本次改动引入的 TS/JSX 问题**

```text
检查 recently edited files 的 diagnostics。
```

**Step 2: 运行最小必要回归**

Run:

```bash
npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx src/utils/format.test.ts
npm run build
```

Expected:

```text
PASS
build completed successfully
```

**Step 3: 人工检查关键结果**

```text
1. 首屏主舞台明显强于下方模块
2. 费用趋势比 Token 趋势更突出
3. 时间切换器选中态明显
4. 整页配色稳定为浅蓝、白、灰
5. loading / error / empty / success 四态仍稳定
```

**Step 4: 整理变更说明**

```text
记录布局变化、测试结果、未覆盖的视觉细节。
```

**Step 5: 提交**

```bash
git add web/src/pages/Dashboard web/src/styles/index.css
git commit -m "feat(web-dashboard): 重构浅蓝指挥舱大屏"
```

提交前先征得用户确认。
