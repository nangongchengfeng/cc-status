# CC Status Web

`cc-status/web` 是 Claude Code 使用量统计的前端仪表板，基于 React + Vite + Tailwind CSS + Recharts 构建。

## 技术栈

- **React 19**：UI 框架
- **Vite**：构建工具
- **TypeScript**：类型安全
- **Tailwind CSS**：原子化样式
- **Recharts**：图表库
- **TanStack Query**：服务端状态管理
- **React Router**：路由
- **Axios**：HTTP 客户端
- **Vitest**：测试框架

## 项目架构

```
src/
├── pages/              # 页面组件（路由级）
│   └── Dashboard/     # 仪表板页面
│       ├── components/ # 页面专属子组件
│       └── DashboardPage.tsx
├── components/         # 可复用 UI 组件
├── hooks/              # 自定义 hooks
├── api/                # API 请求函数
├── types/              # TypeScript 类型定义
├── utils/              # 工具函数
├── constants/          # 常量
└── styles/             # 全局样式
```

## 快速开始

### 安装依赖

在仓库根目录或 web 目录执行：

```bash
cd web
npm install
```

### 开发模式

确保 server 服务已在 `http://127.0.0.1:8080` 启动，然后：

```bash
npm run dev
```

访问 `http://localhost:5173` 查看应用。

Vite 会自动代理 `/api` 请求到后端服务。

### 构建生产版本

```bash
npm run build
```

构建产物输出到 `dist/` 目录。

### 运行测试

```bash
npm run test
```

## 功能特性

- **概览卡片**：总费用、总 token、请求数等核心指标
- **费用趋势图**：按天/周/月展示费用变化
- **Token 趋势图**：输入/输出 token 变化趋势
- **模型排行**：按费用/使用量排序的模型列表
- **客户端排行**：按费用/使用量排序的客户端列表
- **缓存分析**：缓存读写 token 占比
- **最近请求**：原始使用记录列表

## 环境变量

可以在项目根目录创建 `.env` 或 `.env.local` 文件：

```bash
# API 基础地址（默认 /api/v1）
VITE_API_BASE_URL=/api/v1

# 开发服务器端口（默认 5173）
VITE_PORT=5173
```

## API 对接

前端通过 `/api/v1` 与后端 Gin 服务对接：

- `GET /api/v1/stats/dashboard` - 仪表盘数据
- `GET /api/v1/stats/overview` - 总览数据
- `GET /api/v1/stats/trend` - 趋势数据
- `GET /api/v1/logs` - 使用记录列表
- `GET /api/v1/model-pricings` - 模型定价

所有需要认证的接口使用 Bearer Token，通过 `Authorization` 头传递。

## 开发配置

### 代理配置

`vite.config.ts` 中已配置开发时代理：

```typescript
proxy: {
  '/api': {
    target: 'http://127.0.0.1:8080',
    changeOrigin: true,
  },
}
```

### 路径别名

`@` 别名指向 `src/` 目录：

```typescript
import { formatCost } from '@/utils/format';
```

## 与 Server 联调

1. 启动 server 服务：

```bash
# 在仓库根目录
export CC_USAGE_SERVER_AUTH_TOKEN="dev-token"
go run ./server/cmd/server
```

2. 启动 web 开发服务器：

```bash
# 在另一个终端
cd web
npm run dev
```

3. 访问 `http://localhost:5173`

## 开发验证

```bash
# 类型检查
tsc -b

# 运行测试
npm run test

# 构建检查
npm run build
```
