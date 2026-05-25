Status: done

# 编写 UI 集成测试

## What to build

在 server/internal/handler/ 下新增 ui_test.go，使用 httptest 写集成测试，验证路由优先级正确、SPA fallback 工作、静态文件能正确返回。

## Acceptance criteria

- [x] 测试验证 /api/v1/ping 不受 UI 路由拦截。
- [x] 测试验证 / 返回 index.html（使用测试 embed 文件）。
- [x] 测试验证不存在的路径返回 index.html。
- [x] 测试风格与现有 router_test.go 保持一致。

## Blocked by

- 01-setup-ui-handler-framework
- 02-implement-static-file-serving
