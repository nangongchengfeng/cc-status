# React 前端项目

React + Vite + Tailwind CSS + Recharts，功能分层架构，对接 Gin RESTful API。

## 项目架构

```
src/
├── pages/              # 页面组件（路由级）
├── components/         # 可复用 UI 组件
│   └── common/         # 通用基础组件
├── hooks/              # 自定义 hooks
├── api/                # API 请求函数（按资源分文件）
├── stores/             # Zustand 状态管理
├── types/              # TypeScript 类型定义
├── utils/              # 工具函数
├── constants/          # 常量
└── styles/             # 全局样式、主题覆盖
```

**分层约束：**
- `pages/` 对应路由，每个页面一个目录（含页面组件 + 子组件）
- `components/` 只放跨页面复用组件，页面专属组件放页面目录内
- `api/` 按后端资源分文件（如 `user.ts`、`order.ts`），与后端 handler 对应
- `types/` 中前端类型与后端 DTO 对齐，但不直接复制后端结构

## API 对接

**后端响应格式对接：**
- 成功响应 `{ "data": { ... } }`：axios 拦截器提取 `data`，TanStack Query 拿到已解包数据
- 分页响应 `{ "data": [...], "total": 100, "page": 1, "size": 20 }`：保留完整结构

**后端错误格式对接：**
- 错误响应 `{ "code": "NOT_FOUND", "message": "用户不存在" }`
- 响应拦截器捕获非 2xx，用 `message.error()` 展示 `message`
- `code` 用于特殊处理（如 `UNAUTHORIZED` 跳转登录）

**API 函数：**
- `api/` 下每个文件对应一个后端资源
- 函数命名：`getXxx` / `createXxx` / `updateXxx` / `deleteXxx`
- 函数只发请求返回数据，不含业务逻辑
- 基础 URL 从 `VITE_API_BASE_URL` 读取，默认 `/api/v1`

## TanStack Query

- 每个 API 函数对应一个自定义 hook（`hooks/useUsers.ts` 等）
- mutation 用 `useMutation` + `onSuccess` 使相关 query 失效
- query key 数组格式：`['users']` / `['users', id]` / `['users', { page, size }]`
- 全局默认：staleTime 5 分钟，retry 1 次

## Tailwind CSS

- 优先使用 Tailwind 原子类完成布局、间距、颜色、边框和响应式适配
- 主题变量放在 `tailwind.config.*` 或全局样式中统一维护，避免在组件内散落魔法值
- 大屏页面优先使用深色主题、卡片化布局和自由网格，不强行套用后台管理台风格
- 仅在明显复用三次以上时再抽取样式组件，避免过早封装

## Recharts

- 图表统一使用 Recharts，不混用多套图表库
- 趋势图优先 `LineChart`、`AreaChart`
- 排行优先 `BarChart`
- 占比优先 `PieChart`
- Tooltip、Legend、坐标轴格式化逻辑放在页面目录或 `utils/` 中复用

**组件编写：**
- 函数组件 + hooks，禁止 class 组件
- Props 类型用 `interface` 定义，与组件同文件导出
- 文件名 PascalCase（如 `UserList.tsx`）
- 事件处理器 `handleXxx` 命名
- 页面级大屏组件允许在页面目录内拆分子组件，跨页面复用组件再放入 `components/`

**筛选器：**
- 使用原生表单元素或轻量自定义组件，不为简单筛选引入重型表单框架
- 时间范围等展示层预设放 `constants/`，由页面转换为 API 查询参数

**表格：**
- 使用语义化 `table` 或轻量表格组件实现最近请求列表
- 分页参数与后端对齐：`offset` + `limit`
- 列定义抽为常量数组，不内联 JSX

## Zustand

- 全局状态放 `stores/`，每个 store 一个文件
- 命名 `useXxxStore`
- 只管客户端 UI 状态（认证、UI 切换等），服务端状态由 TanStack Query 管理
- 用 `create` + 简单对象，不使用 middleware

## 路由

- React Router v6，定义集中在 `src/router.tsx`
- 页面用 `React.lazy` + `Suspense` 懒加载
- 路由守卫通过 wrapper 组件（如 `<RequireAuth>`）
- 路径风格：小写 + 连字符（与后端 URL 风格一致）

## 测试

- Vitest + React Testing Library
- 测试文件同目录，`*.test.tsx`
- 重点：自定义 hooks、工具函数、组件交互
- 不测试图表库和样式框架本身
- 运行：`npm run test`

## 构建与开发

- 开发端口从 `VITE_PORT` 配置，默认 5173
- Vite `server.proxy` 代理 `/api` 到后端地址
- 构建：`npm run build`
- 环境变量前缀 `VITE_`

## 代码规范

- ESLint + Prettier，提交前自动格式化
- TypeScript strict 模式
- 导入顺序：React → 第三方库 → @/ 别名 → 相对路径
