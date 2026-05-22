# Dashboard Static Replica Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在 `web/` 中把首页 `DashboardPage` 重构为一版按参考图 1:1 静态高保真复刻的大屏页面，并保持可测试、可构建。

**Architecture:** 保留根路由 `/` 和 `DashboardPage` 入口不变，页面第一阶段改为由本地 mock 数据驱动。展示层按视觉区块拆分为页头、核心指标、综合效益、趋势区、排行区、最近请求和底部信息条，复用现有 React、Tailwind CSS 和 Recharts 能力，不接真实查询逻辑。

**Tech Stack:** React 19、TypeScript 5、Vite、Tailwind CSS、Recharts、Vitest、React Testing Library

---

### Task 1: 建立静态复刻的测试锚点和页面骨架

**Files:**
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.test.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.states.test.tsx`

**Step 1: Write the failing test**

- 在 `DashboardPage.test.tsx` 中把断言从旧版文案调整为新版静态页关键锚点：
- 标题 `Claude 用量看板`
- 模块标题 `核心指标`
- 模块标题 `综合效益`
- 模块标题 `最近请求`

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：旧页面无法满足新版标题和结构断言，测试失败。

**Step 3: Write minimal implementation**

- 暂时移除页面对 `useDashboardQuery()`、`useRecentLogsQuery()` 和时间预设切换器的直接依赖。
- 在 `DashboardPage.tsx` 中先搭出新版静态布局骨架：
- 顶部标题区
- 左侧核心指标面板
- 右侧综合效益面板
- 左下趋势与排行区
- 右下最近请求面板
- 底部信息栏

**Step 4: Run test to verify it passes**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：至少新版锚点文案测试通过。

**Step 5: Commit**

- 提交：`refactor(web): 重构dashboard静态复刻骨架`

### Task 2: 创建静态 mock 数据与页面区块组件

**Files:**
- Create: `web/src/pages/Dashboard/dashboardReplica.mock.ts`
- Create: `web/src/pages/Dashboard/components/DashboardHero.tsx`
- Create: `web/src/pages/Dashboard/components/MetricsPanel.tsx`
- Create: `web/src/pages/Dashboard/components/BenefitPanel.tsx`
- Create: `web/src/pages/Dashboard/components/DashboardFooterMeta.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: Write the failing test**

- 在 `DashboardPage.states.test.tsx` 中新增静态页面断言：
- 更新时间胶囊存在
- 四个核心指标卡标题存在
- 三个综合效益卡标题存在

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：因为 mock 数据文件和新组件未创建，测试失败。

**Step 3: Write minimal implementation**

- 在 `dashboardReplica.mock.ts` 中定义：
- `updatedAtLabel`
- `coreMetrics`
- `benefitMetrics`
- `costTrend`
- `tokenTrend`
- `modelRanking`
- `clientRanking`
- `recentRequests`
- `footerMeta`
- 创建 `DashboardHero`、`MetricsPanel`、`BenefitPanel`、`DashboardFooterMeta` 四个组件。
- `DashboardPage.tsx` 改为只负责装配这些组件和布局容器。

**Step 4: Run test to verify it passes**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：静态核心信息渲染稳定通过。

**Step 5: Commit**

- 提交：`feat(web): 补齐dashboard静态复刻基础区块`

### Task 3: 重做全局背景与顶部装饰层

**Files:**
- Modify: `web/src/styles/index.css`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Modify: `web/src/pages/Dashboard/components/DashboardHero.tsx`

**Step 1: Write the failing test**

- 为 `DashboardHero` 增加一个组件测试，断言：
- 主标题存在
- 副标题存在
- 更新时间胶囊存在
- 装饰层容器带有可识别的 `aria-hidden` 或 `data-testid`

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.test.tsx`
- 预期：顶部装饰层尚未按新结构实现，测试失败。

**Step 3: Write minimal implementation**

- 在 `index.css` 中重设页面背景：
- 深蓝黑渐变底
- 径向光晕
- 细噪点覆盖
- 在 `DashboardHero.tsx` 中加入：
- 发光球体
- 横向粒子波浪
- 右上状态胶囊
- 副标题与标题的对齐布局

**Step 4: Run test to verify it passes**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.test.tsx`
- 预期：页头结构测试通过。

**Step 5: Commit**

- 提交：`style(web): 还原dashboard顶部视觉层`

### Task 4: 复刻核心指标与综合效益面板

**Files:**
- Modify: `web/src/pages/Dashboard/components/MetricsPanel.tsx`
- Modify: `web/src/pages/Dashboard/components/BenefitPanel.tsx`
- Modify: `web/src/utils/format.ts`
- Modify: `web/src/utils/format.test.ts`

**Step 1: Write the failing test**

- 在 `format.test.ts` 中补充静态页面格式化断言：
- 大数字千分位格式
- 金额两位小数
- 百分比带正负号文案

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- src/utils/format.test.ts`
- 预期：缺少页面复刻所需的涨跌幅文案格式化，测试失败。

**Step 3: Write minimal implementation**

- 为指标卡增加：
- 彩色圆形图标底盘
- 主值
- 单位
- 辅助标签
- 环比文案
- 为综合效益卡增加：
- 顶部说明
- 三列指标
- 黄/绿/蓝三色高亮数字
- 统一玻璃卡、发光、边框和间距

**Step 4: Run test to verify it passes**

- 运行：`npm run test -- src/utils/format.test.ts src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：格式化和面板渲染测试通过。

**Step 5: Commit**

- 提交：`feat(web): 复刻dashboard指标与综合效益面板`

### Task 5: 复刻趋势分析与双榜单区域

**Files:**
- Create: `web/src/pages/Dashboard/components/TrendPanel.tsx`
- Create: `web/src/pages/Dashboard/components/RankingPanel.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: Write the failing test**

- 在 `DashboardPage.states.test.tsx` 中增加断言：
- `费用趋势（USD）`
- `Token 趋势（万）`
- `模型排行`
- `客户端排行`

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：趋势区和榜单区标题与结构尚未完成，测试失败。

**Step 3: Write minimal implementation**

- `TrendPanel.tsx` 内使用 Recharts 复刻：
- 左侧费用折线图
- 右侧 Token 双面积图
- 两个面板右上角下拉样式胶囊
- `RankingPanel.tsx` 内复刻：
- 左侧紫色模型排行横条
- 右侧蓝色客户端排行横条
- TOP 5 序号、名称、数值布局
- 保持图例和轴线接近参考图密度

**Step 4: Run test to verify it passes**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：关键标题与区块稳定渲染。

**Step 5: Commit**

- 提交：`feat(web): 复刻dashboard趋势与排行区域`

### Task 6: 复刻最近请求表格和底部信息栏

**Files:**
- Create: `web/src/pages/Dashboard/components/RequestsPanel.tsx`
- Modify: `web/src/pages/Dashboard/components/DashboardFooterMeta.tsx`
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`

**Step 1: Write the failing test**

- 在 `DashboardPage.states.test.tsx` 中增加断言：
- 表头 `时间`
- 表头 `模型`
- 表头 `输入`
- 表头 `输出`
- 表头 `费用（USD）`
- 按钮 `查看全部`

**Step 2: Run test to verify it fails**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：最近请求表格与按钮结构未完成，测试失败。

**Step 3: Write minimal implementation**

- 使用 mock 数据创建高密度深色表格。
- 表头、行高、分割线、费用高亮色、长模型名和客户端名截断样式按参考图微调。
- 在底部信息栏加入左侧提示图标文本和右侧数据来源信息。

**Step 4: Run test to verify it passes**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.states.test.tsx`
- 预期：表格和底部信息栏断言通过。

**Step 5: Commit**

- 提交：`feat(web): 复刻dashboard请求表格与底部信息栏`

### Task 7: 统一微调并完成最终回归

**Files:**
- Modify: `web/src/pages/Dashboard/DashboardPage.tsx`
- Modify: `web/src/pages/Dashboard/components/DashboardHero.tsx`
- Modify: `web/src/pages/Dashboard/components/MetricsPanel.tsx`
- Modify: `web/src/pages/Dashboard/components/BenefitPanel.tsx`
- Modify: `web/src/pages/Dashboard/components/TrendPanel.tsx`
- Modify: `web/src/pages/Dashboard/components/RankingPanel.tsx`
- Modify: `web/src/pages/Dashboard/components/RequestsPanel.tsx`
- Modify: `web/src/pages/Dashboard/components/DashboardFooterMeta.tsx`
- Modify: `web/src/styles/index.css`

**Step 1: Write the failing test**

- 不新增低价值测试，先整理现有测试断言是否覆盖关键静态锚点。
- 如缺口明显，仅补充一个最小断言，验证底部数据来源文案存在。

**Step 2: Run test to verify current gap**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx src/utils/format.test.ts`
- 预期：若有最终文案或结构未收口，出现失败项。

**Step 3: Write minimal implementation**

- 调整所有模块的：
- 间距
- 边框透明度
- 发光阴影
- 字号层级
- 表格密度
- 图表高度
- 容器宽度比例
- 确保桌面宽屏下整体接近参考图比例。

**Step 4: Run tests**

- 运行：`npm run test -- src/pages/Dashboard/DashboardPage.test.tsx src/pages/Dashboard/DashboardPage.states.test.tsx src/utils/format.test.ts`
- 预期：全部通过。

**Step 5: Final verification**

- 运行：`npm run build`
- 预期：构建通过。

**Step 6: Commit**

- 提交：`style(web): 完成dashboard静态复刻微调`
