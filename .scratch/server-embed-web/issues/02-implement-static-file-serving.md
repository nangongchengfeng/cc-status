Status: done

# 实现静态文件服务与 SPA fallback

## What to build

在 ui.go 中实现完整的静态文件服务逻辑，使用 Gin 的 http.FileServer 和 NoRoute 机制，支持 SPA 路由 fallback 到 index.html。该切片先用本地文件系统测试（暂不 embed），验证路由优先级和 fallback 行为正确。

## Acceptance criteria

- [x] 访问根路径 / 返回 index.html（假设在 internal/ui/dist 下有测试文件）。
- [x] 访问不存在的路径也返回 index.html（SPA fallback）。
- [x] 访问 /api/v1/ping 仍正常返回，不被 UI 路由拦截。
- [x] 访问 /healthz 仍正常返回。

## Blocked by

- 01-setup-ui-handler-framework
