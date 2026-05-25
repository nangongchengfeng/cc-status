Status: done

# 搭建 UI handler 框架与路由注册

## What to build

在 server/internal/handler/ 下新建 ui.go，封装静态文件服务的框架，并且在 router.go 中预留 UI 路由注册位置。该切片不要求真正 serve 静态文件，只需要把框架搭好，把路由注册顺序理清楚（先注册 /healthz，再注册 /api/v1/*，最后预留 UI 位置）。

## Acceptance criteria

- [x] server/internal/handler/ui.go 已创建，包含基础的 UI handler 封装结构。
- [x] router.go 中路由注册顺序正确：/healthz -> /api/v1/* -> UI 路由（目前预留位置）。
- [x] 现有 API 接口不受影响，/healthz 和 /api/v1/ping 仍能正常访问。

## Blocked by

None - can start immediately
