Status: completed

# Server 使用 go:embed 打包 Web 静态资源

## Problem Statement

当前 cc-status 项目的 web 前端和 server 后端是分离运行的：前端需要用 `npm run dev` 启动 Vite 开发服务器，后端需要用 `go run ./server/cmd/server` 启动 Gin 服务。这对于最终交付和演示来说不够方便——用户希望能直接运行一个二进制文件就同时启动后端 API 和前端页面，不需要分别启动两个服务，也不需要部署两个端口。

## Solution

在 server 中使用 Go 1.16+ 的 `go:embed` 功能，把 web 构建后的静态资源（`web/dist`）打包进 server 二进制文件。server 同时提供 API 服务和静态文件服务，访问根路径 `/` 直接返回前端页面，`/api/v1/*` 提供后端 API，`/healthz` 保持健康检查。开发阶段仍支持前后端分离运行（Vite 代理 `/api` 到后端），CI/CD 构建时才把前端产物打包进二进制。

## User Stories

1. 作为最终用户，我希望只运行一个二进制文件就能启动整个应用，以便部署和演示更简单。
2. 作为最终用户，我希望访问根路径 `/` 就能看到前端仪表盘，以便直接进入演示页面。
3. 作为开发者，我希望开发阶段仍能使用 Vite 的热更新功能，以便前端开发体验不受影响。
4. 作为开发者，我希望用一个 Makefile 命令就能完成全量构建，以便 CI/CD 流程更简单。
5. 作为维护者，我希望构建产物不提交到 Git，以便仓库保持整洁。
6. 作为测试者，我希望有简单的集成测试验证路由优先级正确，以便不会回归。
7. 作为用户，我希望前端路由（如 `/dashboard`）刷新时不会 404，以便单页应用能正常工作。
8. 作为用户，我希望 API 接口仍需要认证，以便敏感数据不被未授权访问。
9. 作为开发者，我希望静态文件服务逻辑封装在独立模块，以便代码结构清晰。
10. 作为跨平台用户，我希望构建脚本在 Windows 和 Linux/macOS 都能工作，以便不同环境都能正常构建。

## Implementation Decisions

- **静态资源来源**：web 构建产物默认输出到 `web/dist`，构建时复制到 `server/internal/ui/dist` 再 embed（避免 go:embed 引用父目录的限制）。
- **静态资源路由**：根路径 `/` 直接 serve 前端，优先级低于 `/healthz` 和 `/api/v1/*`。
- **SPA 路由 fallback**：使用 Gin 的 `NoRoute` + `http.FileServer` 组合，找不到静态文件时 fallback 回 `index.html`。
- **Embed 时机**：总是定义 embed.FS，但运行时检查目录是否为空；开发阶段目录为空则不注册 UI 路由，只在 CI/CD 构建时才复制文件进去。
- **构建流程**：根目录提供 Makefile，包含 `build-web`、`build-server`、`build` 三个主要目标；`build-web` 负责构建前端并复制到 server 目录。
- **UI Handler 组织**：在 `server/internal/handler/` 下新建 `ui.go`，封装静态文件服务逻辑，保持 `router.go` 简洁。
- **Git 忽略**：`server/.gitignore` 忽略 `internal/ui/dist/`，`web/.gitignore` 已有 `dist/` 忽略。
- **路由注册顺序**：先注册 `/healthz`，再注册 `/api/v1/*` 组，最后注册 UI 相关路由，保证 API 优先级高于 UI fallback。
- **认证策略**：UI 静态页面不需要认证（仅返回 HTML/CSS/JS，不包含敏感数据），API 接口保持现有的 Bearer Token 认证。
- **Vite 配置**：保持默认 `base` 配置（`./` 或 `/`），无需修改。
- **跨平台复制**：Makefile 检测操作系统，Windows 使用 PowerShell `Copy-Item`，Linux/macOS 使用 `cp -r`，复制前先清理旧目录。

## Testing Decisions

- 好的测试只验证外部可观察行为：路由优先级（API 不被 UI 覆盖）、SPA fallback 工作、静态文件能正确返回，不测试 embed 内部实现细节。
- 测试重点：`/api/v1/ping` 能正常返回（不被 UI 路由拦截）、`/` 返回 `index.html`、不存在的路径也返回 `index.html`。
- 测试策略：使用 Gin 的 `httptest` 包写集成测试，类似现有的 `router_test.go` 风格。
- 测试文件位置：在 `server/internal/handler/` 下新增 `ui_test.go`。

## Out of Scope

- 前端登录认证页面（首版保持写死 token 的演示模式）。
- 运行时动态重新加载前端资源（embed 是编译时确定的）。
- CDN 加速静态资源（首版只考虑二进制内联）。
- 复杂的构建缓存机制（每次构建都重新复制）。
- gzip 压缩静态资源（首版保持简单，后续可再加）。

## Further Notes

- 该 PRD 依赖 web 前端已经可以正常构建（web-dashboard 相关 issue 已完成）。
- 开发阶段的工作流保持不变：同时运行 Vite dev server 和 Go server，Vite 代理 `/api` 到后端。
- 最终交付物是一个单一的 Go 二进制文件，包含完整的后端 API 和前端页面。
- 后续如需优化，可考虑 gzip 压缩、构建缓存等改进，但首版优先保证功能可用和流程简单。
